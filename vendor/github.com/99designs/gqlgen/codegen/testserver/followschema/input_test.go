package followschema

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

func TestInput(t *testing.T) {
	resolvers := &Stub{}
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers}))
	c := client.New(srv)

	t.Run("when function errors on directives", func(t *testing.T) {
		resolvers.QueryResolver.InputSlice = func(ctx context.Context, arg []string) (b bool, e error) {
			return true, nil
		}

		var resp struct {
			DirectiveArg *string
		}

		err := c.Post(`query { inputSlice(arg: ["ok", 1, 2, "ok"]) }`, &resp)

		require.EqualError(t, err, `http 422: {"errors":[{"message":"String cannot represent a non string value: 1","locations":[{"line":1,"column":32}],"extensions":{"code":"GRAPHQL_VALIDATION_FAILED"}},{"message":"String cannot represent a non string value: 2","locations":[{"line":1,"column":35}],"extensions":{"code":"GRAPHQL_VALIDATION_FAILED"}}],"data":null}`)
		require.Nil(t, resp.DirectiveArg)
	})

	t.Run("when input slice nullable", func(t *testing.T) {
		resolvers.QueryResolver.InputNullableSlice = func(ctx context.Context, arg []string) (b bool, e error) {
			return arg == nil, nil
		}

		var resp struct {
			InputNullableSlice bool
		}
		var err error
		err = c.Post(`query { inputNullableSlice(arg: null) }`, &resp)
		require.NoError(t, err)
		require.True(t, resp.InputNullableSlice)

		err = c.Post(`query { inputNullableSlice(arg: []) }`, &resp)
		require.NoError(t, err)
		require.False(t, resp.InputNullableSlice)
	})

	t.Run("coerce single value to slice", func(t *testing.T) {
		check := func(ctx context.Context, arg []string) (b bool, e error) {
			return len(arg) == 1 && arg[0] == "coerced", nil
		}
		resolvers.QueryResolver.InputSlice = check
		resolvers.QueryResolver.InputNullableSlice = check

		var resp struct {
			Coerced bool
		}
		var err error
		err = c.Post(`query { coerced: inputSlice(arg: "coerced") }`, &resp)
		require.NoError(t, err)
		require.True(t, resp.Coerced)

		err = c.Post(`query { coerced: inputNullableSlice(arg: "coerced") }`, &resp)
		require.NoError(t, err)
		require.True(t, resp.Coerced)
	})
}
