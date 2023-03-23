package sqlite

import (
	"context"
	"fmt"

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

// MigrateBlob migrates a blob from the filesystem to the database, or vice versa.
// The target is determined by the UseDatabase and UseFilesystem options.
// If deleteOld is true, the blob is deleted from the source after migration.
func (qb *BlobStore) MigrateBlob(ctx context.Context, checksum string, deleteOld bool) error {
	if !qb.options.UseDatabase && !qb.options.UseFilesystem {
		panic("no blob store configured")
	}

	if qb.options.UseDatabase && qb.options.UseFilesystem {
		panic("both filesystem and database configured")
	}

	if qb.options.Path == "" {
		panic("no blob path configured")
	}

	if qb.options.UseDatabase {
		return qb.migrateBlobDatabase(ctx, checksum, deleteOld)
	}

	return qb.migrateBlobFilesystem(ctx, checksum, deleteOld)
}

// migrateBlobDatabase migrates a blob from the filesystem to the database
func (qb *BlobStore) migrateBlobDatabase(ctx context.Context, checksum string, deleteOld bool) error {
	// ignore if the blob is already present in the database
	// (still delete the old data if requested)
	existing, err := qb.readFromDatabase(ctx, checksum)
	if err != nil {
		return fmt.Errorf("reading from database: %w", err)
	}

	if len(existing) == 0 {
		// find the blob in the filesystem
		blob, err := qb.fsStore.Read(ctx, checksum)
		if err != nil {
			return fmt.Errorf("reading from filesystem: %w", err)
		}

		// write the blob to the database
		if err := qb.update(ctx, checksum, blob); err != nil {
			return fmt.Errorf("writing to database: %w", err)
		}
	}

	if deleteOld {
		// delete the blob from the filesystem after commit
		if err := qb.fsStore.Delete(ctx, checksum); err != nil {
			return fmt.Errorf("deleting from filesystem: %w", err)
		}
	}

	return nil
}

// migrateBlobFilesystem migrates a blob from the database to the filesystem
func (qb *BlobStore) migrateBlobFilesystem(ctx context.Context, checksum string, deleteOld bool) error {
	// find the blob in the database
	blob, err := qb.readFromDatabase(ctx, checksum)
	if err != nil {
		return fmt.Errorf("reading from database: %w", err)
	}

	if len(blob) == 0 {
		// it's possible that the blob is already present in the filesystem
		// just ignore
		return nil
	}

	// write the blob to the filesystem
	if err := qb.fsStore.Write(ctx, checksum, blob); err != nil {
		return fmt.Errorf("writing to filesystem: %w", err)
	}

	if deleteOld {
		// delete the blob from the database row
		if err := qb.update(ctx, checksum, nil); err != nil {
			return err
		}
	}

	return nil
}
