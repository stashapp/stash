httplog
=======

Small but powerful structured logging package for HTTP request logging in Go.

## Example

(see [_example/](./_example/main.go))

```go
package main

import (
  "net/http"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  "github.com/go-chi/httplog"
)

func main() {
  // Logger
  logger := httplog.NewLogger("httplog-example", httplog.Options{
    JSON: true,
  })

  // Service
  r := chi.NewRouter()
  r.Use(httplog.RequestLogger(logger))
  r.Use(middleware.Heartbeat("/ping"))

  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello world"))
  })

  r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
    panic("oh no")
  })

  r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
    oplog := httplog.LogEntry(r.Context())
    w.Header().Add("Content-Type", "text/plain")
    oplog.Info().Msg("info here")
    w.Write([]byte("info here"))
  })

  r.Get("/warn", func(w http.ResponseWriter, r *http.Request) {
    oplog := httplog.LogEntry(r.Context())
    oplog.Warn().Msg("warn here")
    w.WriteHeader(400)
    w.Write([]byte("warn here"))
  })

  r.Get("/err", func(w http.ResponseWriter, r *http.Request) {
    oplog := httplog.LogEntry(r.Context())
    oplog.Error().Msg("err here")
    w.WriteHeader(500)
    w.Write([]byte("err here"))
  })

  http.ListenAndServe(":5555", r)
}

```

## License

MIT
