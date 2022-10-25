package js

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/robertkrimen/otto"
)

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

func gqlRequestFunc(ctx context.Context, vm *otto.Otto, cookie *http.Cookie, gqlHandler http.Handler) func(call otto.FunctionCall) otto.Value {
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

		r, err := http.NewRequestWithContext(ctx, "POST", "/graphql", &body)
		if err != nil {
			throw(vm, "could not make request")
		}
		r.Header.Set("Content-Type", "application/json")

		if cookie != nil {
			r.AddCookie(cookie)
		}

		w := &responseWriter{
			header: make(http.Header),
		}

		gqlHandler.ServeHTTP(w, r)

		if w.statusCode != http.StatusOK && w.statusCode != 0 {
			throw(vm, fmt.Sprintf("graphQL query failed: %d - %s. Query: %s. Variables: %v", w.statusCode, w.r.String(), in.Query, in.Variables))
		}

		output := w.r.String()
		// convert to JSON
		var obj map[string]interface{}
		if err = json.Unmarshal([]byte(output), &obj); err != nil {
			throw(vm, fmt.Sprintf("could not unmarshal object %s: %s", output, err.Error()))
		}

		retErr, hasErr := obj["errors"]

		if hasErr {
			errOut, _ := json.Marshal(retErr)
			throw(vm, fmt.Sprintf("graphql error: %s", string(errOut)))
		}

		v, err := vm.ToValue(obj["data"])
		if err != nil {
			throw(vm, fmt.Sprintf("could not create return value: %s", err.Error()))
		}

		return v
	}
}

func AddGQLAPI(ctx context.Context, vm *otto.Otto, cookie *http.Cookie, gqlHandler http.Handler) error {
	gql, _ := vm.Object("({})")
	if err := gql.Set("Do", gqlRequestFunc(ctx, vm, cookie, gqlHandler)); err != nil {
		return fmt.Errorf("unable to set GraphQL Do function: %w", err)
	}

	if err := vm.Set("gql", gql); err != nil {
		return fmt.Errorf("unable to set gql: %w", err)
	}

	return nil
}
