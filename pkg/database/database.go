package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/gobuffalo/packr/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/jmoiron/sqlx"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

var DB *sqlx.DB
var appSchemaVersion uint = 2

const sqlite3Driver = "sqlite3_regexp"

func init() {
	// register custom driver with regexp function
	registerRegexpFunc()
}

func Initialize(databasePath string) {
	runMigrations(databasePath)

	// https://github.com/mattn/go-sqlite3
	conn, err := sqlx.Open(sqlite3Driver, "file:"+databasePath+"?_fk=true")
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(4)
	if err != nil {
		logger.Fatalf("db.Open(): %q\n", err)
	}
	DB = conn
}

func Reset(databasePath string) error {
	err := DB.Close()

	if err != nil {
		return errors.New("Error closing database: " + err.Error())
	}

	err = os.Remove(databasePath)
	if err != nil {
		return errors.New("Error removing database: " + err.Error())
	}

	Initialize(databasePath)
	return nil
}

// Migrate the database
func runMigrations(databasePath string) {
	migrationsBox := packr.New("Migrations Box", "./migrations")
	packrSource := &Packr2Source{
		Box:        migrationsBox,
		Migrations: source.NewMigrations(),
	}

	databasePath = utils.FixWindowsPath(databasePath)
	s, _ := WithInstance(packrSource)
	m, err := migrate.NewWithSourceInstance(
		"packr2",
		s,
		fmt.Sprintf("sqlite3://%s", "file:"+databasePath),
	)
	if err != nil {
		panic(err.Error())
	}

	databaseSchemaVersion, _, _ := m.Version()
	stepNumber := appSchemaVersion - databaseSchemaVersion
	if stepNumber != 0 {
		err = m.Steps(int(stepNumber))
		if err != nil {
			panic(err.Error())
		}
	}
	m.Close()
}

func registerRegexpFunc() {
	regexFn := func(re, s string) (bool, error) {
		return regexp.MatchString(re, s)
	}

	sql.Register(sqlite3Driver,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				return conn.RegisterFunc("regexp", regexFn, true)
			},
		})
}
