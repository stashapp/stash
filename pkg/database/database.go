package database

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

var DB *sqlx.DB
var WriteMu sync.Mutex
var dbPath string
var appSchemaVersion uint = 31
var databaseSchemaVersion uint

//go:embed migrations/*.sql
var migrationsBox embed.FS

var (
	// ErrMigrationNeeded indicates that a database migration is needed
	// before the database can be initialized
	ErrMigrationNeeded = errors.New("database migration required")

	// ErrDatabaseNotInitialized indicates that the database is not
	// initialized, usually due to an incomplete configuration.
	ErrDatabaseNotInitialized = errors.New("database not initialized")
)

const sqlite3Driver = "sqlite3ex"

// Ready returns an error if the database is not ready to begin transactions.
func Ready() error {
	if DB == nil {
		return ErrDatabaseNotInitialized
	}

	return nil
}

func init() {
	// register custom driver with regexp function
	registerCustomDriver()
}

// Initialize initializes the database. If the database is new, then it
// performs a full migration to the latest schema version. Otherwise, any
// necessary migrations must be run separately using RunMigrations.
// Returns true if the database is new.
func Initialize(databasePath string) error {
	dbPath = databasePath

	if err := getDatabaseSchemaVersion(); err != nil {
		return fmt.Errorf("error getting database schema version: %v", err)
	}

	if databaseSchemaVersion == 0 {
		// new database, just run the migrations
		if err := RunMigrations(); err != nil {
			return fmt.Errorf("error running initial schema migrations: %v", err)
		}
		// RunMigrations calls Initialise. Just return
		return nil
	} else {
		if databaseSchemaVersion > appSchemaVersion {
			panic(fmt.Sprintf("Database schema version %d is incompatible with required schema version %d", databaseSchemaVersion, appSchemaVersion))
		}

		// if migration is needed, then don't open the connection
		if NeedsMigration() {
			logger.Warnf("Database schema version %d does not match required schema version %d.", databaseSchemaVersion, appSchemaVersion)
			return nil
		}
	}

	const disableForeignKeys = false
	DB = open(databasePath, disableForeignKeys)

	if err := runCustomMigrations(); err != nil {
		return err
	}

	return nil
}

func Close() error {
	WriteMu.Lock()
	defer WriteMu.Unlock()

	if DB != nil {
		if err := DB.Close(); err != nil {
			return err
		}

		DB = nil
	}

	return nil
}

func open(databasePath string, disableForeignKeys bool) *sqlx.DB {
	// https://github.com/mattn/go-sqlite3
	url := "file:" + databasePath + "?_journal=WAL&_sync=NORMAL"
	if !disableForeignKeys {
		url += "&_fk=true"
	}

	conn, err := sqlx.Open(sqlite3Driver, url)
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(4)
	conn.SetConnMaxLifetime(30 * time.Second)
	if err != nil {
		logger.Fatalf("db.Open(): %q\n", err)
	}

	return conn
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

	if err := Initialize(databasePath); err != nil {
		return fmt.Errorf("[reset DB] unable to initialize: %w", err)
	}

	return nil
}

// Backup the database. If db is nil, then uses the existing database
// connection.
func Backup(db *sqlx.DB, backupPath string) error {
	if db == nil {
		var err error
		db, err = sqlx.Connect(sqlite3Driver, "file:"+dbPath+"?_fk=true")
		if err != nil {
			return fmt.Errorf("open database %s failed: %v", dbPath, err)
		}
		defer db.Close()
	}

	logger.Infof("Backing up database into: %s", backupPath)
	_, err := db.Exec(`VACUUM INTO "` + backupPath + `"`)
	if err != nil {
		return fmt.Errorf("vacuum failed: %v", err)
	}

	return nil
}

func RestoreFromBackup(backupPath string) error {
	logger.Infof("Restoring backup database %s into %s", backupPath, dbPath)
	return os.Rename(backupPath, dbPath)
}

// Migrate the database
func NeedsMigration() bool {
	return databaseSchemaVersion != appSchemaVersion
}

func AppSchemaVersion() uint {
	return appSchemaVersion
}

func DatabasePath() string {
	return dbPath
}

func DatabaseBackupPath() string {
	return fmt.Sprintf("%s.%d.%s", dbPath, databaseSchemaVersion, time.Now().Format("20060102_150405"))
}

func Version() uint {
	return databaseSchemaVersion
}

func getMigrate() (*migrate.Migrate, error) {
	migrations, err := iofs.New(migrationsBox, "migrations")
	if err != nil {
		panic(err.Error())
	}

	const disableForeignKeys = true
	conn := open(dbPath, disableForeignKeys)

	driver, err := sqlite3mig.WithInstance(conn.DB, &sqlite3mig.Config{})
	if err != nil {
		return nil, err
	}

	// use sqlite3Driver so that migration has access to durationToTinyInt
	return migrate.NewWithInstance(
		"iofs",
		migrations,
		dbPath,
		driver,
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
	defer m.Close()

	databaseSchemaVersion, _, _ = m.Version()
	stepNumber := appSchemaVersion - databaseSchemaVersion
	if stepNumber != 0 {
		logger.Infof("Migrating database from version %d to %d", databaseSchemaVersion, appSchemaVersion)
		err = m.Steps(int(stepNumber))
		if err != nil {
			// migration failed
			return err
		}
	}

	// re-initialise the database
	if err = Initialize(dbPath); err != nil {
		logger.Warnf("Error re-initializing the database: %v", err)
	}

	// run a vacuum on the database
	logger.Info("Performing vacuum on database")
	_, err = DB.Exec("VACUUM")
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
