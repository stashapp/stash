package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/robertkrimen/otto"
	"github.com/stashapp/stash/pkg/plugin/common"
)

type jsTaskBuilder struct{}

func (*jsTaskBuilder) build(task pluginTask) Task {
	return &jsPluginTask{
		pluginTask: task,
	}
}

type jsPluginTask struct {
	pluginTask

	started   bool
	waitGroup sync.WaitGroup
}

type responseWriter struct {
	r          strings.Builder
	header     http.Header
	statusCode int
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.r.Write(b)
}

func throw(vm *otto.Otto, str string) {
	value, _ := vm.Call("new Error", nil, str)
	panic(value)
}

func gqlRequestFunc(vm *otto.Otto, gqlHandler http.HandlerFunc) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 0 {
			throw(vm, "missing argument")
		}

		query := call.Argument(0)
		vars := call.Argument(1)
		var variables map[string]interface{}
		if !vars.IsUndefined() {
			exported, _ := vars.Export()
			variables, _ = exported.(map[string]interface{})
		}

		in := struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables,omitempty"`
		}{
			Query:     query.String(),
			Variables: variables,
		}

		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(in)
		if err != nil {
			throw(vm, err.Error())
		}

		r, err := http.NewRequest("POST", "/graphql", &body)
		if err != nil {
			throw(vm, "could not make request")
		}
		r.Header.Set("Content-Type", "application/json")

		w := &responseWriter{
			header: make(http.Header),
		}

		gqlHandler(w, r)

		if w.statusCode != http.StatusOK && w.statusCode != 0 {
			throw(vm, fmt.Sprintf("graphQL query failed: %d - %s. Query: %s. Variables: %v", w.statusCode, w.r.String(), in.Query, in.Variables))
		}

		output := w.r.String()
		// convert to JSON
		var obj map[string]interface{}
		if err = json.Unmarshal([]byte(output), &obj); err != nil {
			throw(vm, fmt.Sprintf("could not unmarshal object %s: %s", output, err.Error()))
		}

		retErr, hasErr := obj["error"]

		if hasErr {
			throw(vm, fmt.Sprintf("graphql error: %v", retErr))
		}

		v, err := vm.ToValue(obj["data"])
		if err != nil {
			throw(vm, fmt.Sprintf("could not create return value: %s", err.Error()))
		}

		return v
	}
}

func (t *jsPluginTask) makeOutput(o otto.Value) {
	t.result = &common.PluginOutput{}

	asObj := o.Object()
	if asObj == nil {
		return
	}

	t.result.Output, _ = asObj.Get("Output")
	err, _ := asObj.Get("Error")
	if !err.IsUndefined() {
		errStr := err.String()
		t.result.Error = &errStr
	}
}

func (t *jsPluginTask) Start() error {
	if t.started {
		return errors.New("task already started")
	}

	t.started = true

	if len(t.plugin.Exec) == 0 {
		return errors.New("no script specified in exec")
	}

	scriptFile := t.plugin.Exec[0]

	vm := otto.New()
	pluginPath := t.plugin.getConfigPath()
	script, err := vm.Compile(filepath.Join(pluginPath, scriptFile), nil)
	if err != nil {
		return err
	}

	input := t.buildPluginInput()

	vm.Set("input", input)
	vm.Set("gql", gqlRequestFunc(vm, t.gqlHandler))
	// TODO - vm.Set("log")
	output, err := vm.Run(script)

	t.makeOutput(output)

	return err
}

func (t *jsPluginTask) Wait() {
}

func (t *jsPluginTask) Stop() error {
	return nil
}
