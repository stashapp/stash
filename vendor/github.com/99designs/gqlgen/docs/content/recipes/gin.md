---
title: "Using Gin to setup HTTP handlers"
description: Setting up HTTP handlers using Gin, a HTTP web framework written in Go.
linkTitle: Gin
menu: { main: { parent: 'recipes' } }
---

Gin is an excellent alternative for the `net/http` router. From their official [GitHub page](https://github.com/gin-gonic/gin):

> Gin is a web framework written in Go (Golang). It features a martini-like API with much better performance, up to 40 times faster thanks to httprouter. If you need performance and good productivity, you will love Gin.

Here are the steps to setup Gin and gqlgen together:

Install Gin:
```bash
$ go get github.com/gin-gonic/gin
```

In your router file, define the handlers for the GraphQL and Playground endpoints in two different methods and tie them together in the Gin router:

```go
import (
	"github.com/[username]/gqlgen-todos/graph"	// Replace username with your github username
	"github.com/[username]/gqlgen-todos/graph/generated" // Replace username with your github username
	"github.com/gin-gonic/gin"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	// Setting up Gin
	r := gin.Default()
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	r.Run()
}

```

## Accessing gin.Context
At the Resolver level, `gqlgen` gives you access to the `context.Context` object. One way to access the `gin.Context` is to add it to the context and retrieve it again.

First, create a `gin` middleware to add its context to the `context.Context`:
```go
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
```

In the router definition, use the middleware:
```go
r.Use(GinContextToContextMiddleware())
```

Define a function to recover the `gin.Context` from the `context.Context` struct:
```go
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}
```

Lastly, in the Resolver, retrieve the `gin.Context` with the previous defined function:
```go
func (r *resolver) Todo(ctx context.Context) (*Todo, error) {
	gc, err := GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// ...
}
```
