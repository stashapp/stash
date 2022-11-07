package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Yamashou/gqlgenc/graphqljson"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// HTTPRequestOption represents the options applicable to the http client
type HTTPRequestOption func(req *http.Request)

// Client is the http client wrapper
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

// NewClient creates a new http client wrapper
func NewClient(client *http.Client, baseURL string, options ...HTTPRequestOption) *Client {
	return &Client{
		Client:             client,
		BaseURL:            baseURL,
		HTTPRequestOptions: options,
	}
}

func (c *Client) newRequest(ctx context.Context, operationName, query string, vars map[string]interface{}, httpRequestOptions []HTTPRequestOption) (*http.Request, error) {
	r := &Request{
		Query:         query,
		Variables:     vars,
		OperationName: operationName,
	}

	requestBody, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("create request struct failed: %w", err)
	}

	for _, httpRequestOption := range c.HTTPRequestOptions {
		httpRequestOption(req)
	}
	for _, httpRequestOption := range httpRequestOptions {
		httpRequestOption(req)
	}

	return req, nil
}

// GqlErrorList is the struct of a standard graphql error response
type GqlErrorList struct {
	Errors gqlerror.List `json:"errors"`
}

func (e *GqlErrorList) Error() string {
	return e.Errors.Error()
}

// HTTPError is the error when a GqlErrorList cannot be parsed
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse represent an handled error
type ErrorResponse struct {
	// populated when http status code is not OK
	NetworkError *HTTPError `json:"networkErrors"`
	// populated when http status code is OK but the server returned at least one graphql error
	GqlErrors *gqlerror.List `json:"graphqlErrors"`
}

// HasErrors returns true when at least one error is declared
func (er *ErrorResponse) HasErrors() bool {
	return er.NetworkError != nil || er.GqlErrors != nil
}

func (er *ErrorResponse) Error() string {
	content, err := json.Marshal(er)
	if err != nil {
		return err.Error()
	}

	return string(content)
}

// Post sends a http POST request to the graphql endpoint with the given query then unpacks
// the response into the given object.
func (c *Client) Post(ctx context.Context, operationName, query string, respData interface{}, vars map[string]interface{}, httpRequestOptions ...HTTPRequestOption) error {
	req, err := c.newRequest(ctx, operationName, query, vars, httpRequestOptions)
	if err != nil {
		return fmt.Errorf("don't create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return parseResponse(body, resp.StatusCode, respData)
}

func parseResponse(body []byte, httpCode int, result interface{}) error {
	errResponse := &ErrorResponse{}
	isKOCode := httpCode < 200 || 299 < httpCode
	if isKOCode {
		errResponse.NetworkError = &HTTPError{
			Code:    httpCode,
			Message: fmt.Sprintf("Response body %s", string(body)),
		}
	}

	// some servers return a graphql error with a non OK http code, try anyway to parse the body
	if err := unmarshal(body, result); err != nil {
		if gqlErr, ok := err.(*GqlErrorList); ok {
			errResponse.GqlErrors = &gqlErr.Errors
		} else if !isKOCode { // if is KO code there is already the http error, this error should not be returned
			return err
		}
	}

	if errResponse.HasErrors() {
		return errResponse
	}

	return nil
}

// response is a GraphQL layer response from a handler.
type response struct {
	Data   json.RawMessage `json:"data"`
	Errors json.RawMessage `json:"errors"`
}

func unmarshal(data []byte, res interface{}) error {
	resp := response{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("failed to decode data %s: %w", string(data), err)
	}

	if resp.Errors != nil && len(resp.Errors) > 0 {
		// try to parse standard graphql error
		errors := &GqlErrorList{}
		if e := json.Unmarshal(data, errors); e != nil {
			return fmt.Errorf("faild to parse graphql errors. Response content %s - %w ", string(data), e)
		}

		return errors
	}

	if err := graphqljson.UnmarshalData(resp.Data, res); err != nil {
		return fmt.Errorf("failed to decode data into response %s: %w", string(data), err)
	}

	return nil
}
