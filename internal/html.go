package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func GetHTML(rawURL string) (string, error) {
	response, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to get %s: %s", rawURL, err)
	}

	if response.StatusCode >= 400 {
		return "", fmt.Errorf("failed to get %s: status code %d", rawURL, response.StatusCode)
	}

	contentType := response.Header.Get("content-type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf("invalid response content-type from %s. Expected text/html, got: %s", rawURL, contentType)
	}

	html, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err)
	}

	return string(html), nil
}

func GetURLsFromHTML(htmlBody, rawBaseURL string) (urls []string, err error) {
	urls = []string{}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %s", err)
	}

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
					href, err := url.Parse(a.Val)
					if err != nil {
						fmt.Printf("couldn't parse href '%s': %v\n", a.Val, err)
						continue
					}

					absoluteURL := baseURL.ResolveReference(href)
					urls = append(urls, absoluteURL.String())
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
