package database

import (
	"context"
)

type BlobStoreOptions struct {
	// UseFilesystem should be true if blob data should be stored in the filesystem
	UseFilesystem bool
	// UseDatabase should be true if blob data should be stored in the database
	UseDatabase bool
	// Path is the filesystem path to use for storing blobs
	Path string
	// SupplementaryPaths are alternative filesystem paths that will be used to find blobs
	// No changes will be made to these filesystems
	SupplementaryPaths []string
}

type BlobStore interface {
	Count(ctx context.Context) (int, error)
	Delete(ctx context.Context, checksum string) error
	EntryExists(ctx context.Context, checksum string) (bool, error)
	FindBlobs(ctx context.Context, n uint, lastChecksum string) ([]string, error)
	MigrateBlob(ctx context.Context, checksum string, deleteOld bool) error
	Read(ctx context.Context, checksum string) ([]byte, error)
	Write(ctx context.Context, data []byte) (string, error)
}
