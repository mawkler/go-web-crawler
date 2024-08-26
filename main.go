package main

import (
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
	r := strings.NewReader(htmlBody)
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return urls, nil
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()

			if string(tn) == "a" {
				key, url, _ := z.TagAttr()
				if string(key) == "href" {
					absoluteURL := getAbsolutePath(string(url), rawBaseURL)
					urls = append(urls, absoluteURL)
				}
			}
		}
	}
}

func main() {
	println("Hello, World!")
}
