package sqlgen

import (
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/doug-martin/goqu/v9/internal/errors"
	"github.com/doug-martin/goqu/v9/internal/sb"
)

var ErrNoUpdatedValuesProvided = errors.New("no update values provided")

func ErrCTENotSupported(dialect string) error {
	return errors.New("dialect does not support CTE WITH clause [dialect=%s]", dialect)
}

func ErrRecursiveCTENotSupported(dialect string) error {
	return errors.New("dialect does not support CTE WITH RECURSIVE clause [dialect=%s]", dialect)
}

func ErrReturnNotSupported(dialect string) error {
	return errors.New("dialect does not support RETURNING clause [dialect=%s]", dialect)
}

func ErrNotSupportedFragment(sqlType string, f SQLFragmentType) error {
	return errors.New("unsupported %s SQL fragment %s", sqlType, f)
}

type (
	CommonSQLGenerator interface {
		Dialect() string
		DialectOptions() *SQLDialectOptions
		ExpressionSQLGenerator() ExpressionSQLGenerator
		ReturningSQL(b sb.SQLBuilder, returns exp.ColumnListExpression)
		FromSQL(b sb.SQLBuilder, from exp.ColumnListExpression)
		SourcesSQL(b sb.SQLBuilder, from exp.ColumnListExpression)
		WhereSQL(b sb.SQLBuilder, where exp.ExpressionList)
		OrderSQL(b sb.SQLBuilder, order exp.ColumnListExpression)
		OrderWithOffsetFetchSQL(b sb.SQLBuilder, order exp.ColumnListExpression, offset uint, limit interface{})
		LimitSQL(b sb.SQLBuilder, limit interface{})
		UpdateExpressionSQL(b sb.SQLBuilder, updates ...exp.UpdateExpression)
	}
	commonSQLGenerator struct {
		dialect        string
		esg            ExpressionSQLGenerator
		dialectOptions *SQLDialectOptions
	}
)

func NewCommonSQLGenerator(dialect string, do *SQLDialectOptions) CommonSQLGenerator {
	return &commonSQLGenerator{dialect: dialect, esg: NewExpressionSQLGenerator(dialect, do), dialectOptions: do}
}

func (csg *commonSQLGenerator) Dialect() string {
	return csg.dialect
}

func (csg *commonSQLGenerator) DialectOptions() *SQLDialectOptions {
	return csg.dialectOptions
}

func (csg *commonSQLGenerator) ExpressionSQLGenerator() ExpressionSQLGenerator {
	return csg.esg
}

func (csg *commonSQLGenerator) ReturningSQL(b sb.SQLBuilder, returns exp.ColumnListExpression) {
	if returns != nil && len(returns.Columns()) > 0 {
		if csg.dialectOptions.SupportsReturn {
			b.Write(csg.dialectOptions.ReturningFragment)
			csg.esg.Generate(b, returns)
		} else {
			b.SetError(ErrReturnNotSupported(csg.dialect))
		}
	}
}

// Adds the FROM clause and tables to an sql statement
func (csg *commonSQLGenerator) FromSQL(b sb.SQLBuilder, from exp.ColumnListExpression) {
	if from != nil && !from.IsEmpty() {
		b.Write(csg.dialectOptions.FromFragment)
		csg.SourcesSQL(b, from)
	}
}

// Adds the generates the SQL for a column list
func (csg *commonSQLGenerator) SourcesSQL(b sb.SQLBuilder, from exp.ColumnListExpression) {
	b.WriteRunes(csg.dialectOptions.SpaceRune)
	csg.esg.Generate(b, from)
}

// Generates the WHERE clause for an SQL statement
func (csg *commonSQLGenerator) WhereSQL(b sb.SQLBuilder, where exp.ExpressionList) {
	if where != nil && !where.IsEmpty() {
		b.Write(csg.dialectOptions.WhereFragment)
		csg.esg.Generate(b, where)
	}
}

// Generates the ORDER BY clause for an SQL statement
func (csg *commonSQLGenerator) OrderSQL(b sb.SQLBuilder, order exp.ColumnListExpression) {
	if order != nil && len(order.Columns()) > 0 {
		b.Write(csg.dialectOptions.OrderByFragment)
		csg.esg.Generate(b, order)
	}
}

func (csg *commonSQLGenerator) OrderWithOffsetFetchSQL(
	b sb.SQLBuilder,
	order exp.ColumnListExpression,
	offset uint,
	limit interface{},
) {
	if order == nil {
		return
	}

	csg.OrderSQL(b, order)
	if offset > 0 {
		b.Write(csg.dialectOptions.OffsetFragment)
		csg.esg.Generate(b, offset)
		b.Write([]byte(" ROWS"))

		if limit != nil {
			b.Write(csg.dialectOptions.FetchFragment)
			csg.esg.Generate(b, limit)
			b.Write([]byte(" ROWS ONLY"))
		}
	}
}

// Generates the LIMIT clause for an SQL statement
func (csg *commonSQLGenerator) LimitSQL(b sb.SQLBuilder, limit interface{}) {
	if limit != nil {
		b.Write(csg.dialectOptions.LimitFragment)
		if csg.dialectOptions.SurroundLimitWithParentheses {
			b.WriteRunes(csg.dialectOptions.LeftParenRune)
		}
		csg.esg.Generate(b, limit)
		if csg.dialectOptions.SurroundLimitWithParentheses {
			b.WriteRunes(csg.dialectOptions.RightParenRune)
		}
	}
}

func (csg *commonSQLGenerator) UpdateExpressionSQL(b sb.SQLBuilder, updates ...exp.UpdateExpression) {
	if len(updates) == 0 {
		b.SetError(ErrNoUpdatedValuesProvided)
		return
	}
	updateLen := len(updates)
	for i, update := range updates {
		csg.esg.Generate(b, update)
		if i < updateLen-1 {
			b.WriteRunes(csg.dialectOptions.CommaRune)
		}
	}
}
