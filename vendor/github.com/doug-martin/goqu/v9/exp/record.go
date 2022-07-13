package exp

import (
	"reflect"
	"sort"

	"github.com/doug-martin/goqu/v9/internal/util"
)

// Alternative to writing map[string]interface{}. Can be used for Inserts, Updates or Deletes
type Record map[string]interface{}

func (r Record) Cols() []string {
	cols := make([]string, 0, len(r))
	for col := range r {
		cols = append(cols, col)
	}
	sort.Strings(cols)
	return cols
}

func NewRecordFromStruct(i interface{}, forInsert, forUpdate bool) (r Record, err error) {
	value := reflect.ValueOf(i)
	if value.IsValid() {
		cm, err := util.GetColumnMap(value.Interface())
		if err != nil {
			return nil, err
		}
		cols := cm.Cols()
		r = make(map[string]interface{}, len(cols))
		for _, col := range cols {
			f := cm[col]
			if !shouldSkipField(f, forInsert, forUpdate) {
				if ok, fieldVal := getFieldValue(value, f); ok {
					r[f.ColumnName] = fieldVal
				}
			}
		}
	}
	return
}

func shouldSkipField(f util.ColumnData, forInsert, forUpdate bool) bool {
	shouldSkipInsert := forInsert && !f.ShouldInsert
	shouldSkipUpdate := forUpdate && !f.ShouldUpdate
	return shouldSkipInsert || shouldSkipUpdate
}

func getFieldValue(val reflect.Value, f util.ColumnData) (ok bool, fieldVal interface{}) {
	if v, isAvailable := util.SafeGetFieldByIndex(val, f.FieldIndex); !isAvailable {
		return false, nil
	} else if f.DefaultIfEmpty && util.IsEmptyValue(v) {
		return true, Default()
	} else if v.IsValid() {
		return true, v.Interface()
	} else {
		return true, reflect.Zero(f.GoType).Interface()
	}
}
