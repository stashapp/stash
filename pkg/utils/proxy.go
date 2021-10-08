package utils

import (
	"net/http"
	"strings"
)

func GetProxyPrefix(headers http.Header) string {
	prefix := ""
	if headers.Get("X-Forwarded-Prefix") != "" {
		prefix = strings.TrimRight(headers.Get("X-Forwarded-Prefix"), "/")
	}

	return prefix
}
