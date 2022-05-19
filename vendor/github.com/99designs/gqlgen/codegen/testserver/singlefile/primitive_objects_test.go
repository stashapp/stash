package singlefile

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
)

func TestPrimitiveObjects(t *testing.T) {
	resolvers := &Stub{}
	resolvers.QueryResolver.PrimitiveObject = func(ctx context.Context) (out []Primitive, e error) {
		return []Primitive{2, 4}, nil
	}

	resolvers.PrimitiveResolver.Value = func(ctx context.Context, obj *Primitive) (i int, e error) {
		return int(*obj), nil
	}

	c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

	t.Run("can fetch value", func(t *testing.T) {
		var resp struct {
			PrimitiveObject []struct {
				Value   int
				Squared int
			}
		}
		c.MustPost(`query { primitiveObject { value, squared } }`, &resp)

		assert.Equal(t, 2, resp.PrimitiveObject[0].Value)
		assert.Equal(t, 4, resp.PrimitiveObject[0].Squared)
		assert.Equal(t, 4, resp.PrimitiveObject[1].Value)
		assert.Equal(t, 16, resp.PrimitiveObject[1].Squared)
	})
}

func TestPrimitiveStringObjects(t *testing.T) {
	resolvers := &Stub{}
	resolvers.QueryResolver.PrimitiveStringObject = func(ctx context.Context) (out []PrimitiveString, e error) {
		return []PrimitiveString{"hello", "world"}, nil
	}

	resolvers.PrimitiveStringResolver.Value = func(ctx context.Context, obj *PrimitiveString) (i string, e error) {
		return string(*obj), nil
	}

	resolvers.PrimitiveStringResolver.Len = func(ctx context.Context, obj *PrimitiveString) (i int, e error) {
		return len(string(*obj)), nil
	}

	c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

	t.Run("can fetch value", func(t *testing.T) {
		var resp struct {
			PrimitiveStringObject []struct {
				Value   string
				Doubled string
				Len     int
			}
		}
		c.MustPost(`query { primitiveStringObject { value, doubled, len } }`, &resp)

		assert.Equal(t, "hello", resp.PrimitiveStringObject[0].Value)
		assert.Equal(t, "hellohello", resp.PrimitiveStringObject[0].Doubled)
		assert.Equal(t, 5, resp.PrimitiveStringObject[0].Len)
		assert.Equal(t, "world", resp.PrimitiveStringObject[1].Value)
		assert.Equal(t, "worldworld", resp.PrimitiveStringObject[1].Doubled)
		assert.Equal(t, 5, resp.PrimitiveStringObject[1].Len)
	})
}
