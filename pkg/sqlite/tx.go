package sqlite

import (
	"context"
	"database/sql"
	"fmt"
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

type stmt struct {
	*sql.Stmt
	query string
}

func logSQL(start time.Time, query string, args ...interface{}) {
	since := time.Since(start)
	if since >= slowLogTime {
		logger.Debugf("SLOW SQL [%v]: %s, args: %v", since, query, args)
	} else {
		logger.Tracef("SQL [%v]: %s, args: %v", since, query, args)
	}
}

type dbWrapperType struct{}

var dbWrapper = dbWrapperType{}

func sqlError(err error, sql string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("error executing `%s` [%v]: %w", sql, args, err)
}

func (*dbWrapperType) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, err := getDBReader(ctx)
	if err != nil {
		return sqlError(err, query, args...)
	}

	start := time.Now()
	err = tx.Get(dest, query, args...)
	logSQL(start, query, args...)

	return sqlError(err, query, args...)
}

func (*dbWrapperType) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, err := getDBReader(ctx)
	if err != nil {
		return sqlError(err, query, args...)
	}

	start := time.Now()
	err = tx.Select(dest, query, args...)
	logSQL(start, query, args...)

	return sqlError(err, query, args...)
}

func (*dbWrapperType) Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	tx, err := getDBReader(ctx)
	if err != nil {
		return nil, sqlError(err, query, args...)
	}

	start := time.Now()
	ret, err := tx.Queryx(query, args...)
	logSQL(start, query, args...)

	return ret, sqlError(err, query, args...)
}

func (*dbWrapperType) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	tx, err := getDBReader(ctx)
	if err != nil {
		return nil, sqlError(err, query, args...)
	}

	start := time.Now()
	ret, err := tx.QueryxContext(ctx, query, args...)
	logSQL(start, query, args...)

	return ret, sqlError(err, query, args...)
}

func (*dbWrapperType) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, sqlError(err, query, arg)
	}

	start := time.Now()
	ret, err := tx.NamedExec(query, arg)
	logSQL(start, query, arg)

	return ret, sqlError(err, query, arg)
}

func (*dbWrapperType) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, sqlError(err, query, args...)
	}

	start := time.Now()
	ret, err := tx.Exec(query, args...)
	logSQL(start, query, args...)

	return ret, sqlError(err, query, args...)
}

// Prepare creates a prepared statement.
func (*dbWrapperType) Prepare(ctx context.Context, query string, args ...interface{}) (*stmt, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, sqlError(err, query, args...)
	}

	// nolint:sqlclosecheck
	ret, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, sqlError(err, query, args...)
	}

	return &stmt{
		query: query,
		Stmt:  ret,
	}, nil
}

func (*dbWrapperType) ExecStmt(ctx context.Context, stmt *stmt, args ...interface{}) (sql.Result, error) {
	_, err := getTx(ctx)
	if err != nil {
		return nil, sqlError(err, stmt.query, args...)
	}

	start := time.Now()
	ret, err := stmt.ExecContext(ctx, args...)
	logSQL(start, stmt.query, args...)

	return ret, sqlError(err, stmt.query, args...)
}
