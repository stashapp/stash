package fieldset

import (
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/vektah/gqlparser/v2/ast"
)

// Set represents a FieldSet that is used in federation directives @key and @requires.
// Would be happier to reuse FieldSet parsing from gqlparser, but this suits for now.
//
type Set []Field

// Field represents a single field in a FieldSet
//
type Field []string

// New parses a FieldSet string into a TinyFieldSet.
//
func New(raw string, prefix []string) Set {
	if !strings.Contains(raw, "{") {
		return parseUnnestedKeyFieldSet(raw, prefix)
	}

	var (
		ret       = Set{}
		subPrefix = prefix
	)
	before, during, after := extractSubs(raw)

	if before != "" {
		befores := New(before, prefix)
		if len(befores) > 0 {
			subPrefix = befores[len(befores)-1]
			ret = append(ret, befores[:len(befores)-1]...)
		}
	}
	if during != "" {
		ret = append(ret, New(during, subPrefix)...)
	}
	if after != "" {
		ret = append(ret, New(after, prefix)...)
	}
	return ret
}

// FieldDefinition looks up a field in the type.
//
func (f Field) FieldDefinition(schemaType *ast.Definition, schema *ast.Schema) *ast.FieldDefinition {
	objType := schemaType
	def := objType.Fields.ForName(f[0])

	for _, part := range f[1:] {
		if objType.Kind != ast.Object {
			panic(fmt.Sprintf(`invalid sub-field reference "%s" in %v: `, objType.Name, f))
		}
		x := def.Type.Name()
		objType = schema.Types[x]
		if objType == nil {
			panic("invalid schema type: " + x)
		}
		def = objType.Fields.ForName(part)
	}
	if def == nil {
		return nil
	}
	ret := *def // shallow copy
	ret.Name = f.ToGoPrivate()

	return &ret
}

// TypeReference looks up the type of a field.
//
func (f Field) TypeReference(obj *codegen.Object, objects codegen.Objects) *codegen.Field {
	var def *codegen.Field

	for _, part := range f {
		def = fieldByName(obj, part)
		if def == nil {
			panic("unable to find field " + f[0])
		}
		obj = objects.ByName(def.TypeReference.Definition.Name)
	}
	return def
}

// ToGo converts a (possibly nested) field into a proper public Go name.
//
func (f Field) ToGo() string {
	var ret string

	for _, field := range f {
		ret += templates.ToGo(field)
	}
	return ret
}

// ToGoPrivate converts a (possibly nested) field into a proper private Go name.
//
func (f Field) ToGoPrivate() string {
	var ret string

	for i, field := range f {
		if i == 0 {
			ret += templates.ToGoPrivate(field)
			continue
		}
		ret += templates.ToGo(field)
	}
	return ret
}

// Join concatenates the field parts with a string separator between. Useful in templates.
//
func (f Field) Join(str string) string {
	return strings.Join(f, str)
}

// JoinGo concatenates the Go name of field parts with a string separator between. Useful in templates.
//
func (f Field) JoinGo(str string) string {
	strs := []string{}

	for _, s := range f {
		strs = append(strs, templates.ToGo(s))
	}
	return strings.Join(strs, str)
}

func (f Field) LastIndex() int {
	return len(f) - 1
}

// local functions

// parseUnnestedKeyFieldSet // handles simple case where none of the fields are nested.
//
func parseUnnestedKeyFieldSet(raw string, prefix []string) Set {
	ret := Set{}

	for _, s := range strings.Fields(raw) {
		next := append(prefix[:], s) //nolint:gocritic // slicing out on purpose
		ret = append(ret, next)
	}
	return ret
}

// extractSubs splits out and trims sub-expressions from before, inside, and after "{}".
//
func extractSubs(str string) (string, string, string) {
	start := strings.Index(str, "{")
	end := matchingBracketIndex(str, start)

	if start < 0 || end < 0 {
		panic("invalid key fieldSet: " + str)
	}
	return strings.TrimSpace(str[:start]), strings.TrimSpace(str[start+1 : end]), strings.TrimSpace(str[end+1:])
}

// matchingBracketIndex returns the index of the closing bracket, assuming an open bracket at start.
//
func matchingBracketIndex(str string, start int) int {
	if start < 0 || len(str) <= start+1 {
		return -1
	}
	var depth int

	for i, c := range str[start+1:] {
		switch c {
		case '{':
			depth++
		case '}':
			if depth == 0 {
				return start + 1 + i
			}
			depth--
		}
	}
	return -1
}

func fieldByName(obj *codegen.Object, name string) *codegen.Field {
	for _, field := range obj.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}
