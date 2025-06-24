package main

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/nikolaihg/deadlink-scraper-go/linktype"
	"github.com/nikolaihg/deadlink-scraper-go/queue"
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

const (
	RequestTimeout = 5 * time.Second
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Malformed usage, to few or two many arguments!\n\n Usage: go run .\\main.go https://example.com")
	}
	startURL := os.Args[1]

	// Parse starting url
	u, err := url.Parse(startURL)
	if err != nil {
		log.Fatalf("Error parsing base URL: %v", err)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	baseDomain := u.Host
	startURL = u.String()

	// Create http client
	client := &http.Client{
		Timeout: RequestTimeout,
	}

	// Initialize data structures
	toCrawl := queue.New()
	visited := set.New()
	checked := set.New()
	stats := &LinkStats{ByStatusCode: make(map[string]int)}

	// Enqueue start link
	startLink := linktype.Link{
		URL:  startURL,
		Type: linktype.InternalLink,
	}
	toCrawl.Enqueue(startLink)
	visited.Add(startLink)

	// Start crawling
	for !toCrawl.IsEmpty() {
		currentLink := toCrawl.Dequeue()
		crawl(client, currentLink, baseDomain, visited, checked, toCrawl, stats)
	}
	// End
	printStats(*stats)
	log.Printf("Links visisted: %v\n", visited.Values())
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

func crawl(client *http.Client, link linktype.Link, baseDomain string, visited *set.Set, checked *set.Set, q *queue.Queue, stats *LinkStats) {
	currentURL := link.URL

	if !checked.Contains(link) {
		checked.Add(link)
		validateLink(client, link, stats)
	}
	validateLink(client, link, stats)
	log.Printf("[Crawling]: %s\n", currentURL)

	if link.Type != linktype.InternalLink {
		return
	}

	// Fetch page
	resp, err := client.Get(currentURL)
	if err != nil {
		log.Printf("Error fetching %s: %v", currentURL, err)
		return
	}
	defer resp.Body.Close()

	// Process only HTML pages
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return
	}
	baseURL := resp.Request.URL
	linksMap := extractLinks(root, baseURL)

	// Process all links found on page
	for _, newLink := range linksMap {
		switch newLink.Type {
		case linktype.InternalLink, linktype.PageLink:
			if !visited.Contains(newLink) {
				visited.Add(newLink)
				q.Enqueue(newLink)
			}
		case linktype.ExternalLink:
			if !checked.Contains(newLink) {
				checked.Add(newLink)
				validateLink(client, newLink, stats)
			}
		}
	}
}

func validateLink(client *http.Client, link linktype.Link, stats *LinkStats) {
	stats.Total++

	switch link.Type {
	case linktype.PageLink:
		stats.Skipped++
		log.Printf("[SKIP]   %s (Page link)", link.URL)
		return
	case linktype.InternalLink:
		stats.Internal++
	case linktype.ExternalLink:
		stats.External++
	default:
		stats.Skipped++
		log.Printf("[SKIP]   %s (Unknown type)", link.URL)
		return
	}

	if link.URL == "" {
		stats.Dead++
		log.Printf("[DEAD]   (empty URL)")
		return
	}
	status, statusCode, err := fetchStatus(client, link.URL)
	if err != nil {
		stats.Dead++
		log.Printf("[DEAD]   %s (%v)", link.URL, err)
		return
	}

	codeStr := strconv.Itoa(statusCode)
	stats.ByStatusCode[codeStr]++

	switch {
	case statusCode >= 400:
		stats.Dead++
		log.Printf("[DEAD]   %s (%s)", link.URL, status)

	default:
		stats.Alive++
		log.Printf("[ALIVE]  %s (%s)", link.URL, status)
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
