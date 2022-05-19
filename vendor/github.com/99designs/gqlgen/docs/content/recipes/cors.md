---
title: "Setting CORS headers using rs/cors for gqlgen"
description: Use the best of breed rs/cors library to set CORS headers when working with gqlgen
linkTitle: CORS
menu: { main: { parent: "recipes" } }
---

Cross-Origin Resource Sharing (CORS) headers are required when your graphql server lives on a different domain to the one your client code is served. You can read more about CORS in the [MDN docs](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS).

## rs/cors

gqlgen doesn't include a CORS implementation, but it is built to work with all standard http middleware. Here we are going to use the fantastic `chi` and `rs/cors` to build our server.

```go
package main

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/_examples/starwars"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/gorilla/websocket"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)


    srv := handler.NewDefaultServer(starwars.NewExecutableSchema(starwars.NewResolver()))
    srv.AddTransport(&transport.Websocket{
        Upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                // Check against your desired domains here
                 return r.Host == "example.org"
            },
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        },
    })

	router.Handle("/", playground.Handler("Starwars", "/query"))
	router.Handle("/query", srv)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

```
