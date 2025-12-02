package utils

import "regexp"

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
