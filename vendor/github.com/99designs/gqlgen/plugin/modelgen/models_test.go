package modelgen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen/out"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelGeneration(t *testing.T) {
	cfg, err := config.LoadConfig("testdata/gqlgen.yml")
	require.NoError(t, err)
	require.NoError(t, cfg.Init())
	p := Plugin{
		MutateHook: mutateHook,
		FieldHook:  defaultFieldMutateHook,
	}
	require.NoError(t, p.MutateConfig(cfg))

	require.True(t, cfg.Models.UserDefined("MissingTypeNotNull"))
	require.True(t, cfg.Models.UserDefined("MissingTypeNullable"))
	require.True(t, cfg.Models.UserDefined("MissingEnum"))
	require.True(t, cfg.Models.UserDefined("MissingUnion"))
	require.True(t, cfg.Models.UserDefined("MissingInterface"))
	require.True(t, cfg.Models.UserDefined("TypeWithDescription"))
	require.True(t, cfg.Models.UserDefined("EnumWithDescription"))
	require.True(t, cfg.Models.UserDefined("InterfaceWithDescription"))
	require.True(t, cfg.Models.UserDefined("UnionWithDescription"))

	t.Run("no pointer pointers", func(t *testing.T) {
		generated, err := ioutil.ReadFile("./out/generated.go")
		require.NoError(t, err)
		require.NotContains(t, string(generated), "**")
	})

	t.Run("description is generated", func(t *testing.T) {
		node, err := parser.ParseFile(token.NewFileSet(), "./out/generated.go", nil, parser.ParseComments)
		require.NoError(t, err)
		for _, commentGroup := range node.Comments {
			text := commentGroup.Text()
			words := strings.Split(text, " ")
			require.True(t, len(words) > 1, "expected description %q to have more than one word", text)
		}
	})

	t.Run("tags are applied", func(t *testing.T) {
		file, err := ioutil.ReadFile("./out/generated.go")
		require.NoError(t, err)

		fileText := string(file)

		expectedTags := []string{
			`json:"missing2" database:"MissingTypeNotNullmissing2"`,
			`json:"name" database:"MissingInputname"`,
			`json:"missing2" database:"MissingTypeNullablemissing2"`,
			`json:"name" database:"TypeWithDescriptionname"`,
		}

		for _, tag := range expectedTags {
			require.True(t, strings.Contains(fileText, tag))
		}
	})

	t.Run("field hooks are applied", func(t *testing.T) {
		file, err := ioutil.ReadFile("./out/generated.go")
		require.NoError(t, err)

		fileText := string(file)

		expectedTags := []string{
			`json:"name" anotherTag:"tag"`,
			`json:"enum" yetAnotherTag:"12"`,
			`json:"noVal" yaml:"noVal"`,
			`json:"repeated" someTag:"value" repeated:"true"`,
		}

		for _, tag := range expectedTags {
			require.True(t, strings.Contains(fileText, tag))
		}
	})

	t.Run("concrete types implement interface", func(t *testing.T) {
		var _ out.FooBarer = out.FooBarr{}
	})

	t.Run("implemented interfaces", func(t *testing.T) {
		pkg, err := parseAst("out")
		require.NoError(t, err)

		path := filepath.Join("out", "generated.go")
		generated := pkg.Files[path]

		type field struct {
			typ  string
			name string
		}
		cases := []struct {
			name       string
			wantFields []field
		}{
			{
				name: "A",
				wantFields: []field{
					{
						typ:  "method",
						name: "IsA",
					},
				},
			},
			{
				name: "B",
				wantFields: []field{
					{
						typ:  "method",
						name: "IsB",
					},
				},
			},
			{
				name: "C",
				wantFields: []field{
					{
						typ:  "ident",
						name: "A",
					},
					{
						typ:  "method",
						name: "IsC",
					},
				},
			},
			{
				name: "D",
				wantFields: []field{
					{
						typ:  "ident",
						name: "A",
					},
					{
						typ:  "ident",
						name: "B",
					},
					{
						typ:  "method",
						name: "IsD",
					},
				},
			},
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				typeSpec, ok := generated.Scope.Lookup(tc.name).Decl.(*ast.TypeSpec)
				require.True(t, ok)

				fields := typeSpec.Type.(*ast.InterfaceType).Methods.List
				for i, want := range tc.wantFields {
					if want.typ == "ident" {
						ident, ok := fields[i].Type.(*ast.Ident)
						require.True(t, ok)
						assert.Equal(t, want.name, ident.Name)
					}
					if want.typ == "method" {
						require.GreaterOrEqual(t, 1, len(fields[i].Names))
						name := fields[i].Names[0].Name
						assert.Equal(t, want.name, name)
					}
				}
			})
		}
	})

	t.Run("implemented interfaces type CDImplemented", func(t *testing.T) {
		pkg, err := parseAst("out")
		require.NoError(t, err)

		path := filepath.Join("out", "generated.go")
		generated := pkg.Files[path]

		wantMethods := []string{
			"IsA",
			"IsB",
			"IsC",
			"IsD",
		}

		gots := make([]string, 0, len(wantMethods))
		for _, decl := range generated.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				switch funcDecl.Name.Name {
				case "IsA", "IsB", "IsC", "IsD":
					gots = append(gots, funcDecl.Name.Name)
					require.Len(t, funcDecl.Recv.List, 1)
					recvIdent, ok := funcDecl.Recv.List[0].Type.(*ast.Ident)
					require.True(t, ok)
					require.Equal(t, "CDImplemented", recvIdent.Name)
				}
			}
		}

		sort.Strings(gots)
		require.Equal(t, wantMethods, gots)
	})
}

func mutateHook(b *ModelBuild) *ModelBuild {
	for _, model := range b.Models {
		for _, field := range model.Fields {
			field.Tag += ` database:"` + model.Name + field.Name + `"`
		}
	}

	return b
}

func parseAst(path string) (*ast.Package, error) {
	// test setup to parse the types
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	return pkgs["out"], nil
}
