package extension_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/testserver"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/stretchr/testify/require"
)

func TestAPQIntegration(t *testing.T) {
	h := testserver.New()
	h.Use(&extension.AutomaticPersistedQuery{Cache: graphql.MapCache{}})
	h.AddTransport(&transport.POST{})

	var stats *extension.ApqStats
	h.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		stats = extension.GetApqStats(ctx)
		return next(ctx)
	})

	resp := doRequest(h, "POST", "/graphql", `{"query":"{ name }","extensions":{"persistedQuery":{"version":1,"sha256Hash":"30166fc3298853f22709fce1e4a00e98f1b6a3160eaaaf9cb3b7db6a16073b07"}}}`)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.Equal(t, `{"data":{"name":"test"}}`, resp.Body.String())

	require.NotNil(t, stats)
	require.True(t, stats.SentQuery)
	require.Equal(t, "30166fc3298853f22709fce1e4a00e98f1b6a3160eaaaf9cb3b7db6a16073b07", stats.Hash)
}

func TestAPQ(t *testing.T) {
	const query = "{ me { name } }"
	const hash = "b8d9506e34c83b0e53c2aa463624fcea354713bc38f95276e6f0bd893ffb5b88"

	t.Run("with query and no hash", func(t *testing.T) {
		ctx := newOC()
		params := &graphql.RawParams{
			Query: "original query",
		}
		err := extension.AutomaticPersistedQuery{graphql.MapCache{}}.MutateOperationParameters(ctx, params)
		require.Nil(t, err)

		require.Equal(t, "original query", params.Query)
	})

	t.Run("with hash miss and no query", func(t *testing.T) {
		ctx := newOC()
		params := &graphql.RawParams{
			Extensions: map[string]interface{}{
				"persistedQuery": map[string]interface{}{
					"sha256Hash": hash,
					"version":    1,
				},
			},
		}

		err := extension.AutomaticPersistedQuery{graphql.MapCache{}}.MutateOperationParameters(ctx, params)
		require.Equal(t, err.Message, "PersistedQueryNotFound")
	})

	t.Run("with hash miss and query", func(t *testing.T) {
		ctx := newOC()
		params := &graphql.RawParams{
			Query: query,
			Extensions: map[string]interface{}{
				"persistedQuery": map[string]interface{}{
					"sha256Hash": hash,
					"version":    1,
				},
			},
		}
		cache := graphql.MapCache{}
		err := extension.AutomaticPersistedQuery{cache}.MutateOperationParameters(ctx, params)
		require.Nil(t, err)

		require.Equal(t, "{ me { name } }", params.Query)
		require.Equal(t, "{ me { name } }", cache[hash])
	})

	t.Run("with hash miss and query", func(t *testing.T) {
		ctx := newOC()
		params := &graphql.RawParams{
			Query: query,
			Extensions: map[string]interface{}{
				"persistedQuery": map[string]interface{}{
					"sha256Hash": hash,
					"version":    1,
				},
			},
		}
		cache := graphql.MapCache{}
		err := extension.AutomaticPersistedQuery{cache}.MutateOperationParameters(ctx, params)
		require.Nil(t, err)

		require.Equal(t, "{ me { name } }", params.Query)
		require.Equal(t, "{ me { name } }", cache[hash])
	})

	t.Run("with hash hit and no query", func(t *testing.T) {
		ctx := newOC()
		params := &graphql.RawParams{
			Extensions: map[string]interface{}{
				"persistedQuery": map[string]interface{}{
					"sha256Hash": hash,
					"version":    1,
				},
			},
		}
		cache := graphql.MapCache{
			hash: query,
		}
		err := extension.AutomaticPersistedQuery{cache}.MutateOperationParameters(ctx, params)
		require.Nil(t, err)

		require.Equal(t, "{ me { name } }", params.Query)
	})

	t.Run("with malformed extension payload", func(t *testing.T) {
		ctx := newOC()
		params := &graphql.RawParams{
			Extensions: map[string]interface{}{
				"persistedQuery": "asdf",
			},
		}

		err := extension.AutomaticPersistedQuery{graphql.MapCache{}}.MutateOperationParameters(ctx, params)
		require.Equal(t, err.Message, "invalid APQ extension data")
	})

	t.Run("with invalid extension version", func(t *testing.T) {
		ctx := newOC()
		params := &graphql.RawParams{
			Extensions: map[string]interface{}{
				"persistedQuery": map[string]interface{}{
					"version": 2,
				},
			},
		}
		err := extension.AutomaticPersistedQuery{graphql.MapCache{}}.MutateOperationParameters(ctx, params)
		require.Equal(t, err.Message, "unsupported APQ version")
	})

	t.Run("with hash mismatch", func(t *testing.T) {
		ctx := newOC()
		params := &graphql.RawParams{
			Query: query,
			Extensions: map[string]interface{}{
				"persistedQuery": map[string]interface{}{
					"sha256Hash": "badhash",
					"version":    1,
				},
			},
		}

		err := extension.AutomaticPersistedQuery{graphql.MapCache{}}.MutateOperationParameters(ctx, params)
		require.Equal(t, err.Message, "provided APQ hash does not match query")
	})
}

func newOC() context.Context {
	oc := &graphql.OperationContext{}
	return graphql.WithOperationContext(context.Background(), oc)
}
