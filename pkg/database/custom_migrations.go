package database

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
)

func runCustomMigrations() error {
	if err := createImagesChecksumIndex(); err != nil {
		return err
	}

	return nil
}

func createImagesChecksumIndex() error {
	return WithTxn(func(tx *sqlx.Tx) error {
		row := tx.QueryRow("SELECT 1 AS found FROM sqlite_master WHERE type = 'index' AND name = 'images_checksum_unique'")
		err := row.Err()
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == nil {
			var found bool
			row.Scan(&found)
			if found {
				return nil
			}
		}

		_, err = tx.Exec("CREATE UNIQUE INDEX images_checksum_unique ON images (checksum)")
		if err == nil {
			_, err = tx.Exec("DROP INDEX IF EXISTS index_images_checksum")
			if err != nil {
				logger.Errorf("Failed to remove surrogate images.checksum index: %s", err)
			}
			logger.Info("Created unique constraint on images table")
			return nil
		}

		_, err = tx.Exec("CREATE INDEX IF NOT EXISTS index_images_checksum ON images (checksum)")
		if err != nil {
			logger.Errorf("Unable to create index on images.checksum: %s", err)
		}

		var result []struct {
			Checksum string `db:"checksum"`
		}

		err = tx.Select(&result, "SELECT checksum FROM images GROUP BY checksum HAVING COUNT(1) > 1")
		if err != nil && err != sql.ErrNoRows {
			logger.Errorf("Unable to determine non-unique image checksums: %s", err)
			return nil
		}

		checksums := make([]string, len(result))
		for i, res := range result {
			checksums[i] = res.Checksum
		}

		logger.Warnf("The following duplicate image checksums have been found. Please remove the duplicates and restart. %s", strings.Join(checksums, ", "))

		return nil
	})
}
