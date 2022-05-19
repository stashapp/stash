package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/testserver"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
)

func TestServer(t *testing.T) {
	srv := testserver.New()
	srv.AddTransport(&transport.GET{})

	t.Run("returns an error if no transport matches", func(t *testing.T) {
		resp := post(srv, "/foo", "application/json")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, `{"errors":[{"message":"transport not supported"}],"data":null}`, resp.Body.String())
	})

	t.Run("calls query on executable schema", func(t *testing.T) {
		resp := get(srv, "/foo?query={name}")
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, `{"data":{"name":"test"}}`, resp.Body.String())
	})

	t.Run("mutations are forbidden", func(t *testing.T) {
		resp := get(srv, "/foo?query=mutation{name}")
		assert.Equal(t, http.StatusNotAcceptable, resp.Code)
		assert.Equal(t, `{"errors":[{"message":"GET requests only allow query operations"}],"data":null}`, resp.Body.String())
	})

	t.Run("subscriptions are forbidden", func(t *testing.T) {
		resp := get(srv, "/foo?query=subscription{name}")
		assert.Equal(t, http.StatusNotAcceptable, resp.Code)
		assert.Equal(t, `{"errors":[{"message":"GET requests only allow query operations"}],"data":null}`, resp.Body.String())
	})

	t.Run("invokes operation middleware in order", func(t *testing.T) {
		var calls []string
		srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			calls = append(calls, "first")
			return next(ctx)
		})
		srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			calls = append(calls, "second")
			return next(ctx)
		})

		resp := get(srv, "/foo?query={name}")
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, []string{"first", "second"}, calls)
	})

	t.Run("invokes response middleware in order", func(t *testing.T) {
		var calls []string
		srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			calls = append(calls, "first")
			return next(ctx)
		})
		srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			calls = append(calls, "second")
			return next(ctx)
		})

		resp := get(srv, "/foo?query={name}")
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, []string{"first", "second"}, calls)
	})

	t.Run("invokes field middleware in order", func(t *testing.T) {
		var calls []string
		srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			calls = append(calls, "first")
			return next(ctx)
		})
		srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			calls = append(calls, "second")
			return next(ctx)
		})

		resp := get(srv, "/foo?query={name}")
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, []string{"first", "second"}, calls)
	})

	t.Run("get query parse error in AroundResponses", func(t *testing.T) {
		var errors1 gqlerror.List
		var errors2 gqlerror.List
		srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			resp := next(ctx)
			errors1 = graphql.GetErrors(ctx)
			errors2 = resp.Errors
			return resp
		})

		resp := get(srv, "/foo?query=invalid")
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
		assert.Equal(t, 1, len(errors1))
		assert.Equal(t, 1, len(errors2))
	})

	t.Run("query caching", func(t *testing.T) {
		ctx := context.Background()
		cache := &graphql.MapCache{}
		srv.SetQueryCache(cache)
		qry := `query Foo {name}`

		t.Run("cache miss populates cache", func(t *testing.T) {
			resp := get(srv, "/foo?query="+url.QueryEscape(qry))
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, `{"data":{"name":"test"}}`, resp.Body.String())

			cacheDoc, ok := cache.Get(ctx, qry)
			require.True(t, ok)
			require.Equal(t, "Foo", cacheDoc.(*ast.QueryDocument).Operations[0].Name)
		})

		t.Run("cache hits use document from cache", func(t *testing.T) {
			doc, err := parser.ParseQuery(&ast.Source{Input: `query Bar {name}`})
			require.Nil(t, err)
			cache.Add(ctx, qry, doc)

			resp := get(srv, "/foo?query="+url.QueryEscape(qry))
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, `{"data":{"name":"test"}}`, resp.Body.String())

			cacheDoc, ok := cache.Get(ctx, qry)
			require.True(t, ok)
			require.Equal(t, "Bar", cacheDoc.(*ast.QueryDocument).Operations[0].Name)
		})
	})
}

func TestErrorServer(t *testing.T) {
	srv := testserver.NewError()
	srv.AddTransport(&transport.GET{})

	t.Run("get resolver error in AroundResponses", func(t *testing.T) {
		var errors1 gqlerror.List
		var errors2 gqlerror.List
		srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			resp := next(ctx)
			errors1 = graphql.GetErrors(ctx)
			errors2 = resp.Errors
			return resp
		})

		resp := get(srv, "/foo?query={name}")
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, 1, len(errors1))
		assert.Equal(t, 1, len(errors2))
	})
}

type panicTransport struct{}

func (t panicTransport) Supports(r *http.Request) bool {
	return true
}

func (h panicTransport) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	panic(fmt.Errorf("panic in transport"))
}

func TestRecover(t *testing.T) {
	srv := testserver.New()
	srv.AddTransport(&panicTransport{})

	t.Run("recover from panic", func(t *testing.T) {
		resp := get(srv, "/foo?query={name}")

		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
	})
}

func get(handler http.Handler, target string) *httptest.ResponseRecorder {
	r := httptest.NewRequest("GET", target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)
	return w
}

func post(handler http.Handler, target, contentType string) *httptest.ResponseRecorder {
	r := httptest.NewRequest("POST", target, nil)
	r.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)
	return w
}
