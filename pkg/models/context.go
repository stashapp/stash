package models

import "context"

type contextKey int

const (
	setupContextKey contextKey = iota
	userIDContextKey
)

// WithUserID returns a context with the provided user ID set.
// The user ID is used to read and write user-specific data such as ratings and play history.
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

// UserIDFromContext returns the user ID set in the context.
// Returns false if no user ID is set, or if the user is not a stored user (such as the guest user).
func UserIDFromContext(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(userIDContextKey).(int)
	if !ok || v == 0 {
		return 0, false
	}

	return v, true
}

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
