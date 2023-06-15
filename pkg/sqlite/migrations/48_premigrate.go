package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

func pre48(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running pre-migration for schema version 48")

	m := schema48PreMigrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.validateScrapedItems(ctx)
}

type schema48PreMigrator struct {
	migrator
}

func (m *schema48PreMigrator) validateScrapedItems(ctx context.Context) error {
	var count int

	row := m.db.QueryRowx("SELECT COUNT(*) FROM scraped_items")
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	return fmt.Errorf("found %d row(s) in scraped_items table, cannot migrate", count)
}

func init() {
	sqlite.RegisterPreMigration(48, pre48)
}
