package sqlite

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type dbi struct{}

func (*dbi) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	return tx.Get(dest, query, args...)
}

func (*dbi) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	return tx.Select(dest, query, args...)
}

func (*dbi) Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.Queryx(query, args...)
}

func (*dbi) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.NamedExec(query, arg)
}

func (*dbi) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.Exec(query, args...)
}
