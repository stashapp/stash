package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"
)

type schema78Migrator struct {
	migrator
}

func post78(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 78")

	m := schema78Migrator{
		migrator: migrator{
			db: db,
		},
	}

	if err := m.migrateCareerLength(ctx); err != nil {
		return fmt.Errorf("migrating career_length: %w", err)
	}

	if err := m.dropCareerLength(); err != nil {
		return fmt.Errorf("dropping career_length column: %w", err)
	}

	return nil
}

func (m *schema78Migrator) migrateCareerLength(ctx context.Context) error {
	logger.Info("Migrating career_length to career_start/career_end")

	const limit = 1000

	lastID := 0
	parsed := 0
	unparseable := 0

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := `SELECT id, career_length FROM performers
				WHERE career_length IS NOT NULL AND career_length != ''`

			if lastID != 0 {
				query += fmt.Sprintf(" AND id > %d", lastID)
			}

			query += fmt.Sprintf(" ORDER BY id LIMIT %d", limit)

			rows, err := tx.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var (
					id           int
					careerLength string
				)

				if err := rows.Scan(&id, &careerLength); err != nil {
					return err
				}

				lastID = id
				gotSome = true

				start, end, err := utils.ParseYearRangeString(careerLength)
				if err != nil {
					logger.Warnf("Could not parse career_length %q for performer %d: %v — preserving as custom field", careerLength, id, err)

					if err := m.preserveAsCustomField(tx, id, careerLength); err != nil {
						return fmt.Errorf("preserving career_length for performer %d: %w", id, err)
					}
					unparseable++
					continue
				}

				if err := m.updateCareerFields(tx, id, start, end); err != nil {
					return fmt.Errorf("updating career fields for performer %d: %w", id, err)
				}
				parsed++
			}

			return rows.Err()
		}); err != nil {
			return err
		}

		if !gotSome {
			break
		}
	}

	logger.Infof("Career length migration complete: %d parsed, %d unparseable (preserved as custom fields)", parsed, unparseable)
	return nil
}

func (m *schema78Migrator) updateCareerFields(tx *sqlx.Tx, id int, start *int, end *int) error {
	_, err := tx.Exec(
		"UPDATE performers SET career_start = ?, career_end = ? WHERE id = ?",
		start, end, id,
	)
	return err
}

func (m *schema78Migrator) preserveAsCustomField(tx *sqlx.Tx, id int, value string) error {
	// check if a career_length custom field already exists
	var existing sql.NullString
	err := tx.Get(&existing, "SELECT value FROM performer_custom_fields WHERE performer_id = ? AND field = 'career_length'", id)
	if err == nil {
		logger.Debugf("career_length custom field already exists for performer %d, skipping", id)
		return nil
	}

	_, err = tx.Exec(
		"INSERT INTO performer_custom_fields (performer_id, field, value) VALUES (?, 'career_length', ?)",
		id, value,
	)
	return err
}

func (m *schema78Migrator) dropCareerLength() error {
	logger.Info("Dropping career_length column from performers table")
	return m.execAll([]string{
		"ALTER TABLE performers DROP COLUMN career_length",
	})
}

func init() {
	sqlite.RegisterPostMigration(78, post78)
}
