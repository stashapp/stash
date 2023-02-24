package federation

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/99designs/gqlgen/plugin/federation/fieldset"
)

//go:embed federation.gotpl
var federationTemplate string

type federation struct {
	Entities []*Entity
	Version  int
}

// New returns a federation plugin that injects
// federated directives and types into the schema
func New(version int) plugin.Plugin {
	if version == 0 {
		version = 1
	}

	return &federation{Version: version}
}

// Name returns the plugin name
func (f *federation) Name() string {
	return "federation"
}

// MutateConfig mutates the configuration
func (f *federation) MutateConfig(cfg *config.Config) error {
	builtins := config.TypeMap{
		"_Service": {
			Model: config.StringList{
				"github.com/99designs/gqlgen/plugin/federation/fedruntime.Service",
			},
		},
		"_Entity": {
			Model: config.StringList{
				"github.com/99designs/gqlgen/plugin/federation/fedruntime.Entity",
			},
		},
		"Entity": {
			Model: config.StringList{
				"github.com/99designs/gqlgen/plugin/federation/fedruntime.Entity",
			},
		},
		"_Any": {
			Model: config.StringList{"github.com/99designs/gqlgen/graphql.Map"},
		},
	}

	for typeName, entry := range builtins {
		if cfg.Models.Exists(typeName) {
			return fmt.Errorf("%v already exists which must be reserved when Federation is enabled", typeName)
		}
		cfg.Models[typeName] = entry
	}
	cfg.Directives["external"] = config.DirectiveConfig{SkipRuntime: true}
	cfg.Directives["requires"] = config.DirectiveConfig{SkipRuntime: true}
	cfg.Directives["provides"] = config.DirectiveConfig{SkipRuntime: true}
	cfg.Directives["key"] = config.DirectiveConfig{SkipRuntime: true}
	cfg.Directives["extends"] = config.DirectiveConfig{SkipRuntime: true}

	// Federation 2 specific directives
	if f.Version == 2 {
		cfg.Directives["shareable"] = config.DirectiveConfig{SkipRuntime: true}
		cfg.Directives["link"] = config.DirectiveConfig{SkipRuntime: true}
		cfg.Directives["tag"] = config.DirectiveConfig{SkipRuntime: true}
		cfg.Directives["override"] = config.DirectiveConfig{SkipRuntime: true}
		cfg.Directives["inaccessible"] = config.DirectiveConfig{SkipRuntime: true}
	}

	return nil
}

func (f *federation) InjectSourceEarly() *ast.Source {
	input := `
	scalar _Any
	scalar _FieldSet

	directive @external on FIELD_DEFINITION
	directive @requires(fields: _FieldSet!) on FIELD_DEFINITION
	directive @provides(fields: _FieldSet!) on FIELD_DEFINITION
	directive @extends on OBJECT | INTERFACE
`
	// add version-specific changes on key directive, as well as adding the new directives for federation 2
	if f.Version == 1 {
		input += `
	directive @key(fields: _FieldSet!) repeatable on OBJECT | INTERFACE
`
	} else if f.Version == 2 {
		input += `
	directive @key(fields: _FieldSet!, resolvable: Boolean = true) repeatable on OBJECT | INTERFACE
	directive @link(import: [String!], url: String!) repeatable on SCHEMA
	directive @shareable on OBJECT | FIELD_DEFINITION
	directive @tag(name: String!) repeatable on FIELD_DEFINITION | INTERFACE | OBJECT | UNION | ARGUMENT_DEFINITION | SCALAR | ENUM | ENUM_VALUE | INPUT_OBJECT | INPUT_FIELD_DEFINITION
	directive @override(from: String!) on FIELD_DEFINITION
	directive @inaccessible on SCALAR | OBJECT | FIELD_DEFINITION | ARGUMENT_DEFINITION | INTERFACE | UNION | ENUM | ENUM_VALUE | INPUT_OBJECT | INPUT_FIELD_DEFINITION
`
	}
	return &ast.Source{
		Name:    "federation/directives.graphql",
		Input:   input,
		BuiltIn: true,
	}
}

// InjectSources creates a GraphQL Entity type with all
// the fields that had the @key directive
func (f *federation) InjectSourceLate(schema *ast.Schema) *ast.Source {
	f.setEntities(schema)

	var entities, resolvers, entityResolverInputDefinitions string
	for i, e := range f.Entities {
		if i != 0 {
			entities += " | "
		}
		entities += e.Name

		for _, r := range e.Resolvers {
			if e.Multi {
				if entityResolverInputDefinitions != "" {
					entityResolverInputDefinitions += "\n\n"
				}
				entityResolverInputDefinitions += "input " + r.InputType + " {\n"
				for _, keyField := range r.KeyFields {
					entityResolverInputDefinitions += fmt.Sprintf("\t%s: %s\n", keyField.Field.ToGo(), keyField.Definition.Type.String())
				}
				entityResolverInputDefinitions += "}"
				resolvers += fmt.Sprintf("\t%s(reps: [%s!]!): [%s]\n", r.ResolverName, r.InputType, e.Name)
			} else {
				resolverArgs := ""
				for _, keyField := range r.KeyFields {
					resolverArgs += fmt.Sprintf("%s: %s,", keyField.Field.ToGoPrivate(), keyField.Definition.Type.String())
				}
				resolvers += fmt.Sprintf("\t%s(%s): %s!\n", r.ResolverName, resolverArgs, e.Name)
			}
		}
	}

	var blocks []string
	if entities != "" {
		entities = `# a union of all types that use the @key directive
union _Entity = ` + entities
		blocks = append(blocks, entities)
	}

	// resolvers can be empty if a service defines only "empty
	// extend" types.  This should be rare.
	if resolvers != "" {
		if entityResolverInputDefinitions != "" {
			blocks = append(blocks, entityResolverInputDefinitions)
		}
		resolvers = `# fake type to build resolver interfaces for users to implement
type Entity {
	` + resolvers + `
}`
		blocks = append(blocks, resolvers)
	}

	_serviceTypeDef := `type _Service {
  sdl: String
}`
	blocks = append(blocks, _serviceTypeDef)

	var additionalQueryFields string
	// Quote from the Apollo Federation subgraph specification:
	// If no types are annotated with the key directive, then the
	// _Entity union and _entities field should be removed from the schema
	if len(f.Entities) > 0 {
		additionalQueryFields += `  _entities(representations: [_Any!]!): [_Entity]!
`
	}
	// _service field is required in any case
	additionalQueryFields += `  _service: _Service!`

	extendTypeQueryDef := `extend type ` + schema.Query.Name + ` {
` + additionalQueryFields + `
}`
	blocks = append(blocks, extendTypeQueryDef)

	return &ast.Source{
		Name:    "federation/entity.graphql",
		BuiltIn: true,
		Input:   "\n" + strings.Join(blocks, "\n\n") + "\n",
	}
}

func (f *federation) GenerateCode(data *codegen.Data) error {
	if len(f.Entities) > 0 {
		if data.Objects.ByName("Entity") != nil {
			data.Objects.ByName("Entity").Root = true
		}
		for _, e := range f.Entities {
			obj := data.Objects.ByName(e.Def.Name)

			for _, r := range e.Resolvers {
				// fill in types for key fields
				//
				for _, keyField := range r.KeyFields {
					if len(keyField.Field) == 0 {
						fmt.Println(
							"skipping @key field " + keyField.Definition.Name + " in " + r.ResolverName + " in " + e.Def.Name,
						)
						continue
					}
					cgField := keyField.Field.TypeReference(obj, data.Objects)
					keyField.Type = cgField.TypeReference
				}
			}

			// fill in types for requires fields
			//
			for _, reqField := range e.Requires {
				if len(reqField.Field) == 0 {
					fmt.Println("skipping @requires field " + reqField.Name + " in " + e.Def.Name)
					continue
				}
				cgField := reqField.Field.TypeReference(obj, data.Objects)
				reqField.Type = cgField.TypeReference
			}
		}
	}

	return templates.Render(templates.Options{
		PackageName:     data.Config.Federation.Package,
		Filename:        data.Config.Federation.Filename,
		Data:            f,
		GeneratedHeader: true,
		Packages:        data.Config.Packages,
		Template:        federationTemplate,
	})
}

func (f *federation) setEntities(schema *ast.Schema) {
	for _, schemaType := range schema.Types {
		keys, ok := isFederatedEntity(schemaType)
		if !ok {
			continue
		}
		e := &Entity{
			Name:      schemaType.Name,
			Def:       schemaType,
			Resolvers: nil,
			Requires:  nil,
		}

		// Let's process custom entity resolver settings.
		dir := schemaType.Directives.ForName("entityResolver")
		if dir != nil {
			if dirArg := dir.Arguments.ForName("multi"); dirArg != nil {
				if dirVal, err := dirArg.Value.Value(nil); err == nil {
					e.Multi = dirVal.(bool)
				}
			}
		}

		// If our schema has a field with a type defined in
		// another service, then we need to define an "empty
		// extend" of that type in this service, so this service
		// knows what the type is like.  But the graphql-server
		// will never ask us to actually resolve this "empty
		// extend", so we don't require a resolver function for
		// it.  (Well, it will never ask in practice; it's
		// unclear whether the spec guarantees this.  See
		// https://github.com/apollographql/apollo-server/issues/3852
		// ).  Example:
		//    type MyType {
		//       myvar: TypeDefinedInOtherService
		//    }
		//    // Federation needs this type, but
		//    // it doesn't need a resolver for it!
		//    extend TypeDefinedInOtherService @key(fields: "id") {
		//       id: ID @external
		//    }
		if !e.allFieldsAreExternal(f.Version) {
			for _, dir := range keys {
				if len(dir.Arguments) > 2 {
					panic("More than two arguments provided for @key declaration.")
				}
				var arg *ast.Argument

				// since keys are able to now have multiple arguments, we need to check both possible for a possible @key(fields="" fields="")
				for _, a := range dir.Arguments {
					if a.Name == "fields" {
						if arg != nil {
							panic("More than one `fields` provided for @key declaration.")
						}
						arg = a
					}
				}

				keyFieldSet := fieldset.New(arg.Value.Raw, nil)

				keyFields := make([]*KeyField, len(keyFieldSet))
				resolverFields := []string{}
				for i, field := range keyFieldSet {
					def := field.FieldDefinition(schemaType, schema)

					if def == nil {
						panic(fmt.Sprintf("no field for %v", field))
					}

					keyFields[i] = &KeyField{Definition: def, Field: field}
					resolverFields = append(resolverFields, keyFields[i].Field.ToGo())
				}

				resolverFieldsToGo := schemaType.Name + "By" + strings.Join(resolverFields, "And")
				var resolverName string
				if e.Multi {
					resolverFieldsToGo += "s" // Pluralize for better API readability
					resolverName = fmt.Sprintf("findMany%s", resolverFieldsToGo)
				} else {
					resolverName = fmt.Sprintf("find%s", resolverFieldsToGo)
				}

				e.Resolvers = append(e.Resolvers, &EntityResolver{
					ResolverName: resolverName,
					KeyFields:    keyFields,
					InputType:    resolverFieldsToGo + "Input",
				})
			}

			e.Requires = []*Requires{}
			for _, f := range schemaType.Fields {
				dir := f.Directives.ForName("requires")
				if dir == nil {
					continue
				}
				if len(dir.Arguments) != 1 || dir.Arguments[0].Name != "fields" {
					panic("Exactly one `fields` argument needed for @requires declaration.")
				}
				requiresFieldSet := fieldset.New(dir.Arguments[0].Value.Raw, nil)
				for _, field := range requiresFieldSet {
					e.Requires = append(e.Requires, &Requires{
						Name:  field.ToGoPrivate(),
						Field: field,
					})
				}
			}
		}
		f.Entities = append(f.Entities, e)
	}

	// make sure order remains stable across multiple builds
	sort.Slice(f.Entities, func(i, j int) bool {
		return f.Entities[i].Name < f.Entities[j].Name
	})
}

func isFederatedEntity(schemaType *ast.Definition) ([]*ast.Directive, bool) {
	switch schemaType.Kind {
	case ast.Object:
		keys := schemaType.Directives.ForNames("key")
		if len(keys) > 0 {
			return keys, true
		}
	case ast.Interface:
		// TODO: support @key and @extends for interfaces
		if dir := schemaType.Directives.ForName("key"); dir != nil {
			fmt.Printf("@key directive found on \"interface %s\". Will be ignored.\n", schemaType.Name)
		}
		if dir := schemaType.Directives.ForName("extends"); dir != nil {
			panic(
				fmt.Sprintf(
					"@extends directive is not currently supported for interfaces, use \"extend interface %s\" instead.",
					schemaType.Name,
				))
		}
	default:
		// ignore
	}
	return nil, false
}
