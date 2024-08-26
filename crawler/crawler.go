package crawler

import (
	"fmt"
	"net/url"

	"github.com/mawkler/go-web-crawler/internal"
)

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

func CrawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) (map[string]int, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to parse base URL: %s", err)
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to parse current URL: %s", err)
	}

	// We only want to parse pages on the same domain
	if currentURL.Host != baseURL.Host {
		return pages, nil
	}

	normalizedURL, err := internal.NormalizeURL(rawCurrentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize URL %s: %s", rawCurrentURL, err)
	}

	_, exists := pages[normalizedURL]
	if exists {
		pages[normalizedURL] += 1
		return pages, nil
	}

	pages[normalizedURL] = 1

	html, err := internal.GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to get HTML: %s", err)
	}

	urls, err := internal.GetURLsFromHTML(html, rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to get URLs from HTML: %s", err)
	}

	for _, u := range urls {
		// We don't need to merge the returned map with `pages` because
		// CrawlPage() writes to `pages` as a side-effect
		_, err := CrawlPage(rawBaseURL, u, pages)
		if err != nil {
			fmt.Println(err)
			pages[normalizedURL] = -1
		}
	}

	return pages, nil
}
