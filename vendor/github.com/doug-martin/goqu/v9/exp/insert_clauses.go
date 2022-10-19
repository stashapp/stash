package exp

type (
	InsertClauses interface {
		CommonTables() []CommonTableExpression
		CommonTablesAppend(cte CommonTableExpression) InsertClauses

		HasInto() bool
		clone() *insertClauses

		Cols() ColumnListExpression
		HasCols() bool
		ColsAppend(cols ColumnListExpression) InsertClauses
		SetCols(cols ColumnListExpression) InsertClauses

		Into() Expression
		SetInto(cl Expression) InsertClauses

		Returning() ColumnListExpression
		HasReturning() bool
		SetReturning(cl ColumnListExpression) InsertClauses

		From() AppendableExpression
		HasFrom() bool
		SetFrom(ae AppendableExpression) InsertClauses

		Rows() []interface{}
		HasRows() bool
		SetRows(rows []interface{}) InsertClauses

		HasAlias() bool
		Alias() IdentifierExpression
		SetAlias(ie IdentifierExpression) InsertClauses

		Vals() [][]interface{}
		HasVals() bool
		SetVals(vals [][]interface{}) InsertClauses
		ValsAppend(vals [][]interface{}) InsertClauses

		OnConflict() ConflictExpression
		SetOnConflict(expression ConflictExpression) InsertClauses
	}
	insertClauses struct {
		commonTables []CommonTableExpression
		cols         ColumnListExpression
		into         Expression
		returning    ColumnListExpression
		alias        IdentifierExpression
		rows         []interface{}
		values       [][]interface{}
		from         AppendableExpression
		conflict     ConflictExpression
	}
)

func NewInsertClauses() InsertClauses {
	return &insertClauses{}
}

func (ic *insertClauses) HasInto() bool {
	return ic.into != nil
}

func (ic *insertClauses) clone() *insertClauses {
	return &insertClauses{
		commonTables: ic.commonTables,
		cols:         ic.cols,
		into:         ic.into,
		returning:    ic.returning,
		alias:        ic.alias,
		rows:         ic.rows,
		values:       ic.values,
		from:         ic.from,
		conflict:     ic.conflict,
	}
}

func (ic *insertClauses) CommonTables() []CommonTableExpression {
	return ic.commonTables
}

func (ic *insertClauses) CommonTablesAppend(cte CommonTableExpression) InsertClauses {
	ret := ic.clone()
	ret.commonTables = append(ret.commonTables, cte)
	return ret
}

func (ic *insertClauses) Cols() ColumnListExpression {
	return ic.cols
}

func (ic *insertClauses) HasCols() bool {
	return ic.cols != nil && !ic.cols.IsEmpty()
}

func (ic *insertClauses) ColsAppend(cl ColumnListExpression) InsertClauses {
	ret := ic.clone()
	ret.cols = ret.cols.Append(cl.Columns()...)
	return ret
}

func (ic *insertClauses) SetCols(cl ColumnListExpression) InsertClauses {
	ret := ic.clone()
	ret.cols = cl
	return ret
}

func (ic *insertClauses) Into() Expression {
	return ic.into
}

func (ic *insertClauses) SetInto(into Expression) InsertClauses {
	ret := ic.clone()
	ret.into = into
	return ret
}

func (ic *insertClauses) Returning() ColumnListExpression {
	return ic.returning
}

func (ic *insertClauses) HasReturning() bool {
	return ic.returning != nil && !ic.returning.IsEmpty()
}

func (ic *insertClauses) HasAlias() bool {
	return ic.alias != nil
}

func (ic *insertClauses) Alias() IdentifierExpression {
	return ic.alias
}

func (ic *insertClauses) SetAlias(ie IdentifierExpression) InsertClauses {
	ret := ic.clone()
	ret.alias = ie
	return ret
}

func (ic *insertClauses) SetReturning(cl ColumnListExpression) InsertClauses {
	ret := ic.clone()
	ret.returning = cl
	return ret
}

func (ic *insertClauses) From() AppendableExpression {
	return ic.from
}

func (ic *insertClauses) HasFrom() bool {
	return ic.from != nil
}

func (ic *insertClauses) SetFrom(ae AppendableExpression) InsertClauses {
	ret := ic.clone()
	ret.from = ae
	return ret
}

func (ic *insertClauses) Rows() []interface{} {
	return ic.rows
}

func (ic *insertClauses) HasRows() bool {
	return ic.rows != nil && len(ic.rows) > 0
}

func (ic *insertClauses) SetRows(rows []interface{}) InsertClauses {
	ret := ic.clone()
	ret.rows = rows
	return ret
}

func (ic *insertClauses) Vals() [][]interface{} {
	return ic.values
}

func (ic *insertClauses) HasVals() bool {
	return ic.values != nil && len(ic.values) > 0
}

func (ic *insertClauses) SetVals(vals [][]interface{}) InsertClauses {
	ret := ic.clone()
	ret.values = vals
	return ret
}

func (ic *insertClauses) ValsAppend(vals [][]interface{}) InsertClauses {
	ret := ic.clone()
	newVals := make([][]interface{}, 0, len(ic.values)+len(vals))
	newVals = append(newVals, ic.values...)
	newVals = append(newVals, vals...)
	ret.values = newVals
	return ret
}

func (ic *insertClauses) OnConflict() ConflictExpression {
	return ic.conflict
}

func (ic *insertClauses) SetOnConflict(expression ConflictExpression) InsertClauses {
	ret := ic.clone()
	ret.conflict = expression
	return ret
}
