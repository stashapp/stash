package blob

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/hash/md5"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrBlobReferenced = errors.New("blob is referenced by another object")
)

type Reader interface {
	Read(ctx context.Context, checksum string) ([]byte, error)
}

type Store interface {
	Reader
	Write(ctx context.Context, checksum string, data []byte) error
	Delete(ctx context.Context, checksum string) error
}

type ServiceOptions struct {
	// UseFilesystem should be true if blob data should be stored in the filesystem
	UseFilesystem bool
	// UseDatabase should be true if blob data should be stored in the database
	UseDatabase bool

	Path string

	FS       FS
	Database Store
}

type Writer interface {
	Write(ctx context.Context, data []byte) (string, error)
}

type Service interface {
	Writer
	Reader
	Delete(ctx context.Context, checksum string) error
}

type service struct {
	options ServiceOptions
	fsStore Store
	dbStore Store
}

// NewService
func NewService(options ServiceOptions) Service {
	return &service{
		options: options,
		fsStore: NewFilesystemStore(options.Path, options.FS),
		dbStore: options.Database,
	}
}

// Write stores the data and its checksum in enabled stores.
// Always writes at least the checksum to the database.
func (s *service) Write(ctx context.Context, data []byte) (string, error) {
	if !s.options.UseDatabase && !s.options.UseFilesystem {
		panic("no blob store configured")
	}

	if len(data) == 0 {
		return "", fmt.Errorf("cannot write empty data")
	}

	checksum := md5.FromBytes(data)

	// only write blob to the database if UseDatabase is true
	// always at least write the checksum
	var storedData []byte
	if s.options.UseDatabase {
		storedData = data
	}

	if err := s.options.Database.Write(ctx, checksum, storedData); err != nil {
		return "", fmt.Errorf("writing to database: %w", err)
	}

	if s.options.UseFilesystem {
		if err := s.fsStore.Write(ctx, checksum, data); err != nil {
			return "", fmt.Errorf("writing to filesystem: %w", err)
		}
	}

	return checksum, nil
}

// Read reads the data from the database or filesystem, depending on which is enabled.
func (s *service) Read(ctx context.Context, checksum string) ([]byte, error) {
	if !s.options.UseDatabase && !s.options.UseFilesystem {
		panic("no blob store configured")
	}

	// check the database first
	ret, err := s.options.Database.Read(ctx, checksum)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("reading from database: %w", err)
		}

		// not found in the database - does not exist
		return nil, ErrNotFound
	}

	if s.options.UseDatabase && ret != nil {
		return ret, nil
	}

	if s.options.UseFilesystem {
		return s.fsStore.Read(ctx, checksum)
	}

	return nil, fmt.Errorf("unexpected nil blob")
}

// Delete marks a checksum as no longer in use by a single reference.
// If no references remain, the blob is deleted from the database and filesystem.
func (s *service) Delete(ctx context.Context, checksum string) error {
	// try to delete the blob from the database
	if err := s.options.Database.Delete(ctx, checksum); err != nil {
		if errors.Is(err, ErrBlobReferenced) {
			// blob is still referenced - do not delete
			return nil
		}

		// unexpected error
		return fmt.Errorf("deleting from database: %w", err)
	}

	// blob was deleted from the database - delete from filesystem if enabled
	if s.options.UseFilesystem {
		if err := s.fsStore.Delete(ctx, checksum); err != nil {
			return fmt.Errorf("deleting from filesystem: %w", err)
		}
	}

	return nil
}
