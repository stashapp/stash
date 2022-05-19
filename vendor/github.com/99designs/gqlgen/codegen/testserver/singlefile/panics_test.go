package singlefile

import (
	"context"
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/graphql"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestPanics(t *testing.T) {
	resolvers := &Stub{}
	resolvers.QueryResolver.Panics = func(ctx context.Context) (panics *Panics, e error) {
		return &Panics{}, nil
	}
	resolvers.PanicsResolver.ArgUnmarshal = func(ctx context.Context, obj *Panics, u []MarshalPanic) (b bool, e error) {
		return true, nil
	}
	resolvers.PanicsResolver.FieldScalarMarshal = func(ctx context.Context, obj *Panics) (marshalPanic []MarshalPanic, e error) {
		return []MarshalPanic{MarshalPanic("aa"), MarshalPanic("bb")}, nil
	}

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers}))
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) (userMessage error) {
		return fmt.Errorf("panic: %v", err)
	})

	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		return &gqlerror.Error{
			Message: "presented: " + err.Error(),
			Path:    graphql.GetPath(ctx),
		}
	})

	c := client.New(srv)

	t.Run("panics in marshallers will not kill server", func(t *testing.T) {
		var resp interface{}
		err := c.Post(`query { panics { fieldScalarMarshal } }`, &resp)

		require.EqualError(t, err, "http 422: {\"errors\":[{\"message\":\"presented: panic: BOOM\"}],\"data\":null}")
	})

	t.Run("panics in unmarshalers will not kill server", func(t *testing.T) {
		var resp interface{}
		err := c.Post(`query { panics { argUnmarshal(u: ["aa", "bb"]) } }`, &resp)

		require.EqualError(t, err, "[{\"message\":\"presented: input: panics.argUnmarshal panic: BOOM\",\"path\":[\"panics\",\"argUnmarshal\"]}]")
	})

	t.Run("panics in funcs unmarshal return errors", func(t *testing.T) {
		var resp interface{}
		err := c.Post(`query { panics { fieldFuncMarshal(u: ["aa", "bb"]) } }`, &resp)

		require.EqualError(t, err, "[{\"message\":\"presented: input: panics.fieldFuncMarshal panic: BOOM\",\"path\":[\"panics\",\"fieldFuncMarshal\"]}]")
	})

	t.Run("panics in funcs marshal return errors", func(t *testing.T) {
		var resp interface{}
		err := c.Post(`query { panics { fieldFuncMarshal(u: []) } }`, &resp)

		require.EqualError(t, err, "http 422: {\"errors\":[{\"message\":\"presented: panic: BOOM\"}],\"data\":null}")
	})
}
