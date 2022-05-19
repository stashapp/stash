package graphql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFieldInputContext(t *testing.T) {
	require.Nil(t, GetFieldContext(context.Background()))

	rc := &PathContext{}
	require.Equal(t, rc, GetPathContext(WithPathContext(context.Background(), rc)))
}
