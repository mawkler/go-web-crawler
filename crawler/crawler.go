package crawler

import (
	"fmt"
	"net/url"
	"sort"
	"sync"

	"github.com/mawkler/go-web-crawler/internal"
)

type Crawler struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

type page struct {
	url   string
	count int
}

func NewCrawler(baseURL *url.URL, concurrencyControl chan struct{}, wg *sync.WaitGroup, maxPages int) Crawler {
	pages := map[string]int{}
	mu := &sync.Mutex{}
	return Crawler{pages, baseURL, mu, concurrencyControl, wg, maxPages}
}

func (c *Crawler) addPageVisit(normalizedURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.pages[normalizedURL]
	if exists {
		c.pages[normalizedURL]++
	} else {
		c.pages[normalizedURL] = 1
	}

	return exists
}

func (c *Crawler) getSortedPages() []page {
	pages := make([]page, 0, len(c.pages))
	for url, count := range c.pages {
		pages = append(pages, page{url, count})
	}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].count > pages[j].count
	})

	return pages
}

func (c Crawler) String() string {
	title := fmt.Sprintf(" REPORT for %s ", c.baseURL)
	lines := ""

	for range len(title) {
		lines += "="
	}

	string := fmt.Sprintf("%s\n%s\n%s\n\n", lines, title, lines)

	pages := c.getSortedPages()

	for _, p := range pages {
		string += fmt.Sprintf("Found %d internal links to %s\n", p.count, p.url)
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

// Name description
func (c *Crawler) pagesLength() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.pages)
}

func (c *Crawler) CrawlPage(rawCurrentURL string) {
	c.concurrencyControl <- struct{}{}

	defer func() {
		c.wg.Done()
		<-c.concurrencyControl
	}()

	// Stop crawling if crawling limit reached
	if c.pagesLength() >= c.maxPages {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("failed to parse current URL: %s\n", err)
		return
	}

	// We only want to parse pages on the same domain
	if currentURL.Host != c.baseURL.Host {
		return
	}

	normalizedURL, err := internal.NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("failed to normalize URL %s: %s\n", rawCurrentURL, err)
		return
	}

	fmt.Printf("crawling %s\n", normalizedURL)

	firstVisit := c.addPageVisit(normalizedURL)
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
		c.wg.Add(1)
		go c.CrawlPage(u)
	}
}
