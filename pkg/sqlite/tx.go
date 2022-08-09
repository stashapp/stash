package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
)

const (
	slowLogTime = time.Millisecond * 200
)

type dbReader interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

func logSQL(start time.Time, query string, args ...interface{}) {
	since := time.Since(start)
	if since >= slowLogTime {
		logger.Debugf("SLOW SQL [%v]: %s, args: %v", since, query, args)
	} else {
		logger.Tracef("SQL [%v]: %s, args: %v", since, query, args)
	}
}

type dbWrapper struct{}

func (*dbWrapper) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, err := getDBReader(ctx)
	if err != nil {
		return err
	}

	start := time.Now()
	err = tx.Get(dest, query, args...)
	logSQL(start, query, args...)

	return err
}

func (*dbWrapper) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, err := getDBReader(ctx)
	if err != nil {
		return err
	}

	start := time.Now()
	err = tx.Select(dest, query, args...)
	logSQL(start, query, args...)

	return err
}

func (*dbWrapper) Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	tx, err := getDBReader(ctx)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	ret, err := tx.Queryx(query, args...)
	logSQL(start, query, args...)

	return ret, err
}

func (*dbWrapper) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	tx, err := getDBReader(ctx)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	ret, err := tx.QueryxContext(ctx, query, args...)
	logSQL(start, query, args...)

	return ret, err
}

func (*dbWrapper) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	ret, err := tx.NamedExec(query, arg)
	logSQL(start, query, arg)

	return ret, err
}

func (*dbWrapper) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	ret, err := tx.Exec(query, args...)
	logSQL(start, query, args...)

	return ret, err
}
