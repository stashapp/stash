package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
)

const (
	// TODO: Test for optimality
	maxWriteConnections = 5
	maxReadConnections  = 15
	// Idle connection timeout, in seconds
	// Closes a connection after a period of inactivity, which saves on memory and
	// causes the sqlite -wal and -shm files to be automatically deleted.
	dbConnTimeout = 30 * time.Second
)

var appSchemaVersion uint = 12

//go:embed migrations/*.sql
var migrationsBox embed.FS

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

	readDB  *sqlx.DB
	writeDB *sqlx.DB
	dbPath  string

	schemaVersion uint

	lockChan chan struct{}
}

func NewDatabase() *Database {
	fileStore := NewFileStore()
	folderStore := NewFolderStore()
	galleryStore := NewGalleryStore(fileStore, folderStore)
	blobStore := NewBlobStore(database.BlobStoreOptions{})
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

func (db *Database) DatabaseBackend() database.DatabaseType {
	return database.PostgresBackend
}

func (db *Database) SetBlobStoreOptions(options database.BlobStoreOptions) {
	*db.storeRepository.Blobs = *NewBlobStore(options)
}

// Ready returns an error if the database is not ready to begin transactions.
func (db *Database) Ready() error {
	if db.readDB == nil || db.writeDB == nil {
		return database.ErrDatabaseNotInitialized
	}

	return nil
}

// Open initializes the database. If the database is new, then it
// performs a full migration to the latest schema version. Otherwise, any
// necessary migrations must be run separately using RunMigrations.
// Returns true if the database is new.
func (db *Database) Open(dbPath string) error {
	db.lock()
	defer db.unlock()

	db.dbPath, _ = strings.CutPrefix(dbPath, string(database.PostgresBackend)+":")

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
			return &database.MismatchedSchemaVersionError{
				CurrentSchemaVersion:  databaseSchemaVersion,
				RequiredSchemaVersion: appSchemaVersion,
			}
		}

		// if migration is needed, then don't open the connection
		if db.needsMigration() {
			return &database.MigrationNeededError{
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
	conn, err = sqlx.Open("pgx", db.dbPath)

	if err != nil {
		return nil, fmt.Errorf("db.Open(): %w", err)
	}

	if disableForeignKeys {
		_, err = conn.Exec("SET session_replication_role = replica;")

		if err != nil {
			return nil, fmt.Errorf("conn.Exec(): %w", err)
		}
	}
	if !writable {
		_, err = conn.Exec("SET SESSION CHARACTERISTICS AS TRANSACTION READ ONLY;")

		if err != nil {
			return nil, fmt.Errorf("conn.Exec(): %w", err)
		}
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

func (db *Database) Remove() (err error) {
	_, err = db.writeDB.Exec(`
DO $$
DECLARE 
    r record;
BEGIN
    FOR r IN SELECT quote_ident(tablename) AS tablename, quote_ident(schemaname) AS schemaname FROM pg_tables WHERE schemaname = 'public'
    LOOP
        RAISE INFO 'Dropping table %.%', r.schemaname, r.tablename;
        EXECUTE format('DROP TABLE IF EXISTS %I.%I CASCADE', r.schemaname, r.tablename);
    END LOOP;
END$$;
`)

	return err
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
	logger.Warn("Postgres backend detected, ignoring Backup request")
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
	logger.Warn("Postgres backend detected, ignoring RestoreFromBackup request")
	return nil
}

func (db *Database) AppSchemaVersion() uint {
	return appSchemaVersion
}

func (db *Database) DatabasePath() string {
	return db.dbPath
}

func (db *Database) DatabaseBackupPath(backupDirectoryPath string) string {
	logger.Warn("Postgres backend detected, ignoring DatabaseBackupPath request")
	return ""
}

func (db *Database) AnonymousDatabasePath(backupDirectoryPath string) string {
	fn := fmt.Sprintf("%s.anonymous.%d.%s", "postgres", db.schemaVersion, time.Now().Format("20060102_150405"))

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
	_, err := db.writeDB.ExecContext(ctx, "VACUUM (FULL, ANALYZE, VERBOSE)")
	return err
}

// Analyze runs an ANALYZE on the database to improve query performance.
func (db *Database) Analyze(ctx context.Context) error {
	return analyze(ctx, db.writeDB)
}

// analyze runs an ANALYZE on the database to improve query performance.
func analyze(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, "ANALYZE")
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

	return rowsAffected, nil, nil
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

func getBasenameSQL(sql string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("basename(%s, '\\')", sql)
	}

	return fmt.Sprintf("basename(%s)", sql)
}
