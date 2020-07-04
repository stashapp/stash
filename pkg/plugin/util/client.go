package util

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/shurcooL/graphql"

	"github.com/stashapp/stash/pkg/plugin/common"
)

func NewClient(provider common.StashServerProvider) *graphql.Client {
	portStr := strconv.Itoa(provider.GetPort())

	u, _ := url.Parse("http://localhost:" + portStr + "/graphql")
	u.Scheme = provider.GetScheme()

	cookieJar, _ := cookiejar.New(nil)

	cookie := provider.GetSessionCookie()
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
