package transport

import (
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

// Options responds to http OPTIONS and HEAD requests
type Options struct {
	// AllowedMethods is a list of allowed HTTP methods.
	AllowedMethods []string
}

var _ graphql.Transport = Options{}

func (o Options) Supports(r *http.Request) bool {
	return r.Method == "HEAD" || r.Method == "OPTIONS"
}

func (o Options) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	switch r.Method {
	case http.MethodOptions:
		w.Header().Set("Allow", o.allowedMethods())
		w.WriteHeader(http.StatusOK)
	case http.MethodHead:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (o Options) allowedMethods() string {
	if len(o.AllowedMethods) == 0 {
		return "OPTIONS, GET, POST"
	}
	return strings.Join(o.AllowedMethods, ", ")
}
