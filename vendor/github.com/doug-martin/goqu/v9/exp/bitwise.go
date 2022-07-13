package exp

type bitwise struct {
	lhs Expression
	rhs interface{}
	op  BitwiseOperation
}

func NewBitwiseExpression(op BitwiseOperation, lhs Expression, rhs interface{}) BitwiseExpression {
	return bitwise{op: op, lhs: lhs, rhs: rhs}
}

func (b bitwise) Clone() Expression {
	return NewBitwiseExpression(b.op, b.lhs.Clone(), b.rhs)
}

func (b bitwise) RHS() interface{} {
	return b.rhs
}

func (b bitwise) LHS() Expression {
	return b.lhs
}

func (b bitwise) Op() BitwiseOperation {
	return b.op
}

func (b bitwise) Expression() Expression                           { return b }
func (b bitwise) As(val interface{}) AliasedExpression             { return NewAliasExpression(b, val) }
func (b bitwise) Eq(val interface{}) BooleanExpression             { return eq(b, val) }
func (b bitwise) Neq(val interface{}) BooleanExpression            { return neq(b, val) }
func (b bitwise) Gt(val interface{}) BooleanExpression             { return gt(b, val) }
func (b bitwise) Gte(val interface{}) BooleanExpression            { return gte(b, val) }
func (b bitwise) Lt(val interface{}) BooleanExpression             { return lt(b, val) }
func (b bitwise) Lte(val interface{}) BooleanExpression            { return lte(b, val) }
func (b bitwise) Asc() OrderedExpression                           { return asc(b) }
func (b bitwise) Desc() OrderedExpression                          { return desc(b) }
func (b bitwise) Like(i interface{}) BooleanExpression             { return like(b, i) }
func (b bitwise) NotLike(i interface{}) BooleanExpression          { return notLike(b, i) }
func (b bitwise) ILike(i interface{}) BooleanExpression            { return iLike(b, i) }
func (b bitwise) NotILike(i interface{}) BooleanExpression         { return notILike(b, i) }
func (b bitwise) RegexpLike(val interface{}) BooleanExpression     { return regexpLike(b, val) }
func (b bitwise) RegexpNotLike(val interface{}) BooleanExpression  { return regexpNotLike(b, val) }
func (b bitwise) RegexpILike(val interface{}) BooleanExpression    { return regexpILike(b, val) }
func (b bitwise) RegexpNotILike(val interface{}) BooleanExpression { return regexpNotILike(b, val) }
func (b bitwise) In(i ...interface{}) BooleanExpression            { return in(b, i...) }
func (b bitwise) NotIn(i ...interface{}) BooleanExpression         { return notIn(b, i...) }
func (b bitwise) Is(i interface{}) BooleanExpression               { return is(b, i) }
func (b bitwise) IsNot(i interface{}) BooleanExpression            { return isNot(b, i) }
func (b bitwise) IsNull() BooleanExpression                        { return is(b, nil) }
func (b bitwise) IsNotNull() BooleanExpression                     { return isNot(b, nil) }
func (b bitwise) IsTrue() BooleanExpression                        { return is(b, true) }
func (b bitwise) IsNotTrue() BooleanExpression                     { return isNot(b, true) }
func (b bitwise) IsFalse() BooleanExpression                       { return is(b, false) }
func (b bitwise) IsNotFalse() BooleanExpression                    { return isNot(b, false) }
func (b bitwise) Distinct() SQLFunctionExpression                  { return NewSQLFunctionExpression("DISTINCT", b) }
func (b bitwise) Between(val RangeVal) RangeExpression             { return between(b, val) }
func (b bitwise) NotBetween(val RangeVal) RangeExpression          { return notBetween(b, val) }

// used internally to create a Bitwise Inversion BitwiseExpression
func bitwiseInversion(rhs Expression) BitwiseExpression {
	return NewBitwiseExpression(BitwiseInversionOp, nil, rhs)
}

// used internally to create a Bitwise OR BitwiseExpression
func bitwiseOr(lhs Expression, rhs interface{}) BitwiseExpression {
	return NewBitwiseExpression(BitwiseOrOp, lhs, rhs)
}

// used internally to create a Bitwise AND BitwiseExpression
func bitwiseAnd(lhs Expression, rhs interface{}) BitwiseExpression {
	return NewBitwiseExpression(BitwiseAndOp, lhs, rhs)
}

// used internally to create a Bitwise XOR BitwiseExpression
func bitwiseXor(lhs Expression, rhs interface{}) BitwiseExpression {
	return NewBitwiseExpression(BitwiseXorOp, lhs, rhs)
}

// used internally to create a Bitwise LEFT SHIFT BitwiseExpression
func bitwiseLeftShift(lhs Expression, rhs interface{}) BitwiseExpression {
	return NewBitwiseExpression(BitwiseLeftShiftOp, lhs, rhs)
}

// used internally to create a Bitwise RIGHT SHIFT BitwiseExpression
func bitwiseRightShift(lhs Expression, rhs interface{}) BitwiseExpression {
	return NewBitwiseExpression(BitwiseRightShiftOp, lhs, rhs)
}
