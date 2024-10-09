package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

type SQLiteDB Database

func RegisterSqliteDialect() {
	opts := sqlite3.DialectOptions()
	opts.SupportsReturn = true
	goqu.RegisterDialect("sqlite3new", opts)
}

func NewSQLiteDatabase(dbPath string) *Database {
	dialect = goqu.Dialect("sqlite3new")

	db := &SQLiteDB{
		storeRepository: newDatabase(),
		lockChan:        make(chan struct{}, 1),
		dbConfig:        dbPath,
	}
	db.dbInterface = db

	dbWrapper.dbType = SqliteBackend

	return (*Database)(db)
}

func (db *SQLiteDB) DatabaseType() DatabaseType {
	return SqliteBackend
}

func (db *SQLiteDB) AppSchemaVersion() uint {
	return appSchemaVersion
}

func (db *SQLiteDB) DatabasePath() string {
	return (db.dbConfig).(string)
}

func (db *SQLiteDB) open(disableForeignKeys bool, writable bool) (conn *sqlx.DB, err error) {
	// https://github.com/mattn/go-sqlite3
	url := "file:" + db.DatabasePath() + "?_journal=WAL&_sync=NORMAL&_busy_timeout=50"
	if !disableForeignKeys {
		url += "&_fk=true"
	}

	if writable {
		url += "&_txlock=immediate"
	} else {
		url += "&mode=ro"
	}

	// #5155 - set the cache size if the environment variable is set
	// default is -2000 which is 2MB
	if cacheSize := os.Getenv(cacheSizeEnv); cacheSize != "" {
		url += "&_cache_size=" + cacheSize
	}

	conn, err = sqlx.Open(sqlite3Driver, url)

	if err != nil {
		return nil, fmt.Errorf("db.Open(): %w", err)
	}

	return conn, nil
}

func (db *SQLiteDB) Remove() error {
	databasePath := db.DatabasePath()
	err := (*Database)(db).Close()

	if err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}

	err = os.Remove(databasePath)
	if err != nil {
		return fmt.Errorf("error removing database: %w", err)
	}

	// remove the -shm, -wal files ( if they exist )
	walFiles := []string{databasePath + "-shm", databasePath + "-wal"}
	for _, wf := range walFiles {
		if exists, _ := fsutil.FileExists(wf); exists {
			err = os.Remove(wf)
			if err != nil {
				return fmt.Errorf("error removing database: %w", err)
			}
		}
	}

	return nil
}

func (db *SQLiteDB) Reset() error {
	if err := db.Remove(); err != nil {
		return err
	}

	if err := (*Database)(db).Open(); err != nil {
		return fmt.Errorf("[reset DB] unable to initialize: %w", err)
	}

	return nil
}

// Backup the database. If db is nil, then uses the existing database
// connection.
func (db *SQLiteDB) Backup(backupPath string) (err error) {
	thisDB := db.writeDB
	if thisDB == nil {
		thisDB, err = sqlx.Connect(sqlite3Driver, "file:"+db.DatabasePath()+"?_fk=true")
		if err != nil {
			return fmt.Errorf("open database %s failed: %w", db.DatabasePath(), err)
		}
		defer thisDB.Close()
	}

	logger.Infof("Backing up database into: %s", backupPath)
	_, err = thisDB.Exec(`VACUUM INTO "` + backupPath + `"`)
	if err != nil {
		return fmt.Errorf("vacuum failed: %w", err)
	}

	return nil
}

func (db *SQLiteDB) RestoreFromBackup(backupPath string) error {
	logger.Infof("Restoring backup database %s into %s", backupPath, db.DatabasePath())
	return os.Rename(backupPath, db.DatabasePath())
}

func (db *SQLiteDB) DatabaseBackupPath(backupDirectoryPath string) string {
	fn := fmt.Sprintf("%s.%d.%s", filepath.Base(db.DatabasePath()), db.schemaVersion, time.Now().Format("20060102_150405"))

	if backupDirectoryPath != "" {
		return filepath.Join(backupDirectoryPath, fn)
	}

	return fn
}

func (db *SQLiteDB) AnonymousDatabasePath(backupDirectoryPath string) string {
	fn := fmt.Sprintf("%s.anonymous.%d.%s", filepath.Base(db.DatabasePath()), db.schemaVersion, time.Now().Format("20060102_150405"))

	if backupDirectoryPath != "" {
		return filepath.Join(backupDirectoryPath, fn)
	}

	return fn
}
