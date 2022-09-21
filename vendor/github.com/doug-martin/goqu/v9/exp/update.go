package exp

import (
	"reflect"
	"sort"

	"github.com/doug-martin/goqu/v9/internal/errors"
	"github.com/doug-martin/goqu/v9/internal/util"
)

type (
	update struct {
		col IdentifierExpression
		val interface{}
	}
)

func set(col IdentifierExpression, val interface{}) UpdateExpression {
	return update{col: col, val: val}
}

func NewUpdateExpressions(update interface{}) (updates []UpdateExpression, err error) {
	if u, ok := update.(UpdateExpression); ok {
		updates = append(updates, u)
		return updates, nil
	}
	updateValue := reflect.Indirect(reflect.ValueOf(update))
	switch updateValue.Kind() {
	case reflect.Map:
		keys := util.ValueSlice(updateValue.MapKeys())
		sort.Sort(keys)
		for _, key := range keys {
			updates = append(updates, ParseIdentifier(key.String()).Set(updateValue.MapIndex(key).Interface()))
		}
	case reflect.Struct:
		return getUpdateExpressionsStruct(updateValue)
	default:
		return nil, errors.New("unsupported update interface type %+v", updateValue.Type())
	}
	return updates, nil
}

func getUpdateExpressionsStruct(value reflect.Value) (updates []UpdateExpression, err error) {
	r, err := NewRecordFromStruct(value.Interface(), false, true)
	if err != nil {
		return updates, err
	}
	cols := r.Cols()
	for _, col := range cols {
		updates = append(updates, ParseIdentifier(col).Set(r[col]))
	}
	return updates, nil
}

func (u update) Expression() Expression {
	return u
}

func (u update) Clone() Expression {
	return update{col: u.col.Clone().(IdentifierExpression), val: u.val}
}

func (u update) Col() IdentifierExpression {
	return u.col
}

func (u update) Val() interface{} {
	return u.val
}
