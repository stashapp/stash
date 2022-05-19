package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	. "github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"

	"github.com/99designs/gqlgen/graphql"
)

type Tracer struct {
	DisableColor bool
	au           Aurora
	out          io.Writer
}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
} = &Tracer{}

func (a Tracer) ExtensionName() string {
	return "ApolloTracing"
}

func (a *Tracer) Validate(schema graphql.ExecutableSchema) error {
	isTTY := isatty.IsTerminal(os.Stdout.Fd())

	a.au = NewAurora(!a.DisableColor && isTTY)
	a.out = colorable.NewColorableStdout()

	return nil
}

func stringify(value interface{}) string {
	valueJson, err := json.MarshalIndent(value, "  ", "  ")
	if err == nil {
		return string(valueJson)
	}

	return fmt.Sprint(value)
}

func (a Tracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	rctx := graphql.GetOperationContext(ctx)

	fmt.Fprintln(a.out, "GraphQL Request {")
	for _, line := range strings.Split(rctx.RawQuery, "\n") {
		fmt.Fprintln(a.out, " ", Cyan(line))
	}
	for name, value := range rctx.Variables {
		fmt.Fprintf(a.out, "  var %s = %s\n", name, Yellow(stringify(value)))
	}
	resp := next(ctx)

	fmt.Fprintln(a.out, "  resp:", Green(stringify(resp)))
	for _, err := range resp.Errors {
		fmt.Fprintln(a.out, "  error:", Bold(err.Path.String()+":"), Red(err.Message))
	}
	fmt.Fprintln(a.out, "}")
	fmt.Fprintln(a.out)
	return resp
}
