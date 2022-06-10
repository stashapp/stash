package sqlite

import "github.com/jmoiron/sqlx"

type customMigrationFunc func(db *sqlx.DB) error

func RegisterCustomMigration(schemaVersion uint, fn customMigrationFunc) {
	v := customMigrations[schemaVersion]
	v = append(v, fn)
	customMigrations[schemaVersion] = v
}

var customMigrations = make(map[uint][]customMigrationFunc)
