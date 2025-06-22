package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/nikolaihg/deadlink-scraper-go/linktype"
	"golang.org/x/net/html"
)

func main() {
	mainURL := "https://scrape-me.dreamsofcode.io/about"

	doc := fetchAndParse(mainURL)

	baseURL, err := url.Parse(mainURL)
	if err != nil {
		log.Fatalf("Error parsing base URL: %v", err)
	}

	links := make(map[string]linktype.Link)
	findHref(doc, baseURL, links)

	printLinks(links)
}

func fetchAndParse(urlStr string) *html.Node {
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatalf("Error fetching page: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("BaseURL Response status:", resp.Status)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing html: %v", err)
	}

	return doc
}

func findHref(node *html.Node, baseURL *url.URL, links map[string]linktype.Link) {
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

func filterLink(href string, baseURL *url.URL) linktype.Link {
	// Ignore empty
	if href == "" {
		return linktype.Link{}
	}
	// find page links
	if strings.HasPrefix(href, "#") {
		return linktype.Link{
			URL:  baseURL.String() + href,
			Type: linktype.PageLink,
		}
	}

	parsedHref, err := url.Parse(href)
	if err != nil {
		// If href is malformed, return as-is as external to avoid crash
		return linktype.Link{
			URL:  href,
			Type: linktype.ExternalLink,
		}
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

func printLinks(links map[string]linktype.Link) {
	internal := []linktype.Link{}
	external := []linktype.Link{}
	page := []linktype.Link{}

	for _, link := range links {
		switch link.Type {
		case linktype.InternalLink:
			internal = append(internal, link)
		case linktype.ExternalLink:
			external = append(external, link)
		case linktype.PageLink:
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
	fmt.Printf("Internal: %d\n", counts[linktype.InternalLink])
	fmt.Printf("External: %d\n", counts[linktype.ExternalLink])
	fmt.Printf("Page:     %d\n", counts[linktype.PageLink])
}
