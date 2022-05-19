package followschema

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/stretchr/testify/require"
)

func TestIntrospection(t *testing.T) {
	t.Run("disabled when creating your own server", func(t *testing.T) {
		resolvers := &Stub{}

		srv := handler.New(NewExecutableSchema(Config{Resolvers: resolvers}))
		srv.AddTransport(transport.POST{})
		c := client.New(srv)

		var resp interface{}
		err := c.Post(introspection.Query, &resp)
		require.EqualError(t, err, "[{\"message\":\"introspection disabled\",\"path\":[\"__schema\"]}]")
	})

	t.Run("enabled by default", func(t *testing.T) {
		resolvers := &Stub{}

		c := client.New(handler.NewDefaultServer(
			NewExecutableSchema(Config{Resolvers: resolvers}),
		))

		var resp interface{}
		err := c.Post(introspection.Query, &resp)
		require.NoError(t, err)

		t.Run("does not return empty deprecation strings", func(t *testing.T) {
			q := `{
			  __type(name:"InnerObject") {
			    fields {
			      name
			      deprecationReason
			    }
			  }
			}`

			var resp struct {
				Type struct {
					Fields []struct {
						Name              string
						DeprecationReason *string
					}
				} `json:"__type"`
			}
			err := c.Post(q, &resp)
			require.NoError(t, err)

			require.Equal(t, "id", resp.Type.Fields[0].Name)
			require.Nil(t, resp.Type.Fields[0].DeprecationReason)
		})
	})

	t.Run("disabled by middleware", func(t *testing.T) {
		resolvers := &Stub{}

		srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers}))
		srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			graphql.GetOperationContext(ctx).DisableIntrospection = true
			return next(ctx)
		})
		c := client.New(srv)

		var resp interface{}
		err := c.Post(introspection.Query, &resp)
		require.EqualError(t, err, "[{\"message\":\"introspection disabled\",\"path\":[\"__schema\"]}]")
	})
}
