package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type migrator struct {
	db *sqlx.DB
}

func (m *migrator) withTxn(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		}

		if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

func (m *migrator) execAll(stmts []string) error {
	for _, stmt := range stmts {
		if _, err := m.db.Exec(stmt); err != nil {
			return fmt.Errorf("executing statement %s: %w", stmt, err)
		}
	}

	return nil
}
