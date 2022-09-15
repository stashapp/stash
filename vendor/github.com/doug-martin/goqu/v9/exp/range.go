package exp

type (
	ranged struct {
		lhs Expression
		rhs RangeVal
		op  RangeOperation
	}
	rangeVal struct {
		start interface{}
		end   interface{}
	}
)

// used internally to create an BETWEEN comparison RangeExpression
func between(lhs Expression, rhs RangeVal) RangeExpression {
	return NewRangeExpression(BetweenOp, lhs, rhs)
}

// used internally to create an NOT BETWEEN comparison RangeExpression
func notBetween(lhs Expression, rhs RangeVal) RangeExpression {
	return NewRangeExpression(NotBetweenOp, lhs, rhs)
}

func NewRangeExpression(op RangeOperation, lhs Expression, rhs RangeVal) RangeExpression {
	return ranged{op: op, lhs: lhs, rhs: rhs}
}

func (r ranged) Clone() Expression {
	return NewRangeExpression(r.op, r.lhs.Clone(), r.rhs)
}

func (r ranged) Expression() Expression {
	return r
}

func (r ranged) RHS() RangeVal {
	return r.rhs
}

func (r ranged) LHS() Expression {
	return r.lhs
}

func (r ranged) Op() RangeOperation {
	return r.op
}

// Creates a new Range to be used with a Between expression
//    exp.C("col").Between(exp.Range(1, 10))
func NewRangeVal(start, end interface{}) RangeVal {
	return rangeVal{start: start, end: end}
}

func (rv rangeVal) Start() interface{} {
	return rv.start
}

func (rv rangeVal) End() interface{} {
	return rv.end
}
