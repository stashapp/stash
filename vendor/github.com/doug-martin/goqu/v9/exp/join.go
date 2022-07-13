package exp

type (
	joinExpression struct {
		isConditioned bool
		// The JoinType
		joinType JoinType
		// The table expressions (e.g. LEFT JOIN "my_table", ON (....))
		table Expression
	}
	// Container for all joins within a dataset
	conditionedJoin struct {
		joinExpression
		// The condition to join (e.g. USING("a", "b"), ON("my_table"."fkey" = "other_table"."id")
		condition JoinCondition
	}
	JoinExpressions []JoinExpression
)

func NewUnConditionedJoinExpression(joinType JoinType, table Expression) JoinExpression {
	return joinExpression{
		joinType:      joinType,
		table:         table,
		isConditioned: false,
	}
}

func (je joinExpression) Clone() Expression {
	return je
}

func (je joinExpression) Expression() Expression {
	return je
}

func (je joinExpression) IsConditioned() bool {
	return je.isConditioned
}

func (je joinExpression) JoinType() JoinType {
	return je.joinType
}

func (je joinExpression) Table() Expression {
	return je.table
}

func NewConditionedJoinExpression(joinType JoinType, table Expression, condition JoinCondition) ConditionedJoinExpression {
	return conditionedJoin{
		joinExpression: joinExpression{
			joinType:      joinType,
			table:         table,
			isConditioned: true,
		},
		condition: condition,
	}
}

func (je conditionedJoin) Clone() Expression {
	return je
}

func (je conditionedJoin) Expression() Expression {
	return je
}

func (je conditionedJoin) Condition() JoinCondition {
	return je.condition
}

func (je conditionedJoin) IsConditionEmpty() bool {
	return je.condition == nil || je.condition.IsEmpty()
}

func (jes JoinExpressions) Clone() JoinExpressions {
	ret := make(JoinExpressions, 0, len(jes))
	for _, jc := range jes {
		ret = append(ret, jc.Clone().(JoinExpression))
	}
	return ret
}

type (
	JoinConditionType int
	JoinCondition     interface {
		Type() JoinConditionType
		IsEmpty() bool
	}
	JoinOnCondition interface {
		JoinCondition
		On() ExpressionList
	}
	JoinUsingCondition interface {
		JoinCondition
		Using() ColumnListExpression
	}
	joinOnCondition struct {
		on ExpressionList
	}

	joinUsingCondition struct {
		using ColumnListExpression
	}
)

// Creates a new ON clause to be used within a join
//    ds.Join(I("my_table"), On(I("my_table.fkey").Eq(I("other_table.id")))
func NewJoinOnCondition(expressions ...Expression) JoinCondition {
	return joinOnCondition{on: NewExpressionList(AndType, expressions...)}
}

func (joc joinOnCondition) Type() JoinConditionType {
	return OnJoinCondType
}

func (joc joinOnCondition) On() ExpressionList {
	return joc.on
}

func (joc joinOnCondition) IsEmpty() bool {
	return len(joc.on.Expressions()) == 0
}

// Creates a new USING clause to be used within a join
func NewJoinUsingCondition(expressions ...interface{}) JoinCondition {
	return joinUsingCondition{using: NewColumnListExpression(expressions...)}
}

func (juc joinUsingCondition) Type() JoinConditionType {
	return UsingJoinCondType
}

func (juc joinUsingCondition) Using() ColumnListExpression {
	return juc.using
}

func (juc joinUsingCondition) IsEmpty() bool {
	return len(juc.using.Columns()) == 0
}
