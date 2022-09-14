package sqlgen

import (
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/doug-martin/goqu/v9/internal/errors"
	"github.com/doug-martin/goqu/v9/internal/sb"
)

type (
	// An adapter interface to be used by a Dataset to generate SQL for a specific dialect.
	// See DefaultAdapter for a concrete implementation and examples.
	SelectSQLGenerator interface {
		Dialect() string
		Generate(b sb.SQLBuilder, clauses exp.SelectClauses)
	}
	// The default adapter. This class should be used when building a new adapter. When creating a new adapter you can
	// either override methods, or more typically update default values.
	// See (github.com/doug-martin/goqu/dialect/postgres)
	selectSQLGenerator struct {
		CommonSQLGenerator
	}
)

func ErrNotSupportedJoinType(j exp.JoinExpression) error {
	return errors.New("dialect does not support %v", j.JoinType())
}

func ErrJoinConditionRequired(j exp.JoinExpression) error {
	return errors.New("join condition required for conditioned join %v", j.JoinType())
}

func ErrDistinctOnNotSupported(dialect string) error {
	return errors.New("dialect does not support DISTINCT ON clause [dialect=%s]", dialect)
}

func ErrWindowNotSupported(dialect string) error {
	return errors.New("dialect does not support WINDOW clause [dialect=%s]", dialect)
}

var ErrNoWindowName = errors.New("window expresion has no valid name")

func NewSelectSQLGenerator(dialect string, do *SQLDialectOptions) SelectSQLGenerator {
	return &selectSQLGenerator{NewCommonSQLGenerator(dialect, do)}
}

func (ssg *selectSQLGenerator) Generate(b sb.SQLBuilder, clauses exp.SelectClauses) {
	for _, f := range ssg.DialectOptions().SelectSQLOrder {
		if b.Error() != nil {
			return
		}
		switch f {
		case CommonTableSQLFragment:
			ssg.ExpressionSQLGenerator().Generate(b, clauses.CommonTables())
		case SelectSQLFragment:
			ssg.SelectSQL(b, clauses)
		case SelectWithLimitSQLFragment:
			ssg.SelectWithLimitSQL(b, clauses)
		case FromSQLFragment:
			ssg.FromSQL(b, clauses.From())
		case JoinSQLFragment:
			ssg.JoinSQL(b, clauses.Joins())
		case WhereSQLFragment:
			ssg.WhereSQL(b, clauses.Where())
		case GroupBySQLFragment:
			ssg.GroupBySQL(b, clauses.GroupBy())
		case HavingSQLFragment:
			ssg.HavingSQL(b, clauses.Having())
		case WindowSQLFragment:
			ssg.WindowSQL(b, clauses.Windows())
		case CompoundsSQLFragment:
			ssg.CompoundsSQL(b, clauses.Compounds())
		case OrderSQLFragment:
			ssg.OrderSQL(b, clauses.Order())
		case OrderWithOffsetFetchSQLFragment:
			ssg.OrderWithOffsetFetchSQL(b, clauses.Order(), clauses.Offset(), clauses.Limit())
		case LimitSQLFragment:
			ssg.LimitSQL(b, clauses.Limit())
		case OffsetSQLFragment:
			ssg.OffsetSQL(b, clauses.Offset())
		case ForSQLFragment:
			ssg.ForSQL(b, clauses.Lock())
		default:
			b.SetError(ErrNotSupportedFragment("SELECT", f))
		}
	}
}

func (ssg *selectSQLGenerator) selectSQLCommon(b sb.SQLBuilder, clauses exp.SelectClauses) {
	dc := clauses.Distinct()
	if dc != nil {
		b.Write(ssg.DialectOptions().DistinctFragment)
		if !dc.IsEmpty() {
			if ssg.DialectOptions().SupportsDistinctOn {
				b.Write(ssg.DialectOptions().OnFragment).WriteRunes(ssg.DialectOptions().LeftParenRune)
				ssg.ExpressionSQLGenerator().Generate(b, dc)
				b.WriteRunes(ssg.DialectOptions().RightParenRune, ssg.DialectOptions().SpaceRune)
			} else {
				b.SetError(ErrDistinctOnNotSupported(ssg.Dialect()))
				return
			}
		} else {
			b.WriteRunes(ssg.DialectOptions().SpaceRune)
		}
	}

	if cols := clauses.Select(); clauses.IsDefaultSelect() || len(cols.Columns()) == 0 {
		b.WriteRunes(ssg.DialectOptions().StarRune)
	} else {
		ssg.ExpressionSQLGenerator().Generate(b, cols)
	}
}

// Adds the SELECT clause and columns to a sql statement
func (ssg *selectSQLGenerator) SelectSQL(b sb.SQLBuilder, clauses exp.SelectClauses) {
	b.Write(ssg.DialectOptions().SelectClause).WriteRunes(ssg.DialectOptions().SpaceRune)
	ssg.selectSQLCommon(b, clauses)
}

// Adds the SELECT clause along with LIMIT to a SQL statement (e.g. MSSQL dialect: SELECT TOP 10 ...)
func (ssg *selectSQLGenerator) SelectWithLimitSQL(b sb.SQLBuilder, clauses exp.SelectClauses) {
	b.Write(ssg.DialectOptions().SelectClause).WriteRunes(ssg.DialectOptions().SpaceRune)
	if clauses.Offset() == 0 && clauses.Limit() != nil {
		ssg.LimitSQL(b, clauses.Limit())
		b.WriteRunes(ssg.DialectOptions().SpaceRune)
	}
	ssg.selectSQLCommon(b, clauses)
}

// Generates the JOIN clauses for an SQL statement
func (ssg *selectSQLGenerator) JoinSQL(b sb.SQLBuilder, joins exp.JoinExpressions) {
	if len(joins) > 0 {
		for _, j := range joins {
			joinType, ok := ssg.DialectOptions().JoinTypeLookup[j.JoinType()]
			if !ok {
				b.SetError(ErrNotSupportedJoinType(j))
				return
			}
			b.Write(joinType)
			ssg.ExpressionSQLGenerator().Generate(b, j.Table())
			if t, ok := j.(exp.ConditionedJoinExpression); ok {
				if t.IsConditionEmpty() {
					b.SetError(ErrJoinConditionRequired(j))
					return
				}
				ssg.joinConditionSQL(b, t.Condition())
			}
		}
	}
}

// Generates the GROUP BY clause for an SQL statement
func (ssg *selectSQLGenerator) GroupBySQL(b sb.SQLBuilder, groupBy exp.ColumnListExpression) {
	if groupBy != nil && len(groupBy.Columns()) > 0 {
		b.Write(ssg.DialectOptions().GroupByFragment)
		ssg.ExpressionSQLGenerator().Generate(b, groupBy)
	}
}

// Generates the HAVING clause for an SQL statement
func (ssg *selectSQLGenerator) HavingSQL(b sb.SQLBuilder, having exp.ExpressionList) {
	if having != nil && len(having.Expressions()) > 0 {
		b.Write(ssg.DialectOptions().HavingFragment)
		ssg.ExpressionSQLGenerator().Generate(b, having)
	}
}

// Generates the OFFSET clause for an SQL statement
func (ssg *selectSQLGenerator) OffsetSQL(b sb.SQLBuilder, offset uint) {
	if offset > 0 {
		b.Write(ssg.DialectOptions().OffsetFragment)
		ssg.ExpressionSQLGenerator().Generate(b, offset)
	}
}

// Generates the compound sql clause for an SQL statement (e.g. UNION, INTERSECT)
func (ssg *selectSQLGenerator) CompoundsSQL(b sb.SQLBuilder, compounds []exp.CompoundExpression) {
	for _, compound := range compounds {
		ssg.ExpressionSQLGenerator().Generate(b, compound)
	}
}

// Generates the FOR (aka "locking") clause for an SQL statement
func (ssg *selectSQLGenerator) ForSQL(b sb.SQLBuilder, lockingClause exp.Lock) {
	if lockingClause == nil {
		return
	}
	switch lockingClause.Strength() {
	case exp.ForNolock:
		return
	case exp.ForUpdate:
		b.Write(ssg.DialectOptions().ForUpdateFragment)
	case exp.ForNoKeyUpdate:
		b.Write(ssg.DialectOptions().ForNoKeyUpdateFragment)
	case exp.ForShare:
		b.Write(ssg.DialectOptions().ForShareFragment)
	case exp.ForKeyShare:
		b.Write(ssg.DialectOptions().ForKeyShareFragment)
	}

	of := lockingClause.Of()
	if ofLen := len(of); ofLen > 0 {
		if ofFragment := ssg.DialectOptions().OfFragment; len(ofFragment) > 0 {
			b.Write(ofFragment)
			for i, table := range of {
				ssg.ExpressionSQLGenerator().Generate(b, table)
				if i < ofLen-1 {
					b.WriteRunes(ssg.DialectOptions().CommaRune, ssg.DialectOptions().SpaceRune)
				}
			}
			b.WriteRunes(ssg.DialectOptions().SpaceRune)
		}
	}

	// the WAIT case is the default in Postgres, and is what you get if you don't specify NOWAIT or
	// SKIP LOCKED. There's no special syntax for it in PG, so we don't do anything for it here
	switch lockingClause.WaitOption() {
	case exp.Wait:
		return
	case exp.NoWait:
		b.Write(ssg.DialectOptions().NowaitFragment)
	case exp.SkipLocked:
		b.Write(ssg.DialectOptions().SkipLockedFragment)
	}
}

func (ssg *selectSQLGenerator) WindowSQL(b sb.SQLBuilder, windows []exp.WindowExpression) {
	weLen := len(windows)
	if weLen == 0 {
		return
	}
	if !ssg.DialectOptions().SupportsWindowFunction {
		b.SetError(ErrWindowNotSupported(ssg.Dialect()))
		return
	}
	b.Write(ssg.DialectOptions().WindowFragment)
	for i, we := range windows {
		if !we.HasName() {
			b.SetError(ErrNoWindowName)
		}
		ssg.ExpressionSQLGenerator().Generate(b, we)
		if i < weLen-1 {
			b.WriteRunes(ssg.DialectOptions().CommaRune, ssg.DialectOptions().SpaceRune)
		}
	}
}

func (ssg *selectSQLGenerator) joinConditionSQL(b sb.SQLBuilder, jc exp.JoinCondition) {
	switch t := jc.(type) {
	case exp.JoinOnCondition:
		ssg.joinOnConditionSQL(b, t)
	case exp.JoinUsingCondition:
		ssg.joinUsingConditionSQL(b, t)
	}
}

func (ssg *selectSQLGenerator) joinUsingConditionSQL(b sb.SQLBuilder, jc exp.JoinUsingCondition) {
	b.Write(ssg.DialectOptions().UsingFragment).
		WriteRunes(ssg.DialectOptions().LeftParenRune)
	ssg.ExpressionSQLGenerator().Generate(b, jc.Using())
	b.WriteRunes(ssg.DialectOptions().RightParenRune)
}

func (ssg *selectSQLGenerator) joinOnConditionSQL(b sb.SQLBuilder, jc exp.JoinOnCondition) {
	b.Write(ssg.DialectOptions().OnFragment)
	ssg.ExpressionSQLGenerator().Generate(b, jc.On())
}
