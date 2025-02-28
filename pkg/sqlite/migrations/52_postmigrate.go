package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema52Migrator struct {
	migrator
}

func post52(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 52")

	m := schema52Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrate(ctx)
}

func (m *schema52Migrator) migrate(ctx context.Context) error {
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		query := "SELECT `folders`.`id`, `folders`.`path`, `parent_folder`.`path` FROM `folders` " +
			"INNER JOIN `folders` AS `parent_folder` ON `parent_folder`.`id` = `folders`.`parent_folder_id`"

		rows, err := tx.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				id               int
				folderPath       string
				parentFolderPath string
			)

			err := rows.Scan(&id, &folderPath, &parentFolderPath)
			if err != nil {
				return err
			}

			// ensure folder path is correct
			if !strings.HasPrefix(folderPath, parentFolderPath) {
				logger.Debugf("folder path %s does not have prefix %s. Correcting...", folderPath, parentFolderPath)

				// get the basename of the zip folder path and append it to the correct path
				folderBasename := filepath.Base(folderPath)
				correctPath := filepath.Join(parentFolderPath, folderBasename)

				logger.Infof("correcting folder path %s to %s", folderPath, correctPath)

				// ensure the correct path is unique
				var v int
				isEmptyErr := tx.Get(&v, "SELECT 1 FROM folders WHERE path = ?", correctPath)
				if isEmptyErr != nil && !errors.Is(isEmptyErr, sql.ErrNoRows) {
					return fmt.Errorf("error checking if correct path %s is unique: %w", correctPath, isEmptyErr)
				}

				if isEmptyErr == nil {
					// correct path is not unique, log and skip
					logger.Warnf("correct path %s already exists, skipping...", correctPath)
					continue
				}

				if _, err := tx.Exec("UPDATE folders SET path = ? WHERE id = ?", correctPath, id); err != nil {
					return fmt.Errorf("error updating folder path %s to %s: %w", folderPath, correctPath, err)
				}
			}
		}

		return rows.Err()
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(52, post52)
}
