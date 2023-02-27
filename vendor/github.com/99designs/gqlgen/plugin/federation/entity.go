package federation

import (
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/federation/fieldset"
	"github.com/vektah/gqlparser/v2/ast"
)

// Entity represents a federated type
// that was declared in the GQL schema.
type Entity struct {
	Name      string // The same name as the type declaration
	Def       *ast.Definition
	Resolvers []*EntityResolver
	Requires  []*Requires
	Multi     bool
}

type EntityResolver struct {
	ResolverName string      // The resolver name, such as FindUserByID
	KeyFields    []*KeyField // The fields declared in @key.
	InputType    string      // The Go generated input type for multi entity resolvers
}

type KeyField struct {
	Definition *ast.FieldDefinition
	Field      fieldset.Field        // len > 1 for nested fields
	Type       *config.TypeReference // The Go representation of that field type
}

// Requires represents an @requires clause
type Requires struct {
	Name  string                // the name of the field
	Field fieldset.Field        // source Field, len > 1 for nested fields
	Type  *config.TypeReference // The Go representation of that field type
}

func (e *Entity) allFieldsAreExternal(federationVersion int) bool {
	for _, field := range e.Def.Fields {
		if !e.isFieldImplicitlyExternal(field, federationVersion) && field.Directives.ForName("external") == nil {
			return false
		}
	}
	return true
}

// In federation v2, key fields are implicitly external.
func (e *Entity) isFieldImplicitlyExternal(field *ast.FieldDefinition, federationVersion int) bool {
	// Key fields are only implicitly external in Federation 2
	if federationVersion != 2 {
		return false
	}
	// TODO: From the spec, it seems like if an entity is not resolvable then it should not only not have a resolver, but should not appear in the _Entitiy union.
	// The current implementation is a less drastic departure from the previous behavior, but should probably be reviewed.
	// See https://www.apollographql.com/docs/federation/subgraph-spec/
	if e.isResolvable() {
		return false
	}
	// If the field is a key field, it is implicitly external
	if e.isKeyField(field) {
		return true
	}

	return false
}

// Determine if the entity is resolvable.
func (e *Entity) isResolvable() bool {
	key := e.Def.Directives.ForName("key")
	if key == nil {
		// If there is no key directive, the entity is resolvable.
		return true
	}
	resolvable := key.Arguments.ForName("resolvable")
	if resolvable == nil {
		// If there is no resolvable argument, the entity is resolvable.
		return true
	}
	// only if resolvable: false has been set on the @key directive do we consider the entity non-resolvable.
	return resolvable.Value.Raw != "false"
}

// Determine if a field is part of the entities key.
func (e *Entity) isKeyField(field *ast.FieldDefinition) bool {
	for _, keyField := range e.keyFields() {
		if keyField == field.Name {
			return true
		}
	}
	return false
}

// Get the key fields for this entity.
func (e *Entity) keyFields() []string {
	key := e.Def.Directives.ForName("key")
	if key == nil {
		return []string{}
	}
	fields := key.Arguments.ForName("fields")
	if fields == nil {
		return []string{}
	}
	fieldSet := fieldset.New(fields.Value.Raw, nil)
	keyFields := make([]string, len(fieldSet))
	for i, field := range fieldSet {
		keyFields[i] = field[0]
	}
	return keyFields
}
