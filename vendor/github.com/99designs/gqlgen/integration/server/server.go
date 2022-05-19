package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/integration"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	cfg := integration.Config{Resolvers: &integration.Resolver{}}
	cfg.Complexity.Query.Complexity = func(childComplexity, value int) int {
		// Allow the integration client to dictate the complexity, to verify this
		// function is executed.
		return value
	}

	srv := handler.NewDefaultServer(integration.NewExecutableSchema(cfg))
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		var ie *integration.CustomError
		if errors.As(e, &ie) {
			return &gqlerror.Error{
				Message: ie.UserMessage,
				Path:    graphql.GetPath(ctx),
			}
		}
		return graphql.DefaultErrorPresenter(ctx, e)
	})
	srv.Use(extension.FixedComplexityLimit(1000))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
