package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
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

var appSchemaVersion uint = 69

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
	PostgresBackend DatabaseType = "POSTGRESQL"
	SqliteBackend   DatabaseType = "SQLITE"
)

type DBInterface interface {
	Analyze(ctx context.Context) error
	Anonymise(outPath string) error
	AnonymousDatabasePath(backupDirectoryPath string) string
	AppSchemaVersion() uint
	Backup(backupPath string) (err error)
	Begin(ctx context.Context, writable bool) (context.Context, error)
	Close() error
	Commit(ctx context.Context) error
	DatabaseBackupPath(backupDirectoryPath string) string
	DatabasePath() string
	DatabaseType() DatabaseType
	ExecSQL(ctx context.Context, query string, args []interface{}) (*int64, error)
	getDatabaseSchemaVersion() (uint, error)
	GetReadDB() *sqlx.DB
	GetRepo() *storeRepository
	GetWriteDB() *sqlx.DB
	initialise() error
	IsLocked(err error) bool
	lock()
	needsMigration() bool
	Open() error
	open(disableForeignKeys bool, writable bool) (conn *sqlx.DB, err error)
	openReadDB() error
	openWriteDB() error
	Optimise(ctx context.Context) error
	QuerySQL(ctx context.Context, query string, args []interface{}) ([]string, [][]interface{}, error)
	Ready() error
	ReInitialise() error
	Remove() error
	Repository() models.Repository
	Reset() error
	RestoreFromBackup(backupPath string) error
	Rollback(ctx context.Context) error
	RunAllMigrations() error
	SetBlobStoreOptions(options BlobStoreOptions)
	SetSchemaVersion(version uint)
	txnComplete(ctx context.Context)
	unlock()
	Vacuum(ctx context.Context) error
	Version() uint
	WithDatabase(ctx context.Context) (context.Context, error)
	TestMode()
}

type Database struct {
	*storeRepository
	DBInterface

	readDB   *sqlx.DB
	writeDB  *sqlx.DB
	dbConfig string

	schemaVersion uint

	lockChan chan struct{}
}

func newDatabase() *storeRepository {
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

	return r
}

func getDBBoolean(val bool) string {
	switch dbWrapper.dbType {
	case SqliteBackend:
		if val {
			return "1"
		} else {
			return "0"
		}
	default:
		return strconv.FormatBool(val)
	}
}

func getDBValuesFix() (val string) {
	if dbWrapper.dbType == PostgresBackend {
		val = " AS v(column1, column2)"
	}

	return val
}

func getDBMinFunc() string {
	switch dbWrapper.dbType {
	case PostgresBackend:
		return "LEAST"
	default:
		return "MIN"
	}
}

func getDBLike() string {
	switch dbWrapper.dbType {
	case PostgresBackend:
		return "ILIKE"
	default:
		return "LIKE"
	}
}

func getDBGroupConcat(columnName string) string {
	switch dbWrapper.dbType {
	case PostgresBackend:
		return "STRING_AGG(" + columnName + "::TEXT, ',')"
	default:
		return "GROUP_CONCAT(" + columnName + ")"
	}
}

func getDBRowId() string {
	switch dbWrapper.dbType {
	case PostgresBackend:
		return "ctid"
	default:
		return "rowid"
	}
}

func isConstraintError(err error) bool {
	switch dbWrapper.dbType {
	case PostgresBackend:
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Class 23 — Integrity Constraint Violation
			return pgErr.Code[:2] == "23"
		}
	case SqliteBackend:
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			return sqliteError.Code == sqlite3.ErrConstraint
		}
	}
	return false
}

func (db *Database) SetSchemaVersion(version uint) {
	db.schemaVersion = version
}

func (db *Database) GetRepo() *storeRepository {
	return db.storeRepository
}

func (db *Database) SetBlobStoreOptions(options BlobStoreOptions) {
	*db.Blobs = *NewBlobStore(options)
}

func (db *Database) DatabasePath() string {
	return ""
}

// Ready returns an error if the database is not ready to begin transactions.
func (db *Database) Ready() error {
	if db.readDB == nil || db.writeDB == nil {
		return ErrDatabaseNotInitialized
	}

	return nil
}

// Open initializes the database. If the database is new, then it
// performs a full migration to the latest schema version. Otherwise, any
// necessary migrations must be run separately using RunMigrations.
// Returns true if the database is new.
func (db *Database) Open() error {
	db.lock()
	defer db.unlock()

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
		if databaseSchemaVersion > db.AppSchemaVersion() {
			return &MismatchedSchemaVersionError{
				CurrentSchemaVersion:  databaseSchemaVersion,
				RequiredSchemaVersion: db.AppSchemaVersion(),
			}
		}

		// if migration is needed, then don't open the connection
		if db.needsMigration() {
			return &MigrationNeededError{
				CurrentSchemaVersion:  databaseSchemaVersion,
				RequiredSchemaVersion: db.AppSchemaVersion(),
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

func (db *Database) initialise() error {
	if err := db.openReadDB(); err != nil {
		return fmt.Errorf("opening read database: %w", err)
	}
	if err := db.openWriteDB(); err != nil {
		return fmt.Errorf("opening write database: %w", err)
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

func (db *Database) Version() uint {
	return db.schemaVersion
}

func (db *Database) Reset() error {
	if err := db.Remove(); err != nil {
		return err
	}

	if err := db.Open(); err != nil {
		return fmt.Errorf("[reset DB] unable to initialize: %w", err)
	}

	return nil
}

func (db *Database) AnonymousDatabasePath(backupDirectoryPath string) string {
	fn := fmt.Sprintf("%s.anonymous.%d.%s", filepath.Base(db.DatabasePath()), db.schemaVersion, time.Now().Format("20060102_150405"))

	if backupDirectoryPath != "" {
		return filepath.Join(backupDirectoryPath, fn)
	}

	return fn
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

func (db *Database) GetWriteDB() *sqlx.DB {
	return db.writeDB
}

func (db *Database) GetReadDB() *sqlx.DB {
	return db.readDB
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
