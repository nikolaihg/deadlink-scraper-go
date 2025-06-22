package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type LinkType int

const (
	InternalLink LinkType = iota
	ExternalLink
	PageLink
)

type Link struct {
	URL  string
	Type LinkType
}

func main() {
	mainURL := "https://simple.wikipedia.org/wiki/The_Beatles"

	doc := fetchAndParse(mainURL)

	baseURL, err := url.Parse(mainURL)
	if err != nil {
		log.Fatalf("Error parsing base URL: %v", err)
	}

	links := make(map[string]Link)
	findHref(doc, baseURL, links)

	printLinks(links)
}

func fetchAndParse(urlStr string) *html.Node {
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatalf("Error fetching page: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing html: %v", err)
	}

	return doc
}

func findHref(node *html.Node, baseURL *url.URL, links map[string]Link) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				link := filterLink(attr.Val, baseURL)
				// add to map, avoid duplicate
				if _, exists := links[link.URL]; !exists {
					links[link.URL] = link
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findHref(child, baseURL, links)
	}
}

func filterLink(href string, baseURL *url.URL) Link {
	// Ignore empty
	if href == "" {
		return Link{}
	}
	if strings.HasPrefix(href, "#") {
		return Link{
			URL:  baseURL.String() + href,
			Type: PageLink,
		}
	}

	parsedHref, err := url.Parse(href)
	if err != nil {
		// If href is malformed, return as-is as external to avoid crash
		return Link{
			URL:  href,
			Type: ExternalLink,
		}
	}

	absURL := baseURL.ResolveReference(parsedHref)

	var linkType LinkType
	if absURL.Host == baseURL.Host {
		linkType = InternalLink
	} else {
		linkType = ExternalLink
	}

	return Link{
		URL:  absURL.String(),
		Type: linkType,
	}
}

func printLinks(links map[string]Link) {
	internal := []Link{}
	external := []Link{}
	page := []Link{}

	for _, link := range links {
		switch link.Type {
		case InternalLink:
			internal = append(internal, link)
		case ExternalLink:
			external = append(external, link)
		case PageLink:
			page = append(page, link)
		}
	}

	counts := [3]int{}
	for _, link := range links {
		counts[link.Type]++
	}

	fmt.Println("\n======= Internal Links =======")
	for _, link := range internal {
		fmt.Println(link.URL)
	}

	fmt.Println("\n======= External Links =======")
	for _, link := range external {
		fmt.Println(link.URL)
	}

	fmt.Println("\n======= Page Links =======")
	for _, link := range page {
		fmt.Println(link.URL)
	}

	fmt.Println("\n======= Link Counts =======")
	fmt.Printf("Internal: %d\n", counts[InternalLink])
	fmt.Printf("External: %d\n", counts[ExternalLink])
	fmt.Printf("Page:     %d\n", counts[PageLink])
}
