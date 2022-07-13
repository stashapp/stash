package exp

import (
	"fmt"
	"reflect"

	"github.com/doug-martin/goqu/v9/internal/util"
)

type columnList struct {
	columns []Expression
}

func NewColumnListExpression(vals ...interface{}) ColumnListExpression {
	cols := []Expression{}
	for _, val := range vals {
		switch t := val.(type) {
		case nil: // do nothing
		case string:
			cols = append(cols, ParseIdentifier(t))
		case ColumnListExpression:
			cols = append(cols, t.Columns()...)
		case Expression:
			cols = append(cols, t)
		default:
			_, valKind := util.GetTypeInfo(val, reflect.Indirect(reflect.ValueOf(val)))

			if valKind == reflect.Struct {
				cm, err := util.GetColumnMap(val)
				if err != nil {
					panic(err.Error())
				}
				structCols := cm.Cols()
				for _, col := range structCols {
					i := ParseIdentifier(col)
					var sc Expression = i
					if i.IsQualified() {
						sc = i.As(NewIdentifierExpression("", "", col))
					}
					cols = append(cols, sc)
				}
			} else {
				panic(fmt.Sprintf("Cannot created expression from  %+v", val))
			}
		}
	}
	return columnList{columns: cols}
}

func NewOrderedColumnList(vals ...OrderedExpression) ColumnListExpression {
	exps := make([]interface{}, 0, len(vals))
	for _, col := range vals {
		exps = append(exps, col.Expression())
	}
	return NewColumnListExpression(exps...)
}

func (cl columnList) Clone() Expression {
	newExps := make([]Expression, 0, len(cl.columns))
	for _, exp := range cl.columns {
		newExps = append(newExps, exp.Clone())
	}
	return columnList{columns: newExps}
}

func (cl columnList) Expression() Expression {
	return cl
}

func (cl columnList) IsEmpty() bool {
	return len(cl.columns) == 0
}

func (cl columnList) Columns() []Expression {
	return cl.columns
}

func (cl columnList) Append(cols ...Expression) ColumnListExpression {
	ret := columnList{}
	exps := append(ret.columns, cl.columns...)
	exps = append(exps, cols...)
	ret.columns = exps
	return ret
}
