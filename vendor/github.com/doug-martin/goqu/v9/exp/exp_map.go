package exp

import (
	"sort"
	"strings"

	"github.com/doug-martin/goqu/v9/internal/errors"
)

type (
	// A map of expressions to be ANDed together where the keys are string that will be used as Identifiers and values
	// will be used in a boolean operation.
	// The Ex map can be used in tandem with Op map to create more complex expression such as LIKE, GT, LT...
	// See examples.
	Ex map[string]interface{}
	// A map of expressions to be ORed together where the keys are string that will be used as Identifiers and values
	// will be used in a boolean operation.
	// The Ex map can be used in tandem with Op map to create more complex expression such as LIKE, GT, LT...
	// See examples.
	ExOr map[string]interface{}
	// Used in tandem with the Ex map to create complex comparisons such as LIKE, GT, LT... See examples
	Op map[string]interface{}
)

func (e Ex) Expression() Expression {
	return e
}

func (e Ex) Clone() Expression {
	ret := Ex{}
	for key, val := range e {
		ret[key] = val
	}
	return ret
}

func (e Ex) IsEmpty() bool {
	return len(e) == 0
}

func (e Ex) ToExpressions() (ExpressionList, error) {
	return mapToExpressionList(e, AndType)
}

func (eo ExOr) Expression() Expression {
	return eo
}

func (eo ExOr) Clone() Expression {
	ret := ExOr{}
	for key, val := range eo {
		ret[key] = val
	}
	return ret
}

func (eo ExOr) IsEmpty() bool {
	return len(eo) == 0
}

func (eo ExOr) ToExpressions() (ExpressionList, error) {
	return mapToExpressionList(eo, OrType)
}

func getExMapKeys(ex map[string]interface{}) []string {
	keys := make([]string, 0, len(ex))
	for key := range ex {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func mapToExpressionList(ex map[string]interface{}, eType ExpressionListType) (ExpressionList, error) {
	keys := getExMapKeys(ex)
	ret := make([]Expression, 0, len(keys))
	for _, key := range keys {
		lhs := ParseIdentifier(key)
		rhs := ex[key]
		var exp Expression
		if op, ok := rhs.(Op); ok {
			ors, err := createOredExpressionFromMap(lhs, op)
			if err != nil {
				return nil, err
			}
			exp = NewExpressionList(OrType, ors...)
		} else {
			exp = lhs.Eq(rhs)
		}
		ret = append(ret, exp)
	}
	if eType == OrType {
		return NewExpressionList(OrType, ret...), nil
	}
	return NewExpressionList(AndType, ret...), nil
}

func createOredExpressionFromMap(lhs IdentifierExpression, op Op) ([]Expression, error) {
	opKeys := getExMapKeys(op)
	ors := make([]Expression, 0, len(opKeys))
	for _, opKey := range opKeys {
		if exp, err := createExpressionFromOp(lhs, opKey, op); err != nil {
			return nil, err
		} else if exp != nil {
			ors = append(ors, exp)
		}
	}
	return ors, nil
}

// nolint:gocyclo // not complex just long
func createExpressionFromOp(lhs IdentifierExpression, opKey string, op Op) (exp Expression, err error) {
	switch strings.ToLower(opKey) {
	case EqOp.String():
		exp = lhs.Eq(op[opKey])
	case NeqOp.String():
		exp = lhs.Neq(op[opKey])
	case IsOp.String():
		exp = lhs.Is(op[opKey])
	case IsNotOp.String():
		exp = lhs.IsNot(op[opKey])
	case GtOp.String():
		exp = lhs.Gt(op[opKey])
	case GteOp.String():
		exp = lhs.Gte(op[opKey])
	case LtOp.String():
		exp = lhs.Lt(op[opKey])
	case LteOp.String():
		exp = lhs.Lte(op[opKey])
	case InOp.String():
		exp = lhs.In(op[opKey])
	case NotInOp.String():
		exp = lhs.NotIn(op[opKey])
	case LikeOp.String():
		exp = lhs.Like(op[opKey])
	case NotLikeOp.String():
		exp = lhs.NotLike(op[opKey])
	case ILikeOp.String():
		exp = lhs.ILike(op[opKey])
	case NotILikeOp.String():
		exp = lhs.NotILike(op[opKey])
	case RegexpLikeOp.String():
		exp = lhs.RegexpLike(op[opKey])
	case RegexpNotLikeOp.String():
		exp = lhs.RegexpNotLike(op[opKey])
	case RegexpILikeOp.String():
		exp = lhs.RegexpILike(op[opKey])
	case RegexpNotILikeOp.String():
		exp = lhs.RegexpNotILike(op[opKey])
	case betweenStr:
		rangeVal, ok := op[opKey].(RangeVal)
		if ok {
			exp = lhs.Between(rangeVal)
		}
	case "notbetween":
		rangeVal, ok := op[opKey].(RangeVal)
		if ok {
			exp = lhs.NotBetween(rangeVal)
		}
	default:
		err = errors.New("unsupported expression type %s", opKey)
	}
	return exp, err
}
