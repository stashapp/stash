package introspection

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
)

func ParseIntrospectionQuery(query Query) *ast.SchemaDocument {
	var doc ast.SchemaDocument
	typeMap := query.Schema.Types.NameMap()

	doc.Schema = append(doc.Schema, parseSchemaDefinition(query, typeMap))

	for _, typeVale := range typeMap {
		doc.Definitions = append(doc.Definitions, parseTypeSystemDefinition(typeVale))
	}

	for _, directiveValue := range query.Schema.Directives {
		doc.Directives = append(doc.Directives, parseDirectiveDefinition(directiveValue))
	}

	return &doc
}

func parseSchemaDefinition(query Query, typeMap map[string]*FullType) *ast.SchemaDefinition {
	def := ast.SchemaDefinition{}

	def.OperationTypes = append(def.OperationTypes,
		parseOperationTypeDefinitionForQuery(typeMap[*query.Schema.QueryType.Name]),
		parseOperationTypeDefinitionForMutation(typeMap[*query.Schema.MutationType.Name]),
	)

	return &def
}

func parseOperationTypeDefinitionForQuery(fullType *FullType) *ast.OperationTypeDefinition {
	var op ast.OperationTypeDefinition
	op.Operation = ast.Query
	op.Type = *fullType.Name

	return &op
}

func parseOperationTypeDefinitionForMutation(fullType *FullType) *ast.OperationTypeDefinition {
	var op ast.OperationTypeDefinition
	op.Operation = ast.Mutation
	op.Type = *fullType.Name

	return &op
}

func parseDirectiveDefinition(directiveValue *DirectiveType) *ast.DirectiveDefinition {
	args := make(ast.ArgumentDefinitionList, 0, len(directiveValue.Args))
	for _, arg := range directiveValue.Args {
		argumentDefinition := buildInputValue(arg)
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
	}
}

func parseObjectFields(typeVale *FullType) ast.FieldList {
	fieldList := make(ast.FieldList, 0, len(typeVale.Fields))
	for _, field := range typeVale.Fields {
		typ := getType(&field.Type)
		args := make(ast.ArgumentDefinitionList, 0, len(field.Args))
		for _, arg := range field.Args {
			argumentDefinition := buildInputValue(arg)
			args = append(args, argumentDefinition)
		}

		fieldDefinition := &ast.FieldDefinition{
			Description: pointerString(field.Description),
			Name:        field.Name,
			Arguments:   args,
			Type:        typ,
		}
		fieldList = append(fieldList, fieldDefinition)
	}

	return fieldList
}

func parseInputObjectFields(typeVale *FullType) ast.FieldList {
	fieldList := make(ast.FieldList, 0, len(typeVale.InputFields))
	for _, field := range typeVale.InputFields {
		typ := getType(&field.Type)
		fieldDefinition := &ast.FieldDefinition{
			Description: pointerString(field.Description),
			Name:        field.Name,
			Type:        typ,
		}
		fieldList = append(fieldList, fieldDefinition)
	}

	return fieldList
}

func parseObjectTypeDefinition(typeVale *FullType) *ast.Definition {
	fieldList := parseObjectFields(typeVale)
	interfaces := make([]string, 0, len(typeVale.Interfaces))
	for _, intf := range typeVale.Interfaces {
		interfaces = append(interfaces, pointerString(intf.Name))
	}

	enums := make(ast.EnumValueList, 0, len(typeVale.EnumValues))
	for _, enum := range typeVale.EnumValues {
		enumValue := &ast.EnumValueDefinition{
			Description: pointerString(enum.Description),
			Name:        enum.Name,
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
		Position:    nil,
		BuiltIn:     true,
	}
}

func parseInterfaceTypeDefinition(typeVale *FullType) *ast.Definition {
	fieldList := parseObjectFields(typeVale)
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
		Position:    nil,
		BuiltIn:     true,
	}
}

func parseInputObjectTypeDefinition(typeVale *FullType) *ast.Definition {
	fieldList := parseInputObjectFields(typeVale)
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
		Position:    nil,
		BuiltIn:     true,
	}
}

func parseUnionTypeDefinition(typeVale *FullType) *ast.Definition {
	unions := make([]string, 0, len(typeVale.PossibleTypes))
	for _, unionValue := range typeVale.PossibleTypes {
		unions = append(unions, *unionValue.Name)
	}

	return &ast.Definition{
		Kind:        ast.Union,
		Description: pointerString(typeVale.Description),
		Name:        pointerString(typeVale.Name),
		Types:       unions,
		Position:    nil,
		BuiltIn:     true,
	}
}

func parseEnumTypeDefinition(typeVale *FullType) *ast.Definition {
	enums := make(ast.EnumValueList, 0, len(typeVale.EnumValues))
	for _, enum := range typeVale.EnumValues {
		enumValue := &ast.EnumValueDefinition{
			Description: pointerString(enum.Description),
			Name:        enum.Name,
		}
		enums = append(enums, enumValue)
	}

	return &ast.Definition{
		Kind:        ast.Enum,
		Description: pointerString(typeVale.Description),
		Name:        pointerString(typeVale.Name),
		EnumValues:  enums,
		Position:    nil,
		BuiltIn:     true,
	}
}

func parseScalarTypeExtension(typeVale *FullType) *ast.Definition {
	return &ast.Definition{
		Kind:        ast.Scalar,
		Description: pointerString(typeVale.Description),
		Name:        pointerString(typeVale.Name),
		Position:    nil,
		BuiltIn:     true,
	}
}

func parseTypeSystemDefinition(typeVale *FullType) *ast.Definition {
	switch typeVale.Kind {
	case TypeKindScalar:
		return parseScalarTypeExtension(typeVale)
	case TypeKindInterface:
		return parseInterfaceTypeDefinition(typeVale)
	case TypeKindEnum:
		return parseEnumTypeDefinition(typeVale)
	case TypeKindUnion:
		return parseUnionTypeDefinition(typeVale)
	case TypeKindObject:
		return parseObjectTypeDefinition(typeVale)
	case TypeKindInputObject:
		return parseInputObjectTypeDefinition(typeVale)
	case TypeKindList, TypeKindNonNull:
		panic(fmt.Sprintf("not match Kind: %s", typeVale.Kind))
	}

	panic(fmt.Sprintf("not match Kind: %s", typeVale.Kind))
}

func pointerString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func buildInputValue(input *InputValue) *ast.ArgumentDefinition {
	typ := getType(&input.Type)

	var defaultValue *ast.Value
	if input.DefaultValue != nil {
		defaultValue = &ast.Value{
			Raw:  pointerString(input.DefaultValue),
			Kind: ast.Variable,
		}
	}

	return &ast.ArgumentDefinition{
		Description:  pointerString(input.Description),
		Name:         input.Name,
		DefaultValue: defaultValue,
		Type:         typ,
	}
}

func getType(typeRef *TypeRef) *ast.Type {
	if typeRef.Kind == TypeKindList {
		itemRef := typeRef.OfType
		if itemRef == nil {
			panic("Decorated type deeper than introspection query.")
		}

		return ast.ListType(getType(itemRef), nil)
	}

	if typeRef.Kind == TypeKindNonNull {
		nullableRef := typeRef.OfType
		if nullableRef == nil {
			panic("Decorated type deeper than introspection query.")
		}
		nullableType := getType(nullableRef)
		nullableType.NonNull = true

		return nullableType
	}

	return ast.NamedType(pointerString(typeRef.Name), nil)
}
