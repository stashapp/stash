package handler

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

type executableSchemaMock struct {
	MutationFunc func(ctx context.Context, op *ast.OperationDefinition) *graphql.Response
}

var _ graphql.ExecutableSchema = &executableSchemaMock{}

func (e *executableSchemaMock) Schema() *ast.Schema {
	return gqlparser.MustLoadSchema(&ast.Source{Input: `
		schema { query: Query, mutation: Mutation }
		type Query {
			empty: String!
		}
		scalar Upload
        type File {
            id: Int!
        }
        input UploadFile {
            id: Int!
            file: Upload!
        }
        type Mutation {
            singleUpload(file: Upload!): File!
            singleUploadWithPayload(req: UploadFile!): File!
            multipleUpload(files: [Upload!]!): [File!]!
            multipleUploadWithPayload(req: [UploadFile!]!): [File!]!
        }
	`})
}

func (e *executableSchemaMock) Complexity(typeName, field string, childComplexity int, args map[string]interface{}) (int, bool) {
	return 0, false
}

func (e *executableSchemaMock) Query(ctx context.Context, op *ast.OperationDefinition) *graphql.Response {
	return graphql.ErrorResponse(ctx, "queries are not supported")
}

func (e *executableSchemaMock) Mutation(ctx context.Context, op *ast.OperationDefinition) *graphql.Response {
	return e.MutationFunc(ctx, op)
}

func (e *executableSchemaMock) Subscription(ctx context.Context, op *ast.OperationDefinition) func() *graphql.Response {
	return func() *graphql.Response {
		<-ctx.Done()
		return nil
	}
}
