package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema86Migrator struct {
	migrator
}

func post86(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 86")

	m := schema86Migrator{
		migrator: migrator{
			db: db,
		},
	}

	if err := m.migrateMeasurements(ctx); err != nil {
		return fmt.Errorf("migrating measurements: %w", err)
	}

	if err := m.dropMeasurements(); err != nil {
		return fmt.Errorf("dropping measurements column: %w", err)
	}

	return nil
}

func (m *schema86Migrator) migrateMeasurements(ctx context.Context) error {
	logger.Info("Migrating measurements to band_size/cup_size/waist_size/hip_size")

	const limit = 1000

	lastID := 0
	parsed := 0
	skipped := 0

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := `SELECT id, measurements FROM performers
				WHERE measurements IS NOT NULL AND measurements != ''`

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
					measurements string
				)

				if err := rows.Scan(&id, &measurements); err != nil {
					return err
				}

				lastID = id
				gotSome = true

				bandSize, cupSize, waistSize, hipSize, err := models.ParseMeasurementsString(measurements)
				if err != nil {
					logger.Warnf("Could not parse measurements %q for performer %d: %v — skipping", measurements, id, err)
					skipped++
					continue
				}

				if err := m.updateMeasurementFields(tx, id, bandSize, cupSize, waistSize, hipSize); err != nil {
					return fmt.Errorf("updating measurement fields for performer %d: %w", id, err)
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

	logger.Infof("Measurements migration complete: %d parsed, %d skipped (unparseable)", parsed, skipped)
	return nil
}

func (m *schema86Migrator) updateMeasurementFields(tx *sqlx.Tx, id int, bandSize *int, cupSize *string, waistSize *int, hipSize *int) error {
	_, err := tx.Exec(
		"UPDATE performers SET band_size = ?, cup_size = ?, waist_size = ?, hip_size = ? WHERE id = ?",
		bandSize, cupSize, waistSize, hipSize, id,
	)
	return err
}

func (m *schema86Migrator) dropMeasurements() error {
	logger.Info("Dropping measurements column from performers table")
	return m.execAll([]string{
		"ALTER TABLE performers DROP COLUMN measurements",
	})
}

func init() {
	sqlite.RegisterPostMigration(86, post86)
}
