package graphql

import (
	"context"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const fieldInputCtx key = "field_input_context"

type FieldInputContext struct {
	ParentField *FieldContext
	ParentInput *FieldInputContext
	Field       *string
	Index       *int
}

func (fic *FieldInputContext) Path() ast.Path {
	var inputPath ast.Path
	for it := fic; it != nil; it = it.ParentInput {
		if it.Index != nil {
			inputPath = append(inputPath, ast.PathIndex(*it.Index))
		} else if it.Field != nil {
			inputPath = append(inputPath, ast.PathName(*it.Field))
		}
	}

	// because we are walking up the chain, all the elements are backwards, do an inplace flip.
	for i := len(inputPath)/2 - 1; i >= 0; i-- {
		opp := len(inputPath) - 1 - i
		inputPath[i], inputPath[opp] = inputPath[opp], inputPath[i]
	}

	if fic.ParentField != nil {
		fieldPath := fic.ParentField.Path()
		return append(fieldPath, inputPath...)

	}

	return inputPath
}

func NewFieldInputWithField(field string) *FieldInputContext {
	return &FieldInputContext{Field: &field}
}

func NewFieldInputWithIndex(index int) *FieldInputContext {
	return &FieldInputContext{Index: &index}
}

func WithFieldInputContext(ctx context.Context, fic *FieldInputContext) context.Context {
	if fieldContext := GetFieldContext(ctx); fieldContext != nil {
		fic.ParentField = fieldContext
	}
	if fieldInputContext := GetFieldInputContext(ctx); fieldInputContext != nil {
		fic.ParentInput = fieldInputContext
	}

	return context.WithValue(ctx, fieldInputCtx, fic)
}

func GetFieldInputContext(ctx context.Context) *FieldInputContext {
	if val, ok := ctx.Value(fieldInputCtx).(*FieldInputContext); ok {
		return val
	}
	return nil
}

func WrapErrorWithInputPath(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	inputContext := GetFieldInputContext(ctx)
	path := inputContext.Path()
	if gerr, ok := err.(*gqlerror.Error); ok {
		if gerr.Path == nil {
			gerr.Path = path
		}
		return gerr
	} else {
		return gqlerror.WrapPath(path, err)
	}
}
