package codegen

import (
	"fmt"
	"sort"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/99designs/gqlgen/codegen/config"
)

// Data is a unified model of the code to be generated. Plugins may modify this structure to do things like implement
// resolvers or directives automatically (eg grpc, validation)
type Data struct {
	Config *config.Config
	Schema *ast.Schema
	// If a schema is broken up into multiple Data instance, each representing part of the schema,
	// AllDirectives should contain the directives for the entire schema. Directives() can
	// then be used to get the directives that were defined in this Data instance's sources.
	// If a single Data instance is used for the entire schema, AllDirectives and Directives()
	// will be identical.
	// AllDirectives should rarely be used directly.
	AllDirectives   DirectiveList
	Objects         Objects
	Inputs          Objects
	Interfaces      map[string]*Interface
	ReferencedTypes map[string]*config.TypeReference
	ComplexityRoots map[string]*Object

	QueryRoot        *Object
	MutationRoot     *Object
	SubscriptionRoot *Object
}

type builder struct {
	Config     *config.Config
	Schema     *ast.Schema
	Binder     *config.Binder
	Directives map[string]*Directive
}

// Get only the directives which are defined in the config's sources.
func (d *Data) Directives() DirectiveList {
	res := DirectiveList{}
	for k, directive := range d.AllDirectives {
		for _, s := range d.Config.Sources {
			if directive.Position.Src.Name == s.Name {
				res[k] = directive
				break
			}
		}
	}
	return res
}

func BuildData(cfg *config.Config) (*Data, error) {
	// We reload all packages to allow packages to be compared correctly.
	cfg.ReloadAllPackages()

	b := builder{
		Config: cfg,
		Schema: cfg.Schema,
	}

	b.Binder = b.Config.NewBinder()

	var err error
	b.Directives, err = b.buildDirectives()
	if err != nil {
		return nil, err
	}

	dataDirectives := make(map[string]*Directive)
	for name, d := range b.Directives {
		if !d.Builtin {
			dataDirectives[name] = d
		}
	}

	s := Data{
		Config:        cfg,
		AllDirectives: dataDirectives,
		Schema:        b.Schema,
		Interfaces:    map[string]*Interface{},
	}

	for _, schemaType := range b.Schema.Types {
		switch schemaType.Kind {
		case ast.Object:
			obj, err := b.buildObject(schemaType)
			if err != nil {
				return nil, fmt.Errorf("unable to build object definition: %w", err)
			}

			s.Objects = append(s.Objects, obj)
		case ast.InputObject:
			input, err := b.buildObject(schemaType)
			if err != nil {
				return nil, fmt.Errorf("unable to build input definition: %w", err)
			}

			s.Inputs = append(s.Inputs, input)

		case ast.Union, ast.Interface:
			s.Interfaces[schemaType.Name], err = b.buildInterface(schemaType)
			if err != nil {
				return nil, fmt.Errorf("unable to bind to interface: %w", err)
			}
		}
	}

	if s.Schema.Query != nil {
		s.QueryRoot = s.Objects.ByName(s.Schema.Query.Name)
	} else {
		return nil, fmt.Errorf("query entry point missing")
	}

	if s.Schema.Mutation != nil {
		s.MutationRoot = s.Objects.ByName(s.Schema.Mutation.Name)
	}

	if s.Schema.Subscription != nil {
		s.SubscriptionRoot = s.Objects.ByName(s.Schema.Subscription.Name)
	}

	if err := b.injectIntrospectionRoots(&s); err != nil {
		return nil, err
	}

	s.ReferencedTypes = b.buildTypes()

	sort.Slice(s.Objects, func(i, j int) bool {
		return s.Objects[i].Definition.Name < s.Objects[j].Definition.Name
	})

	sort.Slice(s.Inputs, func(i, j int) bool {
		return s.Inputs[i].Definition.Name < s.Inputs[j].Definition.Name
	})

	if b.Binder.SawInvalid {
		// if we have a syntax error, show it
		err := cfg.Packages.Errors()
		if len(err) > 0 {
			return nil, err
		}

		// otherwise show a generic error message
		return nil, fmt.Errorf("invalid types were encountered while traversing the go source code, this probably means the invalid code generated isnt correct. add try adding -v to debug")
	}

	return &s, nil
}

func (b *builder) injectIntrospectionRoots(s *Data) error {
	obj := s.Objects.ByName(b.Schema.Query.Name)
	if obj == nil {
		return fmt.Errorf("root query type must be defined")
	}

	__type, err := b.buildField(obj, &ast.FieldDefinition{
		Name: "__type",
		Type: ast.NamedType("__Type", nil),
		Arguments: []*ast.ArgumentDefinition{
			{
				Name: "name",
				Type: ast.NonNullNamedType("String", nil),
			},
		},
	})
	if err != nil {
		return err
	}

	__schema, err := b.buildField(obj, &ast.FieldDefinition{
		Name: "__schema",
		Type: ast.NamedType("__Schema", nil),
	})
	if err != nil {
		return err
	}

	obj.Fields = append(obj.Fields, __type, __schema)

	return nil
}
