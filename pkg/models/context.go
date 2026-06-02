package models

import "context"

type contextKey int

const (
	setupContextKey contextKey = iota
)

func WithSetupContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, setupContextKey, true)
}

func IsSetupContext(ctx context.Context) bool {
	v := ctx.Value(setupContextKey)
	if v == nil {
		return false
	}

	b, ok := v.(bool)
	return ok && b
}
