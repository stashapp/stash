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

	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

const (
	// Number of database connections to use
	// The same value is used for both the maximum and idle limit,
	// to prevent opening connections on the fly which has a notieable performance penalty.
	// Fewer connections use less memory, more connections increase performance,
	// but have diminishing returns.
	// 10 was found to be a good tradeoff.
	dbConns = 10
	// Idle connection timeout, in seconds
	// Closes a connection after a period of inactivity, which saves on memory and
	// causes the sqlite -wal and -shm files to be automatically deleted.
	dbConnTimeout = 30
)

var appSchemaVersion uint = 64

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

type Database struct {
	*storeRepository

	db     *sqlx.DB
	dbPath string

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
	db.lockNoCtx()
	defer db.unlock()

	db.dbPath = dbPath

	databaseSchemaVersion, err := db.getDatabaseSchemaVersion()
	if err != nil {
		return fmt.Errorf("getting database schema version: %w", err)
	}

	db.schemaVersion = databaseSchemaVersion

	if databaseSchemaVersion == 0 {
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

	// RunMigrations may have opened a connection already
	if db.db == nil {
		const disableForeignKeys = false
		db.db, err = db.open(disableForeignKeys)
		if err != nil {
			return err
		}
	}

	return nil
}

// lock locks the database for writing.
// This method will block until the lock is acquired of the context is cancelled.
func (db *Database) lock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case db.lockChan <- struct{}{}:
		return nil
	}
}

// lock locks the database for writing. This method will block until the lock is acquired.
func (db *Database) lockNoCtx() {
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
	db.lockNoCtx()
	defer db.unlock()

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
	url := "file:" + db.dbPath + "?_journal=WAL&_sync=NORMAL&_busy_timeout=50"
	if !disableForeignKeys {
		url += "&_fk=true"
	}

	conn, err := sqlx.Open(sqlite3Driver, url)
	conn.SetMaxOpenConns(dbConns)
	conn.SetMaxIdleConns(dbConns)
	conn.SetConnMaxIdleTime(dbConnTimeout * time.Second)
	if err != nil {
		return nil, fmt.Errorf("db.Open(): %w", err)
	}

	return conn, nil
}

func (db *Database) Remove() error {
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
	databasePath := db.dbPath
	if err := db.Remove(); err != nil {
		return err
	}

	if err := db.Open(databasePath); err != nil {
		return fmt.Errorf("[reset DB] unable to initialize: %w", err)
	}

	return nil
}

// Backup the database. If db is nil, then uses the existing database
// connection.
func (db *Database) Backup(backupPath string) (err error) {
	thisDB := db.db
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
	_, err := db.db.ExecContext(ctx, "VACUUM")
	return err
}

// Analyze runs an ANALYZE on the database to improve query performance.
func (db *Database) Analyze(ctx context.Context) error {
	_, err := db.db.ExecContext(ctx, "ANALYZE")
	return err
}

func (db *Database) ExecSQL(ctx context.Context, query string, args []interface{}) (*int64, *int64, error) {
	wrapper := dbWrapperType{}

	result, err := wrapper.Exec(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}

	var rowsAffected *int64
	ra, err := result.RowsAffected()
	if err == nil {
		rowsAffected = &ra
	}

	var lastInsertId *int64
	li, err := result.LastInsertId()
	if err == nil {
		lastInsertId = &li
	}

	return rowsAffected, lastInsertId, nil
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
