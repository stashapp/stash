package sqlite

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fvbommel/sortorder"
	"github.com/golang-migrate/migrate/v4"
	sqlite3mig "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	sqlite3 "github.com/mattn/go-sqlite3"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

var appSchemaVersion uint = 31

//go:embed migrations/*.sql
var migrationsBox embed.FS

var (
	// ErrDatabaseNotInitialized indicates that the database is not
	// initialized, usually due to an incomplete configuration.
	ErrDatabaseNotInitialized = errors.New("database not initialized")
)

// ErrMigrationNeeded indicates that a database migration is needed
// before the database can be initialized
type MigrationNeededError struct {
	CurrentSchemaVersion  uint
	RequiredSchemaVersion uint
}

func (e *MigrationNeededError) Error() string {
	return fmt.Sprintf("database schema version %d does not match required schema version %d", e.CurrentSchemaVersion, e.RequiredSchemaVersion)
}

type MismatchedSchemaVersionError struct {
	CurrentSchemaVersion  uint
	RequiredSchemaVersion uint
}

func (e *MismatchedSchemaVersionError) Error() string {
	return fmt.Sprintf("schema version %d is incompatible with required schema version %d", e.CurrentSchemaVersion, e.RequiredSchemaVersion)
}

const sqlite3Driver = "sqlite3ex"

func init() {
	// register custom driver with regexp function
	registerCustomDriver()
}

type Database struct {
	db     *sqlx.DB
	dbPath string

	schemaVersion uint

	writeMu sync.Mutex
}

// Ready returns an error if the database is not ready to begin transactions.
func (db *Database) Ready() error {
	if db.db == nil {
		return ErrDatabaseNotInitialized
	}

	return nil
}

// Open initializes the database. If the database is new, then it
// performs a full migration to the latest schema version. Otherwise, any
// necessary migrations must be run separately using RunMigrations.
// Returns true if the database is new.
func (db *Database) Open(dbPath string) error {
	db.writeMu.Lock()
	defer db.writeMu.Unlock()

	db.dbPath = dbPath

	databaseSchemaVersion, err := db.getDatabaseSchemaVersion()
	if err != nil {
		return fmt.Errorf("getting database schema version: %w", err)
	}

	db.schemaVersion = databaseSchemaVersion

	if databaseSchemaVersion == 0 {
		// new database, just run the migrations
		if err := db.RunMigrations(); err != nil {
			return fmt.Errorf("error running initial schema migrations: %v", err)
		}
	} else {
		if databaseSchemaVersion > appSchemaVersion {
			return &MismatchedSchemaVersionError{
				CurrentSchemaVersion:  databaseSchemaVersion,
				RequiredSchemaVersion: appSchemaVersion,
			}
		}

		// if migration is needed, then don't open the connection
		if db.needsMigration() {
			return &MigrationNeededError{
				CurrentSchemaVersion:  databaseSchemaVersion,
				RequiredSchemaVersion: appSchemaVersion,
			}
		}
	}

	// RunMigrations may have opened a connection already
	if db.db == nil {
		const disableForeignKeys = false
		db.db, err = db.open(disableForeignKeys)
		if err != nil {
			return err
		}
	}

	if err := db.runCustomMigrations(); err != nil {
		return err
	}

	return nil
}

func (db *Database) Close() error {
	db.writeMu.Lock()
	defer db.writeMu.Unlock()

	if db.db != nil {
		if err := db.db.Close(); err != nil {
			return err
		}

		db.db = nil
	}

	return nil
}

func (db *Database) open(disableForeignKeys bool) (*sqlx.DB, error) {
	// https://github.com/mattn/go-sqlite3
	url := "file:" + db.dbPath + "?_journal=WAL&_sync=NORMAL"
	if !disableForeignKeys {
		url += "&_fk=true"
	}

	conn, err := sqlx.Open(sqlite3Driver, url)
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(4)
	conn.SetConnMaxLifetime(30 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("db.Open(): %w", err)
	}

	return conn, nil
}

func (db *Database) Reset() error {
	databasePath := db.dbPath
	err := db.Close()

	if err != nil {
		return errors.New("Error closing database: " + err.Error())
	}

	err = os.Remove(databasePath)
	if err != nil {
		return errors.New("Error removing database: " + err.Error())
	}

	// remove the -shm, -wal files ( if they exist )
	walFiles := []string{databasePath + "-shm", databasePath + "-wal"}
	for _, wf := range walFiles {
		if exists, _ := fsutil.FileExists(wf); exists {
			err = os.Remove(wf)
			if err != nil {
				return errors.New("Error removing database: " + err.Error())
			}
		}
	}

	if err := db.Open(databasePath); err != nil {
		return fmt.Errorf("[reset DB] unable to initialize: %w", err)
	}

	return nil
}

// Backup the database. If db is nil, then uses the existing database
// connection.
func (db *Database) Backup(backupPath string) error {
	thisDB := db.db
	if thisDB == nil {
		var err error
		thisDB, err = sqlx.Connect(sqlite3Driver, "file:"+db.dbPath+"?_fk=true")
		if err != nil {
			return fmt.Errorf("open database %s failed: %v", db.dbPath, err)
		}
		defer thisDB.Close()
	}

	logger.Infof("Backing up database into: %s", backupPath)
	_, err := thisDB.Exec(`VACUUM INTO "` + backupPath + `"`)
	if err != nil {
		return fmt.Errorf("vacuum failed: %v", err)
	}

	return nil
}

func (db *Database) RestoreFromBackup(backupPath string) error {
	logger.Infof("Restoring backup database %s into %s", backupPath, db.dbPath)
	return os.Rename(backupPath, db.dbPath)
}

// Migrate the database
func (db *Database) needsMigration() bool {
	return db.schemaVersion != appSchemaVersion
}

func (db *Database) AppSchemaVersion() uint {
	return appSchemaVersion
}

func (db *Database) DatabasePath() string {
	return db.dbPath
}

func (db *Database) DatabaseBackupPath() string {
	return fmt.Sprintf("%s.%d.%s", db.dbPath, db.schemaVersion, time.Now().Format("20060102_150405"))
}

func (db *Database) Version() uint {
	return db.schemaVersion
}

func (db *Database) getMigrate() (*migrate.Migrate, error) {
	migrations, err := iofs.New(migrationsBox, "migrations")
	if err != nil {
		panic(err.Error())
	}

	const disableForeignKeys = true
	conn, err := db.open(disableForeignKeys)
	if err != nil {
		return nil, err
	}

	driver, err := sqlite3mig.WithInstance(conn.DB, &sqlite3mig.Config{})
	if err != nil {
		return nil, err
	}

	// use sqlite3Driver so that migration has access to durationToTinyInt
	return migrate.NewWithInstance(
		"iofs",
		migrations,
		db.dbPath,
		driver,
	)
}

func (db *Database) getDatabaseSchemaVersion() (uint, error) {
	m, err := db.getMigrate()
	if err != nil {
		return 0, err
	}
	defer m.Close()

	ret, _, _ := m.Version()
	return ret, nil
}

// Migrate the database
func (db *Database) RunMigrations() error {
	m, err := db.getMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	databaseSchemaVersion, _, _ := m.Version()
	stepNumber := appSchemaVersion - databaseSchemaVersion
	if stepNumber != 0 {
		logger.Infof("Migrating database from version %d to %d", databaseSchemaVersion, appSchemaVersion)
		err = m.Steps(int(stepNumber))
		if err != nil {
			// migration failed
			return err
		}
	}

	// update the schema version
	db.schemaVersion, _, _ = m.Version()

	// re-initialise the database
	const disableForeignKeys = false
	db.db, err = db.open(disableForeignKeys)
	if err != nil {
		return fmt.Errorf("re-initializing the database: %w", err)
	}

	// run a vacuum on the database
	logger.Info("Performing vacuum on database")
	_, err = db.db.Exec("VACUUM")
	if err != nil {
		logger.Warnf("error while performing post-migration vacuum: %v", err)
	}

	return nil
}

func registerCustomDriver() {
	sql.Register(sqlite3Driver,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				funcs := map[string]interface{}{
					"regexp":            regexFn,
					"durationToTinyInt": durationToTinyIntFn,
				}

				for name, fn := range funcs {
					if err := conn.RegisterFunc(name, fn, true); err != nil {
						return fmt.Errorf("error registering function %s: %s", name, err.Error())
					}
				}

				// COLLATE NATURAL_CS - Case sensitive natural sort
				err := conn.RegisterCollation("NATURAL_CS", func(s string, s2 string) int {
					if sortorder.NaturalLess(s, s2) {
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
		},
	)
}
