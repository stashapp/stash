package template

import (
	"fmt"
	"strings"

	"github.com/matryer/moq/internal/registry"
)

// Data is the template data used to render the Moq template.
type Data struct {
	PkgName         string
	SrcPkgQualifier string
	Imports         []*registry.Package
	Mocks           []MockData
	StubImpl        bool
	SkipEnsure      bool
}

// MocksSomeMethod returns true of any one of the Mocks has at least 1
// method.
func (d Data) MocksSomeMethod() bool {
	for _, m := range d.Mocks {
		if len(m.Methods) > 0 {
			return true
		}
	}

	return false
}

// MockData is the data used to generate a mock for some interface.
type MockData struct {
	InterfaceName string
	MockName      string
	Methods       []MethodData
}

// MethodData is the data which represents a method on some interface.
type MethodData struct {
	Name    string
	Params  []ParamData
	Returns []ParamData
}

// ArgList is the string representation of method parameters, ex:
// 's string, n int, foo bar.Baz'.
func (m MethodData) ArgList() string {
	params := make([]string, len(m.Params))
	for i, p := range m.Params {
		params[i] = p.MethodArg()
	}
	return strings.Join(params, ", ")
}

// ArgCallList is the string representation of method call parameters,
// ex: 's, n, foo'. In case of a last variadic parameter, it will be of
// the format 's, n, foos...'
func (m MethodData) ArgCallList() string {
	params := make([]string, len(m.Params))
	for i, p := range m.Params {
		params[i] = p.CallName()
	}
	return strings.Join(params, ", ")
}

// ReturnArgTypeList is the string representation of method return
// types, ex: 'bar.Baz', '(string, error)'.
func (m MethodData) ReturnArgTypeList() string {
	params := make([]string, len(m.Returns))
	for i, p := range m.Returns {
		params[i] = p.TypeString()
	}
	if len(m.Returns) > 1 {
		return fmt.Sprintf("(%s)", strings.Join(params, ", "))
	}
	return strings.Join(params, ", ")
}

// ReturnArgNameList is the string representation of values being
// returned from the method, ex: 'foo', 's, err'.
func (m MethodData) ReturnArgNameList() string {
	params := make([]string, len(m.Returns))
	for i, p := range m.Returns {
		params[i] = p.Name()
	}
	return strings.Join(params, ", ")
}

// ParamData is the data which represents a parameter to some method of
// an interface.
type ParamData struct {
	Var      *registry.Var
	Variadic bool
}

// Name returns the name of the parameter.
func (p ParamData) Name() string {
	return p.Var.Name
}

// MethodArg is the representation of the parameter in the function
// signature, ex: 'name a.Type'.
func (p ParamData) MethodArg() string {
	if p.Variadic {
		return fmt.Sprintf("%s ...%s", p.Name(), p.TypeString()[2:])
	}
	return fmt.Sprintf("%s %s", p.Name(), p.TypeString())
}

// CallName returns the string representation of the parameter to be
// used for a method call. For a variadic paramter, it will be of the
// format 'foos...'.
func (p ParamData) CallName() string {
	if p.Variadic {
		return p.Name() + "..."
	}
	return p.Name()
}

// TypeString returns the string representation of the type of the
// parameter.
func (p ParamData) TypeString() string {
	return p.Var.TypeString()
}
