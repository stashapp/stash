package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/jmoiron/sqlx"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

var DB *sqlx.DB
var dbPath string
var appSchemaVersion uint = 4
var databaseSchemaVersion uint

const sqlite3Driver = "sqlite3_regexp"

func init() {
	// register custom driver with regexp function
	registerRegexpFunc()
}

func Initialize(databasePath string) {
	dbPath = databasePath

	if err := getDatabaseSchemaVersion(); err != nil {
		panic(err)
	}

	if databaseSchemaVersion == 0 {
		// new database, just run the migrations
		if err := RunMigrations(); err != nil {
			panic(err)
		}
		// RunMigrations calls Initialise. Just return
		return
	} else {
		if databaseSchemaVersion > appSchemaVersion {
			panic(fmt.Sprintf("Database schema version %d is incompatible with required schema version %d", databaseSchemaVersion, appSchemaVersion))
		}

		// if migration is needed, then don't open the connection
		if NeedsMigration() {
			logger.Warnf("Database schema version %d does not match required schema version %d.", databaseSchemaVersion, appSchemaVersion)
			return
		}
	}

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

// Backup the database if a migration is needed
// or force is set to true
func Backup(backupPath string) error {
	db, err := sqlx.Connect(sqlite3Driver, "file:"+dbPath+"?_fk=true")
	if err != nil {
		return fmt.Errorf("Open database %s failed:%s", dbPath, err)
	}
	defer db.Close()

	_, err = db.Exec(`VACUUM INTO "` + backupPath + `"`)
	if err != nil {
		return fmt.Errorf("Vacuum failed: %s", err)
	}

	return nil
}

// Migrate the database
func NeedsMigration() bool {
	return databaseSchemaVersion != appSchemaVersion
}

func AppSchemaVersion() uint {
	return appSchemaVersion
}

func DatabaseBackupPath() string {
	return fmt.Sprintf("%s.%d.%s", dbPath, databaseSchemaVersion, time.Now().Format("20060102_150405"))
}

func Version() uint {
	return databaseSchemaVersion
}

func getMigrate() (*migrate.Migrate, error) {
	migrationsBox := packr.New("Migrations Box", "./migrations")
	packrSource := &Packr2Source{
		Box:        migrationsBox,
		Migrations: source.NewMigrations(),
	}

	databasePath := utils.FixWindowsPath(dbPath)
	s, _ := WithInstance(packrSource)
	return migrate.NewWithSourceInstance(
		"packr2",
		s,
		fmt.Sprintf("sqlite3://%s", "file:"+databasePath),
	)
}

func getDatabaseSchemaVersion() error {
	m, err := getMigrate()
	if err != nil {
		return err
	}

	databaseSchemaVersion, _, _ = m.Version()
	m.Close()
	return nil
}

// Migrate the database
func RunMigrations() error {
	m, err := getMigrate()
	if err != nil {
		panic(err.Error())
	}

	databaseSchemaVersion, _, _ = m.Version()
	stepNumber := appSchemaVersion - databaseSchemaVersion
	if stepNumber != 0 {
		err = m.Steps(int(stepNumber))
		if err != nil {
			return err
		}
	}
	m.Close()

	// re-initialise the database
	Initialize(dbPath)

	return nil
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
