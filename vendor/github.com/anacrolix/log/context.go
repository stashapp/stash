package log

import (
	"context"
)

var loggerContextKey interface{} = (*Logger)(nil)

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func ContextLogger(ctx context.Context) Logger {
	return ctx.Value(loggerContextKey).(Logger)
}
