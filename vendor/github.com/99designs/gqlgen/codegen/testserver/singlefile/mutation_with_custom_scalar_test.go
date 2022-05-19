package singlefile

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

func TestErrorInsideMutationArgument(t *testing.T) {
	resolvers := &Stub{}
	resolvers.MutationResolver.UpdateSomething = func(_ context.Context, input SpecialInput) (s string, err error) {
		return "Hello world", nil
	}

	c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

	t.Run("mutation with correct input doesn't return error", func(t *testing.T) {
		var resp map[string]interface{}
		input := map[string]interface{}{
			"nesting": map[string]interface{}{
				"field": "email@example.com",
			},
		}
		err := c.Post(
			`mutation TestMutation($input: SpecialInput!) { updateSomething(input: $input) }`,
			&resp,
			client.Var("input", input),
		)
		require.Equal(t, resp["updateSomething"], "Hello world")
		require.NoError(t, err)
	})

	t.Run("mutation with incorrect input returns full path", func(t *testing.T) {
		var resp map[string]interface{}
		input := map[string]interface{}{
			"nesting": map[string]interface{}{
				"field": "not-an-email",
			},
		}
		err := c.Post(
			`mutation TestMutation($input: SpecialInput!) { updateSomething(input: $input) }`,
			&resp,
			client.Var("input", input),
		)
		require.EqualError(t, err, `[{"message":"invalid email format","path":["updateSomething","input","nesting","field"]}]`)
	})
}
