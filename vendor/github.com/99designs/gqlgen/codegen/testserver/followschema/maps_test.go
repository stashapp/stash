package followschema

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

func TestMaps(t *testing.T) {
	resolver := &Stub{}
	resolver.QueryResolver.MapStringInterface = func(ctx context.Context, in map[string]interface{}) (i map[string]interface{}, e error) {
		return in, nil
	}
	resolver.QueryResolver.MapNestedStringInterface = func(ctx context.Context, in *NestedMapInput) (i map[string]interface{}, e error) {
		if in == nil {
			return nil, nil
		}
		return in.Map, nil
	}

	c := client.New(handler.NewDefaultServer(
		NewExecutableSchema(Config{Resolvers: resolver}),
	))
	t.Run("unset", func(t *testing.T) {
		var resp struct {
			MapStringInterface map[string]interface{}
		}
		err := c.Post(`query { mapStringInterface { a, b } }`, &resp)
		require.NoError(t, err)
		require.Nil(t, resp.MapStringInterface)
	})

	t.Run("nil", func(t *testing.T) {
		var resp struct {
			MapStringInterface map[string]interface{}
		}
		err := c.Post(`query { mapStringInterface(in: null) { a, b } }`, &resp)
		require.NoError(t, err)
		require.Nil(t, resp.MapStringInterface)
	})

	t.Run("values", func(t *testing.T) {
		var resp struct {
			MapStringInterface map[string]interface{}
		}
		err := c.Post(`query { mapStringInterface(in: { a: "a", b: null }) { a, b } }`, &resp)
		require.NoError(t, err)
		require.Equal(t, "a", resp.MapStringInterface["a"])
		require.Nil(t, resp.MapStringInterface["b"])
	})

	t.Run("nested", func(t *testing.T) {
		var resp struct {
			MapNestedStringInterface map[string]interface{}
		}
		err := c.Post(`query { mapNestedStringInterface(in: { map: { a: "a", b: null } }) { a, b } }`, &resp)
		require.NoError(t, err)
		require.Equal(t, "a", resp.MapNestedStringInterface["a"])
		require.Nil(t, resp.MapNestedStringInterface["b"])
	})

	t.Run("nested nil", func(t *testing.T) {
		var resp struct {
			MapNestedStringInterface map[string]interface{}
		}
		err := c.Post(`query { mapNestedStringInterface(in: { map: null }) { a, b } }`, &resp)
		require.NoError(t, err)
		require.Nil(t, resp.MapNestedStringInterface)
	})
}
