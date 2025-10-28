package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

func pre48(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running pre-migration for schema version 48")

	m := schema48PreMigrator{
		migrator: migrator{
			db: db,
		},
	}

	if err := m.validateScrapedItems(ctx); err != nil {
		return err
	}

	if err := m.fixStudioNames(ctx); err != nil {
		return err
	}

	return nil
}

type schema48PreMigrator struct {
	migrator
}

func (m *schema48PreMigrator) validateScrapedItems(ctx context.Context) error {
	var count int

	row := m.db.QueryRowx("SELECT COUNT(*) FROM scraped_items")
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	return fmt.Errorf("found %d row(s) in scraped_items table, cannot migrate", count)
}

func (m *schema48PreMigrator) fixStudioNames(ctx context.Context) error {
	// First remove NULL names
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.Exec("UPDATE studios SET name = 'NULL' WHERE name IS NULL")
		return err
	}); err != nil {
		return err
	}

	// Then remove duplicate names

	dupes := make(map[string][]int)

	// collect names
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		rows, err := tx.Query("SELECT id, name FROM studios ORDER BY name, id")
		if err != nil {
			return err
		}
		defer rows.Close()

		first := true
		var lastName string

		for rows.Next() {
			var (
				id   int
				name string
			)

			err := rows.Scan(&id, &name)
			if err != nil {
				return err
			}

			if first {
				first = false
				lastName = name
				continue
			}

			if lastName == name {
				dupes[name] = append(dupes[name], id)
			} else {
				lastName = name
			}
		}

		return rows.Err()
	}); err != nil {
		return err
	}

	// rename them
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		for name, ids := range dupes {
			i := 0
			for _, id := range ids {
				var newName string
				for j := 0; ; j++ {
					i++
					newName = fmt.Sprintf("%s (%d)", name, i)

					var count int

					row := tx.QueryRowx("SELECT COUNT(*) FROM studios WHERE name = ?", newName)
					err := row.Scan(&count)
					if err != nil {
						return err
					}

					if count == 0 {
						break
					}

					// try up to 100 times to find a unique name
					if j == 100 {
						return fmt.Errorf("cannot make unique studio name for %s", name)
					}
				}

				logger.Infof("Renaming duplicate studio id %d to %s", id, newName)
				_, err := tx.Exec("UPDATE studios SET name = ? WHERE id = ?", newName, id)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	sqlite.RegisterPreMigration(48, pre48)
}
