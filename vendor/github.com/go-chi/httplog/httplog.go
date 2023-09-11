package httplog

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger(serviceName string, opts ...Options) zerolog.Logger {
	if len(opts) > 0 {
		Configure(opts[0])
	} else {
		Configure(DefaultOptions)
	}
	logger := log.With().Str("service", strings.ToLower(serviceName))
	if !DefaultOptions.Concise && len(DefaultOptions.Tags) > 0 {
		logger = logger.Fields(map[string]interface{}{
			"tags": DefaultOptions.Tags,
		})
	}
	return logger.Logger()
}

// RequestLogger is an http middleware to log http requests and responses.
//
// NOTE: for simplicity, RequestLogger automatically makes use of the chi RequestID and
// Recoverer middleware.
func RequestLogger(logger zerolog.Logger, skipPaths ...[]string) func(next http.Handler) http.Handler {
	return chi.Chain(
		middleware.RequestID,
		Handler(logger, skipPaths...),
		middleware.Recoverer,
	).Handler
}

func Handler(logger zerolog.Logger, optSkipPaths ...[]string) func(next http.Handler) http.Handler {
	var f middleware.LogFormatter = &requestLogger{logger}

	skipPaths := map[string]struct{}{}
	if len(optSkipPaths) > 0 {
		for _, path := range optSkipPaths[0] {
			skipPaths[path] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Skip the logger if the path is in the skip list
			if len(skipPaths) > 0 {
				_, skip := skipPaths[r.URL.Path]
				if skip {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Log the request
			entry := f.NewLogEntry(r)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			buf := newLimitBuffer(512)
			ww.Tee(buf)

			t1 := time.Now()
			defer func() {
				var respBody []byte
				if ww.Status() >= 400 {
					respBody, _ = ioutil.ReadAll(buf)
				}
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), respBody)
			}()

			next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))
		}
		return http.HandlerFunc(fn)
	}
}

type requestLogger struct {
	Logger zerolog.Logger
}

func (l *requestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &RequestLoggerEntry{}
	msg := fmt.Sprintf("Request: %s %s", r.Method, r.URL.Path)
	entry.Logger = l.Logger.With().Fields(requestLogFields(r, true)).Logger()
	if !DefaultOptions.Concise {
		entry.Logger.Info().Fields(requestLogFields(r, DefaultOptions.Concise)).Msgf(msg)
	}
	return entry
}

type RequestLoggerEntry struct {
	Logger zerolog.Logger
	msg    string
}

func (l *RequestLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	msg := fmt.Sprintf("Response: %d %s", status, statusLabel(status))
	if l.msg != "" {
		msg = fmt.Sprintf("%s - %s", msg, l.msg)
	}

	responseLog := map[string]interface{}{
		"status":  status,
		"bytes":   bytes,
		"elapsed": float64(elapsed.Nanoseconds()) / 1000000.0, // in milliseconds
	}

	if !DefaultOptions.Concise {
		// Include response header, as well for error status codes (>400) we include
		// the response body so we may inspect the log message sent back to the client.
		if status >= 400 {
			body, _ := extra.([]byte)
			responseLog["body"] = string(body)
		}
		if len(header) > 0 {
			responseLog["header"] = headerLogField(header)
		}
	}

	l.Logger.WithLevel(statusLevel(status)).Fields(map[string]interface{}{
		"httpResponse": responseLog,
	}).Msgf(msg)
}

func (l *RequestLoggerEntry) Panic(v interface{}, stack []byte) {
	stacktrace := "#"
	if DefaultOptions.JSON {
		stacktrace = string(stack)
	}

	l.Logger = l.Logger.With().
		Str("stacktrace", stacktrace).
		Str("panic", fmt.Sprintf("%+v", v)).
		Logger()

	l.msg = fmt.Sprintf("%+v", v)

	if !DefaultOptions.JSON {
		middleware.PrintPrettyStack(v)
	}
}

func requestLogFields(r *http.Request, concise bool) map[string]interface{} {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	requestURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	requestFields := map[string]interface{}{
		"requestURL":    requestURL,
		"requestMethod": r.Method,
		"requestPath":   r.URL.Path,
		"remoteIP":      r.RemoteAddr,
		"proto":         r.Proto,
	}
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		requestFields["requestID"] = reqID
	}

	if concise {
		return map[string]interface{}{
			"httpRequest": requestFields,
		}
	}

	requestFields["scheme"] = scheme

	if len(r.Header) > 0 {
		requestFields["header"] = headerLogField(r.Header)
	}

	return map[string]interface{}{
		"httpRequest": requestFields,
	}
}

func headerLogField(header http.Header) map[string]string {
	headerField := map[string]string{}
	for k, v := range header {
		k = strings.ToLower(k)
		switch {
		case len(v) == 0:
			continue
		case len(v) == 1:
			headerField[k] = v[0]
		default:
			headerField[k] = fmt.Sprintf("[%s]", strings.Join(v, "], ["))
		}
		if k == "authorization" || k == "cookie" || k == "set-cookie" {
			headerField[k] = "***"
		}

		for _, skip := range DefaultOptions.SkipHeaders {
			if k == skip {
				headerField[k] = "***"
				break
			}
		}
	}
	return headerField
}

func statusLevel(status int) zerolog.Level {
	switch {
	case status <= 0:
		return zerolog.WarnLevel
	case status < 400: // for codes in 100s, 200s, 300s
		return zerolog.InfoLevel
	case status >= 400 && status < 500:
		return zerolog.WarnLevel
	case status >= 500:
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func statusLabel(status int) string {
	switch {
	case status >= 100 && status < 300:
		return "OK"
	case status >= 300 && status < 400:
		return "Redirect"
	case status >= 400 && status < 500:
		return "Client Error"
	case status >= 500:
		return "Server Error"
	default:
		return "Unknown"
	}
}

// Helper methods used by the application to get the request-scoped
// logger entry and set additional fields between handlers.
//
// This is a useful pattern to use to set state on the entry as it
// passes through the handler chain, which at any point can be logged
// with a call to .Print(), .Info(), etc.

func LogEntry(ctx context.Context) zerolog.Logger {
	entry, ok := ctx.Value(middleware.LogEntryCtxKey).(*RequestLoggerEntry)
	if !ok || entry == nil {
		return zerolog.Nop()
	} else {
		return entry.Logger
	}
}

func LogEntrySetField(ctx context.Context, key, value string) {
	if entry, ok := ctx.Value(middleware.LogEntryCtxKey).(*RequestLoggerEntry); ok {
		entry.Logger = entry.Logger.With().Str(key, value).Logger()
	}
}

func LogEntrySetFields(ctx context.Context, fields map[string]interface{}) {
	if entry, ok := ctx.Value(middleware.LogEntryCtxKey).(*RequestLoggerEntry); ok {
		entry.Logger = entry.Logger.With().Fields(fields).Logger()
	}
}
