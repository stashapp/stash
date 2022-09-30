package goqu

import (
	"github.com/doug-martin/goqu/v9/exec"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/doug-martin/goqu/v9/internal/sb"
)

type TruncateDataset struct {
	dialect      SQLDialect
	clauses      exp.TruncateClauses
	isPrepared   prepared
	queryFactory exec.QueryFactory
	err          error
}

// used internally by database to create a database with a specific adapter
func newTruncateDataset(d string, queryFactory exec.QueryFactory) *TruncateDataset {
	return &TruncateDataset{
		clauses:      exp.NewTruncateClauses(),
		dialect:      GetDialect(d),
		queryFactory: queryFactory,
	}
}

func Truncate(table ...interface{}) *TruncateDataset {
	return newTruncateDataset("default", nil).Table(table...)
}

// Sets the adapter used to serialize values and create the SQL statement
func (td *TruncateDataset) WithDialect(dl string) *TruncateDataset {
	ds := td.copy(td.GetClauses())
	ds.dialect = GetDialect(dl)
	return ds
}

// Set the parameter interpolation behavior. See examples
//
// prepared: If true the dataset WILL NOT interpolate the parameters.
func (td *TruncateDataset) Prepared(prepared bool) *TruncateDataset {
	ret := td.copy(td.clauses)
	ret.isPrepared = preparedFromBool(prepared)
	return ret
}

func (td *TruncateDataset) IsPrepared() bool {
	return td.isPrepared.Bool()
}

// Returns the current adapter on the dataset
func (td *TruncateDataset) Dialect() SQLDialect {
	return td.dialect
}

// Returns the current adapter on the dataset
func (td *TruncateDataset) SetDialect(dialect SQLDialect) *TruncateDataset {
	cd := td.copy(td.GetClauses())
	cd.dialect = dialect
	return cd
}

func (td *TruncateDataset) Expression() exp.Expression {
	return td
}

// Clones the dataset
func (td *TruncateDataset) Clone() exp.Expression {
	return td.copy(td.clauses)
}

// Returns the current clauses on the dataset.
func (td *TruncateDataset) GetClauses() exp.TruncateClauses {
	return td.clauses
}

// used interally to copy the dataset
func (td *TruncateDataset) copy(clauses exp.TruncateClauses) *TruncateDataset {
	return &TruncateDataset{
		dialect:      td.dialect,
		clauses:      clauses,
		isPrepared:   td.isPrepared,
		queryFactory: td.queryFactory,
		err:          td.err,
	}
}

// Adds a FROM clause. This return a new dataset with the original sources replaced. See examples.
// You can pass in the following.
//   string: Will automatically be turned into an identifier
//   IdentifierExpression
//   LiteralExpression: (See Literal) Will use the literal SQL
func (td *TruncateDataset) Table(table ...interface{}) *TruncateDataset {
	return td.copy(td.clauses.SetTable(exp.NewColumnListExpression(table...)))
}

// Adds a CASCADE clause
func (td *TruncateDataset) Cascade() *TruncateDataset {
	opts := td.clauses.Options()
	opts.Cascade = true
	return td.copy(td.clauses.SetOptions(opts))
}

// Clears the CASCADE clause
func (td *TruncateDataset) NoCascade() *TruncateDataset {
	opts := td.clauses.Options()
	opts.Cascade = false
	return td.copy(td.clauses.SetOptions(opts))
}

// Adds a RESTRICT clause
func (td *TruncateDataset) Restrict() *TruncateDataset {
	opts := td.clauses.Options()
	opts.Restrict = true
	return td.copy(td.clauses.SetOptions(opts))
}

// Clears the RESTRICT clause
func (td *TruncateDataset) NoRestrict() *TruncateDataset {
	opts := td.clauses.Options()
	opts.Restrict = false
	return td.copy(td.clauses.SetOptions(opts))
}

// Add a IDENTITY clause (e.g. RESTART)
func (td *TruncateDataset) Identity(identity string) *TruncateDataset {
	opts := td.clauses.Options()
	opts.Identity = identity
	return td.copy(td.clauses.SetOptions(opts))
}

// Get any error that has been set or nil if no error has been set.
func (td *TruncateDataset) Error() error {
	return td.err
}

// Set an error on the dataset if one has not already been set. This error will be returned by a future call to Error
// or as part of ToSQL. This can be used by end users to record errors while building up queries without having to
// track those separately.
func (td *TruncateDataset) SetError(err error) *TruncateDataset {
	if td.err == nil {
		td.err = err
	}

	return td
}

// Generates a TRUNCATE sql statement, if Prepared has been called with true then the parameters will not be interpolated.
// See examples.
//
// Errors:
//  * There is an error generating the SQL
func (td *TruncateDataset) ToSQL() (sql string, params []interface{}, err error) {
	return td.truncateSQLBuilder().ToSQL()
}

// Generates the TRUNCATE sql, and returns an Exec struct with the sql set to the TRUNCATE statement
//    db.From("test").Truncate().Executor().Exec()
func (td *TruncateDataset) Executor() exec.QueryExecutor {
	return td.queryFactory.FromSQLBuilder(td.truncateSQLBuilder())
}

func (td *TruncateDataset) truncateSQLBuilder() sb.SQLBuilder {
	buf := sb.NewSQLBuilder(td.isPrepared.Bool())
	if td.err != nil {
		return buf.SetError(td.err)
	}
	td.dialect.ToTruncateSQL(buf, td.clauses)
	return buf
}
