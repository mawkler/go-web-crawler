package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func getAbsolutePath(url, baseURL string) string {
	if strings.HasPrefix(url, "/") {
		return baseURL + url
	} else {
		return url
	}
}

func getURLsFromHTML(htmlBody, rawBaseURL string) (urls []string, err error) {
	urls = []string{}

	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("failed to parase HTML body: %s", err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					url := getAbsolutePath(a.Val, rawBaseURL)
					urls = append(urls, url)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return urls, nil
}

func main() {
	println("Hello, World!")
}
