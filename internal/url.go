package internal

import (
	"fmt"
	netUrl "net/url"
	"strings"
)

func NormalizeURL(url string) (string, error) {
	parsedURL, err := netUrl.Parse(url)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL %s: %s", url, err)
	}

	path := parsedURL.Host + parsedURL.Path
	normalizedPath := strings.TrimSuffix(strings.ToLower(path), "/")

	return normalizedPath, nil
}
