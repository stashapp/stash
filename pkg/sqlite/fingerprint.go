package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/file"
	"gopkg.in/guregu/null.v4"
)

const (
	fingerprintTable    = "fingerprints"
	fingerprintIDColumn = "fingerprint_id"
)

type fingerprintRow struct {
	ID          int         `db:"id" goqu:"skipinsert"`
	Type        string      `db:"type"`
	Fingerprint interface{} `db:"fingerprint"`
}

func (f *fingerprintRow) fromFingerprint(fp file.Fingerprint) {
	f.Type = fp.Type
	f.Fingerprint = fp.Fingerprint
}

type fingerprintQueryRow struct {
	Type        null.String `db:"fingerprint_type"`
	Fingerprint interface{} `db:"fingerprint"`
}

func (r *fingerprintQueryRow) resolve() file.Fingerprint {
	return file.Fingerprint{
		Type:        r.Type.String,
		Fingerprint: r.Fingerprint,
	}
}

type fingerprintQueryBuilder struct {
	repository

	tableMgr *table
}

var FingerprintReaderWriter = &fingerprintQueryBuilder{
	repository: repository{
		tableName: fingerprintTable,
		idColumn:  idColumn,
	},

	tableMgr: fingerprintTableMgr,
}

func (qb *fingerprintQueryBuilder) getOrCreate(ctx context.Context, f file.Fingerprint) (*int, error) {
	id, err := qb.getID(ctx, f)
	if err != nil {
		return nil, err
	}

	if id != nil {
		return id, nil
	}

	return qb.Create(ctx, f)
}

func (qb *fingerprintQueryBuilder) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *fingerprintQueryBuilder) getID(ctx context.Context, f file.Fingerprint) (*int, error) {
	table := qb.table()
	q := dialect.From(table).Select(table.Col(idColumn)).Where(table.Col("type").Eq(f.Type), table.Col("fingerprint").Eq(f.Fingerprint))

	var id *int
	const single = true
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var v int
		if err := rows.Scan(&v); err != nil {
			return err
		}

		id = &v
		return nil
	}); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return id, nil
}

func (qb *fingerprintQueryBuilder) Create(ctx context.Context, f file.Fingerprint) (*int, error) {
	var r fingerprintRow
	r.fromFingerprint(f)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
