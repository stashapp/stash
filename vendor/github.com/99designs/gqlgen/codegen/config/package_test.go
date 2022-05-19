package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageConfig(t *testing.T) {
	t.Run("when given just a filename", func(t *testing.T) {
		p := PackageConfig{Filename: "testdata/example.go"}
		require.True(t, p.IsDefined())

		require.NoError(t, p.Check())

		require.Equal(t, p.Package, "config_test_data")
		require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.ImportPath())

		require.Equal(t, "config_test_data", p.Pkg().Name())
		require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.Pkg().Path())

		require.Contains(t, filepath.ToSlash(p.Filename), "codegen/config/testdata/example.go")
		require.Contains(t, filepath.ToSlash(p.Dir()), "codegen/config/testdata")
	})

	t.Run("when given both", func(t *testing.T) {
		p := PackageConfig{Filename: "testdata/example.go", Package: "wololo"}
		require.True(t, p.IsDefined())

		require.NoError(t, p.Check())

		require.Equal(t, p.Package, "wololo")
		require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.ImportPath())

		require.Equal(t, "wololo", p.Pkg().Name())
		require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.Pkg().Path())

		require.Contains(t, filepath.ToSlash(p.Filename), "codegen/config/testdata/example.go")
		require.Contains(t, filepath.ToSlash(p.Dir()), "codegen/config/testdata")
	})

	t.Run("when given nothing", func(t *testing.T) {
		p := PackageConfig{}
		require.False(t, p.IsDefined())

		require.EqualError(t, p.Check(), "filename must be specified")

		require.Equal(t, "", p.Package)
		require.Equal(t, "", p.ImportPath())

		require.Nil(t, p.Pkg())

		require.Equal(t, "", p.Filename)
		require.Equal(t, "", p.Dir())
	})

	t.Run("when given invalid filename", func(t *testing.T) {
		p := PackageConfig{Filename: "wololo.sql"}
		require.True(t, p.IsDefined())

		require.EqualError(t, p.Check(), "filename should be path to a go source file")
	})

	t.Run("when package includes a filename", func(t *testing.T) {
		p := PackageConfig{Filename: "foo.go", Package: "foo/foo.go"}
		require.True(t, p.IsDefined())

		require.EqualError(t, p.Check(), "package should be the output package name only, do not include the output filename")
	})
}
