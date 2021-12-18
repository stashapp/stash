# [Chi](https://github.com/go-chi/chi) logrus-logger middleware

logrus-logger is a request logging middleware for Chi using [Logrus](https://github.com/sirupsen/logrus) logging library

[![Documentation](https://godoc.org/github.com/chi-middleware/logrus-logger?status.svg)](https://pkg.go.dev/github.com/chi-middleware/logrus-logger)
[![codecov](https://codecov.io/gh/chi-middleware/logrus-logger/branch/master/graph/badge.svg)](https://codecov.io/gh/chi-middleware/logrus-logger)
[![Go Report Card](https://goreportcard.com/badge/github.com/chi-middleware/logrus-logger)](https://goreportcard.com/report/github.com/chi-middleware/logrus-logger)
[![Build Status](https://cloud.drone.io/api/badges/chi-middleware/logrus-logger/status.svg?ref=refs/heads/master)](https://cloud.drone.io/chi-middleware/logrus-logger)

## Usage

Import using:

```go
import "github.com/chi-middleware/logrus-logger"
```

Use middleware:

```go
    log := logrus.New()

    r := chi.NewRouter()
    r.Use(logger.Logger("router", log))
```
