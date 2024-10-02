package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
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

const (
	maxWriteConnections = 1
	// Number of database read connections to use
	// The same value is used for both the maximum and idle limit,
	// to prevent opening connections on the fly which has a notieable performance penalty.
	// Fewer connections use less memory, more connections increase performance,
	// but have diminishing returns.
	// 10 was found to be a good tradeoff.
	maxReadConnections = 10
	// Idle connection timeout, in seconds
	// Closes a connection after a period of inactivity, which saves on memory and
	// causes the sqlite -wal and -shm files to be automatically deleted.
	dbConnTimeout = 30 * time.Second

	// environment variable to set the cache size
	cacheSizeEnv = "STASH_SQLITE_CACHE_SIZE"
)

var appSchemaVersion uint = 67

//go:embed migrations/*.sql migrationsPostgres/*.sql
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

type storeRepository struct {
	Blobs          *BlobStore
	File           *FileStore
	Folder         *FolderStore
	Image          *ImageStore
	Gallery        *GalleryStore
	GalleryChapter *GalleryChapterStore
	Scene          *SceneStore
	SceneMarker    *SceneMarkerStore
	Performer      *PerformerStore
	SavedFilter    *SavedFilterStore
	Studio         *StudioStore
	Tag            *TagStore
	Group          *GroupStore
}

type DatabaseType string

const (
	PostgresBackend DatabaseType = "POSTGRES"
	SqliteBackend   DatabaseType = "SQLITE"
)

type Database struct {
	*storeRepository

	readDB   *sqlx.DB
	writeDB  *sqlx.DB
	dbPath   string
	dbType   DatabaseType
	dbString string

	schemaVersion uint

	lockChan chan struct{}
}

func NewDatabase() *Database {
	fileStore := NewFileStore()
	folderStore := NewFolderStore()
	galleryStore := NewGalleryStore(fileStore, folderStore)
	blobStore := NewBlobStore(BlobStoreOptions{})
	performerStore := NewPerformerStore(blobStore)
	studioStore := NewStudioStore(blobStore)
	tagStore := NewTagStore(blobStore)

	r := &storeRepository{}
	*r = storeRepository{
		Blobs:          blobStore,
		File:           fileStore,
		Folder:         folderStore,
		Scene:          NewSceneStore(r, blobStore),
		SceneMarker:    NewSceneMarkerStore(),
		Image:          NewImageStore(r),
		Gallery:        galleryStore,
		GalleryChapter: NewGalleryChapterStore(),
		Performer:      performerStore,
		Studio:         studioStore,
		Tag:            tagStore,
		Group:          NewGroupStore(blobStore),
		SavedFilter:    NewSavedFilterStore(),
	}

	ret := &Database{
		storeRepository: r,
		lockChan:        make(chan struct{}, 1),
	}

	return ret
}

func (db *Database) SetBlobStoreOptions(options BlobStoreOptions) {
	*db.Blobs = *NewBlobStore(options)
}

// Ready returns an error if the database is not ready to begin transactions.
func (db *Database) Ready() error {
	if db.readDB == nil || db.writeDB == nil {
		return ErrDatabaseNotInitialized
	}

	return nil
}

func (db *Database) OpenPostgres(dbConnector string) error {
	db.dbType = PostgresBackend
	db.dbString = dbConnector

	dialect = goqu.Dialect("postgres")

	return db.OpenGeneric()
}

func RegisterSqliteDialect() {
	opts := sqlite3.DialectOptions()
	opts.SupportsReturn = true
	goqu.RegisterDialect("sqlite3new", opts)
}

func (db *Database) OpenSqlite(dbPath string) error {
	db.dbType = SqliteBackend
	db.dbPath = dbPath

	dialect = goqu.Dialect("sqlite3new")

	return db.OpenGeneric()
}

// Open initializes the database. If the database is new, then it
// performs a full migration to the latest schema version. Otherwise, any
// necessary migrations must be run separately using RunMigrations.
// Returns true if the database is new.
func (db *Database) OpenGeneric() error {
	goqu.SetDefaultPrepared(false)

	databaseSchemaVersion, err := db.getDatabaseSchemaVersion()
	if err != nil {
		return fmt.Errorf("getting database schema version: %w", err)
	}

	db.schemaVersion = databaseSchemaVersion

	isNew := databaseSchemaVersion == 0

	if isNew {
		// new database, just run the migrations
		if err := db.RunAllMigrations(); err != nil {
			return fmt.Errorf("error running initial schema migrations: %w", err)
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

	if err := db.initialise(); err != nil {
		return err
	}

	if isNew {
		// optimize database after migration
		err = db.Optimise(context.Background())
		if err != nil {
			logger.Warnf("error while performing post-migration optimisation: %v", err)
		}
	}

	return nil
}

// lock locks the database for writing. This method will block until the lock is acquired.
func (db *Database) lock() {
	db.lockChan <- struct{}{}
}

// unlock unlocks the database
func (db *Database) unlock() {
	// will block the caller if the lock is not held, so check first
	select {
	case <-db.lockChan:
		return
	default:
		panic("database is not locked")
	}
}

func (db *Database) Close() error {
	db.lock()
	defer db.unlock()

	if db.readDB != nil {
		if err := db.readDB.Close(); err != nil {
			return err
		}

		db.readDB = nil
	}
	if db.writeDB != nil {
		if err := db.writeDB.Close(); err != nil {
			return err
		}

		db.writeDB = nil
	}

	return nil
}

func (db *Database) open(disableForeignKeys bool, writable bool) (conn *sqlx.DB, err error) {
	// Fail-safe
	err = errors.New("missing backend type")

	// https://github.com/mattn/go-sqlite3
	if db.dbType == SqliteBackend {
		url := "file:" + db.dbPath + "?_journal=WAL&_sync=NORMAL&_busy_timeout=50"
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
	}
	if db.dbType == PostgresBackend {
		conn, err = sqlx.Open("postgres", db.dbString)
		if err == nil {
			if disableForeignKeys {
				conn.Exec("SET session_replication_role = replica;")
			}
			if !writable {
				conn.Exec("SET default_transaction_read_only = ON;")
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("db.Open(): %w", err)
	}

	return conn, nil
}

func (db *Database) initialise() error {
	if err := db.openReadDB(); err != nil {
		return fmt.Errorf("opening read database: %w", err)
	}
	if err := db.openWriteDB(); err != nil {
		return fmt.Errorf("opening write database: %w", err)
	}

	return nil
}

func (db *Database) openReadDB() error {
	const (
		disableForeignKeys = false
		writable           = false
	)
	var err error
	db.readDB, err = db.open(disableForeignKeys, writable)
	db.readDB.SetMaxOpenConns(maxReadConnections)
	db.readDB.SetMaxIdleConns(maxReadConnections)
	db.readDB.SetConnMaxIdleTime(dbConnTimeout)
	return err
}

func (db *Database) openWriteDB() error {
	const (
		disableForeignKeys = false
		writable           = true
	)
	var err error
	db.writeDB, err = db.open(disableForeignKeys, writable)
	db.writeDB.SetMaxOpenConns(maxWriteConnections)
	db.writeDB.SetMaxIdleConns(maxWriteConnections)
	db.writeDB.SetConnMaxIdleTime(dbConnTimeout)
	return err
}

func (db *Database) Remove() error {
	if db.dbType == PostgresBackend {
		logger.Warn("Postgres backend detected, ignoring Remove request")
		return nil
	}

	databasePath := db.dbPath
	err := db.Close()

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

func (db *Database) Reset() error {
	if db.dbType == PostgresBackend {
		logger.Warn("Postgres backend detected, ignoring Reset request")
		return nil
	}

	if err := db.Remove(); err != nil {
		return err
	}

	if err := db.OpenSqlite(db.dbPath); err != nil {
		return fmt.Errorf("[reset DB] unable to initialize: %w", err)
	}

	return nil
}

// Backup the database. If db is nil, then uses the existing database
// connection.
func (db *Database) Backup(backupPath string) (err error) {
	if db.dbType == PostgresBackend {
		logger.Warn("Postgres backend detected, ignoring Backup request")
		return nil
	}

	thisDB := db.writeDB
	if thisDB == nil {
		thisDB, err = sqlx.Connect(sqlite3Driver, "file:"+db.dbPath+"?_fk=true")
		if err != nil {
			return fmt.Errorf("open database %s failed: %w", db.dbPath, err)
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

func (db *Database) Anonymise(outPath string) error {
	anon, err := NewAnonymiser(db, outPath)

	if err != nil {
		return err
	}

	return anon.Anonymise(context.Background())
}

func (db *Database) RestoreFromBackup(backupPath string) error {
	if db.dbType == PostgresBackend {
		logger.Warn("Postgres backend detected, ignoring RestoreFromBackup request")
		return nil
	}

	logger.Infof("Restoring backup database %s into %s", backupPath, db.dbPath)
	return os.Rename(backupPath, db.dbPath)
}

func (db *Database) AppSchemaVersion() uint {
	return appSchemaVersion
}

func (db *Database) DatabasePath() string {
	return db.dbPath
}

func (db *Database) DatabaseBackupPath(backupDirectoryPath string) string {
	fn := fmt.Sprintf("%s.%d.%s", filepath.Base(db.dbPath), db.schemaVersion, time.Now().Format("20060102_150405"))

	if backupDirectoryPath != "" {
		return filepath.Join(backupDirectoryPath, fn)
	}

	return fn
}

func (db *Database) AnonymousDatabasePath(backupDirectoryPath string) string {
	fn := fmt.Sprintf("%s.anonymous.%d.%s", filepath.Base(db.dbPath), db.schemaVersion, time.Now().Format("20060102_150405"))

	if backupDirectoryPath != "" {
		return filepath.Join(backupDirectoryPath, fn)
	}

	return fn
}

func (db *Database) Version() uint {
	return db.schemaVersion
}

func (db *Database) Optimise(ctx context.Context) error {
	logger.Info("Optimising database")

	err := db.Analyze(ctx)
	if err != nil {
		return fmt.Errorf("performing optimization: %w", err)
	}

	err = db.Vacuum(ctx)
	if err != nil {
		return fmt.Errorf("performing vacuum: %w", err)
	}

	return nil
}

// Vacuum runs a VACUUM on the database, rebuilding the database file into a minimal amount of disk space.
func (db *Database) Vacuum(ctx context.Context) error {
	_, err := db.writeDB.ExecContext(ctx, "VACUUM")
	return err
}

// Analyze runs an ANALYZE on the database to improve query performance.
func (db *Database) Analyze(ctx context.Context) error {
	_, err := db.writeDB.ExecContext(ctx, "ANALYZE")
	return err
}

func (db *Database) ExecSQL(ctx context.Context, query string, args []interface{}) (*int64, error) {
	wrapper := dbWrapperType{}

	result, err := wrapper.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var rowsAffected *int64
	ra, err := result.RowsAffected()
	if err == nil {
		rowsAffected = &ra
	}

	return rowsAffected, nil
}

func (db *Database) QuerySQL(ctx context.Context, query string, args []interface{}) ([]string, [][]interface{}, error) {
	wrapper := dbWrapperType{}

	rows, err := wrapper.QueryxContext(ctx, query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var ret [][]interface{}

	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			return nil, nil, err
		}
		ret = append(ret, row)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return cols, ret, nil
}
