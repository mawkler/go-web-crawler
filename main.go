package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/mawkler/go-web-crawler/crawler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := os.Args[1]

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("failed to parse base URL: %s\n", err)
		return
	}

	ch := make(chan struct{})
	wg := &sync.WaitGroup{}
	cr := crawler.NewCrawler(baseURL, ch, wg)

	pages, err := cr.CrawlPage(rawBaseURL)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println(crawler.PagesToString(pages))
}
