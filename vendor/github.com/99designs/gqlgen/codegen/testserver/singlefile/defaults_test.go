package singlefile

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

func assertDefaults(t *testing.T, ret *DefaultParametersMirror) {
	require.NotNil(t, ret)
	require.NotNil(t, ret.FalsyBoolean)
	require.Equal(t, *ret.FalsyBoolean, false)
	require.NotNil(t, ret.TruthyBoolean)
	require.Equal(t, *ret.TruthyBoolean, true)
}

func TestDefaults(t *testing.T) {
	resolvers := &Stub{}
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers}))
	c := client.New(srv)

	t.Run("default field parameters", func(t *testing.T) {
		resolvers.QueryResolver.DefaultParameters = func(
			ctx context.Context,
			falsyBoolean, truthyBoolean *bool,
		) (*DefaultParametersMirror, error) {
			return &DefaultParametersMirror{
				FalsyBoolean:  falsyBoolean,
				TruthyBoolean: truthyBoolean,
			}, nil
		}

		var resp struct{ DefaultParameters *DefaultParametersMirror }
		err := c.Post(`query {
			defaultParameters {
				falsyBoolean
				truthyBoolean
			}
		}`, &resp)
		require.NoError(t, err)
		assertDefaults(t, resp.DefaultParameters)
	})

	t.Run("default input fields", func(t *testing.T) {
		resolvers.MutationResolver.DefaultInput = func(
			ctx context.Context,
			input DefaultInput,
		) (*DefaultParametersMirror, error) {
			return &DefaultParametersMirror{
				FalsyBoolean:  input.FalsyBoolean,
				TruthyBoolean: input.TruthyBoolean,
			}, nil
		}

		var resp struct{ DefaultInput *DefaultParametersMirror }
		err := c.Post(`mutation {
			defaultInput(input: {}) {
				falsyBoolean
				truthyBoolean
			}
		}`, &resp)
		require.NoError(t, err)
		assertDefaults(t, resp.DefaultInput)
	})
}
