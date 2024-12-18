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
	"github.com/stashapp/stash/pkg/logger"
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
	// SupplementaryPaths are alternative filesystem paths that will be used to find blobs
	// No changes will be made to these filesystems
	SupplementaryPaths []string
}

type BlobStore struct {
	repository

	tableMgr *table

	fsStore *blob.FilesystemStore
	// supplementary stores
	otherStores []blob.FilesystemReader
	options     BlobStoreOptions
}

func NewBlobStore(options BlobStoreOptions) *BlobStore {
	fs := &file.OsFS{}

	ret := &BlobStore{
		repository: repository{
			tableName: blobTable,
			idColumn:  blobChecksumColumn,
		},

		tableMgr: blobTableMgr,

		fsStore: blob.NewFilesystemStore(options.Path, fs),
		options: options,
	}

	for _, otherPath := range options.SupplementaryPaths {
		ret.otherStores = append(ret.otherStores, *blob.NewReadonlyFilesystemStore(otherPath, fs))
	}

	return ret
}

type blobRow struct {
	Checksum string `db:"checksum"`
	Blob     []byte `db:"blob"`
}

func (qb *BlobStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *BlobStore) Count(ctx context.Context) (int, error) {
	table := qb.table()
	q := dialect.From(table).Select(goqu.COUNT(table.Col(blobChecksumColumn)))

	var ret int
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
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

func (qb *BlobStore) update(ctx context.Context, checksum string, data []byte) error {
	table := qb.table()
	q := dialect.Update(table).Prepared(true).Set(goqu.Record{
		"blob": data,
	}).Where(goqu.C(blobChecksumColumn).Eq(checksum))

	_, err := exec(ctx, q)
	if err != nil {
		return fmt.Errorf("updating %s: %w", table, err)
	}

	return nil
}

type ChecksumNotFoundError struct {
	Checksum string
}

func (e *ChecksumNotFoundError) Error() string {
	return fmt.Sprintf("checksum %s does not exist", e.Checksum)
}

type ChecksumBlobNotExistError struct {
	Checksum string
}

func (e *ChecksumBlobNotExistError) Error() string {
	return fmt.Sprintf("blob for checksum %s does not exist", e.Checksum)
}

func (qb *BlobStore) readSQL(ctx context.Context, querySQL string, args ...interface{}) ([]byte, string, error) {
	if !qb.options.UseDatabase && !qb.options.UseFilesystem {
		panic("no blob store configured")
	}

	// always try to get from the database first, even if set to use filesystem
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

	checksum := row.Checksum

	if row.Blob != nil {
		return row.Blob, checksum, nil
	}

	// don't use the filesystem if not configured to do so
	if qb.options.UseFilesystem {
		ret, err := qb.readFromFilesystem(ctx, checksum)
		if err != nil {
			return nil, checksum, err
		}

		return ret, checksum, nil
	}

	return nil, checksum, &ChecksumBlobNotExistError{
		Checksum: checksum,
	}
}

func (qb *BlobStore) readFromFilesystem(ctx context.Context, checksum string) ([]byte, error) {
	// try to read from primary store first, then supplementaries
	fsStores := append([]blob.FilesystemReader{qb.fsStore.FilesystemReader}, qb.otherStores...)

	for _, fsStore := range fsStores {
		ret, err := fsStore.Read(ctx, checksum)
		if err == nil {
			return ret, nil
		}

		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("reading from filesystem: %w", err)
		}
	}

	// blob not found - should not happen
	return nil, &ChecksumBlobNotExistError{
		Checksum: checksum,
	}
}

func (qb *BlobStore) EntryExists(ctx context.Context, checksum string) (bool, error) {
	q := dialect.From(qb.table()).Select(goqu.COUNT("*")).Where(qb.tableMgr.byID(checksum))

	var found int
	if err := querySimple(ctx, q, &found); err != nil {
		return false, fmt.Errorf("querying %s: %w", qb.table(), err)
	}

	return found != 0, nil
}

// Read reads the data from the database or filesystem, depending on which is enabled.
func (qb *BlobStore) Read(ctx context.Context, checksum string) ([]byte, error) {
	if !qb.options.UseDatabase && !qb.options.UseFilesystem {
		panic("no blob store configured")
	}

	// always try to get from the database first, even if set to use filesystem
	ret, err := qb.readFromDatabase(ctx, checksum)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("reading from database: %w", err)
		}

		// not found in the database - does not exist
		return nil, &ChecksumNotFoundError{
			Checksum: checksum,
		}
	}

	if ret != nil {
		return ret, nil
	}

	// don't use the filesystem if not configured to do so
	if qb.options.UseFilesystem {
		return qb.readFromFilesystem(ctx, checksum)
	}

	// blob not found - should not happen
	return nil, &ChecksumBlobNotExistError{
		Checksum: checksum,
	}
}

func (qb *BlobStore) readFromDatabase(ctx context.Context, checksum string) ([]byte, error) {
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
			logger.Debugf("Blob %s is still referenced - not deleting", checksum)
			return nil
		}

		// unexpected error
		return fmt.Errorf("deleting from database: %w", err)
	}

	// blob was deleted from the database - delete from filesystem if enabled
	if qb.options.UseFilesystem {
		logger.Debugf("Deleting blob %s from filesystem", checksum)
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
	repository repository
	blobStore  *BlobStore

	joinTable string
}

func (qb *blobJoinQueryBuilder) GetImage(ctx context.Context, id int, blobCol string) ([]byte, error) {
	sqlQuery := utils.StrFormat(`
SELECT blobs.checksum, blobs.blob FROM {joinTable} INNER JOIN blobs ON {joinTable}.{joinCol} = blobs.checksum
WHERE {joinTable}.id = ?
`, utils.StrFormatMap{
		"joinTable": qb.joinTable,
		"joinCol":   blobCol,
	})

	ret, _, err := qb.blobStore.readSQL(ctx, sqlQuery, id)
	return ret, err
}

func (qb *blobJoinQueryBuilder) UpdateImage(ctx context.Context, id int, blobCol string, image []byte) error {
	if len(image) == 0 {
		return qb.DestroyImage(ctx, id, blobCol)
	}

	oldChecksum, err := qb.getChecksum(ctx, id, blobCol)
	if err != nil {
		return err
	}

	checksum, err := qb.blobStore.Write(ctx, image)
	if err != nil {
		return err
	}

	sqlQuery := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", qb.joinTable, blobCol)
	if _, err := dbWrapper.Exec(ctx, sqlQuery, checksum, id); err != nil {
		return err
	}

	// #3595 - delete the old blob if the checksum is different
	if oldChecksum != nil && *oldChecksum != checksum {
		if err := qb.blobStore.Delete(ctx, *oldChecksum); err != nil {
			return err
		}
	}

	return nil
}

func (qb *blobJoinQueryBuilder) getChecksum(ctx context.Context, id int, blobCol string) (*string, error) {
	sqlQuery := utils.StrFormat(`
SELECT {joinTable}.{joinCol} FROM {joinTable} WHERE {joinTable}.id = ?
`, utils.StrFormatMap{
		"joinTable": qb.joinTable,
		"joinCol":   blobCol,
	})

	var checksum null.String
	err := qb.repository.querySimple(ctx, sqlQuery, []interface{}{id}, &checksum)
	if err != nil {
		return nil, err
	}

	if !checksum.Valid {
		return nil, nil
	}

	return &checksum.String, nil
}

func (qb *blobJoinQueryBuilder) DestroyImage(ctx context.Context, id int, blobCol string) error {
	checksum, err := qb.getChecksum(ctx, id, blobCol)
	if err != nil {
		return err
	}

	if checksum == nil {
		// no image to delete
		return nil
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET %s = NULL WHERE id = ?", qb.joinTable, blobCol)
	if _, err = dbWrapper.Exec(ctx, updateQuery, id); err != nil {
		return err
	}

	return qb.blobStore.Delete(ctx, *checksum)
}

func (qb *blobJoinQueryBuilder) HasImage(ctx context.Context, id int, blobCol string) (bool, error) {
	stmt := utils.StrFormat("SELECT COUNT(*) as count FROM (SELECT {joinCol} FROM {joinTable} WHERE id = ? AND {joinCol} IS NOT NULL LIMIT 1)", utils.StrFormatMap{
		"joinTable": qb.joinTable,
		"joinCol":   blobCol,
	})

	c, err := qb.repository.runCountQuery(ctx, stmt, []interface{}{id})
	if err != nil {
		return false, err
	}

	return c == 1, nil
}
