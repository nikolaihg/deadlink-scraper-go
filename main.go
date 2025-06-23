package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nikolaihg/deadlink-scraper-go/linktype"
	"github.com/nikolaihg/deadlink-scraper-go/queue"
	"github.com/nikolaihg/deadlink-scraper-go/set"
	"golang.org/x/net/html"
)

const (
	RequestTimeout = 5 * time.Second
)

var httpClient = &http.Client{
	Timeout: RequestTimeout,
}

type LinkStats struct {
	Total        int
	Internal     int
	External     int
	Alive        int
	Dead         int
	Skipped      int
	ByStatusCode map[string]int
}

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("Malformed usage, to few or two many arguments!")
	}
	mainURL := os.Args[1]

	visited := set.New()
	myQueue := queue.New()
	stats := &LinkStats{ByStatusCode: make(map[string]int)}

	doc, err := fetchBase(mainURL)
	if err != nil {
		log.Fatalf("Failed to load base page: %v", err)
	}

	baseURL, err := url.Parse(mainURL)
	if err != nil {
		log.Fatalf("Error parsing base URL: %v", err)
	}

	log.Printf("Scanning base page: %s\n", mainURL)
	links := findLinks(doc, baseURL)
	for _, link := range links {
		if link.Type == linktype.InternalLink {
			myQueue.Enqueue(link)
		}
	}

	log.Printf("Starting validations:\n")
	for _, link := range links {
		if visited.Contains(link) {
			continue
		}
		visited.Add(link)

		validateLink(link, stats)
	}

	printStats(*stats)

	log.Printf("Pages added to queue:\n")
	log.Print(myQueue.String())

	log.Printf("Links visisted: %v\n", visited.Values())
}

func validateLink(link linktype.Link, stats *LinkStats) {
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

	status, statusCode, err := fetch(link.URL)
	if err != nil {
		stats.Dead++
		log.Printf("[DEAD]   %s (%v)", link.URL, err)
		return
	}

	codeStr := strconv.Itoa(statusCode)
	stats.ByStatusCode[codeStr]++

	if statusCode >= 400 {
		stats.Dead++
		log.Printf("[DEAD]   %s (%s)", link.URL, status)
	} else {
		stats.Alive++
		log.Printf("[ALIVE]  %s (%s)", link.URL, status)
	}
}

func fetch(urlStr string) (string, int, error) {
	resp, err := httpClient.Get(urlStr)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("error fetching page: %w", err)
	}
	defer resp.Body.Close()
	return resp.Status, resp.StatusCode, nil
}

func fetchBase(urlStr string) (*html.Node, error) {
	resp, err := httpClient.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error fetching page: %w", err)
	}
	defer resp.Body.Close()

	log.Println("BaseURL Response status:", resp.Status)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing html: %w", err)
	}

	return doc, nil
}

func findLinks(node *html.Node, baseURL *url.URL) map[string]linktype.Link {
	links := make(map[string]linktype.Link)

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key != "href" {
					continue
				}

				link := filterLink(attr.Val, baseURL)
				if link.URL == "" {
					continue
				}
				if _, exists := links[link.URL]; exists {
					continue
				}

				links[link.URL] = link
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
	if strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") || strings.HasPrefix(href, "javascript:") {
		return linktype.Link{}
	}
	if strings.HasPrefix(href, "#") {
		return linktype.Link{
			URL:  baseURL.String() + href,
			Type: linktype.PageLink,
		}
	}

	parsedHref, err := url.Parse(href)
	if err != nil {
		log.Printf("Skipping malformed URL: %s (%v)", href, err)
		return linktype.Link{}
	}

	absURL := baseURL.ResolveReference(parsedHref)

	var linkType linktype.LinkType
	if absURL.Host == baseURL.Host {
		linkType = linktype.InternalLink
	} else {
		linkType = linktype.ExternalLink
	}

	return linktype.Link{
		URL:  absURL.String(),
		Type: linkType,
	}
}

func printStats(stats LinkStats) {
	log.Println("Scan complete:")
	log.Printf("Total:    %d\n", stats.Total)
	log.Printf("Internal: %d\n", stats.Internal)
	log.Printf("External: %d\n", stats.External)
	log.Printf("Alive:    %d\n", stats.Alive)
	log.Printf("Dead:     %d\n", stats.Dead)
	log.Printf("Skipped:  %d\n", stats.Skipped)
	log.Println("Status codes breakdown:")
	for code, count := range stats.ByStatusCode {
		log.Printf("  %s: %d\n", code, count)
	}
}
