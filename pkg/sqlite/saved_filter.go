package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

const (
	savedFilterTable       = "saved_filters"
	savedFilterDefaultName = ""
)

type savedFilterRow struct {
	ID     int    `db:"id" goqu:"skipinsert"`
	Mode   string `db:"mode"`
	Name   string `db:"name"`
	Filter string `db:"filter"`
}

func (r *savedFilterRow) fromSavedFilter(o models.SavedFilter) {
	r.ID = o.ID
	r.Mode = string(o.Mode)
	r.Name = o.Name
	r.Filter = o.Filter
}

func (r *savedFilterRow) resolve() *models.SavedFilter {
	ret := &models.SavedFilter{
		ID:     r.ID,
		Name:   r.Name,
		Mode:   models.FilterMode(r.Mode),
		Filter: r.Filter,
	}

	return ret
}

type SavedFilterStore struct {
	repository

	tableMgr *table
}

func NewSavedFilterStore() *SavedFilterStore {
	return &SavedFilterStore{
		repository: repository{
			tableName: savedFilterTable,
			idColumn:  idColumn,
		},
		tableMgr: savedFilterTableMgr,
	}
}

func (qb *SavedFilterStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *SavedFilterStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *SavedFilterStore) Create(ctx context.Context, newObject *models.SavedFilter) error {
	var r savedFilterRow
	r.fromSavedFilter(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	updated, err := qb.find(ctx, id)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	*newObject = *updated

	return nil
}

func (qb *SavedFilterStore) Update(ctx context.Context, updatedObject *models.SavedFilter) error {
	var r savedFilterRow
	r.fromSavedFilter(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *SavedFilterStore) SetDefault(ctx context.Context, obj *models.SavedFilter) error {
	// find the existing default
	existing, err := qb.FindDefault(ctx, obj.Mode)
	if err != nil {
		return err
	}

	obj.Name = savedFilterDefaultName

	if existing != nil {
		obj.ID = existing.ID
		return qb.Update(ctx, obj)
	}

	return qb.Create(ctx, obj)
}

func (qb *SavedFilterStore) Destroy(ctx context.Context, id int) error {
	return qb.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *SavedFilterStore) Find(ctx context.Context, id int) (*models.SavedFilter, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *SavedFilterStore) FindMany(ctx context.Context, ids []int, ignoreNotFound bool) ([]*models.SavedFilter, error) {
	ret := make([]*models.SavedFilter, len(ids))

	table := qb.table()
	q := qb.selectDataset().Prepared(true).Where(table.Col(idColumn).In(ids))
	unsorted, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	for _, s := range unsorted {
		i := intslice.IntIndex(ids, s.ID)
		ret[i] = s
	}

	if !ignoreNotFound {
		for i := range ret {
			if ret[i] == nil {
				return nil, fmt.Errorf("filter with id %d not found", ids[i])
			}
		}
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *SavedFilterStore) find(ctx context.Context, id int) (*models.SavedFilter, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *SavedFilterStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.SavedFilter, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *SavedFilterStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.SavedFilter, error) {
	const single = false
	var ret []*models.SavedFilter
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f savedFilterRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		s := f.resolve()

		ret = append(ret, s)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *SavedFilterStore) FindByMode(ctx context.Context, mode models.FilterMode) ([]*models.SavedFilter, error) {
	// SELECT * FROM %s WHERE mode = ? AND name != ? ORDER BY name ASC
	table := qb.table()
	sq := qb.selectDataset().Prepared(true).Where(
		table.Col("mode").Eq(mode),
		table.Col("name").Neq(savedFilterDefaultName),
	).Order(table.Col("name").Asc())
	ret, err := qb.getMany(ctx, sq)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *SavedFilterStore) FindDefault(ctx context.Context, mode models.FilterMode) (*models.SavedFilter, error) {
	// SELECT * FROM saved_filters WHERE mode = ? AND name = ?
	table := qb.table()
	sq := qb.selectDataset().Prepared(true).Where(
		table.Col("mode").Eq(mode),
		table.Col("name").Eq(savedFilterDefaultName),
	)

	ret, err := qb.get(ctx, sq)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return ret, nil
}

func (qb *SavedFilterStore) All(ctx context.Context) ([]*models.SavedFilter, error) {
	return qb.getMany(ctx, qb.selectDataset())
}
