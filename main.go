package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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
	mainURL := "https://scrape-me.dreamsofcode.io/locations"

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

	fmt.Printf("Scanning base page: %s\n", mainURL)
	links := findLinks(doc, baseURL)

	fmt.Println("\nStarting validations: ")
	for _, link := range links {
		if visited.Contains(link) {
			continue
		}
		visited.Add(link)

		validateLink(link, stats)
	}

	fmt.Println("\nPages added to queue: ")
	myQueue.Print()

	printStats(*stats)
}

func validateLink(link linktype.Link, stats *LinkStats) {
	stats.Total++

	switch link.Type {
	case linktype.PageLink:
		stats.Skipped++
		fmt.Printf("[SKIP]   %s (Page link)\n", link.URL)
		return
	case linktype.InternalLink:
		stats.Internal++
	case linktype.ExternalLink:
		stats.External++
	default:
		stats.Skipped++
		fmt.Printf("[SKIP]   %s (Unknown type)\n", link.URL)
		return
	}

	status, err := fetch(link.URL)
	if err != nil {
		stats.Dead++
		fmt.Printf("[DEAD]   %s (%v)\n", link.URL, err)
		return
	}

	statusCode := strings.Split(status, " ")[0]
	stats.ByStatusCode[statusCode]++

	if strings.HasPrefix(statusCode, "4") || strings.HasPrefix(statusCode, "5") {
		stats.Dead++
		fmt.Printf("[DEAD]   %s (%s)\n", link.URL, status)
	} else {
		stats.Alive++
		fmt.Printf("[ALIVE]  %s (%s)\n", link.URL, status)
	}
}

func fetch(urlStr string) (string, error) {
	resp, err := httpClient.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("error fetching page: %w", err)
	}
	defer resp.Body.Close()
	return resp.Status, nil
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
	fmt.Println("\nScan complete:")
	fmt.Printf("Total:    %d\n", stats.Total)
	fmt.Printf("Internal: %d\n", stats.Internal)
	fmt.Printf("External: %d\n", stats.External)
	fmt.Printf("Alive:    %d\n", stats.Alive)
	fmt.Printf("Dead:     %d\n", stats.Dead)
	fmt.Printf("Skipped:  %d\n", stats.Skipped)
	fmt.Println("Status codes breakdown:")
	for code, count := range stats.ByStatusCode {
		fmt.Printf("  %s: %d\n", code, count)
	}
}
