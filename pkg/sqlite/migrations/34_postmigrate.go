package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema34Migrator struct {
	migrator
}

func post34(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 34")

	m := schema34Migrator{
		migrator: migrator{
			db: db,
		},
	}

	if err := m.migrateObjects(ctx, "scenes"); err != nil {
		return fmt.Errorf("migrating scenes: %w", err)
	}
	if err := m.migrateObjects(ctx, "images"); err != nil {
		return fmt.Errorf("migrating images: %w", err)
	}
	if err := m.migrateObjects(ctx, "galleries"); err != nil {
		return fmt.Errorf("migrating galleries: %w", err)
	}
	if err := m.migrateObjects(ctx, "files"); err != nil {
		return fmt.Errorf("migrating files: %w", err)
	}
	if err := m.migrateObjects(ctx, "folders"); err != nil {
		return fmt.Errorf("migrating folders: %w", err)
	}

	return nil
}

func (m *schema34Migrator) migrateObjects(ctx context.Context, table string) error {
	logger.Infof("Migrating %s table", table)

	const (
		limit    = 1000
		logEvery = 10000
	)

	lastID := 0
	count := 0

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := fmt.Sprintf("SELECT `id`, `created_at`, `updated_at` FROM `%s` WHERE `created_at` like '%% %%' OR `updated_at` like '%% %%'", table)

			if lastID != 0 {
				query += fmt.Sprintf("AND `id` > %d ", lastID)
			}

			query += fmt.Sprintf("ORDER BY `id` LIMIT %d", limit)

			rows, err := m.db.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var (
					id        int
					createdAt time.Time
					updatedAt time.Time
				)

				err := rows.Scan(&id, &createdAt, &updatedAt)
				if err != nil {
					return err
				}

				lastID = id
				gotSome = true
				count++

				// convert incorrect timestamp string to correct one
				// based on models.SQLTimestamp
				fixedCreated := createdAt.Format(time.RFC3339)
				fixedUpdated := updatedAt.Format(time.RFC3339)

				updateSQL := fmt.Sprintf("UPDATE `%s` SET `created_at` = ?, `updated_at` = ? WHERE `id` = ?", table)

				_, err = m.db.Exec(updateSQL, fixedCreated, fixedUpdated, id)
				if err != nil {
					return err
				}
			}

			return rows.Err()
		}); err != nil {
			return err
		}

		if !gotSome {
			break
		}

		if count%logEvery == 0 {
			logger.Infof("Migrated %d rows", count, table)
		}
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(34, post34)
}
