package introspection

import (
	"strings"

	"github.com/vektah/gqlparser/ast"
)

type Schema struct {
	schema *ast.Schema
}

func (s *Schema) Types() []Type {
	var types []Type
	for _, typ := range s.schema.Types {
		if strings.HasPrefix(typ.Name, "__") {
			continue
		}
		types = append(types, *WrapTypeFromDef(s.schema, typ))
	}
	return types
}

func (s *Schema) QueryType() *Type {
	return WrapTypeFromDef(s.schema, s.schema.Query)
}

func (s *Schema) MutationType() *Type {
	return WrapTypeFromDef(s.schema, s.schema.Mutation)
}

func (s *Schema) SubscriptionType() *Type {
	return WrapTypeFromDef(s.schema, s.schema.Subscription)
}

func (s *Schema) Directives() []Directive {
	var res []Directive

	for _, d := range s.schema.Directives {
		res = append(res, s.directiveFromDef(d))
	}

	return res
}

func (s *Schema) directiveFromDef(d *ast.DirectiveDefinition) Directive {
	var locs []string
	for _, loc := range d.Locations {
		locs = append(locs, string(loc))
	}

	var args []InputValue
	for _, arg := range d.Arguments {
		args = append(args, InputValue{
			Name:         arg.Name,
			Description:  arg.Description,
			DefaultValue: defaultValue(arg.DefaultValue),
			Type:         WrapTypeFromType(s.schema, arg.Type),
		})
	}

	return Directive{
		Name:        d.Name,
		Description: d.Description,
		Locations:   locs,
		Args:        args,
	}
}
