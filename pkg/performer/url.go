package performer

import (
	"regexp"
)

var (
	twitterURLRE   = regexp.MustCompile(`^https?:\/\/(?:www\.)?twitter\.com\/`)
	instagramURLRE = regexp.MustCompile(`^https?:\/\/(?:www\.)?instagram\.com\/`)
)

func IsTwitterURL(url string) bool {
	return twitterURLRE.MatchString(url)
}

func IsInstagramURL(url string) bool {
	return instagramURLRE.MatchString(url)
}
