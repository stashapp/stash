package javascript

import (
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/dop251/goja"
	"github.com/stashapp/stash/pkg/logger"
)

type VM struct {
	*goja.Runtime

	Progress      chan float64
	SessionCookie *http.Cookie
	GQLHandler    http.Handler
}

// optionalFieldNameMapper wraps a goja.FieldNameMapper and returns the field name if the wrapped mapper returns an empty string.
type optionalFieldNameMapper struct {
	mapper goja.FieldNameMapper
}

func (tfm optionalFieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string {
	if ret := tfm.mapper.FieldName(t, f); ret != "" {
		return ret
	}

	return f.Name
}

func (tfm optionalFieldNameMapper) MethodName(t reflect.Type, m reflect.Method) string {
	return tfm.mapper.MethodName(t, m)
}

func NewVM() *VM {
	r := goja.New()

	// enable console for backwards compatibility
	c := console{
		Log{
			Logger: logger.Logger,
		},
	}

	// there should not be any reason for this to fail
	_ = c.AddToVM("console", &VM{Runtime: r})

	r.SetFieldNameMapper(optionalFieldNameMapper{goja.TagFieldNameMapper("json", true)})
	return &VM{Runtime: r}
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
