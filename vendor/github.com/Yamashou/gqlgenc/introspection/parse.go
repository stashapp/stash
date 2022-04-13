package introspection

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
)

func ParseIntrospectionQuery(url string, query Query) *ast.SchemaDocument {
	parser := parser{
		sharedPosition: &ast.Position{Src: &ast.Source{
			Name:    "remote",
			BuiltIn: false,
		}},
	}

	if url != "" {
		parser.sharedPosition.Src.Name = url
	}

	return parser.parseIntrospectionQuery(query)
}

type parser struct {
	sharedPosition *ast.Position
}

func (p parser) parseIntrospectionQuery(query Query) *ast.SchemaDocument {
	var doc ast.SchemaDocument
	typeMap := query.Schema.Types.NameMap()

	doc.Schema = append(doc.Schema, p.parseSchemaDefinition(query, typeMap))
	doc.Position = p.sharedPosition

	for _, typeVale := range typeMap {
		doc.Definitions = append(doc.Definitions, p.parseTypeSystemDefinition(typeVale))
	}

	for _, directiveValue := range query.Schema.Directives {
		doc.Directives = append(doc.Directives, p.parseDirectiveDefinition(directiveValue))
	}

	return &doc
}

func (p parser) parseSchemaDefinition(query Query, typeMap map[string]*FullType) *ast.SchemaDefinition {
	def := ast.SchemaDefinition{}
	def.Position = p.sharedPosition

	if query.Schema.QueryType.Name != nil {
		def.OperationTypes = append(def.OperationTypes,
			p.parseOperationTypeDefinitionForQuery(typeMap[*query.Schema.QueryType.Name]),
		)
	}

	if query.Schema.MutationType != nil {
		def.OperationTypes = append(def.OperationTypes,
			p.parseOperationTypeDefinitionForMutation(typeMap[*query.Schema.MutationType.Name]),
		)
	}

	return &def
}

func (p parser) parseOperationTypeDefinitionForQuery(fullType *FullType) *ast.OperationTypeDefinition {
	var op ast.OperationTypeDefinition
	op.Operation = ast.Query
	op.Type = *fullType.Name
	op.Position = p.sharedPosition

	return &op
}

func (p parser) parseOperationTypeDefinitionForMutation(fullType *FullType) *ast.OperationTypeDefinition {
	var op ast.OperationTypeDefinition
	op.Operation = ast.Mutation
	op.Type = *fullType.Name
	op.Position = p.sharedPosition

	return &op
}

func (p parser) parseDirectiveDefinition(directiveValue *DirectiveType) *ast.DirectiveDefinition {
	args := make(ast.ArgumentDefinitionList, 0, len(directiveValue.Args))
	for _, arg := range directiveValue.Args {
		argumentDefinition := p.buildInputValue(arg)
		args = append(args, argumentDefinition)
	}
	locations := make([]ast.DirectiveLocation, 0, len(directiveValue.Locations))
	for _, locationValue := range directiveValue.Locations {
		locations = append(locations, ast.DirectiveLocation(locationValue))
	}

	return &ast.DirectiveDefinition{
		Description: pointerString(directiveValue.Description),
		Name:        directiveValue.Name,
		Arguments:   args,
		Locations:   locations,
		Position:    p.sharedPosition,
	}
}

func (p parser) parseObjectFields(typeVale *FullType) ast.FieldList {
	fieldList := make(ast.FieldList, 0, len(typeVale.Fields))
	for _, field := range typeVale.Fields {
		typ := p.getType(&field.Type)
		args := make(ast.ArgumentDefinitionList, 0, len(field.Args))
		for _, arg := range field.Args {
			argumentDefinition := p.buildInputValue(arg)
			args = append(args, argumentDefinition)
		}

		fieldDefinition := &ast.FieldDefinition{
			Description: pointerString(field.Description),
			Name:        field.Name,
			Arguments:   args,
			Type:        typ,
			Position:    p.sharedPosition,
		}
		fieldList = append(fieldList, fieldDefinition)
	}

	return fieldList
}

func (p parser) parseInputObjectFields(typeVale *FullType) ast.FieldList {
	fieldList := make(ast.FieldList, 0, len(typeVale.InputFields))
	for _, field := range typeVale.InputFields {
		typ := p.getType(&field.Type)
		fieldDefinition := &ast.FieldDefinition{
			Description: pointerString(field.Description),
			Name:        field.Name,
			Type:        typ,
			Position:    p.sharedPosition,
		}
		fieldList = append(fieldList, fieldDefinition)
	}

	return fieldList
}

func (p parser) parseObjectTypeDefinition(typeVale *FullType) *ast.Definition {
	fieldList := p.parseObjectFields(typeVale)
	interfaces := make([]string, 0, len(typeVale.Interfaces))
	for _, intf := range typeVale.Interfaces {
		interfaces = append(interfaces, pointerString(intf.Name))
	}

	enums := make(ast.EnumValueList, 0, len(typeVale.EnumValues))
	for _, enum := range typeVale.EnumValues {
		enumValue := &ast.EnumValueDefinition{
			Description: pointerString(enum.Description),
			Name:        enum.Name,
			Position:    p.sharedPosition,
		}
		enums = append(enums, enumValue)
	}

	return &ast.Definition{
		Kind:        ast.Object,
		Description: pointerString(typeVale.Description),
		Name:        pointerString(typeVale.Name),
		Interfaces:  interfaces,
		Fields:      fieldList,
		EnumValues:  enums,
		Position:    p.sharedPosition,
		BuiltIn:     true,
	}
}

func (p parser) parseInterfaceTypeDefinition(typeVale *FullType) *ast.Definition {
	fieldList := p.parseObjectFields(typeVale)
	interfaces := make([]string, 0, len(typeVale.Interfaces))
	for _, intf := range typeVale.Interfaces {
		interfaces = append(interfaces, pointerString(intf.Name))
	}

	return &ast.Definition{
		Kind:        ast.Interface,
		Description: pointerString(typeVale.Description),
		Name:        pointerString(typeVale.Name),
		Interfaces:  interfaces,
		Fields:      fieldList,
		Position:    p.sharedPosition,
		BuiltIn:     true,
	}
}

func (p parser) parseInputObjectTypeDefinition(typeVale *FullType) *ast.Definition {
	fieldList := p.parseInputObjectFields(typeVale)
	interfaces := make([]string, 0, len(typeVale.Interfaces))
	for _, intf := range typeVale.Interfaces {
		interfaces = append(interfaces, pointerString(intf.Name))
	}

	return &ast.Definition{
		Kind:        ast.InputObject,
		Description: pointerString(typeVale.Description),
		Name:        pointerString(typeVale.Name),
		Interfaces:  interfaces,
		Fields:      fieldList,
		Position:    p.sharedPosition,
		BuiltIn:     true,
	}
}

func (p parser) parseUnionTypeDefinition(typeVale *FullType) *ast.Definition {
	unions := make([]string, 0, len(typeVale.PossibleTypes))
	for _, unionValue := range typeVale.PossibleTypes {
		unions = append(unions, *unionValue.Name)
	}

	return &ast.Definition{
		Kind:        ast.Union,
		Description: pointerString(typeVale.Description),
		Name:        pointerString(typeVale.Name),
		Types:       unions,
		Position:    p.sharedPosition,
		BuiltIn:     true,
	}
}

func (p parser) parseEnumTypeDefinition(typeVale *FullType) *ast.Definition {
	enums := make(ast.EnumValueList, 0, len(typeVale.EnumValues))
	for _, enum := range typeVale.EnumValues {
		enumValue := &ast.EnumValueDefinition{
			Description: pointerString(enum.Description),
			Name:        enum.Name,
			Position:    p.sharedPosition,
		}
		enums = append(enums, enumValue)
	}

	return &ast.Definition{
		Kind:        ast.Enum,
		Description: pointerString(typeVale.Description),
		Name:        pointerString(typeVale.Name),
		EnumValues:  enums,
		Position:    p.sharedPosition,
		BuiltIn:     true,
	}
}

func (p parser) parseScalarTypeExtension(typeVale *FullType) *ast.Definition {
	return &ast.Definition{
		Kind:        ast.Scalar,
		Description: pointerString(typeVale.Description),
		Name:        pointerString(typeVale.Name),
		Position:    p.sharedPosition,
		BuiltIn:     true,
	}
}

func (p parser) parseTypeSystemDefinition(typeVale *FullType) *ast.Definition {
	switch typeVale.Kind {
	case TypeKindScalar:
		return p.parseScalarTypeExtension(typeVale)
	case TypeKindInterface:
		return p.parseInterfaceTypeDefinition(typeVale)
	case TypeKindEnum:
		return p.parseEnumTypeDefinition(typeVale)
	case TypeKindUnion:
		return p.parseUnionTypeDefinition(typeVale)
	case TypeKindObject:
		return p.parseObjectTypeDefinition(typeVale)
	case TypeKindInputObject:
		return p.parseInputObjectTypeDefinition(typeVale)
	case TypeKindList, TypeKindNonNull:
		panic(fmt.Sprintf("not match Kind: %s", typeVale.Kind))
	}

	panic(fmt.Sprintf("not match Kind: %s", typeVale.Kind))
}

func (p parser) buildInputValue(input *InputValue) *ast.ArgumentDefinition {
	typ := p.getType(&input.Type)

	var defaultValue *ast.Value
	if input.DefaultValue != nil {
		defaultValue = &ast.Value{
			Raw:      pointerString(input.DefaultValue),
			Kind:     ast.Variable,
			Position: p.sharedPosition,
		}
	}

	return &ast.ArgumentDefinition{
		Description:  pointerString(input.Description),
		Name:         input.Name,
		DefaultValue: defaultValue,
		Type:         typ,
		Position:     p.sharedPosition,
	}
}

func (p parser) getType(typeRef *TypeRef) *ast.Type {
	if typeRef.Kind == TypeKindList {
		itemRef := typeRef.OfType
		if itemRef == nil {
			panic("Decorated type deeper than introspection query.")
		}

		return ast.ListType(p.getType(itemRef), p.sharedPosition)
	}

	if typeRef.Kind == TypeKindNonNull {
		nullableRef := typeRef.OfType
		if nullableRef == nil {
			panic("Decorated type deeper than introspection query.")
		}
		nullableType := p.getType(nullableRef)
		nullableType.NonNull = true

		return nullableType
	}

	return ast.NamedType(pointerString(typeRef.Name), p.sharedPosition)
}

func pointerString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
