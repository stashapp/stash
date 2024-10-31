package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const (
	savedFilterTable       = "saved_filters"
	savedFilterDefaultName = ""
)

type savedFilterRow struct {
	ID           int               `db:"id" goqu:"skipinsert"`
	Mode         models.FilterMode `db:"mode"`
	Name         string            `db:"name"`
	FindFilter   string            `db:"find_filter"`
	ObjectFilter string            `db:"object_filter"`
	UIOptions    string            `db:"ui_options"`
}

func encodeJSONOrEmpty(v interface{}) string {
	if v == nil {
		return ""
	}

	encoded, err := json.Marshal(v)
	if err != nil {
		logger.Errorf("error encoding json %v: %v", v, err)
	}

	return string(encoded)
}

func decodeJSON(s string, v interface{}) {
	if s == "" {
		return
	}

	if err := json.Unmarshal([]byte(s), v); err != nil {
		logger.Errorf("error decoding json %q: %v", s, err)
	}
}

func (r *savedFilterRow) fromSavedFilter(o models.SavedFilter) {
	r.ID = o.ID
	r.Mode = o.Mode
	r.Name = o.Name

	// encode the filters as json
	r.FindFilter = encodeJSONOrEmpty(o.FindFilter)
	r.ObjectFilter = encodeJSONOrEmpty(o.ObjectFilter)
	r.UIOptions = encodeJSONOrEmpty(o.UIOptions)
}

func (r *savedFilterRow) resolve() *models.SavedFilter {
	ret := &models.SavedFilter{
		ID:   r.ID,
		Mode: r.Mode,
		Name: r.Name,
	}

	// decode the filters from json
	if r.FindFilter != "" {
		ret.FindFilter = &models.FindFilterType{}
		decodeJSON(r.FindFilter, &ret.FindFilter)
	}
	if r.ObjectFilter != "" {
		ret.ObjectFilter = make(map[string]interface{})
		decodeJSON(r.ObjectFilter, &ret.ObjectFilter)
	}
	if r.UIOptions != "" {
		ret.UIOptions = make(map[string]interface{})
		decodeJSON(r.UIOptions, &ret.UIOptions)
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

	updated, err := qb.Find(ctx, id)
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
		i := slices.Index(ids, s.ID)
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

	// TODO - querying on groups needs to include movies
	// remove this when we migrate to remove the movies filter mode in the database
	var whereClause exp.Expression

	if mode == models.FilterModeGroups || mode == models.FilterModeMovies {
		whereClause = goqu.Or(
			table.Col("mode").Eq(models.FilterModeGroups),
			table.Col("mode").Eq(models.FilterModeMovies),
		)
	} else {
		whereClause = table.Col("mode").Eq(mode)
	}

	sq := qb.selectDataset().Prepared(true).Where(whereClause).Order(table.Col("name").Asc())
	ret, err := qb.getMany(ctx, sq)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *SavedFilterStore) All(ctx context.Context) ([]*models.SavedFilter, error) {
	return qb.getMany(ctx, qb.selectDataset())
}
