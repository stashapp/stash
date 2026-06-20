package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

type Store struct {
	DB     *sqlite.Database
	Repo   models.Repository
	Schema uint
}

type BlobStorage string

const (
	BlobStorageDatabase   BlobStorage = "database"
	BlobStorageFilesystem BlobStorage = "filesystem"
)

type Options struct {
	BlobStorage BlobStorage
	BlobsPath   string
}

func Open(path string, options Options) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}

	info, statErr := os.Stat(path)
	if statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return nil, fmt.Errorf("database file does not exist; run Stash server setup first and point database_path at its sqlite file: %w", statErr)
		}
		return nil, fmt.Errorf("read database file: %w", statErr)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("database path is a directory: %s", path)
	}

	db := sqlite.NewDatabase()
	blobOptions, err := sqliteBlobOptions(options)
	if err != nil {
		return nil, err
	}
	db.SetBlobStoreOptions(blobOptions)
	if err := db.Open(path); err != nil {
		return nil, classifyOpenError(err)
	}

	return &Store{
		DB:     db,
		Repo:   db.Repository(),
		Schema: db.AppSchemaVersion(),
	}, nil
}

func sqliteBlobOptions(options Options) (sqlite.BlobStoreOptions, error) {
	switch options.BlobStorage {
	case "", BlobStorageDatabase:
		return sqlite.BlobStoreOptions{UseDatabase: true}, nil
	case BlobStorageFilesystem:
		if options.BlobsPath == "" {
			return sqlite.BlobStoreOptions{}, errors.New("blobs_path is required when blobs_storage=filesystem")
		}
		return sqlite.BlobStoreOptions{UseFilesystem: true, Path: options.BlobsPath}, nil
	default:
		return sqlite.BlobStoreOptions{}, fmt.Errorf("blobs_storage must be one of database, filesystem")
	}
}

func (s *Store) Close() error {
	if s == nil || s.DB == nil {
		return nil
	}
	return s.DB.Close()
}

func classifyOpenError(err error) error {
	var migrationNeeded *sqlite.MigrationNeededError
	if errors.As(err, &migrationNeeded) {
		return fmt.Errorf("database migration required: current schema=%d, required schema=%d; use Stash server or the migration flow to upgrade first: %w", migrationNeeded.CurrentSchemaVersion, migrationNeeded.RequiredSchemaVersion, err)
	}

	var mismatched *sqlite.MismatchedSchemaVersionError
	if errors.As(err, &mismatched) {
		return fmt.Errorf("database schema is too new: current schema=%d, supported schema=%d: %w", mismatched.CurrentSchemaVersion, mismatched.RequiredSchemaVersion, err)
	}

	return fmt.Errorf("open database failed: %w", err)
}
