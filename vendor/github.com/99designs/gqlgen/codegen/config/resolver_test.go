package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolverConfig(t *testing.T) {
	t.Run("single-file", func(t *testing.T) {
		t.Run("when given just a filename", func(t *testing.T) {
			p := ResolverConfig{Filename: "testdata/example.go"}
			require.True(t, p.IsDefined())

			require.NoError(t, p.Check())
			require.NoError(t, p.Check())

			require.Equal(t, p.Package, "config_test_data")
			require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.ImportPath())

			require.Equal(t, "config_test_data", p.Pkg().Name())
			require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.Pkg().Path())

			require.Contains(t, filepath.ToSlash(p.Filename), "codegen/config/testdata/example.go")
			require.Contains(t, filepath.ToSlash(p.Dir()), "codegen/config/testdata")
		})

		t.Run("when given both", func(t *testing.T) {
			p := ResolverConfig{Filename: "testdata/example.go", Package: "wololo"}
			require.True(t, p.IsDefined())

			require.NoError(t, p.Check())
			require.NoError(t, p.Check())

			require.Equal(t, p.Package, "wololo")
			require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.ImportPath())

			require.Equal(t, "wololo", p.Pkg().Name())
			require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.Pkg().Path())

			require.Contains(t, filepath.ToSlash(p.Filename), "codegen/config/testdata/example.go")
			require.Contains(t, filepath.ToSlash(p.Dir()), "codegen/config/testdata")
		})

		t.Run("when given nothing", func(t *testing.T) {
			p := ResolverConfig{}
			require.False(t, p.IsDefined())

			require.EqualError(t, p.Check(), "filename must be specified with layout=single-file")

			require.Equal(t, "", p.Package)
			require.Equal(t, "", p.ImportPath())

			require.Nil(t, p.Pkg())

			require.Equal(t, "", p.Filename)
			require.Equal(t, "", p.Dir())
		})

		t.Run("when given invalid filename", func(t *testing.T) {
			p := ResolverConfig{Filename: "wololo.sql"}
			require.True(t, p.IsDefined())

			require.EqualError(t, p.Check(), "filename should be path to a go source file with layout=single-file")
		})

		t.Run("when package includes a filename", func(t *testing.T) {
			p := ResolverConfig{Filename: "foo.go", Package: "foo/foo.go"}
			require.True(t, p.IsDefined())

			require.EqualError(t, p.Check(), "package should be the output package name only, do not include the output filename")
		})
	})

	t.Run("follow-schema", func(t *testing.T) {
		t.Run("when given just a dir", func(t *testing.T) {
			p := ResolverConfig{Layout: LayoutFollowSchema, DirName: "testdata"}
			require.True(t, p.IsDefined())

			require.NoError(t, p.Check())
			require.NoError(t, p.Check())

			require.Equal(t, p.Package, "config_test_data")
			require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.ImportPath())

			require.Equal(t, "config_test_data", p.Pkg().Name())
			require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.Pkg().Path())

			require.Contains(t, filepath.ToSlash(p.Filename), "codegen/config/testdata/resolver.go")
			require.Contains(t, p.Dir(), "codegen/config/testdata")
		})

		t.Run("when given dir and package name", func(t *testing.T) {
			p := ResolverConfig{Layout: LayoutFollowSchema, DirName: "testdata", Package: "wololo"}
			require.True(t, p.IsDefined())

			require.NoError(t, p.Check())
			require.NoError(t, p.Check())

			require.Equal(t, p.Package, "wololo")
			require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.ImportPath())

			require.Equal(t, "wololo", p.Pkg().Name())
			require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.Pkg().Path())

			require.Contains(t, filepath.ToSlash(p.Filename), "codegen/config/testdata/resolver.go")
			require.Contains(t, p.Dir(), "codegen/config/testdata")
		})

		t.Run("when given a filename", func(t *testing.T) {
			p := ResolverConfig{Layout: LayoutFollowSchema, DirName: "testdata", Filename: "testdata/asdf.go"}
			require.True(t, p.IsDefined())

			require.NoError(t, p.Check())
			require.NoError(t, p.Check())

			require.Equal(t, p.Package, "config_test_data")
			require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.ImportPath())

			require.Equal(t, "config_test_data", p.Pkg().Name())
			require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata", p.Pkg().Path())

			require.Contains(t, filepath.ToSlash(p.Filename), "codegen/config/testdata/asdf.go")
			require.Contains(t, p.Dir(), "codegen/config/testdata")
		})

		t.Run("when given nothing", func(t *testing.T) {
			p := ResolverConfig{Layout: LayoutFollowSchema}
			require.False(t, p.IsDefined())

			require.EqualError(t, p.Check(), "dirname must be specified with layout=follow-schema")

			require.Equal(t, "", p.Package)
			require.Equal(t, "", p.ImportPath())

			require.Nil(t, p.Pkg())

			require.Equal(t, "", p.Filename)
			require.Equal(t, "", p.Dir())
		})
	})

	t.Run("invalid layout", func(t *testing.T) {
		p := ResolverConfig{Layout: "pies", Filename: "asdf.go"}
		require.True(t, p.IsDefined())

		require.EqualError(t, p.Check(), "invalid layout pies. must be single-file or follow-schema")
	})
}
