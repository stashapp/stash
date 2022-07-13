package exp

type sqlWindowExpression struct {
	name          IdentifierExpression
	parent        IdentifierExpression
	partitionCols ColumnListExpression
	orderCols     ColumnListExpression
}

func NewWindowExpression(window, parent IdentifierExpression, partitionCols, orderCols ColumnListExpression) WindowExpression {
	if partitionCols == nil {
		partitionCols = NewColumnListExpression()
	}
	if orderCols == nil {
		orderCols = NewColumnListExpression()
	}
	return sqlWindowExpression{
		name:          window,
		parent:        parent,
		partitionCols: partitionCols,
		orderCols:     orderCols,
	}
}

func (we sqlWindowExpression) clone() sqlWindowExpression {
	return sqlWindowExpression{
		name:          we.name,
		parent:        we.parent,
		partitionCols: we.partitionCols.Clone().(ColumnListExpression),
		orderCols:     we.orderCols.Clone().(ColumnListExpression),
	}
}

func (we sqlWindowExpression) Clone() Expression {
	return we.clone()
}

func (we sqlWindowExpression) Expression() Expression {
	return we
}

func (we sqlWindowExpression) Name() IdentifierExpression {
	return we.name
}

func (we sqlWindowExpression) HasName() bool {
	return we.name != nil
}

func (we sqlWindowExpression) Parent() IdentifierExpression {
	return we.parent
}

func (we sqlWindowExpression) HasParent() bool {
	return we.parent != nil
}

func (we sqlWindowExpression) PartitionCols() ColumnListExpression {
	return we.partitionCols
}

func (we sqlWindowExpression) HasPartitionBy() bool {
	return we.partitionCols != nil && !we.partitionCols.IsEmpty()
}

func (we sqlWindowExpression) OrderCols() ColumnListExpression {
	return we.orderCols
}

func (we sqlWindowExpression) HasOrder() bool {
	return we.orderCols != nil && !we.orderCols.IsEmpty()
}

func (we sqlWindowExpression) PartitionBy(cols ...interface{}) WindowExpression {
	ret := we.clone()
	ret.partitionCols = NewColumnListExpression(cols...)
	return ret
}

func (we sqlWindowExpression) OrderBy(cols ...interface{}) WindowExpression {
	ret := we.clone()
	ret.orderCols = NewColumnListExpression(cols...)
	return ret
}

func (we sqlWindowExpression) Inherit(parent string) WindowExpression {
	ret := we.clone()
	ret.parent = ParseIdentifier(parent)
	return ret
}
