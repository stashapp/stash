package migrations

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema72Migrator struct {
	migrator
}

func post72(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 65")

	m := schema72Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrate()
}

func (m *schema72Migrator) migrate() error {
	// TODO: write to "type" field in custom fields
	return nil
}

func init() {
	sqlite.RegisterPostMigration(72, post72)
}
