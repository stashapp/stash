package performer

import "strings"

func IsTwitterURL(url string) bool {
	return strings.HasPrefix(url, "https://twitter.com/")
}

func IsInstagramURL(url string) bool {
	return strings.HasPrefix(url, "https://instagram.com/")
}
