package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

const savedFilterTable = "saved_filters"
const savedFilterDefaultName = ""

type savedFilterQueryBuilder struct {
	repository
}

var SavedFilterReaderWriter = &savedFilterQueryBuilder{
	repository{
		tableName: savedFilterTable,
		idColumn:  idColumn,
	},
}

func (qb *savedFilterQueryBuilder) Create(ctx context.Context, newObject models.SavedFilter) (*models.SavedFilter, error) {
	var ret models.SavedFilter
	if err := qb.insertObject(ctx, newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *savedFilterQueryBuilder) Update(ctx context.Context, updatedObject models.SavedFilter) (*models.SavedFilter, error) {
	const partial = false
	if err := qb.update(ctx, updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	var ret models.SavedFilter
	if err := qb.getByID(ctx, updatedObject.ID, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *savedFilterQueryBuilder) SetDefault(ctx context.Context, obj models.SavedFilter) (*models.SavedFilter, error) {
	// find the existing default
	existing, err := qb.FindDefault(ctx, obj.Mode)

	if err != nil {
		return nil, err
	}

	obj.Name = savedFilterDefaultName

	if existing != nil {
		obj.ID = existing.ID
		return qb.Update(ctx, obj)
	}

	return qb.Create(ctx, obj)
}

func (qb *savedFilterQueryBuilder) Destroy(ctx context.Context, id int) error {
	return qb.destroyExisting(ctx, []int{id})
}

func (qb *savedFilterQueryBuilder) Find(ctx context.Context, id int) (*models.SavedFilter, error) {
	var ret models.SavedFilter
	if err := qb.getByID(ctx, id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *savedFilterQueryBuilder) FindMany(ctx context.Context, ids []int, ignoreNotFound bool) ([]*models.SavedFilter, error) {
	var filters []*models.SavedFilter
	for _, id := range ids {
		filter, err := qb.Find(ctx, id)
		if err != nil {
			return nil, err
		}

		if filter == nil && !ignoreNotFound {
			return nil, fmt.Errorf("filter with id %d not found", id)
		}

		filters = append(filters, filter)
	}

	return filters, nil
}

func (qb *savedFilterQueryBuilder) FindByMode(ctx context.Context, mode models.FilterMode) ([]*models.SavedFilter, error) {
	// exclude empty-named filters - these are the internal default filters

	query := fmt.Sprintf(`SELECT * FROM %s WHERE mode = ? AND name != ?`, savedFilterTable)

	var ret models.SavedFilters
	if err := qb.query(ctx, query, []interface{}{mode, savedFilterDefaultName}, &ret); err != nil {
		return nil, err
	}

	return []*models.SavedFilter(ret), nil
}

func (qb *savedFilterQueryBuilder) FindDefault(ctx context.Context, mode models.FilterMode) (*models.SavedFilter, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE mode = ? AND name = ?`, savedFilterTable)

	var ret models.SavedFilters
	if err := qb.query(ctx, query, []interface{}{mode, savedFilterDefaultName}, &ret); err != nil {
		return nil, err
	}

	if len(ret) > 0 {
		return ret[0], nil
	}

	return nil, nil
}

func (qb *savedFilterQueryBuilder) All(ctx context.Context) ([]*models.SavedFilter, error) {
	var ret models.SavedFilters
	if err := qb.query(ctx, selectAll(savedFilterTable), nil, &ret); err != nil {
		return nil, err
	}

	return []*models.SavedFilter(ret), nil
}
