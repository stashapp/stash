package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type key int

const (
	txnKey key = iota + 1
)

type Database struct {
	DB *sqlx.DB
}

func (db *Database) Begin(ctx context.Context) (context.Context, error) {
	if tx, _ := getTx(ctx); tx != nil {
		return nil, fmt.Errorf("already in transaction")
	}

	tx, err := db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}

	return context.WithValue(ctx, txnKey, tx), nil
}

func (db *Database) Commit(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *Database) Rollback(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}
	return tx.Rollback()
}

func getTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(txnKey).(*sqlx.Tx)
	if !ok || tx == nil {
		return nil, fmt.Errorf("not in transaction")
	}
	return tx, nil
}
