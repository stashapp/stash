package testutil

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

// WithGQLContext injects a fake gqlgen operation+field context so that
// changesetTranslator works in unit tests without a real GraphQL request.
// inputMap keys are the snake_case field names being treated as "set",
// values are ignored, prefer using `nil`
func WithGQLContext(ctx context.Context, inputMap map[string]interface{}) context.Context {
	const varName = "_testInput"
	if inputMap == nil {
		inputMap = map[string]interface{}{}
	}
	opCtx := &graphql.OperationContext{
		Variables: map[string]interface{}{varName: inputMap},
	}
	fieldCtx := &graphql.FieldContext{
		Field: graphql.CollectedField{
			Field: &ast.Field{
				Definition: &ast.FieldDefinition{
					Arguments: ast.ArgumentDefinitionList{{Name: "input"}},
				},
				Arguments: ast.ArgumentList{
					{
						Name:  "input",
						Value: &ast.Value{Kind: ast.Variable, Raw: varName},
					},
				},
			},
		},
	}
	ctx = graphql.WithOperationContext(ctx, opCtx)
	ctx = graphql.WithFieldContext(ctx, fieldCtx)
	return ctx
}
