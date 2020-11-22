package introspection

type TypeKind string

const (
	TypeKindScalar      TypeKind = "SCALAR"
	TypeKindObject      TypeKind = "OBJECT"
	TypeKindInterface   TypeKind = "INTERFACE"
	TypeKindUnion       TypeKind = "UNION"
	TypeKindEnum        TypeKind = "ENUM"
	TypeKindInputObject TypeKind = "INPUT_OBJECT"
	TypeKindList        TypeKind = "LIST"
	TypeKindNonNull     TypeKind = "NON_NULL"
)

type FullTypes []*FullType

func (fs FullTypes) NameMap() map[string]*FullType {
	typeMap := make(map[string]*FullType)
	for _, typ := range fs {
		typeMap[*typ.Name] = typ
	}

	return typeMap
}

type FullType struct {
	Kind        TypeKind
	Name        *string
	Description *string
	Fields      []*FieldValue
	InputFields []*InputValue
	Interfaces  []*TypeRef
	EnumValues  []*struct {
		Name              string
		Description       *string
		IsDeprecated      bool
		DeprecationReason *string
	}
	PossibleTypes []*TypeRef
}

type FieldValue struct {
	Name              string
	Description       *string
	Args              []*InputValue
	Type              TypeRef
	IsDeprecated      bool
	DeprecationReason *string
}

type InputValue struct {
	Name         string
	Description  *string
	Type         TypeRef
	DefaultValue *string
}

type TypeRef struct {
	Kind   TypeKind
	Name   *string
	OfType *TypeRef
}

type Query struct {
	Schema struct {
		QueryType        struct{ Name *string }
		MutationType     *struct{ Name *string }
		SubscriptionType *struct{ Name *string }
		Types            FullTypes
		Directives       []*DirectiveType
	} `graphql:"__schema"`
}

type DirectiveType struct {
	Name        string
	Description *string
	Locations   []string
	Args        []*InputValue
}
