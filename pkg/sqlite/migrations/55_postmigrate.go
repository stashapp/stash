package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema55Migrator struct {
	migrator
}

func post55(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 55")

	m := schema55Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrate(ctx)
}

func (m *schema55Migrator) migrate(ctx context.Context) error {
	// the last_played_at column was storing in a different format than the rest of the timestamps
	// convert the play history date to the correct format
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		query := "SELECT DISTINCT `scene_id`, `view_date` FROM `scenes_view_dates`"

		rows, err := tx.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				id       int
				viewDate sqlite.Timestamp
			)

			err := rows.Scan(&id, &viewDate)
			if err != nil {
				return err
			}

			utcTimestamp := sqlite.UTCTimestamp{
				Timestamp: viewDate,
			}

			// convert the timestamp to the correct format
			if _, err := tx.Exec("UPDATE scenes_view_dates SET view_date = ? WHERE view_date = ?", utcTimestamp, viewDate.Timestamp); err != nil {
				return fmt.Errorf("error correcting view date %s to %s: %w", viewDate.Timestamp, viewDate, err)
			}
		}

		return rows.Err()
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(55, post55)
}
