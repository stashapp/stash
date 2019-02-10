package handler

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

type executableSchemaStub struct {
	NextResp chan struct{}
}

var _ graphql.ExecutableSchema = &executableSchemaStub{}

func (e *executableSchemaStub) Schema() *ast.Schema {
	return gqlparser.MustLoadSchema(&ast.Source{Input: `
		schema { query: Query }
		type Query {
			me: User!
			user(id: Int): User!
		}
		type User { name: String! }
	`})
}

func (e *executableSchemaStub) Complexity(typeName, field string, childComplexity int, args map[string]interface{}) (int, bool) {
	return 0, false
}

func (e *executableSchemaStub) Query(ctx context.Context, op *ast.OperationDefinition) *graphql.Response {
	return &graphql.Response{Data: []byte(`{"name":"test"}`)}
}

func (e *executableSchemaStub) Mutation(ctx context.Context, op *ast.OperationDefinition) *graphql.Response {
	return graphql.ErrorResponse(ctx, "mutations are not supported")
}

func (e *executableSchemaStub) Subscription(ctx context.Context, op *ast.OperationDefinition) func() *graphql.Response {
	return func() *graphql.Response {
		select {
		case <-ctx.Done():
			return nil
		case <-e.NextResp:
			return &graphql.Response{
				Data: []byte(`{"name":"test"}`),
			}
		}
	}
}
