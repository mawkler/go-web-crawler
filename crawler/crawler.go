package crawler

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/mawkler/go-web-crawler/internal"
)

type Crawler struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func NewCrawler(baseURL *url.URL, concurrencyControl chan struct{}, wg *sync.WaitGroup) Crawler {
	pages := map[string]int{}
	mu := &sync.Mutex{}
	return Crawler{pages, baseURL, mu, concurrencyControl, wg}
}

func (cfg *Crawler) GetPages() map[string]int {
	return cfg.pages
}

func (cfg *Crawler) addPageVisit(normalizedURL string) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	_, exists := cfg.pages[normalizedURL]
	if exists {
		cfg.pages[normalizedURL]++
	} else {
		cfg.pages[normalizedURL] = 1
	}

	return exists
}

func PagesToString(pages map[string]int) string {
	string := ""

	for key, value := range pages {
		string += fmt.Sprintf("%s: %d\n", key, value)
	}

	return string
}

func mergeMaps(map1, map2 map[string]int) map[string]int {
	mergedMap := make(map[string]int)

	for key, value := range map1 {
		mergedMap[key] = value
	}

	for key, value := range map2 {
		mergedMap[key] += value
	}

	return mergedMap
}

func (cfg *Crawler) CrawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}

	defer func() {
		cfg.wg.Done()
		<-cfg.concurrencyControl
	}()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("failed to parse current URL: %s\n", err)
		return
	}

	// We only want to parse pages on the same domain
	if currentURL.Host != cfg.baseURL.Host {
		return
	}

	normalizedURL, err := internal.NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("failed to normalize URL %s: %s\n", rawCurrentURL, err)
		return
	}

	firstVisit := cfg.addPageVisit(normalizedURL)
	if firstVisit {
		return
	}

	html, err := internal.GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("failed to get HTML: %s\n", err)
		return
	}

	urls, err := internal.GetURLsFromHTML(html, rawCurrentURL)
	if err != nil {
		fmt.Printf("failed to get URLs from HTML: %s\n", err)
		return
	}

	for _, u := range urls {
		cfg.wg.Add(1)
		go cfg.CrawlPage(u)
	}
}
