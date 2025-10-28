package sqlite

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/guregu/null.v4"
)

const (
	fingerprintTable = "files_fingerprints"
)

type fingerprintQueryRow struct {
	Type        null.String `db:"fingerprint_type"`
	Fingerprint interface{} `db:"fingerprint"`
}

func (r fingerprintQueryRow) valid() bool {
	return r.Type.Valid
}

func (r *fingerprintQueryRow) resolve() models.Fingerprint {
	return models.Fingerprint{
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
		idColumn:  fileIDColumn,
	},

	tableMgr: fingerprintTableMgr,
}

func (qb *fingerprintQueryBuilder) insert(ctx context.Context, fileID models.FileID, f models.Fingerprint) error {
	table := qb.table()
	q := dialect.Insert(table).Cols(fileIDColumn, "type", "fingerprint").Vals(
		goqu.Vals{fileID, f.Type, f.Fingerprint},
	)
	_, err := exec(ctx, q)
	if err != nil {
		return fmt.Errorf("inserting into %s: %w", table.GetTable(), err)
	}

	return nil
}

func (qb *fingerprintQueryBuilder) insertJoins(ctx context.Context, fileID models.FileID, f []models.Fingerprint) error {
	for _, ff := range f {
		if err := qb.insert(ctx, fileID, ff); err != nil {
			return err
		}
	}

	return nil
}

func (qb *fingerprintQueryBuilder) upsertJoins(ctx context.Context, fileID models.FileID, f []models.Fingerprint) error {
	types := make([]string, len(f))
	for i, ff := range f {
		types[i] = ff.Type
	}

	if err := qb.destroyJoins(ctx, fileID, types); err != nil {
		return err
	}

	for _, ff := range f {
		if err := qb.insert(ctx, fileID, ff); err != nil {
			return err
		}
	}

	return nil
}

func (qb *fingerprintQueryBuilder) replaceJoins(ctx context.Context, fileID models.FileID, f []models.Fingerprint) error {
	if err := qb.destroy(ctx, []int{int(fileID)}); err != nil {
		return err
	}

	return qb.insertJoins(ctx, fileID, f)
}

func (qb *fingerprintQueryBuilder) destroyJoins(ctx context.Context, fileID models.FileID, types []string) error {
	table := qb.table()
	q := dialect.Delete(table).Where(
		table.Col(fileIDColumn).Eq(fileID),
		table.Col("type").In(types),
	)

	_, err := exec(ctx, q)
	if err != nil {
		return fmt.Errorf("deleting from %s: %w", table.GetTable(), err)
	}

	return nil
}

func (qb *fingerprintQueryBuilder) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}
