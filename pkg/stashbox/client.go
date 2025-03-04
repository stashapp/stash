// Package stashbox provides a client interface to a stash-box server instance.
package stashbox

import (
	"context"
	"net/http"
	"regexp"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/stashbox/graphql"
)

// Client represents the client interface to a stash-box server instance.
type Client struct {
	client *graphql.Client
	box    models.StashBox

	// tag patterns to be excluded
	excludeTagRE []*regexp.Regexp
}

// NewClient returns a new instance of a stash-box client.
func NewClient(box models.StashBox, excludeTagPatterns []string) *Client {
	authHeader := func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		req.Header.Set("ApiKey", box.APIKey)
		return next(ctx, req, gqlInfo, res)
	}

	client := &graphql.Client{
		Client: clientv2.NewClient(http.DefaultClient, box.Endpoint, nil, authHeader),
	}

	return &Client{
		client:       client,
		box:          box,
		excludeTagRE: scraper.CompileExclusionRegexps(excludeTagPatterns),
	}
}

func (c Client) getHTTPClient() *http.Client {
	return c.client.Client.Client
}

func (c Client) GetUser(ctx context.Context) (*graphql.Me, error) {
	return c.client.Me(ctx)
}
