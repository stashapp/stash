package sqlite

import (
	"context"
	"fmt"
	"regexp"

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
	default:
		return nil
	}

	return s.setCustomFields(ctx, id, valMap, partial)
}

func (s *customFieldsStore) setCustomFields(ctx context.Context, id int, values map[string]interface{}, partial bool) error {
	if !partial {
		// delete existing custom fields
		if err := s.deleteForID(ctx, id); err != nil {
			return err
		}
	}

	if len(values) == 0 {
		return nil
	}

	conflictKey := s.fk.GetCol().(string) + ", field"
	// upsert new custom fields
	q := dialect.Insert(s.table).Prepared(true).Cols(s.fk, "field", "value").
		OnConflict(goqu.DoUpdate(conflictKey, goqu.Record{"value": goqu.I("excluded.value")}))
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
		var value interface{}
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
		var value interface{}
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

type customFieldsFilterHandler struct {
	table string
	fkCol string
	c     []models.CustomFieldCriterionInput
	idCol string
}

func (h *customFieldsFilterHandler) innerJoin(f *filterBuilder, as string, field string) {
	joinOn := fmt.Sprintf("%s = %s.%s AND %s.field = ?", h.idCol, as, h.fkCol, as)
	f.addInnerJoin(h.table, as, joinOn, field)
}

func (h *customFieldsFilterHandler) leftJoin(f *filterBuilder, as string, field string) {
	joinOn := fmt.Sprintf("%s = %s.%s AND %s.field = ?", h.idCol, as, h.fkCol, as)
	f.addLeftJoin(h.table, as, joinOn, field)
}

func (h *customFieldsFilterHandler) handleCriterion(f *filterBuilder, joinAs string, cc models.CustomFieldCriterionInput) {
	switch cc.Modifier {
	case models.CriterionModifierEquals:
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%[1]s.value IN %s", joinAs, getInBinding(len(cc.Value))), cc.Value...)
	case models.CriterionModifierNotEquals:
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%[1]s.value NOT IN %s", joinAs, getInBinding(len(cc.Value))), cc.Value...)
	case models.CriterionModifierIncludes:
		clauses := make([]sqlClause, len(cc.Value))
		for i, v := range cc.Value {
			clauses[i] = makeClause(fmt.Sprintf("%s.value LIKE ?", joinAs), fmt.Sprintf("%%%v%%", v))
		}
		h.innerJoin(f, joinAs, cc.Field)
		f.whereClauses = append(f.whereClauses, clauses...)
	case models.CriterionModifierExcludes:
		for _, v := range cc.Value {
			f.addWhere(fmt.Sprintf("%[1]s.value NOT LIKE ?", joinAs), fmt.Sprintf("%%%v%%", v))
		}
		h.leftJoin(f, joinAs, cc.Field)
	case models.CriterionModifierMatchesRegex:
		for _, v := range cc.Value {
			vs, ok := v.(string)
			if !ok {
				f.setError(fmt.Errorf("unsupported custom field criterion value type: %T", v))
			}
			if _, err := regexp.Compile(vs); err != nil {
				f.setError(err)
				return
			}
			f.addWhere(fmt.Sprintf("(%s.value regexp ?)", joinAs), v)
		}
		h.innerJoin(f, joinAs, cc.Field)
	case models.CriterionModifierNotMatchesRegex:
		for _, v := range cc.Value {
			vs, ok := v.(string)
			if !ok {
				f.setError(fmt.Errorf("unsupported custom field criterion value type: %T", v))
			}
			if _, err := regexp.Compile(vs); err != nil {
				f.setError(err)
				return
			}
			f.addWhere(fmt.Sprintf("(%s.value IS NULL OR %[1]s.value NOT regexp ?)", joinAs), v)
		}
		h.leftJoin(f, joinAs, cc.Field)
	case models.CriterionModifierIsNull:
		h.leftJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%s.value IS NULL OR TRIM(%[1]s.value) = ''", joinAs))
	case models.CriterionModifierNotNull:
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("TRIM(%[1]s.value) != ''", joinAs))
	case models.CriterionModifierBetween:
		if len(cc.Value) != 2 {
			f.setError(fmt.Errorf("expected 2 values for custom field criterion modifier BETWEEN, got %d", len(cc.Value)))
			return
		}
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%s.value BETWEEN ? AND ?", joinAs), cc.Value[0], cc.Value[1])
	case models.CriterionModifierNotBetween:
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%s.value NOT BETWEEN ? AND ?", joinAs), cc.Value[0], cc.Value[1])
	case models.CriterionModifierLessThan:
		if len(cc.Value) != 1 {
			f.setError(fmt.Errorf("expected 1 value for custom field criterion modifier LESS_THAN, got %d", len(cc.Value)))
			return
		}
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%s.value < ?", joinAs), cc.Value[0])
	case models.CriterionModifierGreaterThan:
		if len(cc.Value) != 1 {
			f.setError(fmt.Errorf("expected 1 value for custom field criterion modifier LESS_THAN, got %d", len(cc.Value)))
			return
		}
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%s.value > ?", joinAs), cc.Value[0])
	default:
		f.setError(fmt.Errorf("unsupported custom field criterion modifier: %s", cc.Modifier))
	}
}

func (h *customFieldsFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	if len(h.c) == 0 {
		return
	}

	for i, cc := range h.c {
		join := fmt.Sprintf("custom_fields_%d", i)
		h.handleCriterion(f, join, cc)
	}
}
