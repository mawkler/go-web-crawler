package main

import (
	"fmt"
	"os"

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

	url := os.Args[1]

	fmt.Printf("starting crawl of: %s\n", url)

	pages, err := crawler.CrawlPage(url, url, map[string]int{})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println(crawler.PagesToString(pages))
}
