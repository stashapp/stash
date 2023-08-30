package manager

import (
	"net/http"
	"strings"

	"github.com/stashapp/stash/internal/manager/config"
)

func GetProxyPrefix(r *http.Request) string {
	return strings.TrimRight(r.Header.Get("X-Forwarded-Prefix"), "/")
}

// Returns stash's baseurl
func GetBaseURL(r *http.Request) string {
	scheme := "http"
	if strings.Compare("https", r.URL.Scheme) == 0 || r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	prefix := GetProxyPrefix(r)

	baseURL := scheme + "://" + r.Host + prefix

	externalHost := config.GetInstance().GetExternalHost()
	if externalHost != "" {
		baseURL = externalHost + prefix
	}

	return baseURL
}
