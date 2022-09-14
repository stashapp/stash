package sqlite

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type customMigrationFunc func(ctx context.Context, db *sqlx.DB) error

func RegisterPostMigration(schemaVersion uint, fn customMigrationFunc) {
	v := postMigrations[schemaVersion]
	v = append(v, fn)
	postMigrations[schemaVersion] = v
}

func RegisterPreMigration(schemaVersion uint, fn customMigrationFunc) {
	v := preMigrations[schemaVersion]
	v = append(v, fn)
	preMigrations[schemaVersion] = v
}

var postMigrations = make(map[uint][]customMigrationFunc)
var preMigrations = make(map[uint][]customMigrationFunc)
