package sqlite

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

const maxCustomFieldNameLength = 64

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

	if valMap != nil {
		if err := s.validateCustomFields(valMap, values.Remove); err != nil {
			return err
		}

		if err := s.setCustomFields(ctx, id, valMap, partial); err != nil {
			return err
		}
	}

	if err := s.deleteCustomFields(ctx, id, values.Remove); err != nil {
		return err
	}

	return nil
}

func (s *customFieldsStore) validateCustomFields(values map[string]interface{}, deleteKeys []string) error {
	// if values is nil, nothing to validate
	if values == nil {
		return nil
	}

	// ensure that custom field names are valid
	// no leading or trailing whitespace, no empty strings
	for k := range values {
		if err := s.validateCustomFieldName(k); err != nil {
			return fmt.Errorf("custom field name %q: %w", k, err)
		}
	}

	// ensure delete keys are not also in values
	for _, k := range deleteKeys {
		if _, ok := values[k]; ok {
			return fmt.Errorf("custom field name %q cannot be in both values and delete keys", k)
		}
	}

	return nil
}

func (s *customFieldsStore) validateCustomFieldName(fieldName string) error {
	// ensure that custom field names are valid
	// no leading or trailing whitespace, no empty strings
	if strings.TrimSpace(fieldName) == "" {
		return fmt.Errorf("custom field name cannot be empty")
	}
	if fieldName != strings.TrimSpace(fieldName) {
		return fmt.Errorf("custom field name cannot have leading or trailing whitespace")
	}
	if len(fieldName) > maxCustomFieldNameLength {
		return fmt.Errorf("custom field name must be less than %d characters", maxCustomFieldNameLength+1)
	}
	return nil
}

func getSQLValueFromCustomFieldInput(input interface{}) (interface{}, error) {
	switch v := input.(type) {
	case []interface{}, map[string]interface{}:
		// TODO - in future it would be nice to convert to a JSON string
		// however, we would need some way to differentiate between a JSON string and a regular string
		// for now, we will not support objects and arrays
		return nil, fmt.Errorf("unsupported custom field value type: %T", input)
	default:
		return v, nil
	}
}

func (s *customFieldsStore) sqlValueToValue(value interface{}) interface{} {
	// TODO - if we ever support objects and arrays we will need to add support here
	return value
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
		v, err := getSQLValueFromCustomFieldInput(value)
		if err != nil {
			return fmt.Errorf("getting SQL value for field %q: %w", key, err)
		}
		r[i] = goqu.Record{"field": key, "value": v, s.fk.GetCol().(string): id}
		i++
	}

	if _, err := exec(ctx, q.Rows(r...)); err != nil {
		return fmt.Errorf("inserting custom fields: %w", err)
	}

	return nil
}

func (s *customFieldsStore) deleteCustomFields(ctx context.Context, id int, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	q := dialect.Delete(s.table).
		Where(s.fk.Eq(id)).
		Where(goqu.I("field").In(keys))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("deleting custom fields: %w", err)
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
		ret[field] = s.sqlValueToValue(value)
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

		m[field] = s.sqlValueToValue(value)
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
	// convert values
	cv := make([]interface{}, len(cc.Value))
	for i, v := range cc.Value {
		var err error
		cv[i], err = getSQLValueFromCustomFieldInput(v)
		if err != nil {
			f.setError(err)
			return
		}
	}

	switch cc.Modifier {
	case models.CriterionModifierEquals:
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%[1]s.value IN %s", joinAs, getInBinding(len(cv))), cv...)
	case models.CriterionModifierNotEquals:
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%[1]s.value NOT IN %s", joinAs, getInBinding(len(cv))), cv...)
	case models.CriterionModifierIncludes:
		clauses := make([]sqlClause, len(cv))
		for i, v := range cv {
			clauses[i] = makeClause(fmt.Sprintf("%s.value LIKE ?", joinAs), fmt.Sprintf("%%%v%%", v))
		}
		h.innerJoin(f, joinAs, cc.Field)
		f.whereClauses = append(f.whereClauses, clauses...)
	case models.CriterionModifierExcludes:
		for _, v := range cv {
			f.addWhere(fmt.Sprintf("%[1]s.value NOT LIKE ?", joinAs), fmt.Sprintf("%%%v%%", v))
		}
		h.leftJoin(f, joinAs, cc.Field)
	case models.CriterionModifierMatchesRegex:
		for _, v := range cv {
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
		for _, v := range cv {
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
		if len(cv) != 2 {
			f.setError(fmt.Errorf("expected 2 values for custom field criterion modifier BETWEEN, got %d", len(cv)))
			return
		}
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%s.value BETWEEN ? AND ?", joinAs), cv[0], cv[1])
	case models.CriterionModifierNotBetween:
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%s.value NOT BETWEEN ? AND ?", joinAs), cv[0], cv[1])
	case models.CriterionModifierLessThan:
		if len(cv) != 1 {
			f.setError(fmt.Errorf("expected 1 value for custom field criterion modifier LESS_THAN, got %d", len(cv)))
			return
		}
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%s.value < ?", joinAs), cv[0])
	case models.CriterionModifierGreaterThan:
		if len(cv) != 1 {
			f.setError(fmt.Errorf("expected 1 value for custom field criterion modifier LESS_THAN, got %d", len(cv)))
			return
		}
		h.innerJoin(f, joinAs, cc.Field)
		f.addWhere(fmt.Sprintf("%s.value > ?", joinAs), cv[0])
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
