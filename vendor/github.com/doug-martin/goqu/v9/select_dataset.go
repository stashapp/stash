package goqu

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9/exec"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/doug-martin/goqu/v9/internal/errors"
	"github.com/doug-martin/goqu/v9/internal/sb"
)

// Dataset for creating and/or executing SELECT SQL statements.
type SelectDataset struct {
	dialect      SQLDialect
	clauses      exp.SelectClauses
	isPrepared   prepared
	queryFactory exec.QueryFactory
	err          error
}

var ErrQueryFactoryNotFoundError = errors.New(
	"unable to execute query did you use goqu.Database#From to create the dataset",
)

// used internally by database to create a database with a specific adapter
func newDataset(d string, queryFactory exec.QueryFactory) *SelectDataset {
	return &SelectDataset{
		clauses:      exp.NewSelectClauses(),
		dialect:      GetDialect(d),
		queryFactory: queryFactory,
	}
}

func From(table ...interface{}) *SelectDataset {
	return newDataset("default", nil).From(table...)
}

func Select(cols ...interface{}) *SelectDataset {
	return newDataset("default", nil).Select(cols...)
}

// Sets the adapter used to serialize values and create the SQL statement
func (sd *SelectDataset) WithDialect(dl string) *SelectDataset {
	ds := sd.copy(sd.GetClauses())
	ds.dialect = GetDialect(dl)
	return ds
}

// Set the parameter interpolation behavior. See examples
//
// prepared: If true the dataset WILL NOT interpolate the parameters.
func (sd *SelectDataset) Prepared(prepared bool) *SelectDataset {
	ret := sd.copy(sd.clauses)
	ret.isPrepared = preparedFromBool(prepared)
	return ret
}

func (sd *SelectDataset) IsPrepared() bool {
	return sd.isPrepared.Bool()
}

// Returns the current adapter on the dataset
func (sd *SelectDataset) Dialect() SQLDialect {
	return sd.dialect
}

// Returns the current adapter on the dataset
func (sd *SelectDataset) SetDialect(dialect SQLDialect) *SelectDataset {
	cd := sd.copy(sd.GetClauses())
	cd.dialect = dialect
	return cd
}

func (sd *SelectDataset) Expression() exp.Expression {
	return sd
}

// Clones the dataset
func (sd *SelectDataset) Clone() exp.Expression {
	return sd.copy(sd.clauses)
}

// Returns the current clauses on the dataset.
func (sd *SelectDataset) GetClauses() exp.SelectClauses {
	return sd.clauses
}

// used interally to copy the dataset
func (sd *SelectDataset) copy(clauses exp.SelectClauses) *SelectDataset {
	return &SelectDataset{
		dialect:      sd.dialect,
		clauses:      clauses,
		isPrepared:   sd.isPrepared,
		queryFactory: sd.queryFactory,
		err:          sd.err,
	}
}

// Creates a new UpdateDataset using the FROM of this dataset. This method will also copy over the `WITH`, `WHERE`,
// `ORDER , and `LIMIT`
func (sd *SelectDataset) Update() *UpdateDataset {
	u := newUpdateDataset(sd.dialect.Dialect(), sd.queryFactory).
		Prepared(sd.isPrepared.Bool())
	if sd.clauses.HasSources() {
		u = u.Table(sd.GetClauses().From().Columns()[0])
	}
	c := u.clauses
	for _, ce := range sd.clauses.CommonTables() {
		c = c.CommonTablesAppend(ce)
	}
	if sd.clauses.Where() != nil {
		c = c.WhereAppend(sd.clauses.Where())
	}
	if sd.clauses.HasLimit() {
		c = c.SetLimit(sd.clauses.Limit())
	}
	if sd.clauses.HasOrder() {
		for _, oe := range sd.clauses.Order().Columns() {
			c = c.OrderAppend(oe.(exp.OrderedExpression))
		}
	}
	u.clauses = c
	return u
}

// Creates a new InsertDataset using the FROM of this dataset. This method will also copy over the `WITH` clause to the
// insert.
func (sd *SelectDataset) Insert() *InsertDataset {
	i := newInsertDataset(sd.dialect.Dialect(), sd.queryFactory).
		Prepared(sd.isPrepared.Bool())
	if sd.clauses.HasSources() {
		i = i.Into(sd.GetClauses().From().Columns()[0])
	}
	c := i.clauses
	for _, ce := range sd.clauses.CommonTables() {
		c = c.CommonTablesAppend(ce)
	}
	i.clauses = c
	return i
}

// Creates a new DeleteDataset using the FROM of this dataset. This method will also copy over the `WITH`, `WHERE`,
// `ORDER , and `LIMIT`
func (sd *SelectDataset) Delete() *DeleteDataset {
	d := newDeleteDataset(sd.dialect.Dialect(), sd.queryFactory).
		Prepared(sd.isPrepared.Bool())
	if sd.clauses.HasSources() {
		d = d.From(sd.clauses.From().Columns()[0])
	}
	c := d.clauses
	for _, ce := range sd.clauses.CommonTables() {
		c = c.CommonTablesAppend(ce)
	}
	if sd.clauses.Where() != nil {
		c = c.WhereAppend(sd.clauses.Where())
	}
	if sd.clauses.HasLimit() {
		c = c.SetLimit(sd.clauses.Limit())
	}
	if sd.clauses.HasOrder() {
		for _, oe := range sd.clauses.Order().Columns() {
			c = c.OrderAppend(oe.(exp.OrderedExpression))
		}
	}
	d.clauses = c
	return d
}

// Creates a new TruncateDataset using the FROM of this dataset.
func (sd *SelectDataset) Truncate() *TruncateDataset {
	td := newTruncateDataset(sd.dialect.Dialect(), sd.queryFactory)
	if sd.clauses.HasSources() {
		td = td.Table(sd.clauses.From())
	}
	return td
}

// Creates a WITH clause for a common table expression (CTE).
//
// The name will be available to SELECT from in the associated query; and can optionally
// contain a list of column names "name(col1, col2, col3)".
//
// The name will refer to the results of the specified subquery.
func (sd *SelectDataset) With(name string, subquery exp.Expression) *SelectDataset {
	return sd.copy(sd.clauses.CommonTablesAppend(exp.NewCommonTableExpression(false, name, subquery)))
}

// Creates a WITH RECURSIVE clause for a common table expression (CTE)
//
// The name will be available to SELECT from in the associated query; and must
// contain a list of column names "name(col1, col2, col3)" for a recursive clause.
//
// The name will refer to the results of the specified subquery. The subquery for
// a recursive query will always end with a UNION or UNION ALL with a clause that
// refers to the CTE by name.
func (sd *SelectDataset) WithRecursive(name string, subquery exp.Expression) *SelectDataset {
	return sd.copy(sd.clauses.CommonTablesAppend(exp.NewCommonTableExpression(true, name, subquery)))
}

// Adds columns to the SELECT clause. See examples
// You can pass in the following.
//   string: Will automatically be turned into an identifier
//   Dataset: Will use the SQL generated from that Dataset. If the dataset is aliased it will use that alias as the
//   column name.
//   LiteralExpression: (See Literal) Will use the literal SQL
//   SQLFunction: (See Func, MIN, MAX, COUNT....)
//   Struct: If passing in an instance of a struct, we will parse the struct for the column names to select.
//   See examples
func (sd *SelectDataset) Select(selects ...interface{}) *SelectDataset {
	if len(selects) == 0 {
		return sd.ClearSelect()
	}
	return sd.copy(sd.clauses.SetSelect(exp.NewColumnListExpression(selects...)))
}

// Adds columns to the SELECT DISTINCT clause. See examples
// You can pass in the following.
//   string: Will automatically be turned into an identifier
//   Dataset: Will use the SQL generated from that Dataset. If the dataset is aliased it will use that alias as the
//   column name.
//   LiteralExpression: (See Literal) Will use the literal SQL
//   SQLFunction: (See Func, MIN, MAX, COUNT....)
//   Struct: If passing in an instance of a struct, we will parse the struct for the column names to select.
//   See examples
// Deprecated: Use Distinct() instead.
func (sd *SelectDataset) SelectDistinct(selects ...interface{}) *SelectDataset {
	if len(selects) == 0 {
		cleared := sd.ClearSelect()
		return cleared.copy(cleared.clauses.SetDistinct(nil))
	}
	return sd.copy(sd.clauses.SetSelect(exp.NewColumnListExpression(selects...)).SetDistinct(exp.NewColumnListExpression()))
}

// Resets to SELECT *. If the SelectDistinct or Distinct was used the returned Dataset will have the the dataset set to SELECT *.
// See examples.
func (sd *SelectDataset) ClearSelect() *SelectDataset {
	return sd.copy(sd.clauses.SetSelect(exp.NewColumnListExpression(exp.Star())).SetDistinct(nil))
}

// Adds columns to the SELECT clause. See examples
// You can pass in the following.
//   string: Will automatically be turned into an identifier
//   Dataset: Will use the SQL generated from that Dataset. If the dataset is aliased it will use that alias as the
//   column name.
//   LiteralExpression: (See Literal) Will use the literal SQL
//   SQLFunction: (See Func, MIN, MAX, COUNT....)
func (sd *SelectDataset) SelectAppend(selects ...interface{}) *SelectDataset {
	return sd.copy(sd.clauses.SelectAppend(exp.NewColumnListExpression(selects...)))
}

func (sd *SelectDataset) Distinct(on ...interface{}) *SelectDataset {
	return sd.copy(sd.clauses.SetDistinct(exp.NewColumnListExpression(on...)))
}

// Adds a FROM clause. This return a new dataset with the original sources replaced. See examples.
// You can pass in the following.
//   string: Will automatically be turned into an identifier
//   Dataset: Will be added as a sub select. If the Dataset is not aliased it will automatically be aliased
//   LiteralExpression: (See Literal) Will use the literal SQL
func (sd *SelectDataset) From(from ...interface{}) *SelectDataset {
	var sources []interface{}
	numSources := 0
	for _, source := range from {
		if ds, ok := source.(*SelectDataset); ok && !ds.clauses.HasAlias() {
			numSources++
			sources = append(sources, ds.As(fmt.Sprintf("t%d", numSources)))
		} else {
			sources = append(sources, source)
		}
	}
	return sd.copy(sd.clauses.SetFrom(exp.NewColumnListExpression(sources...)))
}

// Returns a new Dataset with the current one as an source. If the current Dataset is not aliased (See Dataset#As) then
// it will automatically be aliased. See examples.
func (sd *SelectDataset) FromSelf() *SelectDataset {
	return sd.copy(exp.NewSelectClauses()).From(sd)
}

// Alias to InnerJoin. See examples.
func (sd *SelectDataset) Join(table exp.Expression, condition exp.JoinCondition) *SelectDataset {
	return sd.InnerJoin(table, condition)
}

// Adds an INNER JOIN clause. See examples.
func (sd *SelectDataset) InnerJoin(table exp.Expression, condition exp.JoinCondition) *SelectDataset {
	return sd.joinTable(exp.NewConditionedJoinExpression(exp.InnerJoinType, table, condition))
}

// Adds a FULL OUTER JOIN clause. See examples.
func (sd *SelectDataset) FullOuterJoin(table exp.Expression, condition exp.JoinCondition) *SelectDataset {
	return sd.joinTable(exp.NewConditionedJoinExpression(exp.FullOuterJoinType, table, condition))
}

// Adds a RIGHT OUTER JOIN clause. See examples.
func (sd *SelectDataset) RightOuterJoin(table exp.Expression, condition exp.JoinCondition) *SelectDataset {
	return sd.joinTable(exp.NewConditionedJoinExpression(exp.RightOuterJoinType, table, condition))
}

// Adds a LEFT OUTER JOIN clause. See examples.
func (sd *SelectDataset) LeftOuterJoin(table exp.Expression, condition exp.JoinCondition) *SelectDataset {
	return sd.joinTable(exp.NewConditionedJoinExpression(exp.LeftOuterJoinType, table, condition))
}

// Adds a FULL JOIN clause. See examples.
func (sd *SelectDataset) FullJoin(table exp.Expression, condition exp.JoinCondition) *SelectDataset {
	return sd.joinTable(exp.NewConditionedJoinExpression(exp.FullJoinType, table, condition))
}

// Adds a RIGHT JOIN clause. See examples.
func (sd *SelectDataset) RightJoin(table exp.Expression, condition exp.JoinCondition) *SelectDataset {
	return sd.joinTable(exp.NewConditionedJoinExpression(exp.RightJoinType, table, condition))
}

// Adds a LEFT JOIN clause. See examples.
func (sd *SelectDataset) LeftJoin(table exp.Expression, condition exp.JoinCondition) *SelectDataset {
	return sd.joinTable(exp.NewConditionedJoinExpression(exp.LeftJoinType, table, condition))
}

// Adds a NATURAL JOIN clause. See examples.
func (sd *SelectDataset) NaturalJoin(table exp.Expression) *SelectDataset {
	return sd.joinTable(exp.NewUnConditionedJoinExpression(exp.NaturalJoinType, table))
}

// Adds a NATURAL LEFT JOIN clause. See examples.
func (sd *SelectDataset) NaturalLeftJoin(table exp.Expression) *SelectDataset {
	return sd.joinTable(exp.NewUnConditionedJoinExpression(exp.NaturalLeftJoinType, table))
}

// Adds a NATURAL RIGHT JOIN clause. See examples.
func (sd *SelectDataset) NaturalRightJoin(table exp.Expression) *SelectDataset {
	return sd.joinTable(exp.NewUnConditionedJoinExpression(exp.NaturalRightJoinType, table))
}

// Adds a NATURAL FULL JOIN clause. See examples.
func (sd *SelectDataset) NaturalFullJoin(table exp.Expression) *SelectDataset {
	return sd.joinTable(exp.NewUnConditionedJoinExpression(exp.NaturalFullJoinType, table))
}

// Adds a CROSS JOIN clause. See examples.
func (sd *SelectDataset) CrossJoin(table exp.Expression) *SelectDataset {
	return sd.joinTable(exp.NewUnConditionedJoinExpression(exp.CrossJoinType, table))
}

// Joins this Datasets table with another
func (sd *SelectDataset) joinTable(join exp.JoinExpression) *SelectDataset {
	return sd.copy(sd.clauses.JoinsAppend(join))
}

// Adds a WHERE clause. See examples.
func (sd *SelectDataset) Where(expressions ...exp.Expression) *SelectDataset {
	return sd.copy(sd.clauses.WhereAppend(expressions...))
}

// Removes the WHERE clause. See examples.
func (sd *SelectDataset) ClearWhere() *SelectDataset {
	return sd.copy(sd.clauses.ClearWhere())
}

// Adds a FOR UPDATE clause. See examples.
func (sd *SelectDataset) ForUpdate(waitOption exp.WaitOption, of ...exp.IdentifierExpression) *SelectDataset {
	return sd.withLock(exp.ForUpdate, waitOption, of...)
}

// Adds a FOR NO KEY UPDATE clause. See examples.
func (sd *SelectDataset) ForNoKeyUpdate(waitOption exp.WaitOption, of ...exp.IdentifierExpression) *SelectDataset {
	return sd.withLock(exp.ForNoKeyUpdate, waitOption, of...)
}

// Adds a FOR KEY SHARE clause. See examples.
func (sd *SelectDataset) ForKeyShare(waitOption exp.WaitOption, of ...exp.IdentifierExpression) *SelectDataset {
	return sd.withLock(exp.ForKeyShare, waitOption, of...)
}

// Adds a FOR SHARE clause. See examples.
func (sd *SelectDataset) ForShare(waitOption exp.WaitOption, of ...exp.IdentifierExpression) *SelectDataset {
	return sd.withLock(exp.ForShare, waitOption, of...)
}

func (sd *SelectDataset) withLock(strength exp.LockStrength, option exp.WaitOption, of ...exp.IdentifierExpression) *SelectDataset {
	return sd.copy(sd.clauses.SetLock(exp.NewLock(strength, option, of...)))
}

// Adds a GROUP BY clause. See examples.
func (sd *SelectDataset) GroupBy(groupBy ...interface{}) *SelectDataset {
	return sd.copy(sd.clauses.SetGroupBy(exp.NewColumnListExpression(groupBy...)))
}

// Adds more columns to the current GROUP BY clause. See examples.
func (sd *SelectDataset) GroupByAppend(groupBy ...interface{}) *SelectDataset {
	return sd.copy(sd.clauses.GroupByAppend(exp.NewColumnListExpression(groupBy...)))
}

// Adds a HAVING clause. See examples.
func (sd *SelectDataset) Having(expressions ...exp.Expression) *SelectDataset {
	return sd.copy(sd.clauses.HavingAppend(expressions...))
}

// Adds a ORDER clause. If the ORDER is currently set it replaces it. See examples.
func (sd *SelectDataset) Order(order ...exp.OrderedExpression) *SelectDataset {
	return sd.copy(sd.clauses.SetOrder(order...))
}

// Adds a more columns to the current ORDER BY clause. If no order has be previously specified it is the same as
// calling Order. See examples.
func (sd *SelectDataset) OrderAppend(order ...exp.OrderedExpression) *SelectDataset {
	return sd.copy(sd.clauses.OrderAppend(order...))
}

// Adds a more columns to the beginning of the current ORDER BY clause. If no order has be previously specified it is the same as
// calling Order. See examples.
func (sd *SelectDataset) OrderPrepend(order ...exp.OrderedExpression) *SelectDataset {
	return sd.copy(sd.clauses.OrderPrepend(order...))
}

// Removes the ORDER BY clause. See examples.
func (sd *SelectDataset) ClearOrder() *SelectDataset {
	return sd.copy(sd.clauses.ClearOrder())
}

// Adds a LIMIT clause. If the LIMIT is currently set it replaces it. See examples.
func (sd *SelectDataset) Limit(limit uint) *SelectDataset {
	if limit > 0 {
		return sd.copy(sd.clauses.SetLimit(limit))
	}
	return sd.copy(sd.clauses.ClearLimit())
}

// Adds a LIMIT ALL clause. If the LIMIT is currently set it replaces it. See examples.
func (sd *SelectDataset) LimitAll() *SelectDataset {
	return sd.copy(sd.clauses.SetLimit(L("ALL")))
}

// Removes the LIMIT clause.
func (sd *SelectDataset) ClearLimit() *SelectDataset {
	return sd.copy(sd.clauses.ClearLimit())
}

// Adds an OFFSET clause. If the OFFSET is currently set it replaces it. See examples.
func (sd *SelectDataset) Offset(offset uint) *SelectDataset {
	return sd.copy(sd.clauses.SetOffset(offset))
}

// Removes the OFFSET clause from the Dataset
func (sd *SelectDataset) ClearOffset() *SelectDataset {
	return sd.copy(sd.clauses.ClearOffset())
}

// Creates an UNION statement with another dataset.
// If this or the other dataset has a limit or offset it will use that dataset as a subselect in the FROM clause.
// See examples.
func (sd *SelectDataset) Union(other *SelectDataset) *SelectDataset {
	return sd.withCompound(exp.UnionCompoundType, other.CompoundFromSelf())
}

// Creates an UNION ALL statement with another dataset.
// If this or the other dataset has a limit or offset it will use that dataset as a subselect in the FROM clause.
// See examples.
func (sd *SelectDataset) UnionAll(other *SelectDataset) *SelectDataset {
	return sd.withCompound(exp.UnionAllCompoundType, other.CompoundFromSelf())
}

// Creates an INTERSECT statement with another dataset.
// If this or the other dataset has a limit or offset it will use that dataset as a subselect in the FROM clause.
// See examples.
func (sd *SelectDataset) Intersect(other *SelectDataset) *SelectDataset {
	return sd.withCompound(exp.IntersectCompoundType, other.CompoundFromSelf())
}

// Creates an INTERSECT ALL statement with another dataset.
// If this or the other dataset has a limit or offset it will use that dataset as a subselect in the FROM clause.
// See examples.
func (sd *SelectDataset) IntersectAll(other *SelectDataset) *SelectDataset {
	return sd.withCompound(exp.IntersectAllCompoundType, other.CompoundFromSelf())
}

func (sd *SelectDataset) withCompound(ct exp.CompoundType, other exp.AppendableExpression) *SelectDataset {
	ce := exp.NewCompoundExpression(ct, other)
	ret := sd.CompoundFromSelf()
	ret.clauses = ret.clauses.CompoundsAppend(ce)
	return ret
}

// Used internally to determine if the dataset needs to use iteself as a source.
// If the dataset has an order or limit it will select from itself
func (sd *SelectDataset) CompoundFromSelf() *SelectDataset {
	if sd.clauses.HasOrder() || sd.clauses.HasLimit() {
		return sd.FromSelf()
	}
	return sd.copy(sd.clauses)
}

// Sets the alias for this dataset. This is typically used when using a Dataset as a subselect. See examples.
func (sd *SelectDataset) As(alias string) *SelectDataset {
	return sd.copy(sd.clauses.SetAlias(T(alias)))
}

// Returns the alias value as an identiier expression
func (sd *SelectDataset) GetAs() exp.IdentifierExpression {
	return sd.clauses.Alias()
}

// Sets the WINDOW clauses
func (sd *SelectDataset) Window(ws ...exp.WindowExpression) *SelectDataset {
	return sd.copy(sd.clauses.SetWindows(ws))
}

// Sets the WINDOW clauses
func (sd *SelectDataset) WindowAppend(ws ...exp.WindowExpression) *SelectDataset {
	return sd.copy(sd.clauses.WindowsAppend(ws...))
}

// Sets the WINDOW clauses
func (sd *SelectDataset) ClearWindow() *SelectDataset {
	return sd.copy(sd.clauses.ClearWindows())
}

// Get any error that has been set or nil if no error has been set.
func (sd *SelectDataset) Error() error {
	return sd.err
}

// Set an error on the dataset if one has not already been set. This error will be returned by a future call to Error
// or as part of ToSQL. This can be used by end users to record errors while building up queries without having to
// track those separately.
func (sd *SelectDataset) SetError(err error) *SelectDataset {
	if sd.err == nil {
		sd.err = err
	}

	return sd
}

// Generates a SELECT sql statement, if Prepared has been called with true then the parameters will not be interpolated.
// See examples.
//
// Errors:
//  * There is an error generating the SQL
func (sd *SelectDataset) ToSQL() (sql string, params []interface{}, err error) {
	return sd.selectSQLBuilder().ToSQL()
}

// Generates the SELECT sql, and returns an Exec struct with the sql set to the SELECT statement
//    db.From("test").Select("col").Executor()
//
// See Dataset#ToUpdateSQL for arguments
func (sd *SelectDataset) Executor() exec.QueryExecutor {
	return sd.queryFactory.FromSQLBuilder(sd.selectSQLBuilder())
}

// Appends this Dataset's SELECT statement to the SQLBuilder
// This is used internally for sub-selects by the dialect
func (sd *SelectDataset) AppendSQL(b sb.SQLBuilder) {
	if sd.err != nil {
		b.SetError(sd.err)
		return
	}
	sd.dialect.ToSelectSQL(b, sd.GetClauses())
}

func (sd *SelectDataset) ReturnsColumns() bool {
	return true
}

// Generates the SELECT sql for this dataset and uses Exec#ScanStructs to scan the results into a slice of structs.
//
// ScanStructs will only select the columns that can be scanned in to the struct unless you have explicitly selected
// certain columns. See examples.
//
// i: A pointer to a slice of structs
func (sd *SelectDataset) ScanStructs(i interface{}) error {
	return sd.ScanStructsContext(context.Background(), i)
}

// Generates the SELECT sql for this dataset and uses Exec#ScanStructsContext to scan the results into a slice of
// structs.
//
// ScanStructsContext will only select the columns that can be scanned in to the struct unless you have explicitly
// selected certain columns. See examples.
//
// i: A pointer to a slice of structs
func (sd *SelectDataset) ScanStructsContext(ctx context.Context, i interface{}) error {
	if sd.queryFactory == nil {
		return ErrQueryFactoryNotFoundError
	}
	ds := sd
	if sd.GetClauses().IsDefaultSelect() {
		ds = sd.Select(i)
	}
	return ds.Executor().ScanStructsContext(ctx, i)
}

// Generates the SELECT sql for this dataset and uses Exec#ScanStruct to scan the result into a slice of structs
//
// ScanStruct will only select the columns that can be scanned in to the struct unless you have explicitly selected
// certain columns. See examples.
//
// i: A pointer to a structs
func (sd *SelectDataset) ScanStruct(i interface{}) (bool, error) {
	return sd.ScanStructContext(context.Background(), i)
}

// Generates the SELECT sql for this dataset and uses Exec#ScanStructContext to scan the result into a slice of structs
//
// ScanStructContext will only select the columns that can be scanned in to the struct unless you have explicitly
// selected certain columns. See examples.
//
// i: A pointer to a structs
func (sd *SelectDataset) ScanStructContext(ctx context.Context, i interface{}) (bool, error) {
	if sd.queryFactory == nil {
		return false, ErrQueryFactoryNotFoundError
	}
	ds := sd
	if sd.GetClauses().IsDefaultSelect() {
		ds = sd.Select(i)
	}
	return ds.Limit(1).Executor().ScanStructContext(ctx, i)
}

// Generates the SELECT sql for this dataset and uses Exec#ScanVals to scan the results into a slice of primitive values
//
// i: A pointer to a slice of primitive values
func (sd *SelectDataset) ScanVals(i interface{}) error {
	return sd.ScanValsContext(context.Background(), i)
}

// Generates the SELECT sql for this dataset and uses Exec#ScanValsContext to scan the results into a slice of primitive
// values
//
// i: A pointer to a slice of primitive values
func (sd *SelectDataset) ScanValsContext(ctx context.Context, i interface{}) error {
	if sd.queryFactory == nil {
		return ErrQueryFactoryNotFoundError
	}
	return sd.Executor().ScanValsContext(ctx, i)
}

// Generates the SELECT sql for this dataset and uses Exec#ScanVal to scan the result into a primitive value
//
// i: A pointer to a primitive value
func (sd *SelectDataset) ScanVal(i interface{}) (bool, error) {
	return sd.ScanValContext(context.Background(), i)
}

// Generates the SELECT sql for this dataset and uses Exec#ScanValContext to scan the result into a primitive value
//
// i: A pointer to a primitive value
func (sd *SelectDataset) ScanValContext(ctx context.Context, i interface{}) (bool, error) {
	if sd.queryFactory == nil {
		return false, ErrQueryFactoryNotFoundError
	}
	return sd.Limit(1).Executor().ScanValContext(ctx, i)
}

// Generates the SELECT COUNT(*) sql for this dataset and uses Exec#ScanVal to scan the result into an int64.
func (sd *SelectDataset) Count() (int64, error) {
	return sd.CountContext(context.Background())
}

// Generates the SELECT COUNT(*) sql for this dataset and uses Exec#ScanValContext to scan the result into an int64.
func (sd *SelectDataset) CountContext(ctx context.Context) (int64, error) {
	var count int64
	_, err := sd.Select(COUNT(Star()).As("count")).ScanValContext(ctx, &count)
	return count, err
}

// Generates the SELECT sql only selecting the passed in column and uses Exec#ScanVals to scan the result into a slice
// of primitive values.
//
// i: A slice of primitive values
//
// col: The column to select when generative the SQL
func (sd *SelectDataset) Pluck(i interface{}, col string) error {
	return sd.PluckContext(context.Background(), i, col)
}

// Generates the SELECT sql only selecting the passed in column and uses Exec#ScanValsContext to scan the result into a
// slice of primitive values.
//
// i: A slice of primitive values
//
// col: The column to select when generative the SQL
func (sd *SelectDataset) PluckContext(ctx context.Context, i interface{}, col string) error {
	return sd.Select(col).ScanValsContext(ctx, i)
}

func (sd *SelectDataset) selectSQLBuilder() sb.SQLBuilder {
	buf := sb.NewSQLBuilder(sd.isPrepared.Bool())
	if sd.err != nil {
		return buf.SetError(sd.err)
	}
	sd.dialect.ToSelectSQL(buf, sd.GetClauses())
	return buf
}
