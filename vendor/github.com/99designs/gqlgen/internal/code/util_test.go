package code

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeVendor(t *testing.T) {
	require.Equal(t, "bar/baz", NormalizeVendor("foo/vendor/bar/baz"))
	require.Equal(t, "[]bar/baz", NormalizeVendor("[]foo/vendor/bar/baz"))
	require.Equal(t, "*bar/baz", NormalizeVendor("*foo/vendor/bar/baz"))
	require.Equal(t, "*[]*bar/baz", NormalizeVendor("*[]*foo/vendor/bar/baz"))
	require.Equal(t, "[]*bar/baz", NormalizeVendor("[]*foo/vendor/bar/baz"))
}
