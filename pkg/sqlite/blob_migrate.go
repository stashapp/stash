package sqlite

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/jmoiron/sqlx"
)

func (qb *BlobStore) FindBlobs(ctx context.Context, n uint, lastChecksum string) ([]string, error) {
	table := qb.table()
	q := dialect.From(table).Select(table.Col(blobChecksumColumn)).Order(table.Col(blobChecksumColumn).Asc()).Limit(n)

	if lastChecksum != "" {
		q = q.Where(table.Col(blobChecksumColumn).Gt(lastChecksum))
	}

	const single = false
	var checksums []string
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var checksum string
		if err := rows.Scan(&checksum); err != nil {
			return err
		}
		checksums = append(checksums, checksum)
		return nil
	}); err != nil {
		return nil, err
	}

	return checksums, nil
}

// MigrateBlob migrates a blob to the currently active store from any of the
// other stores. The target is determined by the UseDatabase, UseFilesystem and
// UseS3 options. If deleteOld is true, the blob is deleted from the other
// stores after migration.
func (qb *BlobStore) MigrateBlob(ctx context.Context, checksum string, deleteOld bool) error {
	targets := 0
	for _, enabled := range []bool{qb.options.UseDatabase, qb.options.UseFilesystem, qb.options.UseS3} {
		if enabled {
			targets++
		}
	}

	if targets == 0 {
		panic("no blob store configured")
	}
	if targets > 1 {
		panic("multiple blob stores configured")
	}

	switch {
	case qb.options.UseDatabase:
		return qb.migrateBlobDatabase(ctx, checksum, deleteOld)
	case qb.options.UseFilesystem:
		if qb.options.Path == "" {
			panic("no blob path configured")
		}
		return qb.migrateBlobFilesystem(ctx, checksum, deleteOld)
	default:
		return qb.migrateBlobS3(ctx, checksum, deleteOld)
	}
}

// readFromNonDatabaseStores reads a blob from the filesystem or s3 stores,
// regardless of which store is active. Returns an error wrapping
// ChecksumBlobNotExistError if the blob is in none of them.
func (qb *BlobStore) readFromNonDatabaseStores(ctx context.Context, checksum string) ([]byte, error) {
	fsData, err := qb.readFromFilesystem(ctx, checksum)
	if err == nil {
		return fsData, nil
	}

	var notExist *ChecksumBlobNotExistError
	if !errors.As(err, &notExist) {
		return nil, err
	}

	if qb.s3Store != nil {
		s3Data, s3err := qb.s3Store.Read(ctx, checksum)
		if s3err == nil {
			return s3Data, nil
		}
		if !errors.Is(s3err, fs.ErrNotExist) {
			return nil, fmt.Errorf("reading from s3: %w", s3err)
		}
	}

	return nil, err
}

// deleteFromNonTargetStores removes the blob data from all stores other than
// the active one. Missing blobs are ignored.
func (qb *BlobStore) deleteFromNonTargetStores(ctx context.Context, checksum string) error {
	if !qb.options.UseDatabase {
		// clear the blob column, keep the checksum row
		if err := qb.update(ctx, checksum, nil); err != nil {
			return fmt.Errorf("clearing database blob: %w", err)
		}
	}

	if !qb.options.UseFilesystem && qb.options.Path != "" {
		if err := qb.fsStore.Delete(ctx, checksum); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("deleting from filesystem: %w", err)
		}
	}

	if !qb.options.UseS3 && qb.s3Store != nil {
		if err := qb.s3Store.Delete(ctx, checksum); err != nil {
			return fmt.Errorf("deleting from s3: %w", err)
		}
	}

	return nil
}

// migrateBlobDatabase migrates a blob from the filesystem or s3 to the database
func (qb *BlobStore) migrateBlobDatabase(ctx context.Context, checksum string, deleteOld bool) error {
	// ignore if the blob is already present in the database
	// (still delete the old data if requested)
	existing, err := qb.readFromDatabase(ctx, checksum)
	if err != nil {
		return fmt.Errorf("reading from database: %w", err)
	}

	if len(existing) == 0 {
		// find the blob in the other stores
		blob, err := qb.readFromNonDatabaseStores(ctx, checksum)
		if err != nil {
			return err
		}

		// write the blob to the database
		if err := qb.update(ctx, checksum, blob); err != nil {
			return fmt.Errorf("writing to database: %w", err)
		}
	}

	if deleteOld {
		return qb.deleteFromNonTargetStores(ctx, checksum)
	}

	return nil
}

// migrateBlobFilesystem migrates a blob from the database or s3 to the filesystem
func (qb *BlobStore) migrateBlobFilesystem(ctx context.Context, checksum string, deleteOld bool) error {
	// find the blob in the database
	blobData, err := qb.readFromDatabase(ctx, checksum)
	if err != nil {
		return fmt.Errorf("reading from database: %w", err)
	}

	if len(blobData) == 0 && qb.s3Store != nil {
		s3Data, s3err := qb.s3Store.Read(ctx, checksum)
		if s3err == nil {
			blobData = s3Data
		} else if !errors.Is(s3err, fs.ErrNotExist) {
			return fmt.Errorf("reading from s3: %w", s3err)
		}
	}

	if len(blobData) == 0 {
		// it's possible that the blob is already present in the filesystem
		// just ignore
		return nil
	}

	// write the blob to the filesystem
	if err := qb.fsStore.Write(ctx, checksum, blobData); err != nil {
		return fmt.Errorf("writing to filesystem: %w", err)
	}

	if deleteOld {
		return qb.deleteFromNonTargetStores(ctx, checksum)
	}

	return nil
}

// migrateBlobS3 migrates a blob from the database or filesystem to s3
func (qb *BlobStore) migrateBlobS3(ctx context.Context, checksum string, deleteOld bool) error {
	// find the blob in the database
	blobData, err := qb.readFromDatabase(ctx, checksum)
	if err != nil {
		return fmt.Errorf("reading from database: %w", err)
	}

	if len(blobData) == 0 {
		fsData, fsErr := qb.readFromFilesystem(ctx, checksum)
		if fsErr == nil {
			blobData = fsData
		} else {
			var notExist *ChecksumBlobNotExistError
			if !errors.As(fsErr, &notExist) {
				return fsErr
			}
		}
	}

	if len(blobData) == 0 {
		// it's possible that the blob is already present in s3
		// just ignore
		return nil
	}

	// write the blob to s3
	if err := qb.s3Store.Write(ctx, checksum, blobData); err != nil {
		return fmt.Errorf("writing to s3: %w", err)
	}

	if deleteOld {
		return qb.deleteFromNonTargetStores(ctx, checksum)
	}

	return nil
}
