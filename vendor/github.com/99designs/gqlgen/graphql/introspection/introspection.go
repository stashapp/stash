// introspection implements the spec defined in https://github.com/facebook/graphql/blob/master/spec/Section%204%20--%20Introspection.md#schema-introspection
package introspection

import "github.com/vektah/gqlparser/v2/ast"

type (
	Directive struct {
		Name         string
		description  string
		Locations    []string
		Args         []InputValue
		IsRepeatable bool
	}

	EnumValue struct {
		Name        string
		description string
		deprecation *ast.Directive
	}

	Field struct {
		Name        string
		description string
		Type        *Type
		Args        []InputValue
		deprecation *ast.Directive
	}

	InputValue struct {
		Name         string
		description  string
		DefaultValue *string
		Type         *Type
	}
)

func WrapSchema(schema *ast.Schema) *Schema {
	return &Schema{schema: schema}
}

func (f *EnumValue) Description() *string {
	if f.description == "" {
		return nil
	}
	return &f.description
}

func (f *EnumValue) IsDeprecated() bool {
	return f.deprecation != nil
}

func (f *EnumValue) DeprecationReason() *string {
	if f.deprecation == nil {
		return nil
	}

	reason := f.deprecation.Arguments.ForName("reason")
	if reason == nil {
		return nil
	}

	return &reason.Value.Raw
}

func (f *Field) Description() *string {
	if f.description == "" {
		return nil
	}
	return &f.description
}

func (f *Field) IsDeprecated() bool {
	return f.deprecation != nil
}

func (f *Field) DeprecationReason() *string {
	if f.deprecation == nil {
		return nil
	}

	reason := f.deprecation.Arguments.ForName("reason")
	if reason == nil {
		return nil
	}

	return &reason.Value.Raw
}

func (f *InputValue) Description() *string {
	if f.description == "" {
		return nil
	}
	return &f.description
}

func (f *Directive) Description() *string {
	if f.description == "" {
		return nil
	}
	return &f.description
}
