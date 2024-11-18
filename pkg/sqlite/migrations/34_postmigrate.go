package migrations

import (
	"context"
	"fmt"
	"strings"
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

	objectCols := []string{
		"created_at",
		"updated_at",
	}

	filesystemCols := objectCols
	filesystemCols = append(filesystemCols, "mod_time")

	if err := m.migrateObjects(ctx, "scenes", objectCols); err != nil {
		return fmt.Errorf("migrating scenes: %w", err)
	}
	if err := m.migrateObjects(ctx, "images", objectCols); err != nil {
		return fmt.Errorf("migrating images: %w", err)
	}
	if err := m.migrateObjects(ctx, "galleries", objectCols); err != nil {
		return fmt.Errorf("migrating galleries: %w", err)
	}
	if err := m.migrateObjects(ctx, "files", filesystemCols); err != nil {
		return fmt.Errorf("migrating files: %w", err)
	}
	if err := m.migrateObjects(ctx, "folders", filesystemCols); err != nil {
		return fmt.Errorf("migrating folders: %w", err)
	}

	return nil
}

func (m *schema34Migrator) migrateObjects(ctx context.Context, table string, cols []string) error {
	logger.Infof("Migrating %s table", table)

	quotedCols := make([]string, len(cols)+1)
	quotedCols[0] = "`id`"
	whereClauses := make([]string, len(cols))
	updateClauses := make([]string, len(cols))
	for i, v := range cols {
		quotedCols[i+1] = "`" + v + "`"
		whereClauses[i] = "`" + v + "` like '% %'"
		updateClauses[i] = "`" + v + "` = ?"
	}

	colList := strings.Join(quotedCols, ", ")
	clauseList := strings.Join(whereClauses, " OR ")
	updateList := strings.Join(updateClauses, ", ")

	const (
		limit    = 1000
		logEvery = 10000
	)

	lastID := 0
	count := 0

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := fmt.Sprintf("SELECT %s FROM `%s` WHERE (%s)", colList, table, clauseList)

			if lastID != 0 {
				query += fmt.Sprintf(" AND `id` > %d ", lastID)
			}

			query += fmt.Sprintf(" ORDER BY `id` LIMIT %d", limit)

			rows, err := tx.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var (
					id int
				)

				timeValues := make([]interface{}, len(cols)+1)
				timeValues[0] = &id
				for i := range cols {
					v := time.Time{}
					timeValues[i+1] = &v
				}

				err := rows.Scan(timeValues...)
				if err != nil {
					return err
				}

				lastID = id
				gotSome = true
				count++

				// convert incorrect timestamp string to correct one
				// based on models.SQLTimestamp
				args := make([]interface{}, len(cols)+1)
				for i := range cols {
					tv := timeValues[i+1].(*time.Time)
					args[i] = tv.Format(time.RFC3339)
				}
				args[len(cols)] = id

				updateSQL := fmt.Sprintf("UPDATE `%s` SET %s WHERE `id` = ?", table, updateList)

				_, err = tx.Exec(updateSQL, args...)
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
			logger.Infof("Migrated %d rows", count)
		}
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(34, post34)
}
