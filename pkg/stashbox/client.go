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

	"golang.org/x/time/rate"
)

// DefaultMaxRequestsPerMinute is the default maximum number of requests per minute.
const DefaultMaxRequestsPerMinute = 240

// Client represents the client interface to a stash-box server instance.
type Client struct {
	client     *graphql.Client
	httpClient *http.Client
	box        models.StashBox

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
		if n > 0 {
			c.maxRequestsPerMinute = n
		}
	}
}

func setApiKeyHeader(apiKey string) clientv2.RequestInterceptor {
	return func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		req.Header.Set("ApiKey", apiKey)
		return next(ctx, req, gqlInfo, res)
	}
}

func rateLimit(n int) clientv2.RequestInterceptor {
	perSec := float64(n) / 60
	limiter := rate.NewLimiter(rate.Limit(perSec), 1)

	return func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		if err := limiter.Wait(ctx); err != nil {
			// should only happen if the context is canceled
			return err
		}

		return next(ctx, req, gqlInfo, res)
	}
}

// NewClient returns a new instance of a stash-box client.
func NewClient(box models.StashBox, options ...ClientOption) *Client {
	ret := &Client{
		box:                  box,
		maxRequestsPerMinute: DefaultMaxRequestsPerMinute,
		httpClient:           http.DefaultClient,
	}

	if box.MaxRequestsPerMinute > 0 {
		ret.maxRequestsPerMinute = box.MaxRequestsPerMinute
	}

	for _, option := range options {
		option(ret)
	}

	authHeader := setApiKeyHeader(box.APIKey)
	limitRequests := rateLimit(ret.maxRequestsPerMinute)

	client := &graphql.Client{
		Client: clientv2.NewClient(ret.httpClient, box.Endpoint, nil, authHeader, limitRequests),
	}

	ret.client = client

	return ret
}

func (c Client) GetUser(ctx context.Context) (*graphql.Me, error) {
	return c.client.Me(ctx)
}
