package imports

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/internal/code"
	"github.com/stretchr/testify/require"
)

func TestPrune(t *testing.T) {
	// prime the packages cache so that it's not considered uninitialized

	b, err := Prune("testdata/unused.go", mustReadFile("testdata/unused.go"), &code.Packages{})
	require.NoError(t, err)
	require.Equal(t, strings.ReplaceAll(string(mustReadFile("testdata/unused.expected.go")), "\r\n", "\n"), string(b))
}

func mustReadFile(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return b
}
