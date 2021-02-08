package sqlite

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

type sqlClause struct {
	sql  string
	args []interface{}
}

func makeClause(sql string, args ...interface{}) sqlClause {
	return sqlClause{
		sql:  sql,
		args: args,
	}
}

type criterionHandler interface {
	handle(f *filterBuilder)
}

type criterionHandlerFunc func(f *filterBuilder)

type join struct {
	table    string
	as       string
	onClause string
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
	if j.as != "" && j.as != j.table {
		asStr = " AS " + j.as + " "
	}

	return fmt.Sprintf("LEFT JOIN %s%s on %s", j.table, asStr, j.onClause)
}

type joins []join

func (j *joins) add(newJoins ...join) {
	// only add if not already joined
	for _, newJoin := range newJoins {
		for _, jj := range *j {
			if jj.equals(newJoin) {
				return
			}
		}

		*j = append(*j, newJoin)
	}
}

func (j *joins) toSQL() string {
	var ret []string
	for _, jj := range *j {
		ret = append(ret, jj.toSQL())
	}

	return strings.Join(ret, " ")
}

type filterBuilder struct {
	andFilter *filterBuilder
	orFilter  *filterBuilder
	notFilter *filterBuilder

	joins         joins
	whereClauses  []sqlClause
	havingClauses []sqlClause

	err error
}

var errAndOrFilter error = errors.New(`cannot use "and" and "or" within the same filter`)
var errAndNotFilter error = errors.New(`cannot use "and" and "not" within the same filter`)
var errOrNotFilter error = errors.New(`cannot use "or" and "not" within the same filter`)

func (f *filterBuilder) and(a *filterBuilder) {
	if f.orFilter != nil {
		f.setError(errAndOrFilter)
	}
	if f.notFilter != nil {
		f.setError(errAndNotFilter)
	}

	f.andFilter = a
}

func (f *filterBuilder) or(o *filterBuilder) {
	if f.andFilter != nil {
		f.setError(errAndOrFilter)
	}
	if f.notFilter != nil {
		f.setError(errOrNotFilter)
	}

	f.orFilter = o
}

func (f *filterBuilder) not(n *filterBuilder) {
	if f.andFilter != nil {
		f.setError(errAndNotFilter)
	}
	if f.orFilter != nil {
		f.setError(errOrNotFilter)
	}

	f.notFilter = n
}

func (f *filterBuilder) addJoin(table, as, onClause string) {
	newJoin := join{
		table:    table,
		as:       as,
		onClause: onClause,
	}

	f.joins.add(newJoin)
}

func (f *filterBuilder) addWhere(sql string, args ...interface{}) {
	if sql == "" {
		return
	}
	f.whereClauses = append(f.whereClauses, makeClause(sql, args...))
}

func (f *filterBuilder) addHaving(sql string, args ...interface{}) {
	if sql == "" {
		return
	}
	f.havingClauses = append(f.havingClauses, makeClause(sql, args...))
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
		c := "(" + strings.Join(clauses, " AND ") + ")"
		return c, args
	}

	return "", nil
}

func (f *filterBuilder) subFilter() (s *filterBuilder, op string) {
	if f.andFilter != nil {
		s = f.andFilter
		op = "AND"
	} else if f.orFilter != nil {
		s = f.orFilter
		op = "OR"
	} else if f.notFilter != nil {
		s = f.notFilter
		op = "AND NOT"
	}

	return
}

func (f *filterBuilder) generateWhereClauses() (string, []interface{}) {
	clause, args := f.andClauses(f.whereClauses)

	subFilter, op := f.subFilter()
	if subFilter != nil {
		c, a := subFilter.generateWhereClauses()
		if c != "" {
			clause += " " + op + " " + c
			if len(a) > 0 {
				args = append(args, a...)
			}
		}
	}

	return clause, args
}

func (f *filterBuilder) generateHavingClauses() (string, []interface{}) {
	clause, args := f.andClauses(f.havingClauses)

	subFilter, op := f.subFilter()
	if subFilter != nil {
		c, a := subFilter.generateHavingClauses()
		if c != "" {
			clause += " " + op + " " + c
			if len(a) > 0 {
				args = append(args, a...)
			}
		}
	}

	return clause, args
}

func (f *filterBuilder) getAllJoins() joins {
	var ret joins
	ret.add(f.joins...)
	subFilter, _ := f.subFilter()
	if subFilter != nil {
		subJoins := subFilter.getAllJoins()
		if len(subJoins) > 0 {
			ret.add(subJoins...)
		}
	}

	return ret
}

func (f *filterBuilder) getError() error {
	if f.err != nil {
		return f.err
	}

	subFilter, _ := f.subFilter()
	if subFilter != nil {
		return subFilter.getError()
	}

	return nil
}

func (f *filterBuilder) addToQueryBuilder(qb *queryBuilder) error {
	err := f.getError()
	if err != nil {
		return err
	}

	clause, args := f.generateWhereClauses()
	if len(clause) > 0 {
		qb.addWhere(clause)
	}

	if len(args) > 0 {
		qb.addArg(args...)
	}

	clause, args = f.generateHavingClauses()
	if len(clause) > 0 {
		qb.addHaving(clause)
	}

	if len(args) > 0 {
		qb.addArg(args...)
	}

	qb.addJoins(f.getAllJoins()...)

	return nil
}

func (f *filterBuilder) handleCriterion(handler criterionHandler) {
	f.handleCriterionFunc(func(h *filterBuilder) {
		handler.handle(h)
	})
}

func (f *filterBuilder) handleCriterionFunc(handler criterionHandlerFunc) {
	handler(f)
}

func stringCriterionHandler(c *models.StringCriterionInput, column string) criterionHandlerFunc {
	return func(f *filterBuilder) {
		if c != nil {
			if modifier := c.Modifier; c.Modifier.IsValid() {
				switch modifier {
				case models.CriterionModifierIncludes:
					clause, thisArgs := getSearchBinding([]string{column}, c.Value, false)
					f.addWhere(clause, thisArgs...)
				case models.CriterionModifierExcludes:
					clause, thisArgs := getSearchBinding([]string{column}, c.Value, true)
					f.addWhere(clause, thisArgs...)
				case models.CriterionModifierEquals:
					f.addWhere(column+" LIKE ?", c.Value)
				case models.CriterionModifierNotEquals:
					f.addWhere(column+" NOT LIKE ?", c.Value)
				case models.CriterionModifierMatchesRegex:
					if _, err := regexp.Compile(c.Value); err != nil {
						f.setError(err)
						return
					}
					f.addWhere(column+" regexp ?", c.Value)
				case models.CriterionModifierNotMatchesRegex:
					if _, err := regexp.Compile(c.Value); err != nil {
						f.setError(err)
						return
					}
					f.addWhere(column+" NOT regexp ?", c.Value)
				default:
					clause, count := getSimpleCriterionClause(modifier, "?")

					if count == 1 {
						f.addWhere(column+" "+clause, c.Value)
					} else {
						f.addWhere(column + " " + clause)
					}
				}
			}
		}
	}
}

func intCriterionHandler(c *models.IntCriterionInput, column string) criterionHandlerFunc {
	return func(f *filterBuilder) {
		if c != nil {
			clause, count := getIntCriterionWhereClause(column, *c)

			if count == 1 {
				f.addWhere(clause, c.Value)
			} else {
				f.addWhere(clause)
			}
		}
	}
}

func boolCriterionHandler(c *bool, column string) criterionHandlerFunc {
	return func(f *filterBuilder) {
		if c != nil {
			var v string
			if *c {
				v = "1"
			} else {
				v = "0"
			}

			f.addWhere(column + " = " + v)
		}
	}
}

func stringLiteralCriterionHandler(v *string, column string) criterionHandlerFunc {
	return func(f *filterBuilder) {
		if v != nil {
			f.addWhere(column+" = ?", v)
		}
	}
}
