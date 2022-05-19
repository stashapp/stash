package singlefile

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
)

func TestOk(t *testing.T) {
	resolver := &Stub{}
	resolver.QueryResolver.VOkCaseValue = func(ctx context.Context) (*VOkCaseValue, error) {
		return &VOkCaseValue{}, nil
	}
	resolver.QueryResolver.VOkCaseNil = func(ctx context.Context) (*VOkCaseNil, error) {
		return &VOkCaseNil{}, nil
	}

	c := client.New(handler.NewDefaultServer(
		NewExecutableSchema(Config{Resolvers: resolver}),
	))

	t.Run("v ok case value", func(t *testing.T) {
		var resp struct {
			VOkCaseValue struct {
				Value string
			}
		}
		err := c.Post(`query { vOkCaseValue { value } }`, &resp)
		require.NoError(t, err)
		require.Equal(t, resp.VOkCaseValue.Value, "hi")
	})

	t.Run("v ok case nil", func(t *testing.T) {
		var resp struct {
			VOkCaseNil struct {
				Value *string
			}
		}
		err := c.Post(`query { vOkCaseNil { value } }`, &resp)
		require.NoError(t, err)
		require.Equal(t, true, resp.VOkCaseNil.Value == nil)
	})
}
