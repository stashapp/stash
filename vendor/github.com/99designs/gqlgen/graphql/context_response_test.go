package graphql

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestAddError(t *testing.T) {
	ctx := WithResponseContext(context.Background(), DefaultErrorPresenter, nil)

	root := &FieldContext{
		Field: CollectedField{
			Field: &ast.Field{
				Alias: "foo",
			},
		},
	}
	ctx = WithFieldContext(ctx, root)
	AddError(ctx, errors.New("foo1"))
	AddError(ctx, errors.New("foo2"))

	index := 1
	child := &FieldContext{
		Parent: root,
		Index:  &index,
	}
	userProvidedPath := &FieldContext{
		Parent: child,
		Field: CollectedField{
			Field: &ast.Field{
				Alias: "works",
			},
		},
	}

	ctx = WithFieldContext(ctx, child)
	AddError(ctx, errors.New("bar"))
	AddError(ctx, &gqlerror.Error{
		Message: "foo3",
		Path:    append(child.Path(), ast.PathName("works")),
	})

	specs := []struct {
		Name     string
		RCtx     *FieldContext
		Messages []string
	}{
		{
			Name:     "with root FieldContext",
			RCtx:     root,
			Messages: []string{"foo1", "foo2"},
		},
		{
			Name:     "with child FieldContext",
			RCtx:     child,
			Messages: []string{"bar"},
		},
		{
			Name:     "with user provided path",
			RCtx:     userProvidedPath,
			Messages: []string{"foo3"},
		},
	}

	for _, spec := range specs {
		t.Run(spec.Name, func(t *testing.T) {
			errList := GetFieldErrors(ctx, spec.RCtx)
			require.Len(t, errList, len(spec.Messages))

			for idx, err := range errList {
				assert.Equal(t, spec.Messages[idx], err.Message)
			}
		})
	}
}

func TestGetErrorFromPresenter(t *testing.T) {
	ctx := WithResponseContext(context.Background(), func(ctx context.Context, err error) *gqlerror.Error {
		errs := GetErrors(ctx)

		// because we are still presenting the error it is not expected to be returned, but this should not deadlock.
		require.Len(t, errs, 0)
		return DefaultErrorPresenter(ctx, err)
	}, nil)

	ctx = WithFieldContext(ctx, &FieldContext{})
	AddError(ctx, errors.New("foo1"))
}
