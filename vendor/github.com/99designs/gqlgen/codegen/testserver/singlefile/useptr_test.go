package singlefile

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserPtr(t *testing.T) {
	s := &Stub{}
	r := reflect.TypeOf(s.QueryResolver.OptionalUnion)
	require.True(t, r.Out(0).Kind() == reflect.Interface)
}
