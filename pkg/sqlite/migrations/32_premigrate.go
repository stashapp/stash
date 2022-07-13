package migrations

import (
	"context"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

func pre32(ctx context.Context, db *sqlx.DB) error {
	// verify that folder-based galleries (those with zip = 0 and path is not null) are
	// not zip-based. If they are zip based then set zip to 1
	// we could still miss some if the path does not exist, but this is the best we can do

	logger.Info("Running pre-migration for schema version 32")

	mm := schema32PreMigrator{
		migrator: migrator{
			db: db,
		},
	}

	return mm.migrate(ctx)
}

type schema32PreMigrator struct {
	migrator
}

func (m *schema32PreMigrator) migrate(ctx context.Context) error {
	// query for galleries with zip = 0 and path not null
	result := struct {
		Count int `db:"count"`
	}{0}

	if err := m.db.Get(&result, "SELECT COUNT(*) AS count FROM `galleries` WHERE `zip` = '0' AND `path` IS NOT NULL"); err != nil {
		return err
	}

	if result.Count == 0 {
		return nil
	}

	logger.Infof("Checking %d galleries for incorrect zip value...", result.Count)

	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		const query = "SELECT `id`, `path` FROM `galleries` WHERE `zip` = '0' AND `path` IS NOT NULL ORDER BY `id`"
		rows, err := m.db.Query(query)
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

			// if path does not exist, assume that it is a file and not a folder
			// if it does exist and is a folder, then we ignore it
			// otherwise set zip to 1
			info, err := os.Stat(p)
			if err != nil {
				logger.Warnf("unable to verify if %q is a folder due to error %v. Not migrating.", p, err)
				continue
			}

			if info.IsDir() {
				// ignore it
				continue
			}

			logger.Infof("Correcting %q gallery to be zip-based.", p)

			_, err = m.db.Exec("UPDATE `galleries` SET `zip` = '1' WHERE `id` = ?", id)
			if err != nil {
				return err
			}
		}

		return rows.Err()
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	sqlite.RegisterPreMigration(32, pre32)
}
