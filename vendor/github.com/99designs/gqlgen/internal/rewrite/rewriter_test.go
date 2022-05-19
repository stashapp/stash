package rewrite

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRewriter(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		r, err := New("testdata")
		require.NoError(t, err)

		body := r.GetMethodBody("Foo", "Method")
		require.Equal(t, `
	// leading comment

	// field comment
	m.Field++

	// trailing comment
`, strings.ReplaceAll(body, "\r\n", "\n"))

		imps := r.ExistingImports("testdata/example.go")
		require.Len(t, imps, 2)
		assert.Equal(t, []Import{
			{
				Alias:      "",
				ImportPath: "fmt",
			},
			{
				Alias:      "lol",
				ImportPath: "bytes",
			},
		}, imps)
	})

	t.Run("out of scope dir", func(t *testing.T) {
		_, err := New("../../../out-of-gomod/package")
		require.Error(t, err)
	})
}
