package javascript

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dop251/goja"
)

type VM struct {
	*goja.Runtime

	Progress      chan float64
	SessionCookie *http.Cookie
	GQLHandler    http.Handler
}

func NewVM() *VM {
	return &VM{Runtime: goja.New()}
}

type APIAdder interface {
	AddToVM(globalName string, vm *VM) error
}

type ObjectValueDef struct {
	Name  string
	Value interface{}
}

type setter interface {
	Set(name string, value interface{}) error
}

func Compile(path string) (*goja.Program, error) {
	js, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return goja.Compile(path, string(js), true)
}

func CompileScript(name, script string) (*goja.Program, error) {
	return goja.Compile(name, string(script), true)
}

func SetAll(s setter, defs ...ObjectValueDef) error {
	for _, def := range defs {
		if err := s.Set(def.Name, def.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", def.Name, err)
		}
	}
	return nil
}

func (v *VM) Throw(err error) {
	e, newErr := v.New(v.Get("Error"), v.ToValue(err))
	if newErr != nil {
		panic(newErr)
	}

	panic(e)
}
