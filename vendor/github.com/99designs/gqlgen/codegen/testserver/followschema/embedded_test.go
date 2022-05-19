package followschema

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

type fakeUnexportedEmbeddedInterface struct{}

func (*fakeUnexportedEmbeddedInterface) UnexportedEmbeddedInterfaceExportedMethod() string {
	return "UnexportedEmbeddedInterfaceExportedMethod"
}

func TestEmbedded(t *testing.T) {
	resolver := &Stub{}
	resolver.QueryResolver.EmbeddedCase1 = func(ctx context.Context) (*EmbeddedCase1, error) {
		return &EmbeddedCase1{}, nil
	}
	resolver.QueryResolver.EmbeddedCase2 = func(ctx context.Context) (*EmbeddedCase2, error) {
		return &EmbeddedCase2{&unexportedEmbeddedPointer{}}, nil
	}
	resolver.QueryResolver.EmbeddedCase3 = func(ctx context.Context) (*EmbeddedCase3, error) {
		return &EmbeddedCase3{&fakeUnexportedEmbeddedInterface{}}, nil
	}

	c := client.New(handler.NewDefaultServer(
		NewExecutableSchema(Config{Resolvers: resolver}),
	))

	t.Run("embedded case 1", func(t *testing.T) {
		var resp struct {
			EmbeddedCase1 struct {
				ExportedEmbeddedPointerExportedMethod string
			}
		}
		err := c.Post(`query { embeddedCase1 { exportedEmbeddedPointerExportedMethod } }`, &resp)
		require.NoError(t, err)
		require.Equal(t, resp.EmbeddedCase1.ExportedEmbeddedPointerExportedMethod, "ExportedEmbeddedPointerExportedMethodResponse")
	})

	t.Run("embedded case 2", func(t *testing.T) {
		var resp struct {
			EmbeddedCase2 struct {
				UnexportedEmbeddedPointerExportedMethod string
			}
		}
		err := c.Post(`query { embeddedCase2 { unexportedEmbeddedPointerExportedMethod } }`, &resp)
		require.NoError(t, err)
		require.Equal(t, resp.EmbeddedCase2.UnexportedEmbeddedPointerExportedMethod, "UnexportedEmbeddedPointerExportedMethodResponse")
	})

	t.Run("embedded case 3", func(t *testing.T) {
		var resp struct {
			EmbeddedCase3 struct {
				UnexportedEmbeddedInterfaceExportedMethod string
			}
		}
		err := c.Post(`query { embeddedCase3 { unexportedEmbeddedInterfaceExportedMethod } }`, &resp)
		require.NoError(t, err)
		require.Equal(t, resp.EmbeddedCase3.UnexportedEmbeddedInterfaceExportedMethod, "UnexportedEmbeddedInterfaceExportedMethod")
	})
}
