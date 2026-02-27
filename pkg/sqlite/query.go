package sqlite

import (
	"context"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

type queryBuilder struct {
	repository *repository

	columns []string
	from    string

	joins         joins
	whereClauses  []string
	havingClauses []string
	withClauses   []string
	recursiveWith bool

	withArgs   []interface{}
	joinArgs   []interface{}
	whereArgs  []interface{}
	havingArgs []interface{}

	sortAndPagination string
}

func (qb queryBuilder) allArgs() []interface{} {
	var args []interface{}
	args = append(args, qb.withArgs...)
	args = append(args, qb.joinArgs...)
	args = append(args, qb.whereArgs...)
	args = append(args, qb.havingArgs...)
	return args
}

func (qb queryBuilder) body(includeSortPagination bool) string {
	return fmt.Sprintf("SELECT %s FROM %s%s", strings.Join(qb.columns, ", "), qb.from, qb.joins.toSQL(includeSortPagination))
}

func (qb *queryBuilder) addColumn(column string) {
	qb.columns = append(qb.columns, column)
}

func (qb queryBuilder) toSQL(includeSortPagination bool) string {
	body := qb.body(includeSortPagination)

	withClause := ""
	if len(qb.withClauses) > 0 {
		var recursive string
		if qb.recursiveWith {
			recursive = " RECURSIVE "
		}
		withClause = "WITH " + recursive + strings.Join(qb.withClauses, ", ") + " "
	}

	body = withClause + qb.repository.buildQueryBody(body, qb.whereClauses, qb.havingClauses)
	if includeSortPagination {
		body += qb.sortAndPagination
	}

	return body
}

func (qb queryBuilder) findIDs(ctx context.Context) ([]int, error) {
	const includeSortPagination = true
	sql := qb.toSQL(includeSortPagination)
	return qb.repository.runIdsQuery(ctx, sql, qb.allArgs())
}

func (qb queryBuilder) executeFind(ctx context.Context) ([]int, int, error) {
	const includeSortPagination = true
	body := qb.body(includeSortPagination)
	return qb.repository.executeFindQuery(ctx, body, qb.allArgs(), qb.sortAndPagination, qb.whereClauses, qb.havingClauses, qb.withClauses, qb.recursiveWith)
}

func (qb queryBuilder) executeCount(ctx context.Context) (int, error) {
	const includeSortPagination = false
	body := qb.body(includeSortPagination)

	withClause := ""
	if len(qb.withClauses) > 0 {
		var recursive string
		if qb.recursiveWith {
			recursive = " RECURSIVE "
		}
		withClause = "WITH " + recursive + strings.Join(qb.withClauses, ", ") + " "
	}

	body = qb.repository.buildQueryBody(body, qb.whereClauses, qb.havingClauses)
	countQuery := withClause + qb.repository.buildCountQuery(body)
	return qb.repository.runCountQuery(ctx, countQuery, qb.allArgs())
}

func (qb *queryBuilder) addWhere(clauses ...string) {
	for _, clause := range clauses {
		if len(clause) > 0 {
			qb.whereClauses = append(qb.whereClauses, clause)
		}
	}
}

func (qb *queryBuilder) addHaving(clauses ...string) {
	for _, clause := range clauses {
		if len(clause) > 0 {
			qb.havingClauses = append(qb.havingClauses, clause)
		}
	}
}

func (qb *queryBuilder) addWith(recursive bool, clauses ...string) {
	for _, clause := range clauses {
		if len(clause) > 0 {
			qb.withClauses = append(qb.withClauses, clause)
		}
	}

	qb.recursiveWith = qb.recursiveWith || recursive
}

func (qb *queryBuilder) addArg(args ...interface{}) {
	qb.whereArgs = append(qb.whereArgs, args...)
}

func (qb *queryBuilder) addHavingArg(args ...interface{}) {
	qb.havingArgs = append(qb.havingArgs, args...)
}

func (qb *queryBuilder) hasJoin(alias string) bool {
	for _, j := range qb.joins {
		if j.alias() == alias {
			return true
		}
	}

	return false
}

func (qb *queryBuilder) join(table, as, onClause string) {
	newJoin := join{
		table:    table,
		as:       as,
		onClause: onClause,
		joinType: "LEFT",
	}

	qb.joins.add(newJoin)
}

func (qb *queryBuilder) joinSort(table, as, onClause string) {
	newJoin := join{
		sort:     true,
		table:    table,
		as:       as,
		onClause: onClause,
		joinType: "LEFT",
	}

	qb.joins.add(newJoin)
}

func (qb *queryBuilder) addJoins(joins ...join) {
	for _, j := range joins {
		if qb.joins.addUnique(j) {
			qb.joinArgs = append(qb.joinArgs, j.args...)
		}
	}
}

func (qb *queryBuilder) addFilter(f *filterBuilder) error {
	err := f.getError()
	if err != nil {
		return err
	}

	clause, args := f.generateWithClauses()
	if len(clause) > 0 {
		qb.addWith(f.recursiveWith, clause)
	}
	if len(args) > 0 {
		qb.withArgs = append(qb.withArgs, args...)
	}

	qb.addJoins(f.getAllJoins()...)

	clause, args = f.generateWhereClauses()
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
		qb.addHavingArg(args...)
	}

	return nil
}

func (qb *queryBuilder) parseQueryString(columns []string, q string) {
	specs := models.ParseSearchString(q)

	for _, t := range specs.MustHave {
		var clauses []string

		for _, column := range columns {
			clauses = append(clauses, column+" LIKE ?")
			qb.addArg(like(t))
		}

		qb.addWhere("(" + strings.Join(clauses, " OR ") + ")")
	}

	for _, t := range specs.MustNot {
		for _, column := range columns {
			qb.addWhere(coalesce(column) + " NOT LIKE ?")
			qb.addArg(like(t))
		}
	}

	for _, set := range specs.AnySets {
		var clauses []string

		for _, column := range columns {
			for _, v := range set {
				clauses = append(clauses, column+" LIKE ?")
				qb.addArg(like(v))
			}
		}

		qb.addWhere("(" + strings.Join(clauses, " OR ") + ")")
	}
}
