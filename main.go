package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nikolaihg/deadlink-scraper-go/linktype"
	"github.com/nikolaihg/deadlink-scraper-go/set"
	"golang.org/x/net/html"
)

type LinkStats struct {
	Total        int
	Internal     int
	External     int
	Alive        int
	Dead         int
	Skipped      int
	ByStatusCode map[string]int
}

var (
	visitedMu sync.Mutex
	checkedMu sync.Mutex
	statsMu   sync.Mutex
)

func main() {
	concurrency := flag.Int("concurrency", 10, "number of concurrent workers")
	timeout := flag.Duration("timeout", 5*time.Second, "HTTP request timeout (e.g. 5s, 1m)")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalf("Usage: go run . <start-url> [-concurrency N] [-timeout duration]")
	}
	startURL := flag.Arg(0)

	u, err := url.Parse(startURL)
	if err != nil {
		log.Fatalf("Error parsing start URL: %v", err)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	baseDomain := u.Host
	startURL = u.String()

	client := &http.Client{Timeout: *timeout}

	visited := set.New()
	checked := set.New()
	stats := &LinkStats{ByStatusCode: make(map[string]int)}

	// job channel & waitgroups
	jobs := make(chan linktype.Link)
	var taskWg sync.WaitGroup   // counts outstanding links to process
	var workerWg sync.WaitGroup // waits for workers to exit

	// start N workers
	for i := 0; i < *concurrency; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for link := range jobs {
				// crawl this link, possibly enqueueing more
				crawl(client, link, baseDomain, visited, checked, stats, jobs, &taskWg)
				taskWg.Done()
			}
		}()
	}

	// seed first job
	startLink := linktype.Link{URL: startURL, Type: linktype.InternalLink}
	visitedMu.Lock()
	visited.Add(startLink)
	visitedMu.Unlock()

	taskWg.Add(1)
	jobs <- startLink

	// once all tasks are done, close jobs to shut down workers
	go func() {
		taskWg.Wait()
		close(jobs)
	}()

	// wait for workers to finish
	workerWg.Wait()

	// done
	printStats(*stats)
	log.Printf("Links visited: %v\n", visited.Values())
}

func crawl(client *http.Client, link linktype.Link, baseDomain string, visited, checked *set.Set, stats *LinkStats, jobs chan<- linktype.Link, taskWg *sync.WaitGroup) {
	checkedMu.Lock()
	if !checked.Contains(link) {
		checked.Add(link)
		validateLink(client, link, stats)
	}
	checkedMu.Unlock()

	log.Printf("[Crawling]: %s\n", link.URL)
	if link.Type != linktype.InternalLink {
		return
	}

	// fetch page
	resp, err := client.Get(link.URL)
	if err != nil {
		log.Printf("Error fetching %s: %v\n", link.URL, err)
		return
	}
	defer resp.Body.Close()

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML: %v\n", err)
		return
	}

	baseURL := resp.Request.URL
	linksMap := extractLinks(root, baseURL)

	for _, newLink := range linksMap {
		switch newLink.Type {
		case linktype.InternalLink, linktype.PageLink:
			visitedMu.Lock()
			already := visited.Contains(newLink)
			if !already {
				visited.Add(newLink)
				taskWg.Add(1)   // count the new job
				jobs <- newLink // send to workers
			}
			visitedMu.Unlock()
		case linktype.ExternalLink:
			checkedMu.Lock()
			if !checked.Contains(newLink) {
				checked.Add(newLink)
				validateLink(client, newLink, stats)
			}
			checkedMu.Unlock()
		}
	}
}

func validateLink(client *http.Client, link linktype.Link, stats *LinkStats) {
	statsMu.Lock()
	defer statsMu.Unlock()

	stats.Total++
	switch link.Type {
	case linktype.PageLink:
		stats.Skipped++
		log.Printf("[SKIP]   %s (Page link)\n", link.URL)
		return
	case linktype.InternalLink:
		stats.Internal++
	case linktype.ExternalLink:
		stats.External++
	default:
		stats.Skipped++
		log.Printf("[SKIP]   %s (Unknown type)\n", link.URL)
		return
	}
	if link.URL == "" {
		stats.Dead++
		log.Printf("[DEAD]   (empty URL)\n")
		return
	}
	status, code, err := fetchStatus(client, link.URL)
	if err != nil {
		stats.Dead++
		log.Printf("[DEAD]   %s (%v)\n", link.URL, err)
		return
	}
	codeStr := strconv.Itoa(code)
	stats.ByStatusCode[codeStr]++
	if code >= 400 {
		stats.Dead++
		log.Printf("[DEAD]   %s (%s)\n", link.URL, status)
	} else {
		stats.Alive++
		log.Printf("[ALIVE]  %s (%s)\n", link.URL, status)
	}
}

func extractLinks(node *html.Node, baseURL *url.URL) map[string]linktype.Link {
	links := make(map[string]linktype.Link)

	tagAttr := map[string]string{
		"a":      "href",
		"link":   "href",
		"img":    "src",
		"script": "src",
		"iframe": "src",
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if attr, ok := tagAttr[n.Data]; ok {
				for _, a := range n.Attr {
					if a.Key == attr {
						link := filterLink(a.Val, baseURL)
						if link.URL != "" {
							if _, exists := links[link.URL]; !exists {
								links[link.URL] = link
							}
						}
						break // no need to keep looping attrs
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(node)
	return links
}

func filterLink(href string, baseURL *url.URL) linktype.Link {
	if href == "" {
		return linktype.Link{}
	}

	// Skip non-HTTP links
	switch {
	case strings.HasPrefix(href, "mailto:"),
		strings.HasPrefix(href, "tel:"),
		strings.HasPrefix(href, "javascript:"),
		strings.HasPrefix(href, "ftp:"):
		return linktype.Link{}
	case strings.HasPrefix(href, "#"):
		return linktype.Link{
			URL:  baseURL.String() + href,
			Type: linktype.PageLink,
		}
	}

	// Handle page links
	if strings.HasPrefix(href, "#") {
		return linktype.Link{
			URL:  baseURL.String() + href,
			Type: linktype.PageLink,
		}
	}

	normalized, err := normalizeURL(baseURL, href)
	if err != nil {
		log.Printf("Skipping invalid URL: %s (%v)", href, err)
		return linktype.Link{}
	}

	var linkType linktype.LinkType
	parsed, _ := url.Parse(normalized)
	if parsed.Host == baseURL.Host {
		linkType = linktype.InternalLink
	} else {
		linkType = linktype.ExternalLink
	}

	return linktype.Link{
		URL:  normalized,
		Type: linkType,
	}
}

func normalizeURL(base *url.URL, raw string) (string, error) {
	u, err := base.Parse(raw)
	if err != nil {
		return "", err
	}
	u.Fragment = ""
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New("unsupported scheme")
	}
	u.Host = strings.ToLower(u.Host)
	if (u.Scheme == "http" && u.Port() == "80") || (u.Scheme == "https" && u.Port() == "443") {
		u.Host = u.Hostname()
	}
	u.Path = path.Clean((u.Path))
	return u.String(), nil
}

func fetchStatus(client *http.Client, url string) (string, int, error) {
	// Try HEAD request first, because it is faster
	resp, err := client.Head(url)
	if err == nil {
		defer resp.Body.Close()
		return resp.Status, resp.StatusCode, nil
	}

	// Fallback to GET if HEAD fails
	resp, err = client.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	return resp.Status, resp.StatusCode, nil
}

func printStats(stats LinkStats) {
	log.Println("Scan complete:")
	log.Printf("Total:    %d\n", stats.Total)
	log.Printf("Internal: %d\n", stats.Internal)
	log.Printf("External: %d\n", stats.External)
	log.Printf("Alive:    %d\n", stats.Alive)
	log.Printf("Dead:     %d\n", stats.Dead)
	log.Printf("Skipped:  %d\n", stats.Skipped)
	log.Println("Status codes distribution:")
	for code, count := range stats.ByStatusCode {
		log.Printf("  %s: %d\n", code, count)
	}
}
