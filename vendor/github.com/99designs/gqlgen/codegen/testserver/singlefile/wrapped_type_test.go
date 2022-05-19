package singlefile

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/codegen/testserver/singlefile/otherpkg"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

func TestWrappedTypes(t *testing.T) {
	resolvers := &Stub{}

	c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

	resolvers.QueryResolver.WrappedScalar = func(ctx context.Context) (scalar WrappedScalar, e error) {
		return "hello", nil
	}

	resolvers.QueryResolver.WrappedStruct = func(ctx context.Context) (wrappedStruct *WrappedStruct, e error) {
		wrapped := WrappedStruct(otherpkg.Struct{
			Name: "hello",
		})
		return &wrapped, nil
	}

	resolvers.QueryResolver.WrappedMap = func(ctx context.Context) (wrappedMap WrappedMap, e error) {
		wrapped := WrappedMap(map[string]string{
			"name": "hello",
		})
		return wrapped, nil
	}

	resolvers.QueryResolver.WrappedSlice = func(ctx context.Context) (slice WrappedSlice, err error) {
		wrapped := WrappedSlice([]string{"hello"})
		return wrapped, nil
	}

	resolvers.WrappedMapResolver.Get = func(ctx context.Context, obj WrappedMap, key string) (s string, err error) {
		return obj[key], nil
	}

	resolvers.WrappedSliceResolver.Get = func(ctx context.Context, obj WrappedSlice, idx int) (s string, err error) {
		return obj[idx], nil
	}

	t.Run("wrapped struct", func(t *testing.T) {
		var resp struct {
			WrappedStruct struct {
				Name string
			}
		}

		err := c.Post(`query { wrappedStruct { name } }`, &resp)
		require.NoError(t, err)

		require.Equal(t, "hello", resp.WrappedStruct.Name)
	})

	t.Run("wrapped scalar", func(t *testing.T) {
		var resp struct {
			WrappedScalar string
		}

		err := c.Post(`query { wrappedScalar }`, &resp)
		require.NoError(t, err)

		require.Equal(t, "hello", resp.WrappedScalar)
	})

	t.Run("wrapped map", func(t *testing.T) {
		var resp struct {
			WrappedMap struct {
				Name string
			}
		}

		err := c.Post(`query { wrappedMap { name: get(key: "name") } }`, &resp)
		require.NoError(t, err)

		require.Equal(t, "hello", resp.WrappedMap.Name)
	})

	t.Run("wrapped slice", func(t *testing.T) {
		var resp struct {
			WrappedSlice struct {
				First string
			}
		}

		err := c.Post(`query { wrappedSlice { first: get(idx: 0) } }`, &resp)
		require.NoError(t, err)

		require.Equal(t, "hello", resp.WrappedSlice.First)
	})
}
