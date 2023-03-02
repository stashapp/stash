package modelgen

import (
	"fmt"
	"go/types"
	"sort"
	"strings"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
)

type BuildMutateHook = func(b *ModelBuild) *ModelBuild

type FieldMutateHook = func(td *ast.Definition, fd *ast.FieldDefinition, f *Field) (*Field, error)

// defaultFieldMutateHook is the default hook for the Plugin which applies the GoTagFieldHook.
func defaultFieldMutateHook(td *ast.Definition, fd *ast.FieldDefinition, f *Field) (*Field, error) {
	return GoTagFieldHook(td, fd, f)
}

func defaultBuildMutateHook(b *ModelBuild) *ModelBuild {
	return b
}

type ModelBuild struct {
	PackageName string
	Interfaces  []*Interface
	Models      []*Object
	Enums       []*Enum
	Scalars     []string
}

type Interface struct {
	Description string
	Name        string
	Implements  []string
}

type Object struct {
	Description string
	Name        string
	Fields      []*Field
	Implements  []string
}

type Field struct {
	Description string
	Name        string
	Type        types.Type
	Tag         string
}

type Enum struct {
	Description string
	Name        string
	Values      []*EnumValue
}

type EnumValue struct {
	Description string
	Name        string
}

func New() plugin.Plugin {
	return &Plugin{
		MutateHook: defaultBuildMutateHook,
		FieldHook:  defaultFieldMutateHook,
	}
}

type Plugin struct {
	MutateHook BuildMutateHook
	FieldHook  FieldMutateHook
}

var _ plugin.ConfigMutator = &Plugin{}

func (m *Plugin) Name() string {
	return "modelgen"
}

func (m *Plugin) MutateConfig(cfg *config.Config) error {
	binder := cfg.NewBinder()

	b := &ModelBuild{
		PackageName: cfg.Model.Package,
	}

	for _, schemaType := range cfg.Schema.Types {
		if cfg.Models.UserDefined(schemaType.Name) {
			continue
		}
		switch schemaType.Kind {
		case ast.Interface, ast.Union:
			it := &Interface{
				Description: schemaType.Description,
				Name:        schemaType.Name,
				Implements:  schemaType.Interfaces,
			}

			b.Interfaces = append(b.Interfaces, it)
		case ast.Object, ast.InputObject:
			if schemaType == cfg.Schema.Query || schemaType == cfg.Schema.Mutation || schemaType == cfg.Schema.Subscription {
				continue
			}
			it := &Object{
				Description: schemaType.Description,
				Name:        schemaType.Name,
			}

			// If Interface A implements interface B, and Interface C also implements interface B
			// then both A and C have methods of B.
			// The reason for checking unique is to prevent the same method B from being generated twice.
			uniqueMap := map[string]bool{}
			for _, implementor := range cfg.Schema.GetImplements(schemaType) {
				if !uniqueMap[implementor.Name] {
					it.Implements = append(it.Implements, implementor.Name)
					uniqueMap[implementor.Name] = true
				}
				// for interface implements
				for _, iface := range implementor.Interfaces {
					if !uniqueMap[iface] {
						it.Implements = append(it.Implements, iface)
						uniqueMap[iface] = true
					}
				}
			}

			for _, field := range schemaType.Fields {
				var typ types.Type
				fieldDef := cfg.Schema.Types[field.Type.Name()]

				if cfg.Models.UserDefined(field.Type.Name()) {
					var err error
					typ, err = binder.FindTypeFromName(cfg.Models[field.Type.Name()].Model[0])
					if err != nil {
						return err
					}
				} else {
					switch fieldDef.Kind {
					case ast.Scalar:
						// no user defined model, referencing a default scalar
						typ = types.NewNamed(
							types.NewTypeName(0, cfg.Model.Pkg(), "string", nil),
							nil,
							nil,
						)

					case ast.Interface, ast.Union:
						// no user defined model, referencing a generated interface type
						typ = types.NewNamed(
							types.NewTypeName(0, cfg.Model.Pkg(), templates.ToGo(field.Type.Name()), nil),
							types.NewInterfaceType([]*types.Func{}, []types.Type{}),
							nil,
						)

					case ast.Enum:
						// no user defined model, must reference a generated enum
						typ = types.NewNamed(
							types.NewTypeName(0, cfg.Model.Pkg(), templates.ToGo(field.Type.Name()), nil),
							nil,
							nil,
						)

					case ast.Object, ast.InputObject:
						// no user defined model, must reference a generated struct
						typ = types.NewNamed(
							types.NewTypeName(0, cfg.Model.Pkg(), templates.ToGo(field.Type.Name()), nil),
							types.NewStruct(nil, nil),
							nil,
						)

					default:
						panic(fmt.Errorf("unknown ast type %s", fieldDef.Kind))
					}
				}

				name := field.Name
				if nameOveride := cfg.Models[schemaType.Name].Fields[field.Name].FieldName; nameOveride != "" {
					name = nameOveride
				}

				typ = binder.CopyModifiersFromAst(field.Type, typ)

				if isStruct(typ) && (fieldDef.Kind == ast.Object || fieldDef.Kind == ast.InputObject) {
					typ = types.NewPointer(typ)
				}

				f := &Field{
					Name:        name,
					Type:        typ,
					Description: field.Description,
					Tag:         `json:"` + field.Name + `"`,
				}

				if m.FieldHook != nil {
					mf, err := m.FieldHook(schemaType, field, f)
					if err != nil {
						return fmt.Errorf("generror: field %v.%v: %w", it.Name, field.Name, err)
					}
					f = mf
				}

				it.Fields = append(it.Fields, f)
			}

			b.Models = append(b.Models, it)
		case ast.Enum:
			it := &Enum{
				Name:        schemaType.Name,
				Description: schemaType.Description,
			}

			for _, v := range schemaType.EnumValues {
				it.Values = append(it.Values, &EnumValue{
					Name:        v.Name,
					Description: v.Description,
				})
			}

			b.Enums = append(b.Enums, it)
		case ast.Scalar:
			b.Scalars = append(b.Scalars, schemaType.Name)
		}
	}
	sort.Slice(b.Enums, func(i, j int) bool { return b.Enums[i].Name < b.Enums[j].Name })
	sort.Slice(b.Models, func(i, j int) bool { return b.Models[i].Name < b.Models[j].Name })
	sort.Slice(b.Interfaces, func(i, j int) bool { return b.Interfaces[i].Name < b.Interfaces[j].Name })

	for _, it := range b.Enums {
		cfg.Models.Add(it.Name, cfg.Model.ImportPath()+"."+templates.ToGo(it.Name))
	}
	for _, it := range b.Models {
		cfg.Models.Add(it.Name, cfg.Model.ImportPath()+"."+templates.ToGo(it.Name))
	}
	for _, it := range b.Interfaces {
		cfg.Models.Add(it.Name, cfg.Model.ImportPath()+"."+templates.ToGo(it.Name))
	}
	for _, it := range b.Scalars {
		cfg.Models.Add(it, "github.com/99designs/gqlgen/graphql.String")
	}

	if len(b.Models) == 0 && len(b.Enums) == 0 && len(b.Interfaces) == 0 && len(b.Scalars) == 0 {
		return nil
	}

	if m.MutateHook != nil {
		b = m.MutateHook(b)
	}

	err := templates.Render(templates.Options{
		PackageName:     cfg.Model.Package,
		Filename:        cfg.Model.Filename,
		Data:            b,
		GeneratedHeader: true,
		Packages:        cfg.Packages,
	})
	if err != nil {
		return err
	}

	// We may have generated code in a package we already loaded, so we reload all packages
	// to allow packages to be compared correctly
	cfg.ReloadAllPackages()

	return nil
}

// GoTagFieldHook applies the goTag directive to the generated Field f. When applying the Tag to the field, the field
// name is used when no value argument is present.
func GoTagFieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *Field) (*Field, error) {
	args := make([]string, 0)
	for _, goTag := range fd.Directives.ForNames("goTag") {
		key := ""
		value := fd.Name

		if arg := goTag.Arguments.ForName("key"); arg != nil {
			if k, err := arg.Value.Value(nil); err == nil {
				key = k.(string)
			}
		}

		if arg := goTag.Arguments.ForName("value"); arg != nil {
			if v, err := arg.Value.Value(nil); err == nil {
				value = v.(string)
			}
		}

		args = append(args, key+":\""+value+"\"")
	}

	if len(args) > 0 {
		f.Tag = f.Tag + " " + strings.Join(args, " ")
	}

	return f, nil
}

func isStruct(t types.Type) bool {
	_, is := t.Underlying().(*types.Struct)
	return is
}
