package sqlite

import (
	"database/sql/driver"
	"fmt"
	"strconv"

	"github.com/WithoutPants/sortorder/casefolded"
	modernsqlite "modernc.org/sqlite"
)

const sqliteDriver = "sqlite"

func init() {
	funcs := map[string]struct {
		args int32
		fn   func([]driver.Value) (driver.Value, error)
	}{
		"regexp":            {args: 2, fn: regexpSQLiteFn},
		"durationToTinyInt": {args: 1, fn: durationToTinyIntSQLiteFn},
		"basename":          {args: 1, fn: basenameSQLiteFn},
		"phash_distance":    {args: 2, fn: phashDistanceSQLiteFn},
	}

	for name, fn := range funcs {
		if err := modernsqlite.RegisterScalarFunction(name, fn.args, func(_ *modernsqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			return fn.fn(args)
		}); err != nil {
			panic(fmt.Errorf("registering sqlite function %s: %w", name, err))
		}
	}

	// COLLATE NATURAL_CI - case-insensitive natural sort.
	if err := modernsqlite.RegisterCollationUtf8("NATURAL_CI", func(s, s2 string) int {
		if casefolded.NaturalLess(s, s2) {
			return -1
		}
		if casefolded.NaturalLess(s2, s) {
			return 1
		}
		return 0
	}); err != nil {
		panic(fmt.Errorf("registering natural sort collation: %w", err))
	}
}

func regexpSQLiteFn(args []driver.Value) (driver.Value, error) {
	ret, err := regexFn(sqliteString(args[0]), sqliteString(args[1]))
	return ret, err
}

func durationToTinyIntSQLiteFn(args []driver.Value) (driver.Value, error) {
	return durationToTinyIntFn(sqliteString(args[0]))
}

func basenameSQLiteFn(args []driver.Value) (driver.Value, error) {
	return basenameFn(sqliteString(args[0]))
}

func phashDistanceSQLiteFn(args []driver.Value) (driver.Value, error) {
	return phashDistanceFn(sqliteInt64(args[0]), sqliteInt64(args[1]))
}

func sqliteString(value driver.Value) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprint(v)
	}
}

func sqliteInt64(value driver.Value) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		ret, _ := strconv.ParseInt(v, 10, 64)
		return ret
	case []byte:
		ret, _ := strconv.ParseInt(string(v), 10, 64)
		return ret
	default:
		return 0
	}
}
