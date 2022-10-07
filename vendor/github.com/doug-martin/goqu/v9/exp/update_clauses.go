package exp

type (
	UpdateClauses interface {
		HasTable() bool
		clone() *updateClauses

		CommonTables() []CommonTableExpression
		CommonTablesAppend(cte CommonTableExpression) UpdateClauses

		Table() Expression
		SetTable(table Expression) UpdateClauses

		SetValues() interface{}
		HasSetValues() bool
		SetSetValues(values interface{}) UpdateClauses

		From() ColumnListExpression
		HasFrom() bool
		SetFrom(tables ColumnListExpression) UpdateClauses

		Where() ExpressionList
		ClearWhere() UpdateClauses
		WhereAppend(expressions ...Expression) UpdateClauses

		Order() ColumnListExpression
		HasOrder() bool
		ClearOrder() UpdateClauses
		SetOrder(oes ...OrderedExpression) UpdateClauses
		OrderAppend(...OrderedExpression) UpdateClauses
		OrderPrepend(...OrderedExpression) UpdateClauses

		Limit() interface{}
		HasLimit() bool
		ClearLimit() UpdateClauses
		SetLimit(limit interface{}) UpdateClauses

		Returning() ColumnListExpression
		HasReturning() bool
		SetReturning(cl ColumnListExpression) UpdateClauses
	}
	updateClauses struct {
		commonTables []CommonTableExpression
		table        Expression
		setValues    interface{}
		from         ColumnListExpression
		where        ExpressionList
		order        ColumnListExpression
		limit        interface{}
		returning    ColumnListExpression
	}
)

func NewUpdateClauses() UpdateClauses {
	return &updateClauses{}
}

func (uc *updateClauses) HasTable() bool {
	return uc.table != nil
}

func (uc *updateClauses) clone() *updateClauses {
	return &updateClauses{
		commonTables: uc.commonTables,
		table:        uc.table,
		setValues:    uc.setValues,
		from:         uc.from,
		where:        uc.where,
		order:        uc.order,
		limit:        uc.limit,
		returning:    uc.returning,
	}
}

func (uc *updateClauses) CommonTables() []CommonTableExpression {
	return uc.commonTables
}

func (uc *updateClauses) CommonTablesAppend(cte CommonTableExpression) UpdateClauses {
	ret := uc.clone()
	ret.commonTables = append(ret.commonTables, cte)
	return ret
}

func (uc *updateClauses) Table() Expression {
	return uc.table
}

func (uc *updateClauses) SetTable(table Expression) UpdateClauses {
	ret := uc.clone()
	ret.table = table
	return ret
}

func (uc *updateClauses) SetValues() interface{} {
	return uc.setValues
}

func (uc *updateClauses) HasSetValues() bool {
	return uc.setValues != nil
}

func (uc *updateClauses) SetSetValues(values interface{}) UpdateClauses {
	ret := uc.clone()
	ret.setValues = values
	return ret
}

func (uc *updateClauses) From() ColumnListExpression {
	return uc.from
}

func (uc *updateClauses) HasFrom() bool {
	return uc.from != nil && !uc.from.IsEmpty()
}

func (uc *updateClauses) SetFrom(from ColumnListExpression) UpdateClauses {
	ret := uc.clone()
	ret.from = from
	return ret
}

func (uc *updateClauses) Where() ExpressionList {
	return uc.where
}

func (uc *updateClauses) ClearWhere() UpdateClauses {
	ret := uc.clone()
	ret.where = nil
	return ret
}

func (uc *updateClauses) WhereAppend(expressions ...Expression) UpdateClauses {
	if len(expressions) == 0 {
		return uc
	}
	ret := uc.clone()
	if ret.where == nil {
		ret.where = NewExpressionList(AndType, expressions...)
	} else {
		ret.where = ret.where.Append(expressions...)
	}
	return ret
}

func (uc *updateClauses) Order() ColumnListExpression {
	return uc.order
}

func (uc *updateClauses) HasOrder() bool {
	return uc.order != nil
}

func (uc *updateClauses) ClearOrder() UpdateClauses {
	ret := uc.clone()
	ret.order = nil
	return ret
}

func (uc *updateClauses) SetOrder(oes ...OrderedExpression) UpdateClauses {
	ret := uc.clone()
	ret.order = NewOrderedColumnList(oes...)
	return ret
}

func (uc *updateClauses) OrderAppend(oes ...OrderedExpression) UpdateClauses {
	if uc.order == nil {
		return uc.SetOrder(oes...)
	}
	ret := uc.clone()
	ret.order = ret.order.Append(NewOrderedColumnList(oes...).Columns()...)
	return ret
}

func (uc *updateClauses) OrderPrepend(oes ...OrderedExpression) UpdateClauses {
	if uc.order == nil {
		return uc.SetOrder(oes...)
	}
	ret := uc.clone()
	ret.order = NewOrderedColumnList(oes...).Append(ret.order.Columns()...)
	return ret
}

func (uc *updateClauses) Limit() interface{} {
	return uc.limit
}

func (uc *updateClauses) HasLimit() bool {
	return uc.limit != nil
}

func (uc *updateClauses) ClearLimit() UpdateClauses {
	ret := uc.clone()
	ret.limit = nil
	return ret
}

func (uc *updateClauses) SetLimit(limit interface{}) UpdateClauses {
	ret := uc.clone()
	ret.limit = limit
	return ret
}

func (uc *updateClauses) Returning() ColumnListExpression {
	return uc.returning
}

func (uc *updateClauses) HasReturning() bool {
	return uc.returning != nil && !uc.returning.IsEmpty()
}

func (uc *updateClauses) SetReturning(cl ColumnListExpression) UpdateClauses {
	ret := uc.clone()
	ret.returning = cl
	return ret
}
