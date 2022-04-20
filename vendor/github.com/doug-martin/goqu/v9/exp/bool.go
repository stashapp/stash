package exp

import (
	"reflect"
	"regexp"
)

type boolean struct {
	lhs Expression
	rhs interface{}
	op  BooleanOperation
}

func NewBooleanExpression(op BooleanOperation, lhs Expression, rhs interface{}) BooleanExpression {
	return boolean{op: op, lhs: lhs, rhs: rhs}
}

func (b boolean) Clone() Expression {
	return NewBooleanExpression(b.op, b.lhs.Clone(), b.rhs)
}

func (b boolean) Expression() Expression {
	return b
}

func (b boolean) RHS() interface{} {
	return b.rhs
}

func (b boolean) LHS() Expression {
	return b.lhs
}

func (b boolean) Op() BooleanOperation {
	return b.op
}

func (b boolean) As(val interface{}) AliasedExpression {
	return NewAliasExpression(b, val)
}

// used internally to create an equality BooleanExpression
func eq(lhs Expression, rhs interface{}) BooleanExpression {
	return checkBoolExpType(EqOp, lhs, rhs, false)
}

// used internally to create an in-equality BooleanExpression
func neq(lhs Expression, rhs interface{}) BooleanExpression {
	return checkBoolExpType(EqOp, lhs, rhs, true)
}

// used internally to create an gt comparison BooleanExpression
func gt(lhs Expression, rhs interface{}) BooleanExpression {
	return NewBooleanExpression(GtOp, lhs, rhs)
}

// used internally to create an gte comparison BooleanExpression
func gte(lhs Expression, rhs interface{}) BooleanExpression {
	return NewBooleanExpression(GteOp, lhs, rhs)
}

// used internally to create an lt comparison BooleanExpression
func lt(lhs Expression, rhs interface{}) BooleanExpression {
	return NewBooleanExpression(LtOp, lhs, rhs)
}

// used internally to create an lte comparison BooleanExpression
func lte(lhs Expression, rhs interface{}) BooleanExpression {
	return NewBooleanExpression(LteOp, lhs, rhs)
}

// used internally to create an IN BooleanExpression
func in(lhs Expression, vals ...interface{}) BooleanExpression {
	if len(vals) == 1 && reflect.Indirect(reflect.ValueOf(vals[0])).Kind() == reflect.Slice {
		return NewBooleanExpression(InOp, lhs, vals[0])
	}
	return NewBooleanExpression(InOp, lhs, vals)
}

// used internally to create a NOT IN BooleanExpression
func notIn(lhs Expression, vals ...interface{}) BooleanExpression {
	if len(vals) == 1 && reflect.Indirect(reflect.ValueOf(vals[0])).Kind() == reflect.Slice {
		return NewBooleanExpression(NotInOp, lhs, vals[0])
	}
	return NewBooleanExpression(NotInOp, lhs, vals)
}

// used internally to create an IS BooleanExpression
func is(lhs Expression, val interface{}) BooleanExpression {
	return checkBoolExpType(IsOp, lhs, val, false)
}

// used internally to create an IS NOT BooleanExpression
func isNot(lhs Expression, val interface{}) BooleanExpression {
	return checkBoolExpType(IsOp, lhs, val, true)
}

// used internally to create a LIKE BooleanExpression
func like(lhs Expression, val interface{}) BooleanExpression {
	return checkLikeExp(LikeOp, lhs, val, false)
}

// used internally to create an ILIKE BooleanExpression
func iLike(lhs Expression, val interface{}) BooleanExpression {
	return checkLikeExp(ILikeOp, lhs, val, false)
}

// used internally to create a NOT LIKE BooleanExpression
func notLike(lhs Expression, val interface{}) BooleanExpression {
	return checkLikeExp(LikeOp, lhs, val, true)
}

// used internally to create a NOT ILIKE BooleanExpression
func notILike(lhs Expression, val interface{}) BooleanExpression {
	return checkLikeExp(ILikeOp, lhs, val, true)
}

// used internally to create a LIKE BooleanExpression
func regexpLike(lhs Expression, val interface{}) BooleanExpression {
	return checkLikeExp(RegexpLikeOp, lhs, val, false)
}

// used internally to create an ILIKE BooleanExpression
func regexpILike(lhs Expression, val interface{}) BooleanExpression {
	return checkLikeExp(RegexpILikeOp, lhs, val, false)
}

// used internally to create a NOT LIKE BooleanExpression
func regexpNotLike(lhs Expression, val interface{}) BooleanExpression {
	return checkLikeExp(RegexpLikeOp, lhs, val, true)
}

// used internally to create a NOT ILIKE BooleanExpression
func regexpNotILike(lhs Expression, val interface{}) BooleanExpression {
	return checkLikeExp(RegexpILikeOp, lhs, val, true)
}

// checks an like rhs to create the proper like expression for strings or regexps
func checkLikeExp(op BooleanOperation, lhs Expression, val interface{}, invert bool) BooleanExpression {
	rhs := val

	if t, ok := val.(*regexp.Regexp); ok {
		if op == LikeOp {
			op = RegexpLikeOp
		} else if op == ILikeOp {
			op = RegexpILikeOp
		}
		rhs = t.String()
	}
	if invert {
		op = operatorInversions[op]
	}
	return NewBooleanExpression(op, lhs, rhs)
}

// checks a boolean operation normalizing the operation based on the RHS (e.g. "a" = true vs "a" IS TRUE
func checkBoolExpType(op BooleanOperation, lhs Expression, rhs interface{}, invert bool) BooleanExpression {
	if rhs == nil {
		op = IsOp
	} else {
		switch reflect.Indirect(reflect.ValueOf(rhs)).Kind() {
		case reflect.Bool:
			op = IsOp
		case reflect.Slice:
			// if its a slice of bytes dont treat as an IN
			if _, ok := rhs.([]byte); !ok {
				op = InOp
			}
		case reflect.Struct:
			switch rhs.(type) {
			case SQLExpression:
				op = InOp
			case AppendableExpression:
				op = InOp
			case *regexp.Regexp:
				return checkLikeExp(LikeOp, lhs, rhs, invert)
			}
		default:
		}
	}
	if invert {
		op = operatorInversions[op]
	}
	return NewBooleanExpression(op, lhs, rhs)
}
