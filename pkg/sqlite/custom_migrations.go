package sqlite

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type customMigrationFunc func(ctx context.Context, db *sqlx.DB) error

func RegisterCustomMigration(schemaVersion uint, fn customMigrationFunc) {
	v := customMigrations[schemaVersion]
	v = append(v, fn)
	customMigrations[schemaVersion] = v
}

var customMigrations = make(map[uint][]customMigrationFunc)
