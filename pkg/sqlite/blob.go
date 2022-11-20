package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/blob"
)

type BlobStore struct {
	repository

	tableMgr *table
}

const (
	blobTable          = "blobs"
	blobChecksumColumn = "checksum"
)

func NewBlobStore() *BlobStore {
	return &BlobStore{
		repository: repository{
			tableName: blobTable,
			idColumn:  blobChecksumColumn,
		},

		tableMgr: blobTableMgr,
	}
}

type blobRow struct {
	Checksum string `db:"checksum"`
	Blob     []byte `db:"blob"`
}

func (qb *BlobStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *BlobStore) Write(ctx context.Context, checksum string, data []byte) error {
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

func (qb *BlobStore) Read(ctx context.Context, checksum string) (io.ReadCloser, error) {
	q := dialect.From(qb.table()).Select(qb.table().All()).Where(qb.tableMgr.byID(checksum))

	var row blobRow
	const single = true
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		if err := r.StructScan(&row); err != nil {
			return err
		}

		return nil
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, blob.ErrNotFound
		}

		return nil, fmt.Errorf("querying %s: %w", qb.table(), err)
	}

	if row.Blob == nil {
		return nil, nil
	}

	return io.NopCloser(bytes.NewReader(row.Blob)), nil
}

func (qb *BlobStore) Delete(ctx context.Context, checksum string) error {
	table := qb.table()

	q := dialect.Delete(table).Where(goqu.C(blobChecksumColumn).Eq(checksum))

	_, err := exec(ctx, q)
	if err != nil {
		// TODO - handle checksum in use error
		return fmt.Errorf("deleting from %s: %w", table, err)
	}

	return nil
}
