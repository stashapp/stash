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

// DefaultMaxRequestsPerMinute is the default maximum number of requests per minute.
const DefaultMaxRequestsPerMinute = 240

// Client represents the client interface to a stash-box server instance.
type Client struct {
	client *graphql.Client
	box    models.StashBox

	maxRequestsPerMinute int

	// tag patterns to be excluded
	excludeTagRE []*regexp.Regexp
}

type ClientOption func(*Client)

func ExcludeTagPatterns(patterns []string) ClientOption {
	return func(c *Client) {
		c.excludeTagRE = scraper.CompileExclusionRegexps(patterns)
	}
}

func MaxRequestsPerMinute(n int) ClientOption {
	return func(c *Client) {
		c.maxRequestsPerMinute = n
	}
}

func setApiKeyHeader(apiKey string) clientv2.RequestInterceptor {
	return func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		req.Header.Set("ApiKey", apiKey)
		return next(ctx, req, gqlInfo, res)
	}
}

// NewClient returns a new instance of a stash-box client.
func NewClient(box models.StashBox, options ...ClientOption) *Client {
	authHeader := setApiKeyHeader(box.APIKey)

	client := &graphql.Client{
		Client: clientv2.NewClient(http.DefaultClient, box.Endpoint, nil, authHeader),
	}

	ret := &Client{
		client:               client,
		box:                  box,
		maxRequestsPerMinute: DefaultMaxRequestsPerMinute,
	}

	for _, option := range options {
		option(ret)
	}

	return ret
}

func (c Client) getHTTPClient() *http.Client {
	return c.client.Client.Client
}

func (c Client) GetUser(ctx context.Context) (*graphql.Me, error) {
	return c.client.Me(ctx)
}
