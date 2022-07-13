package sqlite

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type dbReader interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type dbWrapper struct{}

func (*dbWrapper) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, err := getDBReader(ctx)
	if err != nil {
		return err
	}

	return tx.Get(dest, query, args...)
}

func (*dbWrapper) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, err := getDBReader(ctx)
	if err != nil {
		return err
	}

	return tx.Select(dest, query, args...)
}

func (*dbWrapper) Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	tx, err := getDBReader(ctx)
	if err != nil {
		return nil, err
	}

	return tx.Queryx(query, args...)
}

func (*dbWrapper) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.NamedExec(query, arg)
}

func (*dbWrapper) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.Exec(query, args...)
}
