package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
)

// WithTxn executes the provided function within a transaction. It rolls back
// the transaction if the function returns an error, otherwise the transaction
// is committed.
func WithTxn(fn func(tx *sqlx.Tx) error) error {
	ctx := context.TODO()
	tx := DB.MustBeginTx(ctx, nil)

	var err error
	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			if err := tx.Rollback(); err != nil {
				logger.Warnf("failure when performing transaction rollback: %v", err)
			}
			panic(p)
		}

		if err != nil {
			// something went wrong, rollback
			if err := tx.Rollback(); err != nil {
				logger.Warnf("failure when performing transaction rollback: %v", err)
			}
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
