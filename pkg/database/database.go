package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type DatabaseType string

const (
	PostgresBackend DatabaseType = "POSTGRESQL"
	SqliteBackend   DatabaseType = "SQLITE"
)

type Database interface {
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
	DatabaseBackend() DatabaseType
	ExecSQL(ctx context.Context, query string, args []interface{}) (*int64, *int64, error)
	IsLocked(err error) bool
	Open(dbPath string) error
	Optimise(ctx context.Context) error
	QuerySQL(ctx context.Context, query string, args []interface{}) ([]string, [][]interface{}, error)
	ReInitialise() error
	Ready() error
	Remove() error
	Repository() models.Repository
	Reset() error
	RestoreFromBackup(backupPath string) error
	Rollback(ctx context.Context) error
	RunAllMigrations() error
	SetBlobStoreOptions(options BlobStoreOptions)
	Vacuum(ctx context.Context) error
	Version() uint
	WithDatabase(ctx context.Context) (context.Context, error)
	NewMigrator() (MigrateStore, error)

	Blobs() BlobStore
	File() FileStore
	Folder() FolderStore
	Image() ImageStore
	Gallery() GalleryStore
	GalleryChapter() GalleryChapterStore
	Scene() SceneStore
	SceneMarker() SceneMarkerStore
	Performer() PerformerStore
	SavedFilter() SavedFilterStore
	Studio() StudioStore
	Tag() TagStore
	Group() GroupStore
}

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
