package utils

import (
	"regexp"
	"sort"
	"strings"
)

// URLFromHandle adds the site URL to the input if the input is not already a URL
// siteURL must not end with a slash
func URLFromHandle(input string, siteURL string) string {
	// if the input is already a URL, return it
	re := regexp.MustCompile(`^https?://`)
	if re.MatchString(input) {
		return input
	}

	return siteURL + "/" + input
}

// urlSortKey extracts the sortable portion of a URL by removing the protocol and www. prefix
func urlSortKey(url string) string {
	// Remove http:// or https://
	key := strings.TrimPrefix(url, "https://")
	key = strings.TrimPrefix(key, "http://")
	// Remove www. prefix
	key = strings.TrimPrefix(key, "www.")
	return strings.ToLower(key)
}

// SortURLs sorts a slice of URLs alphabetically by their base URL,
// excluding the protocol (http/https) and www. prefix
func SortURLs(urls []string) {
	sort.SliceStable(urls, func(i, j int) bool {
		return urlSortKey(urls[i]) < urlSortKey(urls[j])
	})
}
