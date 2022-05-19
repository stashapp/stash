package graphql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestGetOperationContext(t *testing.T) {
	rc := &OperationContext{}

	t.Run("with operation context", func(t *testing.T) {
		ctx := WithOperationContext(context.Background(), rc)

		require.True(t, HasOperationContext(ctx))
		require.Equal(t, rc, GetOperationContext(ctx))
	})

	t.Run("without operation context", func(t *testing.T) {
		ctx := context.Background()

		require.False(t, HasOperationContext(ctx))
		require.Panics(t, func() {
			GetOperationContext(ctx)
		})
	})
}

func TestCollectAllFields(t *testing.T) {
	t.Run("collect fields", func(t *testing.T) {
		ctx := testContext(ast.SelectionSet{
			&ast.Field{
				Name: "field",
			},
		})
		s := CollectAllFields(ctx)
		require.Equal(t, []string{"field"}, s)
	})

	t.Run("unique field names", func(t *testing.T) {
		ctx := testContext(ast.SelectionSet{
			&ast.Field{
				Name: "field",
			},
			&ast.Field{
				Name:  "field",
				Alias: "field alias",
			},
		})
		s := CollectAllFields(ctx)
		require.Equal(t, []string{"field"}, s)
	})

	t.Run("collect fragments", func(t *testing.T) {
		ctx := testContext(ast.SelectionSet{
			&ast.Field{
				Name: "fieldA",
			},
			&ast.InlineFragment{
				TypeCondition: "ExampleTypeA",
				SelectionSet: ast.SelectionSet{
					&ast.Field{
						Name: "fieldA",
					},
				},
			},
			&ast.InlineFragment{
				TypeCondition: "ExampleTypeB",
				SelectionSet: ast.SelectionSet{
					&ast.Field{
						Name: "fieldB",
					},
				},
			},
		})
		s := CollectAllFields(ctx)
		require.Equal(t, []string{"fieldA", "fieldB"}, s)
	})

	t.Run("collect fragments with same field name on different types", func(t *testing.T) {
		ctx := testContext(ast.SelectionSet{
			&ast.InlineFragment{
				TypeCondition: "ExampleTypeA",
				SelectionSet: ast.SelectionSet{
					&ast.Field{
						Name:             "fieldA",
						ObjectDefinition: &ast.Definition{Name: "ExampleTypeA"},
					},
				},
			},
			&ast.InlineFragment{
				TypeCondition: "ExampleTypeB",
				SelectionSet: ast.SelectionSet{
					&ast.Field{
						Name:             "fieldA",
						ObjectDefinition: &ast.Definition{Name: "ExampleTypeB"},
					},
				},
			},
		})
		resCtx := GetFieldContext(ctx)
		collected := CollectFields(GetOperationContext(ctx), resCtx.Field.Selections, nil)
		require.Len(t, collected, 2)
		require.NotEqual(t, collected[0], collected[1])
		require.Equal(t, collected[0].Name, collected[1].Name)
	})

	t.Run("collect fragments with same field name and different alias", func(t *testing.T) {
		ctx := testContext(ast.SelectionSet{
			&ast.InlineFragment{
				TypeCondition: "ExampleTypeA",
				SelectionSet: ast.SelectionSet{
					&ast.Field{
						Name:             "fieldA",
						Alias:            "fieldA",
						ObjectDefinition: &ast.Definition{Name: "ExampleTypeA"},
					},
					&ast.Field{
						Name:             "fieldA",
						Alias:            "fieldA Alias",
						ObjectDefinition: &ast.Definition{Name: "ExampleTypeA"},
					},
				},
				ObjectDefinition: &ast.Definition{Name: "ExampleType", Kind: ast.Interface},
			},
		})
		resCtx := GetFieldContext(ctx)
		collected := CollectFields(GetOperationContext(ctx), resCtx.Field.Selections, nil)
		require.Len(t, collected, 2)
		require.NotEqual(t, collected[0], collected[1])
		require.Equal(t, collected[0].Name, collected[1].Name)
		require.NotEqual(t, collected[0].Alias, collected[1].Alias)
	})
}
