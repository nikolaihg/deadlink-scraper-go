package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		log.Fatalf("Error fetching the URL: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status: ", resp.Status)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing html: %v", err)
	}

	html.Render(os.Stdout, doc)

	traverseHTML(doc)
}

func traverseHTML(node *html.Node) {
	if node.Type == html.ElementNode {
		fmt.Println("Tag:", node.Data)
	}
	if node.Type == html.TextNode {
		if node.Data != "" {
			fmt.Print("Text:", strings.TrimSpace(node.Data))
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		traverseHTML(child)
	}
}
