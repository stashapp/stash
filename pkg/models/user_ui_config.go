package models

import "context"

// UserUIConfigReaderWriter reads and writes the UI configuration of the current user.
type UserUIConfigReaderWriter interface {
	// Get returns the UI configuration of the current user.
	// Returns nil if the user has not set a UI configuration, or if there is no current user.
	Get(ctx context.Context) (map[string]interface{}, error)
	// Set replaces the UI configuration of the current user.
	// Returns an error if there is no current user.
	Set(ctx context.Context, config map[string]interface{}) error
}
