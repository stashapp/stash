package exp

type compound struct {
	t   CompoundType
	rhs AppendableExpression
}

func NewCompoundExpression(ct CompoundType, rhs AppendableExpression) CompoundExpression {
	return compound{t: ct, rhs: rhs}
}

func (c compound) Expression() Expression { return c }

func (c compound) Clone() Expression {
	return compound{t: c.t, rhs: c.rhs.Clone().(AppendableExpression)}
}

func (c compound) Type() CompoundType        { return c.t }
func (c compound) RHS() AppendableExpression { return c.rhs }
