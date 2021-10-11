package sqlite

import (
	"strings"
)

type queryBuilder struct {
	repository *repository

	body string

	joins         joins
	whereClauses  []string
	havingClauses []string
	args          []interface{}
	withClauses   []string
	recursiveWith bool

	sortAndPagination string

	err error
}

func (qb queryBuilder) executeFind() ([]int, int, error) {
	if qb.err != nil {
		return nil, 0, qb.err
	}

	body := qb.body
	body += qb.joins.toSQL()

	return qb.repository.executeFindQuery(body, qb.args, qb.sortAndPagination, qb.whereClauses, qb.havingClauses, qb.withClauses, qb.recursiveWith)
}

func (qb queryBuilder) executeCount() (int, error) {
	if qb.err != nil {
		return 0, qb.err
	}

	body := qb.body
	body += qb.joins.toSQL()

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
	return qb.repository.runCountQuery(countQuery, qb.args)
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
	qb.args = append(qb.args, args...)
}

func (qb *queryBuilder) join(table, as, onClause string) {
	newJoin := join{
		table:    table,
		as:       as,
		onClause: onClause,
	}

	qb.joins.add(newJoin)
}

func (qb *queryBuilder) addJoins(joins ...join) {
	qb.joins.add(joins...)
}

func (qb *queryBuilder) addFilter(f *filterBuilder) {
	err := f.getError()
	if err != nil {
		qb.err = err
		return
	}

	clause, args := f.generateWithClauses()
	if len(clause) > 0 {
		qb.addWith(f.recursiveWith, clause)
	}

	if len(args) > 0 {
		// WITH clause always comes first and thus precedes alk args
		qb.args = append(args, qb.args...)
	}

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
		qb.addArg(args...)
	}

	qb.addJoins(f.getAllJoins()...)
}
