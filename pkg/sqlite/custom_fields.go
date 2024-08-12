package sqlite

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

type customFieldsStore struct {
	table exp.IdentifierExpression
	fk    exp.IdentifierExpression
}

func (s *customFieldsStore) deleteForID(ctx context.Context, id int) error {
	table := s.table
	q := dialect.Delete(table).Where(s.fk.Eq(id))
	_, err := exec(ctx, q)
	if err != nil {
		return fmt.Errorf("deleting from %s: %w", s.table.GetTable(), err)
	}

	return nil
}

func (s *customFieldsStore) SetCustomFields(ctx context.Context, id int, values models.CustomFieldsInput) error {
	var partial bool
	var valMap map[string]interface{}

	switch {
	case values.Full != nil:
		partial = false
		valMap = values.Full
	case values.Partial != nil:
		partial = true
		valMap = values.Partial
	}

	if len(valMap) == 0 {
		return nil
	}

	return s.setCustomFields(ctx, id, valMap, partial)
}

func (s *customFieldsStore) setCustomFields(ctx context.Context, id int, values map[string]interface{}, partial bool) error {
	if len(values) == 0 {
		return nil
	}

	if !partial {
		// delete existing custom fields
		if err := s.deleteForID(ctx, id); err != nil {
			return err
		}
	}

	conflictKey := s.fk.GetCol().(string) + ", field"
	// upsert new custom fields
	q := dialect.Insert(s.table).Prepared(true).Cols(s.fk, "field", "value").
		OnConflict(goqu.DoUpdate(conflictKey, goqu.Record{"value": goqu.I("excluded").Col("value")}))
	r := make([]interface{}, len(values))
	var i int
	for key, value := range values {
		r[i] = goqu.Record{"field": key, "value": value, s.fk.GetCol().(string): id}
		i++
	}

	if _, err := exec(ctx, q.Rows(r...)); err != nil {
		return fmt.Errorf("inserting custom fields: %w", err)
	}

	return nil
}

func (s *customFieldsStore) GetCustomFields(ctx context.Context, id int) (map[string]interface{}, error) {
	q := dialect.Select("field", "value").From(s.table).Where(s.fk.Eq(id))

	const single = false
	ret := make(map[string]interface{})
	err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var field string
		var value string
		if err := rows.Scan(&field, &value); err != nil {
			return fmt.Errorf("scanning custom fields: %w", err)
		}
		ret[field] = value
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("getting custom fields: %w", err)
	}

	return ret, nil
}

func (s *customFieldsStore) GetCustomFieldsBulk(ctx context.Context, ids []int) ([]models.CustomFieldMap, error) {
	q := dialect.Select(s.fk.As("id"), "field", "value").From(s.table).Where(s.fk.In(ids))

	const single = false
	ret := make([]models.CustomFieldMap, len(ids))

	idi := make(map[int]int, len(ids))
	for i, id := range ids {
		idi[id] = i
	}

	err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var id int
		var field string
		var value string
		if err := rows.Scan(&id, &field, &value); err != nil {
			return fmt.Errorf("scanning custom fields: %w", err)
		}

		i := idi[id]
		m := ret[i]
		if m == nil {
			m = make(map[string]interface{})
			ret[i] = m
		}

		m[field] = value
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("getting custom fields: %w", err)
	}

	return ret, nil
}
