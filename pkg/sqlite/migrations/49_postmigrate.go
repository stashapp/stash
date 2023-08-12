package migrations

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema49Migrator struct {
	migrator
}

func post49(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 49")

	m := schema49Migrator{
		migrator: migrator{
			db: db,
		},
	}

	if err := m.createVideoFiltersTable(ctx, db); err != nil {
		return err
	}

	return nil
}

func (m *schema49Migrator) createVideoFiltersTable(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS scene_filters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			scene_id INTEGER NOT NULL UNIQUE,
			brightness INTEGER NOT NULL,
			contrast INTEGER NOT NULL,
			gamma INTEGER NOT NULL,
			saturate INTEGER NOT NULL,
			hue_rotate INTEGER NOT NULL,
			warmth INTEGER NOT NULL,
			red INTEGER NOT NULL,
			green INTEGER NOT NULL,
			blue INTEGER NOT NULL,
			blur INTEGER NOT NULL,
			rotate REAL NOT NULL,
			scale INTEGER NOT NULL,
			aspect_ratio INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			FOREIGN KEY (scene_id) REFERENCES scenes(id) ON DELETE CASCADE
		)
	`)

	// Create an index on the scene_id column
	_, indexErr := db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS index_scene_filters_on_scene_id
		ON scene_filters(scene_id)
	`)
	if indexErr != nil {
		return indexErr
	}

	return err
}

func init() {
	sqlite.RegisterPostMigration(49, post49)
}
