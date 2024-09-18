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

	maxConcurrency := 50
	maxPages := 1000
	ch := make(chan struct{}, maxConcurrency)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	cr := crawler.NewCrawler(baseURL, ch, wg, maxPages)

	cr.CrawlPage(rawBaseURL)
	wg.Wait()

	fmt.Println(crawler.PagesToString(cr.GetPages()))
}
