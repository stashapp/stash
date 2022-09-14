package exp

import (
	"reflect"
	"sort"

	"github.com/doug-martin/goqu/v9/internal/errors"
	"github.com/doug-martin/goqu/v9/internal/util"
)

type (
	insert struct {
		from AppendableExpression
		cols ColumnListExpression
		vals [][]interface{}
	}
)

func NewInsertExpression(rows ...interface{}) (insertExpression InsertExpression, err error) {
	switch len(rows) {
	case 0:
		return new(insert), nil
	case 1:
		val := reflect.ValueOf(rows[0])
		if val.Kind() == reflect.Slice {
			vals := make([]interface{}, 0, val.Len())
			for i := 0; i < val.Len(); i++ {
				vals = append(vals, val.Index(i).Interface())
			}
			return NewInsertExpression(vals...)
		}
		if ae, ok := rows[0].(AppendableExpression); ok {
			return &insert{from: ae}, nil
		}
	}
	return newInsert(rows...)
}

func (i *insert) Expression() Expression {
	return i
}

func (i *insert) Clone() Expression {
	return i.clone()
}

func (i *insert) clone() *insert {
	return &insert{from: i.from, cols: i.cols, vals: i.vals}
}

func (i *insert) IsEmpty() bool {
	return i.from == nil && (i.cols == nil || i.cols.IsEmpty())
}

func (i *insert) IsInsertFrom() bool {
	return i.from != nil
}

func (i *insert) From() AppendableExpression {
	return i.from
}

func (i *insert) Cols() ColumnListExpression {
	return i.cols
}

func (i *insert) SetCols(cols ColumnListExpression) InsertExpression {
	ci := i.clone()
	ci.cols = cols
	return ci
}

func (i *insert) Vals() [][]interface{} {
	return i.vals
}

func (i *insert) SetVals(vals [][]interface{}) InsertExpression {
	ci := i.clone()
	ci.vals = vals
	return ci
}

// parses the rows gathering and sorting unique columns and values for each record
func newInsert(rows ...interface{}) (insertExp InsertExpression, err error) {
	var mapKeys util.ValueSlice
	rowValue := reflect.Indirect(reflect.ValueOf(rows[0]))
	rowType := rowValue.Type()
	rowKind := rowValue.Kind()
	if rowKind == reflect.Struct {
		return createStructSliceInsert(rows...)
	}
	vals := make([][]interface{}, 0, len(rows))
	var columns ColumnListExpression
	for _, row := range rows {
		if rowType != reflect.Indirect(reflect.ValueOf(row)).Type() {
			return nil, errors.New(
				"rows must be all the same type expected %+v got %+v",
				rowType,
				reflect.Indirect(reflect.ValueOf(row)).Type(),
			)
		}
		newRowValue := reflect.Indirect(reflect.ValueOf(row))
		switch rowKind {
		case reflect.Map:
			if columns == nil {
				mapKeys = util.ValueSlice(newRowValue.MapKeys())
				sort.Sort(mapKeys)
				colKeys := make([]interface{}, 0, len(mapKeys))
				for _, key := range mapKeys {
					colKeys = append(colKeys, key.Interface())
				}
				columns = NewColumnListExpression(colKeys...)
			}
			newMapKeys := util.ValueSlice(newRowValue.MapKeys())
			if len(newMapKeys) != len(mapKeys) {
				return nil, errors.New("rows with different value length expected %d got %d", len(mapKeys), len(newMapKeys))
			}
			if !mapKeys.Equal(newMapKeys) {
				return nil, errors.New("rows with different keys expected %s got %s", mapKeys.String(), newMapKeys.String())
			}
			rowVals := make([]interface{}, 0, len(mapKeys))
			for _, key := range mapKeys {
				rowVals = append(rowVals, newRowValue.MapIndex(key).Interface())
			}
			vals = append(vals, rowVals)
		default:
			return nil, errors.New(
				"unsupported insert must be map, goqu.Record, or struct type got: %T",
				row,
			)
		}
	}
	return &insert{cols: columns, vals: vals}, nil
}

func createStructSliceInsert(rows ...interface{}) (insertExp InsertExpression, err error) {
	rowValue := reflect.Indirect(reflect.ValueOf(rows[0]))
	rowType := rowValue.Type()
	recordRows := make([]interface{}, 0, len(rows))
	for _, row := range rows {
		if rowType != reflect.Indirect(reflect.ValueOf(row)).Type() {
			return nil, errors.New(
				"rows must be all the same type expected %+v got %+v",
				rowType,
				reflect.Indirect(reflect.ValueOf(row)).Type(),
			)
		}
		newRowValue := reflect.Indirect(reflect.ValueOf(row))
		record, err := getFieldsValuesFromStruct(newRowValue)
		if err != nil {
			return nil, err
		}
		recordRows = append(recordRows, record)
	}
	return newInsert(recordRows...)
}

func getFieldsValuesFromStruct(value reflect.Value) (row Record, err error) {
	if value.IsValid() {
		return NewRecordFromStruct(value.Interface(), true, false)
	}
	return
}
