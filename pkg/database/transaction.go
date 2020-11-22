package database

import (
	"context"

	"github.com/jmoiron/sqlx"
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
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
