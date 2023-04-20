package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

const pinnedFilterTable = "pinned_filters"
const pinnedFilterDefaultName = ""

type pinnedFilterQueryBuilder struct {
	repository
}

var PinnedFilterReaderWriter = &pinnedFilterQueryBuilder{
	repository{
		tableName: pinnedFilterTable,
		idColumn:  idColumn,
	},
}

func (qb *pinnedFilterQueryBuilder) Create(ctx context.Context, newObject models.PinnedFilter) (*models.PinnedFilter, error) {
	var ret models.PinnedFilter
	if err := qb.insertObject(ctx, newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *pinnedFilterQueryBuilder) Destroy(ctx context.Context, id int) error {
	return qb.destroyExisting(ctx, []int{id})
}

func (qb *pinnedFilterQueryBuilder) Find(ctx context.Context, id int) (*models.PinnedFilter, error) {
	var ret models.PinnedFilter
	if err := qb.getByID(ctx, id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *pinnedFilterQueryBuilder) FindMany(ctx context.Context, ids []int, ignoreNotFound bool) ([]*models.PinnedFilter, error) {
	var filters []*models.PinnedFilter
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

func (qb *pinnedFilterQueryBuilder) FindByMode(ctx context.Context, mode models.FilterMode) ([]*models.PinnedFilter, error) {
	// exclude empty-named filters - these are the internal default filters

	query := fmt.Sprintf(`SELECT * FROM %s WHERE mode = ? AND name != ? ORDER BY name ASC`, pinnedFilterTable)

	var ret models.PinnedFilters
	if err := qb.query(ctx, query, []interface{}{mode, pinnedFilterDefaultName}, &ret); err != nil {
		return nil, err
	}

	return []*models.PinnedFilter(ret), nil
}

func (qb *pinnedFilterQueryBuilder) All(ctx context.Context) ([]*models.PinnedFilter, error) {
	var ret models.PinnedFilters
	if err := qb.query(ctx, selectAll(pinnedFilterTable), nil, &ret); err != nil {
		return nil, err
	}

	return []*models.PinnedFilter(ret), nil
}
