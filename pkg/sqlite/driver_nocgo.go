//go:build !cgo

package sqlite

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"

	"github.com/WithoutPants/sortorder/casefolded"
	"github.com/jmoiron/sqlx"
	"modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"

	"github.com/stashapp/stash/pkg/logger"
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

	// Define NATURAL_CI collation
	sqlite.MustRegisterCollationUtf8("NATURAL_CI", func(s1, s2 string) int {
		if casefolded.NaturalLess(s1, s2) {
			return -1
		} else {
			return 1
		}
	})
}

func createDBConn(dbPath string, disableForeignKeys bool) (*sqlx.DB, error) {
	// https://pkg.go.dev/modernc.org/sqlite#Driver.Open
	qs := url.Values{}
	qs.Add("_pragma", "busy_timeout(100)")
	qs.Add("_pragma", "journal_mode(WAL)")
	qs.Add("_pragma", "synchronous(NORMAL)")
	if disableForeignKeys {
		qs.Add("_pragma", "foreign_keys(0)")
	} else {
		qs.Add("_pragma", "foreign_keys(1)")
	}
	url := "file:" + dbPath + "?" + qs.Encode()

	logger.Debugf("Connecting to SQLite at '%s' (driver: non-CGo)", url)

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
		// This driver uses extended errors, so we need to grab the last 8 bits only
		// See: http://www.sqlite.org/c3ref/c_abort_rollback.html
		return sqliteErr.Code()&0xFF == sqlitelib.SQLITE_CONSTRAINT
	}
	return false
}
