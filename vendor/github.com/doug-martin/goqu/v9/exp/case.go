package exp

type (
	caseElse struct {
		result interface{}
	}
	caseWhen struct {
		caseElse
		condition interface{}
	}
	caseExpression struct {
		value         interface{}
		whens         []CaseWhen
		elseCondition CaseElse
	}
)

func NewCaseElse(result interface{}) CaseElse {
	return caseElse{result: result}
}

func (ce caseElse) Result() interface{} {
	return ce.result
}

func NewCaseWhen(condition, result interface{}) CaseWhen {
	return caseWhen{caseElse: caseElse{result: result}, condition: condition}
}

func (cw caseWhen) Condition() interface{} {
	return cw.condition
}

func NewCaseExpression() CaseExpression {
	return caseExpression{value: nil, whens: []CaseWhen{}, elseCondition: nil}
}

func (c caseExpression) Expression() Expression {
	return c
}

func (c caseExpression) Clone() Expression {
	return caseExpression{value: c.value, whens: c.whens, elseCondition: c.elseCondition}
}

func (c caseExpression) As(alias interface{}) AliasedExpression {
	return NewAliasExpression(c, alias)
}

func (c caseExpression) GetValue() interface{} {
	return c.value
}

func (c caseExpression) GetWhens() []CaseWhen {
	return c.whens
}

func (c caseExpression) GetElse() CaseElse {
	return c.elseCondition
}

func (c caseExpression) Value(value interface{}) CaseExpression {
	c.value = value
	return c
}

func (c caseExpression) When(condition, result interface{}) CaseExpression {
	c.whens = append(c.whens, NewCaseWhen(condition, result))
	return c
}

func (c caseExpression) Else(result interface{}) CaseExpression {
	c.elseCondition = NewCaseElse(result)
	return c
}

func (c caseExpression) Asc() OrderedExpression  { return asc(c) }
func (c caseExpression) Desc() OrderedExpression { return desc(c) }
