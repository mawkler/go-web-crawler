package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/mawkler/go-web-crawler/crawler"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Not enough arguments provided")
		fmt.Println("usage: crawler <baseURL> <maxConcurrency> <maxPages>")
		os.Exit(1)
	}
	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := os.Args[1]
	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("maxConcurrency not an int")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("maxPages not an int")
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("failed to parse base URL: %s\n", err)
		return
	}

	ch := make(chan struct{}, maxConcurrency)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	cr := crawler.NewCrawler(baseURL, ch, wg, maxPages)

	cr.CrawlPage(rawBaseURL)
	wg.Wait()

	fmt.Println(cr)
}
