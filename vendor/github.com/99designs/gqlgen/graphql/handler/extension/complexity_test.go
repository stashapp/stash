package extension_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/testserver"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/stretchr/testify/require"
)

func TestHandlerComplexity(t *testing.T) {
	h := testserver.New()
	h.Use(&extension.ComplexityLimit{
		Func: func(ctx context.Context, rc *graphql.OperationContext) int {
			if rc.RawQuery == "{ ok: name }" {
				return 4
			}
			return 2
		},
	})
	h.AddTransport(&transport.POST{})
	var stats *extension.ComplexityStats
	h.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		stats = extension.GetComplexityStats(ctx)
		return next(ctx)
	})

	t.Run("below complexity limit", func(t *testing.T) {
		stats = nil
		h.SetCalculatedComplexity(2)
		resp := doRequest(h, "POST", "/graphql", `{"query":"{ name }"}`)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		require.Equal(t, `{"data":{"name":"test"}}`, resp.Body.String())

		require.Equal(t, 2, stats.ComplexityLimit)
		require.Equal(t, 2, stats.Complexity)
	})

	t.Run("above complexity limit", func(t *testing.T) {
		stats = nil
		h.SetCalculatedComplexity(4)
		resp := doRequest(h, "POST", "/graphql", `{"query":"{ name }"}`)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		require.Equal(t, `{"errors":[{"message":"operation has complexity 4, which exceeds the limit of 2","extensions":{"code":"COMPLEXITY_LIMIT_EXCEEDED"}}],"data":null}`, resp.Body.String())

		require.Equal(t, 2, stats.ComplexityLimit)
		require.Equal(t, 4, stats.Complexity)
	})

	t.Run("within dynamic complexity limit", func(t *testing.T) {
		stats = nil
		h.SetCalculatedComplexity(4)
		resp := doRequest(h, "POST", "/graphql", `{"query":"{ ok: name }"}`)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		require.Equal(t, `{"data":{"name":"test"}}`, resp.Body.String())

		require.Equal(t, 4, stats.ComplexityLimit)
		require.Equal(t, 4, stats.Complexity)
	})
}

func TestFixedComplexity(t *testing.T) {
	h := testserver.New()
	h.Use(extension.FixedComplexityLimit(2))
	h.AddTransport(&transport.POST{})

	var stats *extension.ComplexityStats
	h.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		stats = extension.GetComplexityStats(ctx)
		return next(ctx)
	})

	t.Run("below complexity limit", func(t *testing.T) {
		h.SetCalculatedComplexity(2)
		resp := doRequest(h, "POST", "/graphql", `{"query":"{ name }"}`)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		require.Equal(t, `{"data":{"name":"test"}}`, resp.Body.String())

		require.Equal(t, 2, stats.ComplexityLimit)
		require.Equal(t, 2, stats.Complexity)
	})

	t.Run("above complexity limit", func(t *testing.T) {
		h.SetCalculatedComplexity(4)
		resp := doRequest(h, "POST", "/graphql", `{"query":"{ name }"}`)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		require.Equal(t, `{"errors":[{"message":"operation has complexity 4, which exceeds the limit of 2","extensions":{"code":"COMPLEXITY_LIMIT_EXCEEDED"}}],"data":null}`, resp.Body.String())

		require.Equal(t, 2, stats.ComplexityLimit)
		require.Equal(t, 4, stats.Complexity)
	})

	t.Run("bypass __schema field", func(t *testing.T) {
		h.SetCalculatedComplexity(4)
		resp := doRequest(h, "POST", "/graphql", `{ "operationName":"IntrospectionQuery", "query":"query IntrospectionQuery { __schema { queryType { name } mutationType { name }}}"}`)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		require.Equal(t, `{"data":{"name":"test"}}`, resp.Body.String())

		require.Equal(t, 2, stats.ComplexityLimit)
		require.Equal(t, 0, stats.Complexity)
	})
}

func doRequest(handler http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)
	return w
}
