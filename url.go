package main

import (
	"fmt"
	netUrl "net/url"
	"strings"
)

func normalizeURL(url string) (string, error) {
	parsedURL, err := netUrl.Parse(strings.TrimSpace(url))
	if err != nil {
		return "", fmt.Errorf("failed to parse URL %s: %s", url, err)
	}

	return parsedURL.Host + parsedURL.Path, nil
}
