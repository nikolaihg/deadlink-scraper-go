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
	"golang.org/x/net/html"
)

const (
	RequestTimeout = 5 * time.Second
)

var httpClient = &http.Client{
	Timeout: RequestTimeout,
}

func main() {
	mainURL := "https://scrape-me.dreamsofcode.io/about"
	links := make(map[string]linktype.Link)
	myQueue := queue.New()
	//visited := set.New()

	doc := fetchBase(mainURL)

	baseURL, err := url.Parse(mainURL)
	if err != nil {
		log.Fatalf("Error parsing base URL: %v", err)
	}

	findLinks(doc, baseURL, links, myQueue)

	fmt.Printf("Starting scan of %s\n\n", mainURL)

	for _, link := range links {
		switch link.Type {
		case linktype.InternalLink, linktype.ExternalLink:
			status, err := fetch(link.URL)
			if err != nil {
				fmt.Printf("[DEAD]   %s (%v)\n", link.URL, err)
			} else if strings.HasPrefix(status, "4") || strings.HasPrefix(status, "5") {
				fmt.Printf("[DEAD]   %s (%s)\n", link.URL, status)
			} else {
				fmt.Printf("[ALIVE]  %s (%s)\n", link.URL, status)
			}
		case linktype.PageLink:
			// Optional: Skip or log
			fmt.Printf("[SKIP]   %s (Page link)\n", link.URL)
		}
	}

	myQueue.Print()
}

func fetch(urlStr string) (string, error) {
	resp, err := httpClient.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("error fetching page: %w", err)
	}
	defer resp.Body.Close()

	return resp.Status, nil
}

func fetchBase(urlStr string) *html.Node {
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatalf("Error fetching page: %v", err)
	}
	defer resp.Body.Close()

	log.Println("BaseURL Response status:", resp.Status)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing html: %v", err)
	}

	return doc
}

func findLinks(node *html.Node, baseURL *url.URL, links map[string]linktype.Link, q *queue.Queue) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				link := filterLink(attr.Val, baseURL)
				// add to map, avoid duplicate
				if _, exists := links[link.URL]; !exists {
					links[link.URL] = link
					q.Enqueue(link)
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findLinks(child, baseURL, links, q)
	}
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
