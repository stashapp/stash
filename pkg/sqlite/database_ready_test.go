package sqlite

import (
	"context"
	"errors"
	"testing"
)

func TestBeginUninitialisedDatabase(t *testing.T) {
	ctx := context.Background()

	db := NewDatabase()
	if _, err := db.Begin(ctx, false); !errors.Is(err, ErrDatabaseNotInitialized) {
		t.Errorf("Begin() error = %v, want %v", err, ErrDatabaseNotInitialized)
	}
	if _, err := db.WithDatabase(ctx); !errors.Is(err, ErrDatabaseNotInitialized) {
		t.Errorf("WithDatabase() error = %v, want %v", err, ErrDatabaseNotInitialized)
	}

	// simulate Open having found an old schema
	db.schemaVersion = appSchemaVersion - 1
	var migrationErr *MigrationNeededError
	if _, err := db.Begin(ctx, true); !errors.As(err, &migrationErr) {
		t.Errorf("Begin() error = %v, want MigrationNeededError", err)
	}
}
