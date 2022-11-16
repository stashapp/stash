package exp

type (
	expressionList struct {
		operator    ExpressionListType
		expressions []Expression
	}
)

// A list of expressions that should be ORed together
//    Or(I("a").Eq(10), I("b").Eq(11)) //(("a" = 10) OR ("b" = 11))
func NewExpressionList(operator ExpressionListType, expressions ...Expression) ExpressionList {
	el := expressionList{operator: operator}
	exps := make([]Expression, 0, len(el.expressions))
	for _, e := range expressions {
		switch t := e.(type) {
		case ExpressionList:
			if !t.IsEmpty() {
				exps = append(exps, e)
			}
		case Ex:
			if len(t) > 0 {
				exps = append(exps, e)
			}
		case ExOr:
			if len(t) > 0 {
				exps = append(exps, e)
			}
		default:
			exps = append(exps, e)
		}
	}
	el.expressions = exps
	return el
}

func (el expressionList) Clone() Expression {
	newExps := make([]Expression, 0, len(el.expressions))
	for _, exp := range el.expressions {
		newExps = append(newExps, exp.Clone())
	}
	return expressionList{operator: el.operator, expressions: newExps}
}

func (el expressionList) Expression() Expression {
	return el
}

func (el expressionList) IsEmpty() bool {
	return len(el.expressions) == 0
}

func (el expressionList) Type() ExpressionListType {
	return el.operator
}

func (el expressionList) Expressions() []Expression {
	return el.expressions
}

func (el expressionList) Append(expressions ...Expression) ExpressionList {
	exps := make([]Expression, 0, len(el.expressions)+len(expressions))
	exps = append(exps, el.expressions...)
	exps = append(exps, expressions...)
	return NewExpressionList(el.operator, exps...)
}
