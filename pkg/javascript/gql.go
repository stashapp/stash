package javascript

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dop251/goja"
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

type GQL struct {
	Context    context.Context
	Cookie     *http.Cookie
	GQLHandler http.Handler
}

func (g *GQL) gqlRequestFunc(vm *VM) func(query string, variables map[string]interface{}) (goja.Value, error) {
	return func(query string, variables map[string]interface{}) (goja.Value, error) {
		in := struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables,omitempty"`
		}{
			Query:     query,
			Variables: variables,
		}

		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(in)
		if err != nil {
			return nil, err
		}

		r, err := http.NewRequestWithContext(g.Context, "POST", "/graphql", &body)
		if err != nil {
			return nil, fmt.Errorf("could not make request")
		}
		r.Header.Set("Content-Type", "application/json")

		if g.Cookie != nil {
			r.AddCookie(g.Cookie)
		}

		w := &responseWriter{
			header: make(http.Header),
		}

		g.GQLHandler.ServeHTTP(w, r)

		if w.statusCode != http.StatusOK && w.statusCode != 0 {
			vm.Throw(fmt.Errorf("graphQL query failed: %d - %s. Query: %s. Variables: %v", w.statusCode, w.r.String(), in.Query, in.Variables))
		}

		output := w.r.String()
		// convert to JSON
		var obj map[string]interface{}
		if err = json.Unmarshal([]byte(output), &obj); err != nil {
			vm.Throw(fmt.Errorf("could not unmarshal object %s: %s", output, err.Error()))
		}

		retErr, hasErr := obj["errors"]

		if hasErr {
			errOut, _ := json.Marshal(retErr)
			vm.Throw(fmt.Errorf("graphql error: %s", string(errOut)))
		}

		v := vm.ToValue(obj["data"])

		return v, nil
	}
}

func (g *GQL) AddToVM(globalName string, vm *VM) error {
	gql := vm.NewObject()

	if err := gql.Set("Do", g.gqlRequestFunc(vm)); err != nil {
		return fmt.Errorf("unable to set GraphQL Do function: %w", err)
	}

	if err := vm.Set(globalName, gql); err != nil {
		return fmt.Errorf("unable to set gql: %w", err)
	}

	return nil
}
