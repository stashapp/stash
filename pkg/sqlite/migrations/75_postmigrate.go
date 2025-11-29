package migrations

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema73Migrator struct {
	migrator
}

func post73(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 73")

	m := schema73Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrate(ctx)
}

func (m *schema73Migrator) migrate(ctx context.Context) error {
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		query := "SELECT `performer_id`, `value` FROM `performer_custom_fields`"

		rows, err := tx.Queryx(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				performer_id int
				value        interface{}
			)

			err := rows.Scan(&performer_id, &value)
			if err != nil {
				return err
			}

			gotype := reflect.TypeOf(value).String()
			logger.Debugf("setting type for %v to %v for %v", value, gotype, performer_id)
			r, err := tx.Exec("UPDATE performer_custom_fields SET type = ? WHERE performer_id = ?", gotype, performer_id)
			if err != nil {
				return fmt.Errorf("error setting type for %v to %v for %v", value, gotype, performer_id)
			}

			rowsAffected, err := r.RowsAffected()
			if err != nil {
				return err
			}

			if rowsAffected == 0 {
				return fmt.Errorf("no rows affected when updating type %v to %v for %v", value, gotype, performer_id)
			}
		}

		return rows.Err()
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(73, post73)
}
