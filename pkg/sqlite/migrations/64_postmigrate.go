package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

// this is a copy of the 55 post migration
// some non-UTC dates were missed, so we need to correct them

type schema64Migrator struct {
	migrator
}

func post64(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 64")

	m := schema64Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrate(ctx)
}

func (m *schema64Migrator) migrate(ctx context.Context) error {
	// the last_played_at column was storing in a different format than the rest of the timestamps
	// convert the play history date to the correct format
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		query := "SELECT DISTINCT `scene_id`, `view_date` FROM `scenes_view_dates`"

		rows, err := m.db.Query(query)
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

			// skip if already in the correct format
			if viewDate.Timestamp.Location() == time.UTC {
				logger.Debugf("view date %s is already in the correct format", viewDate.Timestamp)
				continue
			}

			utcTimestamp := sqlite.UTCTimestamp{
				Timestamp: viewDate,
			}

			// convert the timestamp to the correct format
			logger.Debugf("correcting view date %q to UTC date %q for scene %d", viewDate.Timestamp, viewDate.Timestamp.UTC(), id)
			r, err := m.db.Exec("UPDATE scenes_view_dates SET view_date = ? WHERE scene_id = ? AND (view_date = ? OR view_date = ?)", utcTimestamp, id, viewDate.Timestamp, viewDate)
			if err != nil {
				return fmt.Errorf("error correcting view date %s to %s: %w", viewDate.Timestamp, viewDate, err)
			}

			rowsAffected, err := r.RowsAffected()
			if err != nil {
				return err
			}

			if rowsAffected == 0 {
				return fmt.Errorf("no rows affected when updating view date %s to %s for scene %d", viewDate.Timestamp, viewDate.Timestamp.UTC(), id)
			}
		}

		return rows.Err()
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(64, post64)
}
