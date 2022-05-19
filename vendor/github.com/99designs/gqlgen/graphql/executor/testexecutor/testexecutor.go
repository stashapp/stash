package testexecutor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/executor"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type MockResponse struct {
	Name string `json:"name"`
}

func (mr *MockResponse) UnmarshalGQL(v interface{}) error {
	return nil
}

func (mr *MockResponse) MarshalGQL(w io.Writer) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(mr)
	if err != nil {
		panic(err)
	}

	ba := bytes.NewBuffer(bytes.TrimRight(buf.Bytes(), "\n"))

	fmt.Fprint(w, ba)
}

// New provides a server for use in tests that isn't relying on generated code. It isnt a perfect reproduction of
// a generated server, but it aims to be good enough to test the handler package without relying on codegen.
func New() *TestExecutor {
	next := make(chan struct{})

	schema := gqlparser.MustLoadSchema(&ast.Source{Input: `
    type Query {
      name: String!
      find(id: Int!): String!
    }
    type Mutation {
      name: String!
    }
    type Subscription {
      name: String!
    }
  `})

	exec := &TestExecutor{
		next: next,
	}

	exec.schema = &graphql.ExecutableSchemaMock{
		ExecFunc: func(ctx context.Context) graphql.ResponseHandler {
			rc := graphql.GetOperationContext(ctx)
			switch rc.Operation.Operation {
			case ast.Query:
				ran := false
				return func(ctx context.Context) *graphql.Response {
					if ran {
						return nil
					}
					ran = true
					// Field execution happens inside the generated code, lets simulate some of it.
					ctx = graphql.WithFieldContext(ctx, &graphql.FieldContext{
						Object: "Query",
						Field: graphql.CollectedField{
							Field: &ast.Field{
								Name:       "name",
								Alias:      "name",
								Definition: schema.Types["Query"].Fields.ForName("name"),
							},
						},
					})
					data := graphql.GetOperationContext(ctx).RootResolverMiddleware(ctx, func(ctx context.Context) graphql.Marshaler {
						res, err := graphql.GetOperationContext(ctx).ResolverMiddleware(ctx, func(ctx context.Context) (interface{}, error) {
							// return &graphql.Response{Data: []byte(`{"name":"test"}`)}, nil
							return &MockResponse{Name: "test"}, nil
						})
						if err != nil {
							panic(err)
						}

						return res.(*MockResponse)
					})

					var buf bytes.Buffer
					data.MarshalGQL(&buf)

					return &graphql.Response{Data: buf.Bytes()}
				}
			case ast.Mutation:
				return graphql.OneShot(graphql.ErrorResponse(ctx, "mutations are not supported"))
			case ast.Subscription:
				return func(context context.Context) *graphql.Response {
					select {
					case <-ctx.Done():
						return nil
					case <-next:
						return &graphql.Response{
							Data: []byte(`{"name":"test"}`),
						}
					}
				}
			default:
				return graphql.OneShot(graphql.ErrorResponse(ctx, "unsupported GraphQL operation"))
			}
		},
		SchemaFunc: func() *ast.Schema {
			return schema
		},
		ComplexityFunc: func(typeName string, fieldName string, childComplexity int, args map[string]interface{}) (i int, b bool) {
			return exec.complexity, true
		},
	}

	exec.Executor = executor.New(exec.schema)
	return exec
}

// NewError provides a server for use in resolver error tests that isn't relying on generated code. It isnt a perfect reproduction of
// a generated server, but it aims to be good enough to test the handler package without relying on codegen.
func NewError() *TestExecutor {
	next := make(chan struct{})

	schema := gqlparser.MustLoadSchema(&ast.Source{Input: `
		type Query {
			name: String!
		}
	`})

	exec := &TestExecutor{
		next: next,
	}

	exec.schema = &graphql.ExecutableSchemaMock{
		ExecFunc: func(ctx context.Context) graphql.ResponseHandler {
			rc := graphql.GetOperationContext(ctx)
			switch rc.Operation.Operation {
			case ast.Query:
				ran := false
				return func(ctx context.Context) *graphql.Response {
					if ran {
						return nil
					}
					ran = true

					graphql.AddError(ctx, fmt.Errorf("resolver error"))

					return &graphql.Response{
						Data: []byte(`null`),
					}
				}
			case ast.Mutation:
				return graphql.OneShot(graphql.ErrorResponse(ctx, "mutations are not supported"))
			case ast.Subscription:
				return graphql.OneShot(graphql.ErrorResponse(ctx, "subscription are not supported"))
			default:
				return graphql.OneShot(graphql.ErrorResponse(ctx, "unsupported GraphQL operation"))
			}
		},
		SchemaFunc: func() *ast.Schema {
			return schema
		},
		ComplexityFunc: func(typeName string, fieldName string, childComplexity int, args map[string]interface{}) (i int, b bool) {
			return exec.complexity, true
		},
	}

	exec.Executor = executor.New(exec.schema)
	return exec
}

type TestExecutor struct {
	*executor.Executor
	schema     graphql.ExecutableSchema
	next       chan struct{}
	complexity int
}

func (e *TestExecutor) Schema() graphql.ExecutableSchema {
	return e.schema
}

func (e *TestExecutor) SendNextSubscriptionMessage() {
	select {
	case e.next <- struct{}{}:
	case <-time.After(1 * time.Second):
		fmt.Println("WARNING: no active subscription")
	}
}

func (e *TestExecutor) SetCalculatedComplexity(complexity int) {
	e.complexity = complexity
}
