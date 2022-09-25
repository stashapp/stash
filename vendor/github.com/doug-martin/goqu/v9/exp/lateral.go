package exp

type (
	lateral struct {
		table AppendableExpression
	}
)

// Creates a new SQL lateral expression
//   L(From("test")) -> LATERAL (SELECT * FROM "tests")
func NewLateralExpression(table AppendableExpression) LateralExpression {
	return lateral{table: table}
}

func (l lateral) Clone() Expression {
	return NewLateralExpression(l.table)
}

func (l lateral) Table() AppendableExpression {
	return l.table
}

func (l lateral) Expression() Expression               { return l }
func (l lateral) As(val interface{}) AliasedExpression { return NewAliasExpression(l, val) }
