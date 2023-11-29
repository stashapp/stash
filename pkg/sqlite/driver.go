package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/WithoutPants/sortorder/casefolded"
	sqlite3 "github.com/mattn/go-sqlite3"
)

const sqlite3Driver = "sqlite3ex"

func init() {
	// register custom driver
	sql.Register(sqlite3Driver, &CustomSQLiteDriver{})
}

type CustomSQLiteDriver struct{}

type CustomSQLiteConn struct {
	*sqlite3.SQLiteConn
}

func (d *CustomSQLiteDriver) Open(dsn string) (driver.Conn, error) {
	sqlite3Driver := &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			funcs := map[string]interface{}{
				"regexp":            regexFn,
				"durationToTinyInt": durationToTinyIntFn,
				"basename":          basenameFn,
				"phash_distance":    phashDistanceFn,
			}

			for name, fn := range funcs {
				if err := conn.RegisterFunc(name, fn, true); err != nil {
					return fmt.Errorf("error registering function %s: %v", name, err)
				}
			}

			// COLLATE NATURAL_CI - Case insensitive natural sort
			err := conn.RegisterCollation("NATURAL_CI", func(s string, s2 string) int {
				if casefolded.NaturalLess(s, s2) {
					return -1
				} else {
					return 1
				}
			})

			if err != nil {
				return fmt.Errorf("error registering natural sort collation: %v", err)
			}

			return nil
		},
	}

	conn, err := sqlite3Driver.Open(dsn)
	if err != nil {
		return nil, err
	}

	return &CustomSQLiteConn{conn.(*sqlite3.SQLiteConn)}, nil
}

func (c *CustomSQLiteConn) Close() error {
	conn := c.SQLiteConn

	_, _ = conn.Exec("PRAGMA analysis_limit=1000; PRAGMA optimize;", []driver.Value{})

	return conn.Close()
}
