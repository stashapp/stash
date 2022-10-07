package exp

type (
	doNothingConflict struct{}
	// ConflictUpdate is the struct that represents the UPDATE fragment of an
	// INSERT ... ON CONFLICT/ON DUPLICATE KEY DO UPDATE statement
	conflictUpdate struct {
		target      string
		update      interface{}
		whereClause ExpressionList
	}
)

// Creates a conflict struct to be passed to InsertConflict to ignore constraint errors
//  InsertConflict(DoNothing(),...) -> INSERT INTO ... ON CONFLICT DO NOTHING
func NewDoNothingConflictExpression() ConflictExpression {
	return &doNothingConflict{}
}

func (c doNothingConflict) Expression() Expression {
	return c
}

func (c doNothingConflict) Clone() Expression {
	return c
}

func (c doNothingConflict) Action() ConflictAction {
	return DoNothingConflictAction
}

// Creates a ConflictUpdate struct to be passed to InsertConflict
// Represents a ON CONFLICT DO UPDATE portion of an INSERT statement (ON DUPLICATE KEY UPDATE for mysql)
//
//  InsertConflict(DoUpdate("target_column", update),...) ->
//  	INSERT INTO ... ON CONFLICT DO UPDATE SET a=b
//  InsertConflict(DoUpdate("target_column", update).Where(Ex{"a": 1},...) ->
//  	INSERT INTO ... ON CONFLICT DO UPDATE SET a=b WHERE a=1
func NewDoUpdateConflictExpression(target string, update interface{}) ConflictUpdateExpression {
	return &conflictUpdate{target: target, update: update}
}

func (c conflictUpdate) Expression() Expression {
	return c
}

func (c conflictUpdate) Clone() Expression {
	return &conflictUpdate{
		target:      c.target,
		update:      c.update,
		whereClause: c.whereClause.Clone().(ExpressionList),
	}
}

func (c conflictUpdate) Action() ConflictAction {
	return DoUpdateConflictAction
}

// Returns the target conflict column. Only necessary for Postgres.
// Will return an error for mysql/sqlite. Will also return an error if missing from a postgres ConflictUpdate.
func (c conflictUpdate) TargetColumn() string {
	return c.target
}

// Returns the Updates which represent the ON CONFLICT DO UPDATE portion of an insert statement. If nil,
// there are no updates.
func (c conflictUpdate) Update() interface{} {
	return c.update
}

// Append to the existing Where clause for an ON CONFLICT DO UPDATE ... WHERE ...
//  InsertConflict(DoNothing(),...) -> INSERT INTO ... ON CONFLICT DO NOTHING
func (c *conflictUpdate) Where(expressions ...Expression) ConflictUpdateExpression {
	if c.whereClause == nil {
		c.whereClause = NewExpressionList(AndType, expressions...)
	} else {
		c.whereClause = c.whereClause.Append(expressions...)
	}
	return c
}

// Append to the existing Where clause for an ON CONFLICT DO UPDATE ... WHERE ...
//  InsertConflict(DoNothing(),...) -> INSERT INTO ... ON CONFLICT DO NOTHING
func (c *conflictUpdate) WhereClause() ExpressionList {
	return c.whereClause
}
