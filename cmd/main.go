package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("https://en.wikipedia.org/wiki/The_Beatles")
	if err != nil {
		log.Fatalf("Error fetching the URL: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status: ", resp.Status)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing html: %v", err)
	}

	findHref(doc)
}

func findHref(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				fmt.Printf("Tag: <%s>, attr: %s, url: %s\n", node.Data, attr.Key, attr.Val)
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findHref(child)
	}
}
