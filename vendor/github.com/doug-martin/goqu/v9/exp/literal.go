package exp

type (
	literal struct {
		literal string
		args    []interface{}
	}
)

// Creates a new SQL literal with the provided arguments.
//   L("a = 1") -> a = 1
// You can also you placeholders. All placeholders within a Literal are represented by '?'
//   L("a = ?", "b") -> a = 'b'
// Literals can also contain placeholders for other expressions
//   L("(? AND ?) OR (?)", I("a").Eq(1), I("b").Eq("b"), I("c").In([]string{"a", "b", "c"}))
func NewLiteralExpression(sql string, args ...interface{}) LiteralExpression {
	return literal{literal: sql, args: args}
}

// Returns a literal for the '*' operator
func Star() LiteralExpression {
	return NewLiteralExpression("*")
}

// Returns a literal for the 'DEFAULT'
func Default() LiteralExpression {
	return NewLiteralExpression("DEFAULT")
}

func (l literal) Clone() Expression {
	return NewLiteralExpression(l.literal, l.args...)
}

func (l literal) Literal() string {
	return l.literal
}

func (l literal) Args() []interface{} {
	return l.args
}

func (l literal) Expression() Expression                           { return l }
func (l literal) As(val interface{}) AliasedExpression             { return NewAliasExpression(l, val) }
func (l literal) Eq(val interface{}) BooleanExpression             { return eq(l, val) }
func (l literal) Neq(val interface{}) BooleanExpression            { return neq(l, val) }
func (l literal) Gt(val interface{}) BooleanExpression             { return gt(l, val) }
func (l literal) Gte(val interface{}) BooleanExpression            { return gte(l, val) }
func (l literal) Lt(val interface{}) BooleanExpression             { return lt(l, val) }
func (l literal) Lte(val interface{}) BooleanExpression            { return lte(l, val) }
func (l literal) Asc() OrderedExpression                           { return asc(l) }
func (l literal) Desc() OrderedExpression                          { return desc(l) }
func (l literal) Between(val RangeVal) RangeExpression             { return between(l, val) }
func (l literal) NotBetween(val RangeVal) RangeExpression          { return notBetween(l, val) }
func (l literal) Like(val interface{}) BooleanExpression           { return like(l, val) }
func (l literal) NotLike(val interface{}) BooleanExpression        { return notLike(l, val) }
func (l literal) ILike(val interface{}) BooleanExpression          { return iLike(l, val) }
func (l literal) NotILike(val interface{}) BooleanExpression       { return notILike(l, val) }
func (l literal) RegexpLike(val interface{}) BooleanExpression     { return regexpLike(l, val) }
func (l literal) RegexpNotLike(val interface{}) BooleanExpression  { return regexpNotLike(l, val) }
func (l literal) RegexpILike(val interface{}) BooleanExpression    { return regexpILike(l, val) }
func (l literal) RegexpNotILike(val interface{}) BooleanExpression { return regexpNotILike(l, val) }
func (l literal) In(vals ...interface{}) BooleanExpression         { return in(l, vals...) }
func (l literal) NotIn(vals ...interface{}) BooleanExpression      { return notIn(l, vals...) }
func (l literal) Is(val interface{}) BooleanExpression             { return is(l, val) }
func (l literal) IsNot(val interface{}) BooleanExpression          { return isNot(l, val) }
func (l literal) IsNull() BooleanExpression                        { return is(l, nil) }
func (l literal) IsNotNull() BooleanExpression                     { return isNot(l, nil) }
func (l literal) IsTrue() BooleanExpression                        { return is(l, true) }
func (l literal) IsNotTrue() BooleanExpression                     { return isNot(l, true) }
func (l literal) IsFalse() BooleanExpression                       { return is(l, false) }
func (l literal) IsNotFalse() BooleanExpression                    { return isNot(l, false) }

func (l literal) BitwiseInversion() BitwiseExpression                { return bitwiseInversion(l) }
func (l literal) BitwiseOr(val interface{}) BitwiseExpression        { return bitwiseOr(l, val) }
func (l literal) BitwiseAnd(val interface{}) BitwiseExpression       { return bitwiseAnd(l, val) }
func (l literal) BitwiseXor(val interface{}) BitwiseExpression       { return bitwiseXor(l, val) }
func (l literal) BitwiseLeftShift(val interface{}) BitwiseExpression { return bitwiseLeftShift(l, val) }
func (l literal) BitwiseRightShift(val interface{}) BitwiseExpression {
	return bitwiseRightShift(l, val)
}
