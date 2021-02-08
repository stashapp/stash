package sqlite

import (
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

func (j join) getTable() string {
	if j.as == "" {
		return j.table
	}

	return j.as
}

func (j join) toSQL() string {
	asStr := ""
	if j.as != "" {
		asStr = " AS " + j.as + " "
	}

	return fmt.Sprintf(" LEFT JOIN %s%s on %s ", j.table, asStr, j.onClause)
}

type filterBuilder struct {
	joins         []join
	whereClauses  []sqlClause
	havingClauses []sqlClause

	err error
}

func (f *filterBuilder) addJoin(table, as, onClause string) {
	// only add if not already joined
	name := as
	if name == "" {
		name = table
	}

	for _, j := range f.joins {
		if j.getTable() == name {
			return
		}
	}

	f.joins = append(f.joins, join{
		table:    table,
		as:       as,
		onClause: onClause,
	})
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

func (f *filterBuilder) joinClauses(input []sqlClause) (string, []interface{}) {
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

func (f *filterBuilder) addToQueryBuilder(qb *queryBuilder) {
	clauses, args := f.joinClauses(f.whereClauses)
	if len(clauses) > 0 {
		qb.addWhere(clauses)
	}

	if len(args) > 0 {
		qb.addArg(args...)
	}

	clauses, args = f.joinClauses(f.havingClauses)
	if len(clauses) > 0 {
		qb.addHaving(clauses)
	}

	if len(args) > 0 {
		qb.addArg(args...)
	}

	for _, j := range f.joins {
		qb.body += j.toSQL()
	}
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
