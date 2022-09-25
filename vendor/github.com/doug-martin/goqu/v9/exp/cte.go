package exp

type commonExpr struct {
	recursive bool
	name      LiteralExpression
	subQuery  Expression
}

// Creates a new WITH common table expression for a SQLExpression, typically Datasets'. This function is used
// internally by Dataset when a CTE is added to another Dataset
func NewCommonTableExpression(recursive bool, name string, subQuery Expression) CommonTableExpression {
	return commonExpr{recursive: recursive, name: NewLiteralExpression(name), subQuery: subQuery}
}

func (ce commonExpr) Expression() Expression { return ce }

func (ce commonExpr) Clone() Expression {
	return commonExpr{recursive: ce.recursive, name: ce.name, subQuery: ce.subQuery.Clone().(SQLExpression)}
}

func (ce commonExpr) IsRecursive() bool       { return ce.recursive }
func (ce commonExpr) Name() LiteralExpression { return ce.name }
func (ce commonExpr) SubQuery() Expression    { return ce.subQuery }
