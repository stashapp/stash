package graphql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRootFieldContext(t *testing.T) {
	require.Nil(t, GetRootFieldContext(context.Background()))

	rc := &RootFieldContext{}
	require.Equal(t, rc, GetRootFieldContext(WithRootFieldContext(context.Background(), rc)))
}
