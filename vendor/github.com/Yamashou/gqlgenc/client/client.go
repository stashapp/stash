package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/xerrors"

	"github.com/Yamashou/gqlgenc/graphqljson"
)

type HTTPRequestOption func(req *http.Request)

type Client struct {
	Client             *http.Client
	BaseURL            string
	HTTPRequestOptions []HTTPRequestOption
}

// Request represents an outgoing GraphQL request
type Request struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

func NewClient(client *http.Client, baseURL string, options ...HTTPRequestOption) *Client {
	return &Client{
		Client:             client,
		BaseURL:            baseURL,
		HTTPRequestOptions: options,
	}
}

func (c *Client) newRequest(ctx context.Context, query string, vars map[string]interface{}, httpRequestOptions []HTTPRequestOption) (*http.Request, error) {
	r := &Request{
		Query:         query,
		Variables:     vars,
		OperationName: "",
	}

	requestBody, err := json.Marshal(r)
	if err != nil {
		return nil, xerrors.Errorf("encode: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, xerrors.Errorf("create request struct failed: %w", err)
	}

	for _, httpRequestOption := range c.HTTPRequestOptions {
		httpRequestOption(req)
	}
	for _, httpRequestOption := range httpRequestOptions {
		httpRequestOption(req)
	}

	return req, nil
}

// Post sends a http POST request to the graphql endpoint with the given query then unpacks
// the response into the given object.
func (c *Client) Post(ctx context.Context, query string, respData interface{}, vars map[string]interface{}, httpRequestOptions ...HTTPRequestOption) error {
	req, err := c.newRequest(ctx, query, vars, httpRequestOptions)
	if err != nil {
		return xerrors.Errorf("don't create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	resp, err := c.Client.Do(req)
	if err != nil {
		return xerrors.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if err := graphqljson.Unmarshal(resp.Body, respData); err != nil {
		return err
	}

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return xerrors.Errorf("http status code: %v", resp.StatusCode)
	}

	return nil
}
