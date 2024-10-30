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
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite/blob"
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
	Checksum string           `db:"checksum"`
	Blob     sql.Null[[]byte] `db:"blob"`
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
	var storedData sql.Null[[]byte]
	if qb.options.UseDatabase {
		storedData.V = data
		storedData.Valid = len(storedData.V) > 0
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

func (qb *BlobStore) write(ctx context.Context, checksum string, data sql.Null[[]byte]) error {
	table := qb.table()
	q := dialect.Insert(table).Rows(blobRow{
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
	q := dialect.Update(table).Set(goqu.Record{
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

func (qb *BlobStore) readSQL(ctx context.Context, querySQL sqler) ([]byte, string, error) {
	if !qb.options.UseDatabase && !qb.options.UseFilesystem {
		panic("no blob store configured")
	}

	query, args, err := querySQL.ToSQL()
	if err != nil {
		return nil, "", fmt.Errorf("reading blob tosql: %w", err)
	}

	// always try to get from the database first, even if set to use filesystem
	var row blobRow
	found := false
	const single = true
	if err := qb.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
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

	if row.Blob.Valid {
		return row.Blob.V, checksum, nil
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

	if ret.Valid {
		return ret.V, nil
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

func (qb *BlobStore) readFromDatabase(ctx context.Context, checksum string) (sql.Null[[]byte], error) {
	q := dialect.From(qb.table()).Select(qb.table().All()).Where(qb.tableMgr.byID(checksum))

	var empty sql.Null[[]byte]
	var row blobRow
	const single = true
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		if err := r.StructScan(&row); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return empty, fmt.Errorf("querying %s: %w", qb.table(), err)
	}

	return row.Blob, nil
}

// Delete marks a checksum as no longer in use by a single reference.
// If no references remain, the blob is deleted from the database and filesystem.
func (qb *BlobStore) Delete(ctx context.Context, checksum string) error {
	// try to delete the blob from the database
	if err := qb.delete(ctx, checksum); err != nil {
		if isConstraintError(err) {
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

func (qb *BlobStore) delete(ctx context.Context, checksum string) error {
	table := qb.table()

	q := dialect.Delete(table).Where(goqu.C(blobChecksumColumn).Eq(checksum))

	err := withSavepoint(ctx, func(ctx context.Context) error {
		_, err := exec(ctx, q)
		return err
	})

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
	sqlQuery := dialect.From(qb.joinTable).
		Join(goqu.I("blobs"), goqu.On(goqu.I(qb.joinTable+"."+blobCol).Eq(goqu.I("blobs.checksum")))).
		Select(goqu.I("blobs.checksum"), goqu.I("blobs.blob")).
		Where(goqu.Ex{"id": id})

	ret, _, err := qb.blobStore.readSQL(ctx, sqlQuery)
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

	sqlQuery := dialect.From(qb.joinTable).Update().
		Set(goqu.Record{blobCol: checksum}).
		Prepared(true).
		Where(goqu.Ex{"id": id})

	query, args, err := sqlQuery.ToSQL()
	if err != nil {
		return err
	}

	if _, err := dbWrapper.Exec(ctx, query, args...); err != nil {
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
	sqlQuery := dialect.From(qb.joinTable).
		Select(blobCol).
		Where(goqu.Ex{"id": id})

	query, args, err := sqlQuery.ToSQL()
	if err != nil {
		return nil, err
	}

	var checksum null.String
	err = qb.repository.querySimple(ctx, query, args, &checksum)
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

	updateQuery := dialect.Update(qb.joinTable).
		Set(goqu.Record{blobCol: nil}).
		Where(goqu.Ex{"id": id})

	query, args, err := updateQuery.ToSQL()
	if err != nil {
		return err
	}

	if _, err = dbWrapper.Exec(ctx, query, args...); err != nil {
		return err
	}

	return qb.blobStore.Delete(ctx, *checksum)
}

func (qb *blobJoinQueryBuilder) HasImage(ctx context.Context, id int, blobCol string) (bool, error) {
	ds := dialect.From(goqu.T(qb.joinTable)).
		Select(goqu.C(blobCol)).
		Where(
			goqu.C("id").Eq(id),
			goqu.C(blobCol).IsNotNull(),
		).
		Limit(1)

	countDs := dialect.From(ds.As("subquery")).Select(goqu.COUNT("*").As("count"))

	sql, params, err := countDs.ToSQL()
	if err != nil {
		return false, err
	}

	c, err := qb.repository.runCountQuery(ctx, sql, params)
	if err != nil {
		return false, err
	}

	return c == 1, nil
}
