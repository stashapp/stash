package codegen

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/stretchr/testify/require"
)

func TestFindField(t *testing.T) {
	input := `
package test

type Std struct {
	Name string
	Value int
}
type Anon struct {
	Name string
	Tags
}
type Tags struct {
	Bar string ` + "`" + `gqlgen:"foo"` + "`" + `
	Foo int    ` + "`" + `gqlgen:"bar"` + "`" + `
}
type Amb struct {
	Bar string ` + "`" + `gqlgen:"foo"` + "`" + `
	Foo int    ` + "`" + `gqlgen:"foo"` + "`" + `
}
type Embed struct {
	Std
	Test string
}
`
	scope, err := parseScope(input, "test")
	require.NoError(t, err)

	std := scope.Lookup("Std").Type().(*types.Named)
	anon := scope.Lookup("Anon").Type().(*types.Named)
	tags := scope.Lookup("Tags").Type().(*types.Named)
	amb := scope.Lookup("Amb").Type().(*types.Named)
	embed := scope.Lookup("Embed").Type().(*types.Named)

	tests := []struct {
		Name        string
		Named       *types.Named
		Field       string
		Tag         string
		Expected    string
		ShouldError bool
	}{
		{"Finds a field by name with no tag", std, "name", "", "Name", false},
		{"Finds a field by name when passed tag but tag not used", std, "name", "gqlgen", "Name", false},
		{"Ignores tags when not passed a tag", tags, "foo", "", "Foo", false},
		{"Picks field with tag over field name when passed a tag", tags, "foo", "gqlgen", "Bar", false},
		{"Errors when ambigious", amb, "foo", "gqlgen", "", true},
		{"Finds a field that is in embedded struct", anon, "bar", "", "Bar", false},
		{"Finds field that is not in embedded struct", embed, "test", "", "Test", false},
	}

	for _, tt := range tests {
		b := builder{Config: &config.Config{StructTag: tt.Tag}}
		target, err := b.findBindTarget(tt.Named, tt.Field)
		if tt.ShouldError {
			require.Nil(t, target, tt.Name)
			require.Error(t, err, tt.Name)
		} else {
			require.NoError(t, err, tt.Name)
			require.Equal(t, tt.Expected, target.Name(), tt.Name)
		}
	}
}

func parseScope(input interface{}, packageName string) (*types.Scope, error) {
	// test setup to parse the types
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", input, 0)
	if err != nil {
		return nil, err
	}

	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check(packageName, fset, []*ast.File{f}, nil)
	if err != nil {
		return nil, err
	}

	return pkg.Scope(), nil
}

func TestEqualFieldName(t *testing.T) {
	tt := []struct {
		Name     string
		Source   string
		Target   string
		Expected bool
	}{
		{Name: "words with same case", Source: "test", Target: "test", Expected: true},
		{Name: "words different case", Source: "test", Target: "tEsT", Expected: true},
		{Name: "different words", Source: "foo", Target: "bar", Expected: false},
		{Name: "separated with underscore", Source: "the_test", Target: "TheTest", Expected: true},
		{Name: "empty values", Source: "", Target: "", Expected: true},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			result := equalFieldName(tc.Source, tc.Target)
			require.Equal(t, tc.Expected, result)
		})
	}
}
