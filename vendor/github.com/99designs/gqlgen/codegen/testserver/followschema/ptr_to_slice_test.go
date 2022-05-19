package followschema

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

func TestPtrToSlice(t *testing.T) {
	resolvers := &Stub{}

	c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

	resolvers.QueryResolver.PtrToSliceContainer = func(ctx context.Context) (wrappedStruct *PtrToSliceContainer, e error) {
		ptrToSliceContainer := PtrToSliceContainer{
			PtrToSlice: &[]string{"hello"},
		}
		return &ptrToSliceContainer, nil
	}

	t.Run("pointer to slice", func(t *testing.T) {
		var resp struct {
			PtrToSliceContainer struct {
				PtrToSlice []string
			}
		}

		err := c.Post(`query { ptrToSliceContainer {  ptrToSlice }}`, &resp)
		require.NoError(t, err)

		require.Equal(t, []string{"hello"}, resp.PtrToSliceContainer.PtrToSlice)
	})
}
