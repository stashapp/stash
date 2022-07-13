package sqlgen

import (
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/doug-martin/goqu/v9/internal/errors"
	"github.com/doug-martin/goqu/v9/internal/sb"
)

type (
	// An adapter interface to be used by a Dataset to generate SQL for a specific dialect.
	// See DefaultAdapter for a concrete implementation and examples.
	DeleteSQLGenerator interface {
		Dialect() string
		Generate(b sb.SQLBuilder, clauses exp.DeleteClauses)
	}
	// The default adapter. This class should be used when building a new adapter. When creating a new adapter you can
	// either override methods, or more typically update default values.
	// See (github.com/doug-martin/goqu/dialect/postgres)
	deleteSQLGenerator struct {
		CommonSQLGenerator
	}
)

var ErrNoSourceForDelete = errors.New("no source found when generating delete sql")

func NewDeleteSQLGenerator(dialect string, do *SQLDialectOptions) DeleteSQLGenerator {
	return &deleteSQLGenerator{NewCommonSQLGenerator(dialect, do)}
}

func (dsg *deleteSQLGenerator) Generate(b sb.SQLBuilder, clauses exp.DeleteClauses) {
	if !clauses.HasFrom() {
		b.SetError(ErrNoSourceForDelete)
		return
	}
	for _, f := range dsg.DialectOptions().DeleteSQLOrder {
		if b.Error() != nil {
			return
		}
		switch f {
		case CommonTableSQLFragment:
			dsg.ExpressionSQLGenerator().Generate(b, clauses.CommonTables())
		case DeleteBeginSQLFragment:
			dsg.DeleteBeginSQL(
				b, exp.NewColumnListExpression(clauses.From()), !(clauses.HasLimit() || clauses.HasOrder()),
			)
		case FromSQLFragment:
			dsg.FromSQL(b, exp.NewColumnListExpression(clauses.From()))
		case WhereSQLFragment:
			dsg.WhereSQL(b, clauses.Where())
		case OrderSQLFragment:
			if dsg.DialectOptions().SupportsOrderByOnDelete {
				dsg.OrderSQL(b, clauses.Order())
			}
		case LimitSQLFragment:
			if dsg.DialectOptions().SupportsLimitOnDelete {
				dsg.LimitSQL(b, clauses.Limit())
			}
		case ReturningSQLFragment:
			dsg.ReturningSQL(b, clauses.Returning())
		default:
			b.SetError(ErrNotSupportedFragment("DELETE", f))
		}
	}
}

// Adds the correct fragment to being an DELETE statement
func (dsg *deleteSQLGenerator) DeleteBeginSQL(b sb.SQLBuilder, from exp.ColumnListExpression, multiTable bool) {
	b.Write(dsg.DialectOptions().DeleteClause)
	if multiTable && dsg.DialectOptions().SupportsDeleteTableHint {
		dsg.SourcesSQL(b, from)
	}
}
