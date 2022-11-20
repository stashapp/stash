package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/sqlite/blob"
	"github.com/stashapp/stash/pkg/utils"
	"gopkg.in/guregu/null.v4"
)

const (
	blobTable          = "blobs"
	blobChecksumColumn = "checksum"
)

type BlobStoreOptions struct {
	// UseFilesystem should be true if blob data should be stored in the filesystem
	UseFilesystem bool
	// UseDatabase should be true if blob data should be stored in the database
	UseDatabase bool
	// Path is the filesystem path to use for storing blobs
	Path string
}

type BlobStore struct {
	repository

	tableMgr *table

	fsStore *blob.FilesystemStore
	options BlobStoreOptions
}

func NewBlobStore(options BlobStoreOptions) *BlobStore {
	return &BlobStore{
		repository: repository{
			tableName: blobTable,
			idColumn:  blobChecksumColumn,
		},

		tableMgr: blobTableMgr,

		fsStore: blob.NewFilesystemStore(options.Path, &file.OsFS{}),
		options: options,
	}
}

type blobRow struct {
	Checksum string `db:"checksum"`
	Blob     []byte `db:"blob"`
}

func (qb *BlobStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

// Write stores the data and its checksum in enabled stores.
// Always writes at least the checksum to the database.
func (qb *BlobStore) Write(ctx context.Context, data []byte) (string, error) {
	if !qb.options.UseDatabase && !qb.options.UseFilesystem {
		panic("no blob store configured")
	}

	if len(data) == 0 {
		return "", fmt.Errorf("cannot write empty data")
	}

	checksum := md5.FromBytes(data)

	// only write blob to the database if UseDatabase is true
	// always at least write the checksum
	var storedData []byte
	if qb.options.UseDatabase {
		storedData = data
	}

	if err := qb.write(ctx, checksum, storedData); err != nil {
		return "", fmt.Errorf("writing to database: %w", err)
	}

	if qb.options.UseFilesystem {
		if err := qb.fsStore.Write(ctx, checksum, data); err != nil {
			return "", fmt.Errorf("writing to filesystem: %w", err)
		}
	}

	return checksum, nil
}

func (qb *BlobStore) write(ctx context.Context, checksum string, data []byte) error {
	table := qb.table()
	q := dialect.Insert(table).Prepared(true).Rows(blobRow{
		Checksum: checksum,
		Blob:     data,
	}).OnConflict(goqu.DoNothing())

	_, err := exec(ctx, q)
	if err != nil {
		return fmt.Errorf("inserting into %s: %w", table, err)
	}

	return nil
}

type ChecksumNotFoundError struct {
	Checksum string
}

func (e *ChecksumNotFoundError) Error() string {
	return fmt.Sprintf("checksum %s does not exist", e.Checksum)
}

type ChecksumFileNotExistError struct {
	Checksum string
}

func (e *ChecksumFileNotExistError) Error() string {
	return fmt.Sprintf("checksum %s does not exist in filesystem", e.Checksum)
}

func (qb *BlobStore) readSQL(ctx context.Context, querySQL string, args ...interface{}) ([]byte, string, error) {
	if !qb.options.UseDatabase && !qb.options.UseFilesystem {
		panic("no blob store configured")
	}

	var row blobRow
	found := false
	const single = true
	if err := qb.queryFunc(ctx, querySQL, args, single, func(r *sqlx.Rows) error {
		found = true
		if err := r.StructScan(&row); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, "", fmt.Errorf("reading from database: %w", err)
	}

	if !found {
		// not found in the database - does not exist
		return nil, "", nil
	}

	if qb.options.UseDatabase && row.Blob != nil {
		return row.Blob, row.Checksum, nil
	}

	checksum := row.Checksum

	if qb.options.UseFilesystem {
		ret, err := qb.fsStore.Read(ctx, checksum)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, checksum, fmt.Errorf("reading from filesystem: %w", err)
			}

			// not found in the filesystem - should not happen
			return nil, checksum, &ChecksumFileNotExistError{
				Checksum: checksum,
			}
		}

		return ret, checksum, nil
	}

	return nil, checksum, fmt.Errorf("unexpected nil blob")
}

// Read reads the data from the database or filesystem, depending on which is enabled.
func (qb *BlobStore) Read(ctx context.Context, checksum string) ([]byte, error) {
	if !qb.options.UseDatabase && !qb.options.UseFilesystem {
		panic("no blob store configured")
	}

	// check the database first
	ret, err := qb.read(ctx, checksum)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("reading from database: %w", err)
		}

		// not found in the database - does not exist
		return nil, &ChecksumNotFoundError{
			Checksum: checksum,
		}
	}

	if qb.options.UseDatabase && ret != nil {
		return ret, nil
	}

	if qb.options.UseFilesystem {
		ret, err := qb.fsStore.Read(ctx, checksum)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("reading from filesystem: %w", err)
			}

			// not found in the filesystem - should not happen
			return nil, &ChecksumFileNotExistError{
				Checksum: checksum,
			}
		}

		return ret, nil
	}

	return nil, fmt.Errorf("unexpected nil blob")
}

func (qb *BlobStore) read(ctx context.Context, checksum string) ([]byte, error) {
	q := dialect.From(qb.table()).Select(qb.table().All()).Where(qb.tableMgr.byID(checksum))

	var row blobRow
	const single = true
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		if err := r.StructScan(&row); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("querying %s: %w", qb.table(), err)
	}

	return row.Blob, nil
}

// Delete marks a checksum as no longer in use by a single reference.
// If no references remain, the blob is deleted from the database and filesystem.
func (qb *BlobStore) Delete(ctx context.Context, checksum string) error {
	// try to delete the blob from the database
	if err := qb.delete(ctx, checksum); err != nil {
		if qb.isConstraintError(err) {
			// blob is still referenced - do not delete
			return nil
		}

		// unexpected error
		return fmt.Errorf("deleting from database: %w", err)
	}

	// blob was deleted from the database - delete from filesystem if enabled
	if qb.options.UseFilesystem {
		if err := qb.fsStore.Delete(ctx, checksum); err != nil {
			return fmt.Errorf("deleting from filesystem: %w", err)
		}
	}

	return nil
}

func (qb *BlobStore) isConstraintError(err error) bool {
	var sqliteError sqlite3.Error
	if errors.As(err, &sqliteError) {
		return sqliteError.Code == sqlite3.ErrConstraint
	}
	return false
}

func (qb *BlobStore) delete(ctx context.Context, checksum string) error {
	table := qb.table()

	q := dialect.Delete(table).Where(goqu.C(blobChecksumColumn).Eq(checksum))

	_, err := exec(ctx, q)
	if err != nil {
		return fmt.Errorf("deleting from %s: %w", table, err)
	}

	return nil
}

type blobJoinQueryBuilder struct {
	repository
	blobStore *BlobStore

	joinTable string
	joinCol   string
}

func (qb *blobJoinQueryBuilder) GetImage(ctx context.Context, id int) ([]byte, error) {
	sqlQuery := utils.StrFormat(`
SELECT blobs.checksum, blobs.blob FROM {joinTable} INNER JOIN blobs ON {joinTable}.{joinCol} = blobs.checksum
WHERE {joinTable}.id = ?
`, utils.StrFormatMap{
		"joinTable": qb.joinTable,
		"joinCol":   qb.joinCol,
	})

	ret, _, err := qb.blobStore.readSQL(ctx, sqlQuery, id)
	return ret, err
}

func (qb *blobJoinQueryBuilder) UpdateImage(ctx context.Context, id int, image []byte) error {
	if len(image) == 0 {
		return qb.DestroyImage(ctx, id)
	}
	checksum, err := qb.blobStore.Write(ctx, image)
	if err != nil {
		return err
	}

	sqlQuery := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", qb.joinTable, qb.joinCol)
	_, err = qb.tx.Exec(ctx, sqlQuery, checksum, id)
	return err
}

func (qb *blobJoinQueryBuilder) DestroyImage(ctx context.Context, id int) error {
	sqlQuery := utils.StrFormat(`
SELECT {joinTable}.{joinCol} FROM {joinTable} WHERE {joinTable}.id = ?
`, utils.StrFormatMap{
		"joinTable": qb.joinTable,
		"joinCol":   qb.joinCol,
	})

	var checksum null.String
	err := qb.repository.querySimple(ctx, sqlQuery, []interface{}{id}, &checksum)
	if err != nil {
		return err
	}

	if !checksum.Valid {
		// no image to delete
		return nil
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET %s = NULL WHERE id = ?", qb.joinTable, qb.joinCol)
	if _, err = qb.tx.Exec(ctx, updateQuery, id); err != nil {
		return err
	}

	return qb.blobStore.Delete(ctx, checksum.String)
}

func (qb *blobJoinQueryBuilder) HasImage(ctx context.Context, id int) (bool, error) {
	stmt := utils.StrFormat("SELECT COUNT(*) as count FROM (SELECT {joinCol} FROM {joinTable} WHERE id = ? LIMIT 1)", utils.StrFormatMap{
		"joinTable": qb.joinTable,
		"joinCol":   qb.joinCol,
	})

	c, err := qb.runCountQuery(ctx, stmt, []interface{}{id})
	if err != nil {
		return false, err
	}

	return c == 1, nil
}
