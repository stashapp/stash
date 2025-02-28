// Package util implements utility and convenience methods for plugins. It is
// not intended for the main stash code to access.
package util

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	graphql "github.com/hasura/go-graphql-client"

	"github.com/stashapp/stash/pkg/plugin/common"
)

// NewClient creates a graphql Client connecting to the stash server using
// the provided server connection details.
// Always connects to the graphql endpoint of the localhost.
func NewClient(provider common.StashServerConnection) *graphql.Client {
	portStr := strconv.Itoa(provider.Port)

	u, _ := url.Parse("http://" + provider.Host + ":" + portStr + "/graphql")
	u.Scheme = provider.Scheme

	cookieJar, _ := cookiejar.New(nil)

	cookie := provider.SessionCookie
	if cookie != nil {
		cookieJar.SetCookies(u, []*http.Cookie{
			cookie,
		})
	}

	httpClient := &http.Client{
		Jar: cookieJar,
	}

	return graphql.NewClient(u.String(), httpClient)
}
