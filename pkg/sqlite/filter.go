package sqlite

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

func illegalFilterCombination(type1, type2 string) error {
	return fmt.Errorf("cannot have %s and %s in the same filter", type1, type2)
}

func validateFilterCombination[T any](sf models.OperatorFilter[T]) error {
	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if sf.And != nil {
		if sf.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if sf.Not != nil {
			return illegalFilterCombination(and, not)
		}
	}

	if sf.Or != nil {
		if sf.Not != nil {
			return illegalFilterCombination(or, not)
		}
	}

	return nil
}

func handleSubFilter[T any](ctx context.Context, handler criterionHandler, f *filterBuilder, subFilter models.OperatorFilter[T]) {
	subQuery := &filterBuilder{}
	handler.handle(ctx, subQuery)

	if subFilter.And != nil {
		f.and(subQuery)
	}
	if subFilter.Or != nil {
		f.or(subQuery)
	}
	if subFilter.Not != nil {
		f.not(subQuery)
	}
}

type sqlClause struct {
	sql  string
	args []interface{}
}

func (c sqlClause) not() sqlClause {
	return sqlClause{
		sql:  "NOT (" + c.sql + ")",
		args: c.args,
	}
}

func makeClause(sql string, args ...interface{}) sqlClause {
	return sqlClause{
		sql:  sql,
		args: args,
	}
}

func joinClauses(joinType string, clauses ...sqlClause) sqlClause {
	var ret []string
	var args []interface{}

	for _, clause := range clauses {
		ret = append(ret, "("+clause.sql+")")
		args = append(args, clause.args...)
	}

	return sqlClause{sql: strings.Join(ret, " "+joinType+" "), args: args}
}

func orClauses(clauses ...sqlClause) sqlClause {
	return joinClauses("OR", clauses...)
}

func andClauses(clauses ...sqlClause) sqlClause {
	return joinClauses("AND", clauses...)
}

type join struct {
	table    string
	as       string
	onClause string
	joinType string
}

// equals returns true if the other join alias/table is equal to this one
func (j join) equals(o join) bool {
	return j.alias() == o.alias()
}

// alias returns the as string, or the table if as is empty
func (j join) alias() string {
	if j.as == "" {
		return j.table
	}

	return j.as
}

func (j join) toSQL() string {
	asStr := ""
	joinStr := j.joinType
	if j.as != "" && j.as != j.table {
		asStr = " AS " + j.as
	}
	if j.joinType == "" {
		joinStr = "LEFT"
	}

	return fmt.Sprintf("%s JOIN %s%s ON %s", joinStr, j.table, asStr, j.onClause)
}

type joins []join

func (j *joins) add(newJoins ...join) {
	// only add if not already joined
	for _, newJoin := range newJoins {
		found := false
		for _, jj := range *j {
			if jj.equals(newJoin) {
				found = true
				break
			}
		}

		if !found {
			*j = append(*j, newJoin)
		}
	}
}

func (j *joins) toSQL() string {
	if len(*j) == 0 {
		return ""
	}

	var ret []string
	for _, jj := range *j {
		ret = append(ret, jj.toSQL())
	}

	return " " + strings.Join(ret, " ")
}

type filterBuilder struct {
	subFilter   *filterBuilder
	subFilterOp string

	joins         joins
	whereClauses  []sqlClause
	havingClauses []sqlClause
	withClauses   []sqlClause
	recursiveWith bool

	err error
}

func (f *filterBuilder) empty() bool {
	return f == nil || (len(f.whereClauses) == 0 && len(f.joins) == 0 && len(f.havingClauses) == 0 && f.subFilter == nil)
}

func filterBuilderFromHandler(ctx context.Context, handler criterionHandler) *filterBuilder {
	f := &filterBuilder{}
	handler.handle(ctx, f)
	return f
}

var errSubFilterAlreadySet = errors.New(`sub-filter already set`)

// sub-filter operator values
var (
	andOp = "AND"
	orOp  = "OR"
	notOp = "AND NOT"
)

// and sets the sub-filter that will be ANDed with this one.
// Sets the error state if sub-filter is already set.
func (f *filterBuilder) and(a *filterBuilder) {
	if f.subFilter != nil {
		f.setError(errSubFilterAlreadySet)
		return
	}

	f.subFilter = a
	f.subFilterOp = andOp
}

// or sets the sub-filter that will be ORed with this one.
// Sets the error state if a sub-filter is already set.
func (f *filterBuilder) or(o *filterBuilder) {
	if f.subFilter != nil {
		f.setError(errSubFilterAlreadySet)
		return
	}

	f.subFilter = o
	f.subFilterOp = orOp
}

// not sets the sub-filter that will be AND NOTed with this one.
// Sets the error state if a sub-filter is already set.
func (f *filterBuilder) not(n *filterBuilder) {
	if f.subFilter != nil {
		f.setError(errSubFilterAlreadySet)
		return
	}

	f.subFilter = n
	f.subFilterOp = notOp
}

// addLeftJoin adds a left join to the filter. The join is expressed in SQL as:
// LEFT JOIN <table> [AS <as>] ON <onClause>
// The AS is omitted if as is empty.
// This method does not add a join if it its alias/table name is already
// present in another existing join.
func (f *filterBuilder) addLeftJoin(table, as, onClause string) {
	newJoin := join{
		table:    table,
		as:       as,
		onClause: onClause,
		joinType: "LEFT",
	}

	f.joins.add(newJoin)
}

// addInnerJoin adds an inner join to the filter. The join is expressed in SQL as:
// INNER JOIN <table> [AS <as>] ON <onClause>
// The AS is omitted if as is empty.
// This method does not add a join if it its alias/table name is already
// present in another existing join.
func (f *filterBuilder) addInnerJoin(table, as, onClause string) {
	newJoin := join{
		table:    table,
		as:       as,
		onClause: onClause,
		joinType: "INNER",
	}

	f.joins.add(newJoin)
}

// addWhere adds a where clause and arguments to the filter. Where clauses
// are ANDed together. Does not add anything if the provided string is empty.
func (f *filterBuilder) addWhere(sql string, args ...interface{}) {
	if sql == "" {
		return
	}
	f.whereClauses = append(f.whereClauses, makeClause(sql, args...))
}

// addHaving adds a where clause and arguments to the filter. Having clauses
// are ANDed together. Does not add anything if the provided string is empty.
func (f *filterBuilder) addHaving(sql string, args ...interface{}) {
	if sql == "" {
		return
	}
	f.havingClauses = append(f.havingClauses, makeClause(sql, args...))
}

// addWith adds a with clause and arguments to the filter
func (f *filterBuilder) addWith(sql string, args ...interface{}) {
	if sql == "" {
		return
	}

	f.withClauses = append(f.withClauses, makeClause(sql, args...))
}

// addRecursiveWith adds a with clause and arguments to the filter, and sets it to recursive
//
//nolint:unused
func (f *filterBuilder) addRecursiveWith(sql string, args ...interface{}) {
	if sql == "" {
		return
	}

	f.addWith(sql, args...)
	f.recursiveWith = true
}

func (f *filterBuilder) getSubFilterClause(clause, subFilterClause string) string {
	ret := clause

	if subFilterClause != "" {
		var op string
		if len(ret) > 0 {
			op = " " + f.subFilterOp + " "
		} else if f.subFilterOp == notOp {
			op = "NOT "
		}

		ret += op + "(" + subFilterClause + ")"
	}

	return ret
}

// generateWhereClauses generates the SQL where clause for this filter.
// All where clauses within the filter are ANDed together. This is combined
// with the sub-filter, which will use the applicable operator (AND/OR/AND NOT).
func (f *filterBuilder) generateWhereClauses() (clause string, args []interface{}) {
	clause, args = f.andClauses(f.whereClauses)

	if f.subFilter != nil {
		c, a := f.subFilter.generateWhereClauses()
		if c != "" {
			clause = f.getSubFilterClause(clause, c)
			if len(a) > 0 {
				args = append(args, a...)
			}
		}
	}

	return
}

// generateHavingClauses generates the SQL having clause for this filter.
// All having clauses within the filter are ANDed together. This is combined
// with the sub-filter, which will use the applicable operator (AND/OR/AND NOT).
func (f *filterBuilder) generateHavingClauses() (string, []interface{}) {
	clause, args := f.andClauses(f.havingClauses)

	if f.subFilter != nil {
		c, a := f.subFilter.generateHavingClauses()
		if c != "" {
			clause = f.getSubFilterClause(clause, c)
			if len(a) > 0 {
				args = append(args, a...)
			}
		}
	}

	return clause, args
}

func (f *filterBuilder) generateWithClauses() (string, []interface{}) {
	var clauses []string
	var args []interface{}
	for _, w := range f.withClauses {
		clauses = append(clauses, w.sql)
		args = append(args, w.args...)
	}

	if len(clauses) > 0 {
		return strings.Join(clauses, ", "), args
	}

	return "", nil
}

// getAllJoins returns all of the joins in this filter and any sub-filter(s).
// Redundant joins will not be duplicated in the return value.
func (f *filterBuilder) getAllJoins() joins {
	var ret joins
	ret.add(f.joins...)
	if f.subFilter != nil {
		subJoins := f.subFilter.getAllJoins()
		if len(subJoins) > 0 {
			ret.add(subJoins...)
		}
	}

	return ret
}

// getError returns the error state on this filter, or on any sub-filter(s) if
// the error state is nil.
func (f *filterBuilder) getError() error {
	if f.err != nil {
		return f.err
	}

	if f.subFilter != nil {
		return f.subFilter.getError()
	}

	return nil
}

// handleCriterion calls the handle function on the provided criterionHandler,
// providing itself.
func (f *filterBuilder) handleCriterion(ctx context.Context, handler criterionHandler) {
	handler.handle(ctx, f)
}

func (f *filterBuilder) setError(e error) {
	if f.err == nil {
		f.err = e
	}
}

func (f *filterBuilder) andClauses(input []sqlClause) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	for _, w := range input {
		clauses = append(clauses, w.sql)
		args = append(args, w.args...)
	}

	if len(clauses) > 0 {
		c := "(" + strings.Join(clauses, ") AND (") + ")"
		if len(clauses) > 1 {
			c = "(" + c + ")"
		}
		return c, args
	}

	return "", nil
}
