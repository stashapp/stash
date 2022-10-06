package exp

type (
	TruncateClauses interface {
		HasTable() bool
		clone() *truncateClauses

		Table() ColumnListExpression
		SetTable(tables ColumnListExpression) TruncateClauses

		Options() TruncateOptions
		SetOptions(opts TruncateOptions) TruncateClauses
	}
	truncateClauses struct {
		tables  ColumnListExpression
		options TruncateOptions
	}
)

func NewTruncateClauses() TruncateClauses {
	return &truncateClauses{}
}

func (tc *truncateClauses) HasTable() bool {
	return tc.tables != nil
}

func (tc *truncateClauses) clone() *truncateClauses {
	return &truncateClauses{
		tables: tc.tables,
	}
}

func (tc *truncateClauses) Table() ColumnListExpression {
	return tc.tables
}
func (tc *truncateClauses) SetTable(tables ColumnListExpression) TruncateClauses {
	ret := tc.clone()
	ret.tables = tables
	return ret
}

func (tc *truncateClauses) Options() TruncateOptions {
	return tc.options
}
func (tc *truncateClauses) SetOptions(opts TruncateOptions) TruncateClauses {
	ret := tc.clone()
	ret.options = opts
	return ret
}
