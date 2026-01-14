package migrations

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

func post84(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 76")

	m := schema84Migrator{
		migrator: migrator{
			db: db,
		},
		folderCache: make(map[string]folderInfo),
	}

	if err := m.migrateFolders(ctx); err != nil {
		return fmt.Errorf("migrating folders: %w", err)
	}

	return nil
}

type schema84Migrator struct {
	migrator
	folderCache map[string]folderInfo
}

func (m *schema84Migrator) migrateFolders(ctx context.Context) error {
	logger.Infof("Migrating folders")

	const (
		limit    = 1000
		logEvery = 10000
	)

	lastID := 0
	count := 0

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := "SELECT `folders`.`id`, `folders`.`path` FROM `folders` "

			if lastID != 0 {
				query += fmt.Sprintf("WHERE `folders`.`id` > %d ", lastID)
			}

			query += fmt.Sprintf("ORDER BY `folders`.`id` LIMIT %d", limit)

			rows, err := tx.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var id int
				var p string

				err := rows.Scan(&id, &p)
				if err != nil {
					return err
				}

				lastID = id
				gotSome = true
				count++

				basename := filepath.Base(p)
				logger.Debugf("Migrating folder %d %q: setting basename to %q", id, p, basename)
				_, err = tx.Exec("UPDATE `folders` SET `basename` = ? WHERE `id` = ?", basename, id)
				if err != nil {
					return fmt.Errorf("error migrating folder %d %q: %w", id, p, err)
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
			logger.Infof("Migrated %d folders", count)
		}
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(84, post84)
}
