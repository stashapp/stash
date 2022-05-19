package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/99designs/gqlgen/internal/code"
)

func TestLoadConfig(t *testing.T) {
	t.Run("config does not exist", func(t *testing.T) {
		_, err := LoadConfig("doesnotexist.yml")
		require.Error(t, err)
	})

	t.Run("malformed config", func(t *testing.T) {
		_, err := LoadConfig("testdata/cfg/malformedconfig.yml")
		require.EqualError(t, err, "unable to parse config: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `asdf` into config.Config")
	})

	t.Run("unknown keys", func(t *testing.T) {
		_, err := LoadConfig("testdata/cfg/unknownkeys.yml")
		require.EqualError(t, err, "unable to parse config: yaml: unmarshal errors:\n  line 2: field unknown not found in type config.Config")
	})

	t.Run("globbed filenames", func(t *testing.T) {
		c, err := LoadConfig("testdata/cfg/glob.yml")
		require.NoError(t, err)

		if runtime.GOOS == "windows" {
			require.Equal(t, c.SchemaFilename[0], `testdata\cfg\glob\bar\bar with spaces.graphql`)
			require.Equal(t, c.SchemaFilename[1], `testdata\cfg\glob\foo\foo.graphql`)
		} else {
			require.Equal(t, c.SchemaFilename[0], "testdata/cfg/glob/bar/bar with spaces.graphql")
			require.Equal(t, c.SchemaFilename[1], "testdata/cfg/glob/foo/foo.graphql")
		}
	})

	t.Run("unwalkable path", func(t *testing.T) {
		_, err := LoadConfig("testdata/cfg/unwalkable.yml")
		if runtime.GOOS == "windows" {
			require.EqualError(t, err, "failed to walk schema at root not_walkable/: CreateFile not_walkable/: The system cannot find the file specified.")
		} else {
			require.EqualError(t, err, "failed to walk schema at root not_walkable/: lstat not_walkable/: no such file or directory")
		}
	})
}

func TestLoadConfigFromDefaultLocation(t *testing.T) {
	testDir, err := os.Getwd()
	require.NoError(t, err)
	var cfg *Config

	t.Run("will find closest match", func(t *testing.T) {
		err = os.Chdir(filepath.Join(testDir, "testdata", "cfg", "subdir"))
		require.NoError(t, err)

		cfg, err = LoadConfigFromDefaultLocations()
		require.NoError(t, err)
		require.Equal(t, StringList{"inner"}, cfg.SchemaFilename)
	})

	t.Run("will find config in parent dirs", func(t *testing.T) {
		err = os.Chdir(filepath.Join(testDir, "testdata", "cfg", "otherdir"))
		require.NoError(t, err)

		cfg, err = LoadConfigFromDefaultLocations()
		require.NoError(t, err)
		require.Equal(t, StringList{"outer"}, cfg.SchemaFilename)
	})

	t.Run("will return error if config doesn't exist", func(t *testing.T) {
		err = os.Chdir(testDir)
		require.NoError(t, err)

		cfg, err = LoadConfigFromDefaultLocations()
		require.True(t, errors.Is(err, fs.ErrNotExist))
	})
}

func TestLoadDefaultConfig(t *testing.T) {
	testDir, err := os.Getwd()
	require.NoError(t, err)
	var cfg *Config

	t.Run("will find the schema", func(t *testing.T) {
		err = os.Chdir(filepath.Join(testDir, "testdata", "defaultconfig"))
		require.NoError(t, err)

		cfg, err = LoadDefaultConfig()
		require.NoError(t, err)
		require.NotEmpty(t, cfg.Sources)
	})

	t.Run("will return error if schema doesn't exist", func(t *testing.T) {
		err = os.Chdir(testDir)
		require.NoError(t, err)

		cfg, err = LoadDefaultConfig()
		require.True(t, errors.Is(err, fs.ErrNotExist))
	})
}

func TestReferencedPackages(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tm := TypeMap{
			"Foo": {Model: StringList{"github.com/test.Foo"}},
			"Bar": {Model: StringList{"github.com/test.Bar"}},
			"Baz": {Model: StringList{"github.com/otherpkg.Baz"}},
			"Map": {Model: StringList{"map[string]interface{}"}},
			"SkipResolver": {
				Fields: map[string]TypeMapField{
					"field": {Resolver: false},
				},
			},
		}

		pkgs := tm.ReferencedPackages()

		assert.Equal(t, []string{"github.com/test", "github.com/otherpkg"}, pkgs)
	})
}

func TestConfigCheck(t *testing.T) {
	for _, execLayout := range []ExecLayout{ExecLayoutSingleFile, ExecLayoutFollowSchema} {
		t.Run(string(execLayout), func(t *testing.T) {
			t.Run("invalid config format due to conflicting package names", func(t *testing.T) {
				config := Config{
					Exec:  ExecConfig{Layout: execLayout, Filename: "generated/exec.go", DirName: "generated", Package: "graphql"},
					Model: PackageConfig{Filename: "generated/models.go"},
				}

				require.EqualError(t, config.check(), "exec and model define the same import path (github.com/99designs/gqlgen/codegen/config/generated) with different package names (graphql vs generated)")
			})

			t.Run("federation must be in exec package", func(t *testing.T) {
				config := Config{
					Exec:       ExecConfig{Layout: execLayout, Filename: "generated/exec.go", DirName: "generated"},
					Federation: PackageConfig{Filename: "anotherpkg/federation.go"},
				}

				require.EqualError(t, config.check(), "federation and exec must be in the same package")
			})

			t.Run("federation must have same package name as exec", func(t *testing.T) {
				config := Config{
					Exec:       ExecConfig{Layout: execLayout, Filename: "generated/exec.go", DirName: "generated"},
					Federation: PackageConfig{Filename: "generated/federation.go", Package: "federation"},
				}

				require.EqualError(t, config.check(), "exec and federation define the same import path (github.com/99designs/gqlgen/codegen/config/generated) with different package names (generated vs federation)")
			})

			t.Run("deprecated federated flag raises an error", func(t *testing.T) {
				config := Config{
					Exec:      ExecConfig{Layout: execLayout, Filename: "generated/exec.go", DirName: "generated"},
					Federated: true,
				}

				require.EqualError(t, config.check(), "federated has been removed, instead use\nfederation:\n    filename: path/to/federated.go")
			})
		})
	}
}

func TestAutobinding(t *testing.T) {
	t.Run("valid paths", func(t *testing.T) {
		cfg := Config{
			Models: TypeMap{},
			AutoBind: []string{
				"github.com/99designs/gqlgen/codegen/config/testdata/autobinding/chat",
				"github.com/99designs/gqlgen/codegen/config/testdata/autobinding/scalars/model",
			},
			Packages: &code.Packages{},
		}

		cfg.Schema = gqlparser.MustLoadSchema(&ast.Source{Name: "TestAutobinding.schema", Input: `
			scalar Banned
			type Message { id: ID }
		`})

		require.NoError(t, cfg.autobind())

		require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata/autobinding/scalars/model.Banned", cfg.Models["Banned"].Model[0])
		require.Equal(t, "github.com/99designs/gqlgen/codegen/config/testdata/autobinding/chat.Message", cfg.Models["Message"].Model[0])
	})

	t.Run("with file path", func(t *testing.T) {
		cfg := Config{
			Models: TypeMap{},
			AutoBind: []string{
				"../chat",
			},
			Packages: &code.Packages{},
		}

		cfg.Schema = gqlparser.MustLoadSchema(&ast.Source{Name: "TestAutobinding.schema", Input: `
			scalar Banned
			type Message { id: ID }
		`})

		require.EqualError(t, cfg.autobind(), "unable to load ../chat - make sure you're using an import path to a package that exists")
	})
}
