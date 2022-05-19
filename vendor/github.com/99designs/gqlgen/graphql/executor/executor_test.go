package executor_test

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor/testexecutor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
)

func TestExecutor(t *testing.T) {
	exec := testexecutor.New()

	t.Run("calls query on executable schema", func(t *testing.T) {
		resp := query(exec, "", "{name}")
		assert.Equal(t, `{"name":"test"}`, string(resp.Data))
	})

	t.Run("invokes operation middleware in order", func(t *testing.T) {
		var calls []string
		exec.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			calls = append(calls, "first")
			return next(ctx)
		})
		exec.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			calls = append(calls, "second")
			return next(ctx)
		})

		resp := query(exec, "", "{name}")
		assert.Equal(t, `{"name":"test"}`, string(resp.Data))
		assert.Equal(t, []string{"first", "second"}, calls)
	})

	t.Run("invokes response middleware in order", func(t *testing.T) {
		var calls []string
		exec.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			calls = append(calls, "first")
			return next(ctx)
		})
		exec.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			calls = append(calls, "second")
			return next(ctx)
		})

		resp := query(exec, "", "{name}")
		assert.Equal(t, `{"name":"test"}`, string(resp.Data))
		assert.Equal(t, []string{"first", "second"}, calls)
	})

	t.Run("invokes root field middleware in order", func(t *testing.T) {
		var calls []string
		exec.AroundRootFields(func(ctx context.Context, next graphql.RootResolver) graphql.Marshaler {
			calls = append(calls, "first")
			return next(ctx)
		})
		exec.AroundRootFields(func(ctx context.Context, next graphql.RootResolver) graphql.Marshaler {
			calls = append(calls, "second")
			return next(ctx)
		})

		resp := query(exec, "", "{name}")
		assert.Equal(t, `{"name":"test"}`, string(resp.Data))
		assert.Equal(t, []string{"first", "second"}, calls)
	})

	t.Run("invokes field middleware in order", func(t *testing.T) {
		var calls []string
		exec.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			calls = append(calls, "first")
			return next(ctx)
		})
		exec.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			calls = append(calls, "second")
			return next(ctx)
		})

		resp := query(exec, "", "{name}")
		assert.Equal(t, `{"name":"test"}`, string(resp.Data))
		assert.Equal(t, []string{"first", "second"}, calls)
	})

	t.Run("invokes operation mutators", func(t *testing.T) {
		var calls []string
		exec.Use(&testParamMutator{
			Mutate: func(ctx context.Context, req *graphql.RawParams) *gqlerror.Error {
				calls = append(calls, "param")
				return nil
			},
		})
		exec.Use(&testCtxMutator{
			Mutate: func(ctx context.Context, rc *graphql.OperationContext) *gqlerror.Error {
				calls = append(calls, "context")
				return nil
			},
		})
		resp := query(exec, "", "{name}")
		assert.Equal(t, `{"name":"test"}`, string(resp.Data))
		assert.Equal(t, []string{"param", "context"}, calls)
	})

	t.Run("get query parse error in AroundResponses", func(t *testing.T) {
		var errors1 gqlerror.List
		var errors2 gqlerror.List
		exec.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			resp := next(ctx)
			errors1 = graphql.GetErrors(ctx)
			errors2 = resp.Errors
			return resp
		})

		resp := query(exec, "", "invalid")
		assert.Equal(t, "", string(resp.Data))
		assert.Equal(t, 1, len(resp.Errors))
		assert.Equal(t, 1, len(errors1))
		assert.Equal(t, 1, len(errors2))
	})

	t.Run("query caching", func(t *testing.T) {
		ctx := context.Background()
		cache := &graphql.MapCache{}
		exec.SetQueryCache(cache)
		qry := `query Foo {name}`

		t.Run("cache miss populates cache", func(t *testing.T) {
			resp := query(exec, "Foo", qry)
			assert.Equal(t, `{"name":"test"}`, string(resp.Data))

			cacheDoc, ok := cache.Get(ctx, qry)
			require.True(t, ok)
			require.Equal(t, "Foo", cacheDoc.(*ast.QueryDocument).Operations[0].Name)
		})

		t.Run("cache hits use document from cache", func(t *testing.T) {
			doc, err := parser.ParseQuery(&ast.Source{Input: `query Bar {name}`})
			require.Nil(t, err)
			cache.Add(ctx, qry, doc)

			resp := query(exec, "Bar", qry)
			assert.Equal(t, `{"name":"test"}`, string(resp.Data))

			cacheDoc, ok := cache.Get(ctx, qry)
			require.True(t, ok)
			require.Equal(t, "Bar", cacheDoc.(*ast.QueryDocument).Operations[0].Name)
		})
	})
}

type testParamMutator struct {
	Mutate func(context.Context, *graphql.RawParams) *gqlerror.Error
}

func (m *testParamMutator) ExtensionName() string {
	return "Operation: Mutate Parameters"
}

func (m *testParamMutator) Validate(s graphql.ExecutableSchema) error {
	return nil
}

func (m *testParamMutator) MutateOperationParameters(ctx context.Context, r *graphql.RawParams) *gqlerror.Error {
	return m.Mutate(ctx, r)
}

type testCtxMutator struct {
	Mutate func(context.Context, *graphql.OperationContext) *gqlerror.Error
}

func (m *testCtxMutator) ExtensionName() string {
	return "Operation: Mutate the Context"
}

func (m *testCtxMutator) Validate(s graphql.ExecutableSchema) error {
	return nil
}

func (m *testCtxMutator) MutateOperationContext(ctx context.Context, rc *graphql.OperationContext) *gqlerror.Error {
	return m.Mutate(ctx, rc)
}

func TestErrorServer(t *testing.T) {
	exec := testexecutor.NewError()

	t.Run("get resolver error in AroundResponses", func(t *testing.T) {
		var errors1 gqlerror.List
		var errors2 gqlerror.List
		exec.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			resp := next(ctx)
			errors1 = graphql.GetErrors(ctx)
			errors2 = resp.Errors
			return resp
		})

		resp := query(exec, "", "{name}")
		assert.Equal(t, "null", string(resp.Data))
		assert.Equal(t, 1, len(errors1))
		assert.Equal(t, 1, len(errors2))
	})
}

func query(exec *testexecutor.TestExecutor, op, q string) *graphql.Response {
	ctx := graphql.StartOperationTrace(context.Background())
	now := graphql.Now()
	rc, err := exec.CreateOperationContext(ctx, &graphql.RawParams{
		Query:         q,
		OperationName: op,
		ReadTime: graphql.TraceTiming{
			Start: now,
			End:   now,
		},
	})
	if err != nil {
		return exec.DispatchError(ctx, err)
	}

	resp, ctx2 := exec.DispatchOperation(ctx, rc)
	return resp(ctx2)
}
