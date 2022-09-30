package exp

type (
	orderedExpression struct {
		sortExpression Expression
		direction      SortDirection
		nullSortType   NullSortType
	}
)

// used internally to create a new SORT_ASC OrderedExpression
func asc(exp Expression) OrderedExpression {
	return NewOrderedExpression(exp, AscDir, NoNullsSortType)
}

// used internally to create a new SORT_DESC OrderedExpression
func desc(exp Expression) OrderedExpression {
	return NewOrderedExpression(exp, DescSortDir, NoNullsSortType)
}

// used internally to create a new SORT_ASC OrderedExpression
func NewOrderedExpression(exp Expression, direction SortDirection, sortType NullSortType) OrderedExpression {
	return orderedExpression{sortExpression: exp, direction: direction, nullSortType: sortType}
}

func (oe orderedExpression) Clone() Expression {
	return NewOrderedExpression(oe.sortExpression, oe.direction, oe.nullSortType)
}

func (oe orderedExpression) Expression() Expression {
	return oe
}

func (oe orderedExpression) SortExpression() Expression {
	return oe.sortExpression
}

func (oe orderedExpression) IsAsc() bool {
	return oe.direction == AscDir
}

func (oe orderedExpression) NullSortType() NullSortType {
	return oe.nullSortType
}

func (oe orderedExpression) NullsFirst() OrderedExpression {
	return NewOrderedExpression(oe.sortExpression, oe.direction, NullsFirstSortType)
}

func (oe orderedExpression) NullsLast() OrderedExpression {
	return NewOrderedExpression(oe.sortExpression, oe.direction, NullsLastSortType)
}
