package sqlite

import (
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

func NewSavedFilterReaderWriter(tx dbi) *savedFilterQueryBuilder {
	return &savedFilterQueryBuilder{
		repository{
			tx:        tx,
			tableName: savedFilterTable,
			idColumn:  idColumn,
		},
	}
}

func (qb *savedFilterQueryBuilder) Create(newObject models.SavedFilter) (*models.SavedFilter, error) {
	var ret models.SavedFilter
	if err := qb.insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *savedFilterQueryBuilder) Update(updatedObject models.SavedFilter) (*models.SavedFilter, error) {
	const partial = false
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	var ret models.SavedFilter
	if err := qb.get(updatedObject.ID, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *savedFilterQueryBuilder) SetDefault(obj models.SavedFilter) (*models.SavedFilter, error) {
	// find the existing default
	existing, err := qb.FindDefault(obj.Mode)

	if err != nil {
		return nil, err
	}

	obj.Name = savedFilterDefaultName

	if existing != nil {
		obj.ID = existing.ID
		return qb.Update(obj)
	}

	return qb.Create(obj)
}

func (qb *savedFilterQueryBuilder) Destroy(id int) error {
	return qb.destroyExisting([]int{id})
}

func (qb *savedFilterQueryBuilder) Find(id int) (*models.SavedFilter, error) {
	var ret models.SavedFilter
	if err := qb.get(id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *savedFilterQueryBuilder) FindByMode(mode models.FilterMode) ([]*models.SavedFilter, error) {
	// exclude empty-named filters - these are the internal default filters

	query := fmt.Sprintf(`SELECT * FROM %s WHERE mode = ? AND name != ?`, savedFilterTable)

	var ret models.SavedFilters
	if err := qb.query(query, []interface{}{mode, savedFilterDefaultName}, &ret); err != nil {
		return nil, err
	}

	return []*models.SavedFilter(ret), nil
}

func (qb *savedFilterQueryBuilder) FindDefault(mode models.FilterMode) (*models.SavedFilter, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE mode = ? AND name = ?`, savedFilterTable)

	var ret models.SavedFilters
	if err := qb.query(query, []interface{}{mode, savedFilterDefaultName}, &ret); err != nil {
		return nil, err
	}

	if len(ret) > 0 {
		return ret[0], nil
	}

	return nil, nil
}
