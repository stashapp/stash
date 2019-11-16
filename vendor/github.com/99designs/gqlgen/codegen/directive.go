package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/ast"
)

type Directive struct {
	Name    string
	Args    []*FieldArgument
	Builtin bool
}

func (b *builder) buildDirectives() (map[string]*Directive, error) {
	directives := make(map[string]*Directive, len(b.Schema.Directives))

	for name, dir := range b.Schema.Directives {
		if _, ok := directives[name]; ok {
			return nil, errors.Errorf("directive with name %s already exists", name)
		}

		var builtin bool
		if name == "skip" || name == "include" || name == "deprecated" {
			builtin = true
		}

		var args []*FieldArgument
		for _, arg := range dir.Arguments {
			tr, err := b.Binder.TypeReference(arg.Type, nil)
			if err != nil {
				return nil, err
			}

			newArg := &FieldArgument{
				ArgumentDefinition: arg,
				TypeReference:      tr,
				VarName:            templates.ToGoPrivate(arg.Name),
			}

			if arg.DefaultValue != nil {
				var err error
				newArg.Default, err = arg.DefaultValue.Value(nil)
				if err != nil {
					return nil, errors.Errorf("default value for directive argument %s(%s) is not valid: %s", dir.Name, arg.Name, err.Error())
				}
			}
			args = append(args, newArg)
		}

		directives[name] = &Directive{
			Name:    name,
			Args:    args,
			Builtin: builtin,
		}
	}

	return directives, nil
}

func (b *builder) getDirectives(list ast.DirectiveList) ([]*Directive, error) {
	dirs := make([]*Directive, len(list))
	for i, d := range list {
		argValues := make(map[string]interface{}, len(d.Arguments))
		for _, da := range d.Arguments {
			val, err := da.Value.Value(nil)
			if err != nil {
				return nil, err
			}
			argValues[da.Name] = val
		}
		def, ok := b.Directives[d.Name]
		if !ok {
			return nil, fmt.Errorf("directive %s not found", d.Name)
		}

		var args []*FieldArgument
		for _, a := range def.Args {
			value := a.Default
			if argValue, ok := argValues[a.Name]; ok {
				value = argValue
			}
			args = append(args, &FieldArgument{
				ArgumentDefinition: a.ArgumentDefinition,
				Value:              value,
				VarName:            a.VarName,
				TypeReference:      a.TypeReference,
			})
		}
		dirs[i] = &Directive{
			Name: d.Name,
			Args: args,
		}

	}

	return dirs, nil
}

func (d *Directive) ArgsFunc() string {
	if len(d.Args) == 0 {
		return ""
	}

	return "dir_" + d.Name + "_args"
}

func (d *Directive) CallArgs() string {
	args := []string{"ctx", "obj", "n"}

	for _, arg := range d.Args {
		args = append(args, "args["+strconv.Quote(arg.Name)+"].("+templates.CurrentImports.LookupType(arg.TypeReference.GO)+")")
	}

	return strings.Join(args, ", ")
}

func (d *Directive) ResolveArgs(obj string, next string) string {
	args := []string{"ctx", obj, next}

	for _, arg := range d.Args {
		dArg := "&" + arg.VarName
		if !arg.TypeReference.IsPtr() {
			if arg.Value != nil {
				dArg = templates.Dump(arg.Value)
			} else {
				dArg = templates.Dump(arg.Default)
			}
		} else if arg.Value == nil && arg.Default == nil {
			dArg = "nil"
		}

		args = append(args, dArg)
	}

	return strings.Join(args, ", ")
}

func (d *Directive) Declaration() string {
	res := ucFirst(d.Name) + " func(ctx context.Context, obj interface{}, next graphql.Resolver"

	for _, arg := range d.Args {
		res += fmt.Sprintf(", %s %s", arg.Name, templates.CurrentImports.LookupType(arg.TypeReference.GO))
	}

	res += ") (res interface{}, err error)"
	return res
}
