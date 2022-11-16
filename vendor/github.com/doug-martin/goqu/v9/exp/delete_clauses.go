package exp

type (
	DeleteClauses interface {
		HasFrom() bool
		clone() *deleteClauses

		CommonTables() []CommonTableExpression
		CommonTablesAppend(cte CommonTableExpression) DeleteClauses

		From() IdentifierExpression
		SetFrom(table IdentifierExpression) DeleteClauses

		Where() ExpressionList
		ClearWhere() DeleteClauses
		WhereAppend(expressions ...Expression) DeleteClauses

		Order() ColumnListExpression
		HasOrder() bool
		ClearOrder() DeleteClauses
		SetOrder(oes ...OrderedExpression) DeleteClauses
		OrderAppend(...OrderedExpression) DeleteClauses
		OrderPrepend(...OrderedExpression) DeleteClauses

		Limit() interface{}
		HasLimit() bool
		ClearLimit() DeleteClauses
		SetLimit(limit interface{}) DeleteClauses

		Returning() ColumnListExpression
		HasReturning() bool
		SetReturning(cl ColumnListExpression) DeleteClauses
	}
	deleteClauses struct {
		commonTables []CommonTableExpression
		from         IdentifierExpression
		where        ExpressionList
		order        ColumnListExpression
		limit        interface{}
		returning    ColumnListExpression
	}
)

func NewDeleteClauses() DeleteClauses {
	return &deleteClauses{}
}

func (dc *deleteClauses) HasFrom() bool {
	return dc.from != nil
}

func (dc *deleteClauses) clone() *deleteClauses {
	return &deleteClauses{
		commonTables: dc.commonTables,
		from:         dc.from,

		where:     dc.where,
		order:     dc.order,
		limit:     dc.limit,
		returning: dc.returning,
	}
}

func (dc *deleteClauses) CommonTables() []CommonTableExpression {
	return dc.commonTables
}

func (dc *deleteClauses) CommonTablesAppend(cte CommonTableExpression) DeleteClauses {
	ret := dc.clone()
	ret.commonTables = append(ret.commonTables, cte)
	return ret
}

func (dc *deleteClauses) From() IdentifierExpression {
	return dc.from
}

func (dc *deleteClauses) SetFrom(table IdentifierExpression) DeleteClauses {
	ret := dc.clone()
	ret.from = table
	return ret
}

func (dc *deleteClauses) Where() ExpressionList {
	return dc.where
}

func (dc *deleteClauses) ClearWhere() DeleteClauses {
	ret := dc.clone()
	ret.where = nil
	return ret
}

func (dc *deleteClauses) WhereAppend(expressions ...Expression) DeleteClauses {
	if len(expressions) == 0 {
		return dc
	}
	ret := dc.clone()
	if ret.where == nil {
		ret.where = NewExpressionList(AndType, expressions...)
	} else {
		ret.where = ret.where.Append(expressions...)
	}
	return ret
}

func (dc *deleteClauses) Order() ColumnListExpression {
	return dc.order
}

func (dc *deleteClauses) HasOrder() bool {
	return dc.order != nil
}

func (dc *deleteClauses) ClearOrder() DeleteClauses {
	ret := dc.clone()
	ret.order = nil
	return ret
}

func (dc *deleteClauses) SetOrder(oes ...OrderedExpression) DeleteClauses {
	ret := dc.clone()
	ret.order = NewOrderedColumnList(oes...)
	return ret
}

func (dc *deleteClauses) OrderAppend(oes ...OrderedExpression) DeleteClauses {
	if dc.order == nil {
		return dc.SetOrder(oes...)
	}
	ret := dc.clone()
	ret.order = ret.order.Append(NewOrderedColumnList(oes...).Columns()...)
	return ret
}

func (dc *deleteClauses) OrderPrepend(oes ...OrderedExpression) DeleteClauses {
	if dc.order == nil {
		return dc.SetOrder(oes...)
	}
	ret := dc.clone()
	ret.order = NewOrderedColumnList(oes...).Append(ret.order.Columns()...)
	return ret
}

func (dc *deleteClauses) Limit() interface{} {
	return dc.limit
}

func (dc *deleteClauses) HasLimit() bool {
	return dc.limit != nil
}

func (dc *deleteClauses) ClearLimit() DeleteClauses {
	ret := dc.clone()
	ret.limit = nil
	return ret
}

func (dc *deleteClauses) SetLimit(limit interface{}) DeleteClauses {
	ret := dc.clone()
	ret.limit = limit
	return ret
}

func (dc *deleteClauses) Returning() ColumnListExpression {
	return dc.returning
}

func (dc *deleteClauses) HasReturning() bool {
	return dc.returning != nil && !dc.returning.IsEmpty()
}

func (dc *deleteClauses) SetReturning(cl ColumnListExpression) DeleteClauses {
	ret := dc.clone()
	ret.returning = cl
	return ret
}
