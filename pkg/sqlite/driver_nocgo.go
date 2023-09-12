//go:build !cgo

package sqlite

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	"modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

const sqlite3Driver = "sqlite"

func init() {
	// regexp
	sqlite.MustRegisterScalarFunction("regexp", 2, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("invalid args length: %d", len(args))
		}

		re, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("arg0 is not string: %T", args[0])
		}
		s, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("arg1 is not string: %T", args[1])
		}

		return regexFn(re, s)
	})

	// durationToTinyInt
	sqlite.MustRegisterDeterministicScalarFunction("durationToTinyInt", 1, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("invalid args length: %d", len(args))
		}

		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("arg0 is not string: %T", args[0])
		}

		return durationToTinyIntFn(str)
	})

	// basename
	sqlite.MustRegisterDeterministicScalarFunction("basename", 1, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("invalid args length: %d", len(args))
		}

		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("arg0 is not string: %T", args[0])
		}

		return basenameFn(str)
	})

	// phash_distance
	sqlite.MustRegisterDeterministicScalarFunction("phash_distance", 2, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("invalid args length: %d", len(args))
		}

		phash1, ok := args[0].(int64)
		if !ok {
			return nil, fmt.Errorf("arg0 is not int64: %T", args[0])
		}
		phash2, ok := args[1].(int64)
		if !ok {
			return nil, fmt.Errorf("arg1 is not int64: %T", args[1])
		}

		return phashDistanceFn(phash1, phash2)
	})

	// TODO: Define NATURAL_CI collation
	// Blocked by https://gitlab.com/cznic/sqlite/-/issues/163
	/*
		err := conn.RegisterCollation("NATURAL_CI", func(s string, s2 string) int {
					if casefolded.NaturalLess(s, s2) {
						return -1
					} else {
						return 1
					}
				})
	*/
}

func createDBConn(dbPath string, disableForeignKeys bool) (*sqlx.DB, error) {
	// https://pkg.go.dev/modernc.org/sqlite#Driver.Open
	var qs url.Values
	qs.Set("_txlock", "immediate")
	qs.Add("_pragma", "busy_timeout(50)")
	qs.Add("_pragma", "journal_mode(WAL)")
	qs.Add("_pragma", "synchronous(NORMAL)")
	if !disableForeignKeys {
		qs.Add("_pragma", "foreign_keys(true)")
	}
	url := "sqlite://file:" + dbPath + "?" + qs.Encode()

	return sqlx.Open("sqlite", url)
}

func IsLockedError(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code() == sqlitelib.SQLITE_BUSY
	}
	return false
}

func IsConstraintError(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code() == sqlitelib.SQLITE_CONSTRAINT
	}
	return false
}
