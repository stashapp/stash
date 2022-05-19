package introspection

import (
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

type Schema struct {
	schema *ast.Schema
}

func (s *Schema) Description() *string {
	if s.schema.Description == "" {
		return nil
	}
	return &s.schema.Description
}

func (s *Schema) Types() []Type {
	typeIndex := map[string]Type{}
	typeNames := make([]string, 0, len(s.schema.Types))
	for _, typ := range s.schema.Types {
		if strings.HasPrefix(typ.Name, "__") {
			continue
		}
		typeNames = append(typeNames, typ.Name)
		typeIndex[typ.Name] = *WrapTypeFromDef(s.schema, typ)
	}
	sort.Strings(typeNames)

	types := make([]Type, len(typeNames))
	for i, t := range typeNames {
		types[i] = typeIndex[t]
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
	dIndex := map[string]Directive{}
	dNames := make([]string, 0, len(s.schema.Directives))

	for _, d := range s.schema.Directives {
		dNames = append(dNames, d.Name)
		dIndex[d.Name] = s.directiveFromDef(d)
	}
	sort.Strings(dNames)

	res := make([]Directive, len(dNames))
	for i, d := range dNames {
		res[i] = dIndex[d]
	}

	return res
}

func (s *Schema) directiveFromDef(d *ast.DirectiveDefinition) Directive {
	locs := make([]string, len(d.Locations))
	for i, loc := range d.Locations {
		locs[i] = string(loc)
	}

	args := make([]InputValue, len(d.Arguments))
	for i, arg := range d.Arguments {
		args[i] = InputValue{
			Name:         arg.Name,
			description:  arg.Description,
			DefaultValue: defaultValue(arg.DefaultValue),
			Type:         WrapTypeFromType(s.schema, arg.Type),
		}
	}

	return Directive{
		Name:         d.Name,
		description:  d.Description,
		Locations:    locs,
		Args:         args,
		IsRepeatable: d.IsRepeatable,
	}
}
