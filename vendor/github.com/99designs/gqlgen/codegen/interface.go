package codegen

import (
	"go/types"

	"github.com/vektah/gqlparser/ast"
)

type Interface struct {
	*ast.Definition
	Type         types.Type
	Implementors []InterfaceImplementor
	InTypemap    bool
}

type InterfaceImplementor struct {
	*ast.Definition

	Interface *Interface
	Type      types.Type
}

func (b *builder) buildInterface(typ *ast.Definition) *Interface {
	obj, err := b.Binder.DefaultUserObject(typ.Name)
	if err != nil {
		panic(err)
	}

	i := &Interface{
		Definition: typ,
		Type:       obj,
		InTypemap:  b.Config.Models.UserDefined(typ.Name),
	}

	for _, implementor := range b.Schema.GetPossibleTypes(typ) {
		obj, err := b.Binder.DefaultUserObject(implementor.Name)
		if err != nil {
			panic(err)
		}

		i.Implementors = append(i.Implementors, InterfaceImplementor{
			Definition: implementor,
			Type:       obj,
			Interface:  i,
		})
	}

	return i
}

func (i *InterfaceImplementor) ValueReceiver() bool {
	interfaceType, err := findGoInterface(i.Interface.Type)
	if interfaceType == nil || err != nil {
		return true
	}

	implementorType, err := findGoNamedType(i.Type)
	if implementorType == nil || err != nil {
		return true
	}

	return types.Implements(implementorType, interfaceType)
}
