package sqlite

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"

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
	handle(ctx context.Context, f *filterBuilder)
}

type criterionHandlerFunc func(ctx context.Context, f *filterBuilder)

func (h criterionHandlerFunc) handle(ctx context.Context, f *filterBuilder) {
	h(ctx, f)
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

func stringCriterionHandler(c *models.StringCriterionInput, column string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
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
					f.addWhere(fmt.Sprintf("(%s IS NOT NULL AND %[1]s regexp ?)", column), c.Value)
				case models.CriterionModifierNotMatchesRegex:
					if _, err := regexp.Compile(c.Value); err != nil {
						f.setError(err)
						return
					}
					f.addWhere(fmt.Sprintf("(%s IS NULL OR %[1]s NOT regexp ?)", column), c.Value)
				case models.CriterionModifierIsNull:
					f.addWhere("(" + column + " IS NULL OR TRIM(" + column + ") = '')")
				case models.CriterionModifierNotNull:
					f.addWhere("(" + column + " IS NOT NULL AND TRIM(" + column + ") != '')")
				default:
					panic("unsupported string filter modifier")
				}
			}
		}
	}
}

func intCriterionHandler(c *models.IntCriterionInput, column string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			clause, args := getIntCriterionWhereClause(column, *c)
			f.addWhere(clause, args...)
		}
	}
}

func boolCriterionHandler(c *bool, column string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
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

// handle for MultiCriterion where there is a join table between the new
// objects
type joinedMultiCriterionHandlerBuilder struct {
	// table containing the primary objects
	primaryTable string
	// table joining primary and foreign objects
	joinTable string
	// alias for join table, if required
	joinAs string
	// foreign key of the primary object on the join table
	primaryFK string
	// foreign key of the foreign object on the join table
	foreignFK string

	addJoinTable func(f *filterBuilder)
}

func (m *joinedMultiCriterionHandlerBuilder) handler(criterion *models.MultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			joinAlias := m.joinAs
			if joinAlias == "" {
				joinAlias = m.joinTable
			}

			if criterion.Modifier == models.CriterionModifierIsNull || criterion.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if criterion.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				m.addJoinTable(f)

				f.addWhere(utils.StrFormat("{table}.{column} IS {not} NULL", utils.StrFormatMap{
					"table":  joinAlias,
					"column": m.foreignFK,
					"not":    notClause,
				}))
				return
			}

			if len(criterion.Value) == 0 {
				return
			}

			var args []interface{}
			for _, tagID := range criterion.Value {
				args = append(args, tagID)
			}

			whereClause := ""
			havingClause := ""

			switch criterion.Modifier {
			case models.CriterionModifierIncludes:
				// includes any of the provided ids
				m.addJoinTable(f)
				whereClause = fmt.Sprintf("%s.%s IN %s", joinAlias, m.foreignFK, getInBinding(len(criterion.Value)))
			case models.CriterionModifierIncludesAll:
				// includes all of the provided ids
				m.addJoinTable(f)
				whereClause = fmt.Sprintf("%s.%s IN %s", joinAlias, m.foreignFK, getInBinding(len(criterion.Value)))
				havingClause = fmt.Sprintf("count(distinct %s.%s) IS %d", joinAlias, m.foreignFK, len(criterion.Value))
			case models.CriterionModifierExcludes:
				// excludes all of the provided ids
				// need to use actual join table name for this
				// <primaryTable>.id NOT IN (select <joinTable>.<primaryFK> from <joinTable> where <joinTable>.<foreignFK> in <values>)
				whereClause = fmt.Sprintf("%[1]s.id NOT IN (SELECT %[3]s.%[2]s from %[3]s where %[3]s.%[4]s in %[5]s)", m.primaryTable, m.primaryFK, m.joinTable, m.foreignFK, getInBinding(len(criterion.Value)))
			}

			f.addWhere(whereClause, args...)
			f.addHaving(havingClause)
		}
	}
}

type multiCriterionHandlerBuilder struct {
	primaryTable string
	foreignTable string
	joinTable    string
	primaryFK    string
	foreignFK    string

	// function that will be called to perform any necessary joins
	addJoinsFunc func(f *filterBuilder)
}

func (m *multiCriterionHandlerBuilder) handler(criterion *models.MultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			if criterion.Modifier == models.CriterionModifierIsNull || criterion.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if criterion.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				table := m.primaryTable
				if m.joinTable != "" {
					table = m.joinTable
					f.addLeftJoin(table, "", fmt.Sprintf("%s.%s = %s.id", table, m.primaryFK, m.primaryTable))
				}

				f.addWhere(fmt.Sprintf("%s.%s IS %s NULL", table, m.foreignFK, notClause))
				return
			}

			if len(criterion.Value) == 0 {
				return
			}

			var args []interface{}
			for _, tagID := range criterion.Value {
				args = append(args, tagID)
			}

			if m.addJoinsFunc != nil {
				m.addJoinsFunc(f)
			}

			whereClause, havingClause := getMultiCriterionClause(m.primaryTable, m.foreignTable, m.joinTable, m.primaryFK, m.foreignFK, criterion)
			f.addWhere(whereClause, args...)
			f.addHaving(havingClause)
		}
	}
}

type countCriterionHandlerBuilder struct {
	primaryTable string
	joinTable    string
	primaryFK    string
}

func (m *countCriterionHandlerBuilder) handler(criterion *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			clause, args := getCountCriterionClause(m.primaryTable, m.joinTable, m.primaryFK, *criterion)

			f.addWhere(clause, args...)
		}
	}
}

// handler for StringCriterion for string list fields
type stringListCriterionHandlerBuilder struct {
	// table joining primary and foreign objects
	joinTable string
	// string field on the join table
	stringColumn string

	addJoinTable func(f *filterBuilder)
}

func (m *stringListCriterionHandlerBuilder) handler(criterion *models.StringCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil && len(criterion.Value) > 0 {
			m.addJoinTable(f)

			stringCriterionHandler(criterion, m.joinTable+"."+m.stringColumn)(ctx, f)
		}
	}
}

type hierarchicalMultiCriterionHandlerBuilder struct {
	tx dbi

	primaryTable string
	foreignTable string
	foreignFK    string

	derivedTable   string
	parentFK       string
	relationsTable string
}

func getHierarchicalValues(ctx context.Context, tx dbi, values []string, table, relationsTable, parentFK string, depth *int) string {
	var args []interface{}

	depthVal := 0
	if depth != nil {
		depthVal = *depth
	}

	if depthVal == 0 {
		valid := true
		var valuesClauses []string
		for _, value := range values {
			id, err := strconv.Atoi(value)
			// In case of invalid value just run the query.
			// Building VALUES() based on provided values just saves a query when depth is 0.
			if err != nil {
				valid = false
				break
			}

			valuesClauses = append(valuesClauses, fmt.Sprintf("(%d,%d)", id, id))
		}

		if valid {
			return "VALUES" + strings.Join(valuesClauses, ",")
		}
	}

	for _, value := range values {
		args = append(args, value)
	}
	inCount := len(args)

	var depthCondition string
	if depthVal != -1 {
		depthCondition = fmt.Sprintf("WHERE depth < %d", depthVal)
	}

	withClauseMap := utils.StrFormatMap{
		"table":           table,
		"relationsTable":  relationsTable,
		"inBinding":       getInBinding(inCount),
		"recursiveSelect": "",
		"parentFK":        parentFK,
		"depthCondition":  depthCondition,
		"unionClause":     "",
	}

	if relationsTable != "" {
		withClauseMap["recursiveSelect"] = utils.StrFormat(`SELECT p.root_id, c.child_id, depth + 1 FROM {relationsTable} AS c
INNER JOIN items as p ON c.parent_id = p.item_id
`, withClauseMap)
	} else {
		withClauseMap["recursiveSelect"] = utils.StrFormat(`SELECT p.root_id, c.id, depth + 1 FROM {table} as c
INNER JOIN items as p ON c.{parentFK} = p.item_id
`, withClauseMap)
	}

	if depthVal != 0 {
		withClauseMap["unionClause"] = utils.StrFormat(`
UNION {recursiveSelect} {depthCondition}
`, withClauseMap)
	}

	withClause := utils.StrFormat(`items AS (
SELECT id as root_id, id as item_id, 0 as depth FROM {table}
WHERE id in {inBinding}
{unionClause})
`, withClauseMap)

	query := fmt.Sprintf("WITH RECURSIVE %s SELECT 'VALUES' || GROUP_CONCAT('(' || root_id || ', ' || item_id || ')') AS val FROM items", withClause)

	var valuesClause string
	err := tx.Get(ctx, &valuesClause, query, args...)
	if err != nil {
		logger.Error(err)
		// return record which never matches so we don't have to handle error here
		return "VALUES(NULL, NULL)"
	}

	return valuesClause
}

func addHierarchicalConditionClauses(f *filterBuilder, criterion *models.HierarchicalMultiCriterionInput, table, idColumn string) {
	switch criterion.Modifier {
	case models.CriterionModifierIncludes:
		f.addWhere(fmt.Sprintf("%s.%s IS NOT NULL", table, idColumn))
	case models.CriterionModifierIncludesAll:
		f.addWhere(fmt.Sprintf("%s.%s IS NOT NULL", table, idColumn))
		f.addHaving(fmt.Sprintf("count(distinct %s.%s) IS %d", table, idColumn, len(criterion.Value)))
	case models.CriterionModifierExcludes:
		f.addWhere(fmt.Sprintf("%s.%s IS NULL", table, idColumn))
	}
}

func (m *hierarchicalMultiCriterionHandlerBuilder) handler(criterion *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			if criterion.Modifier == models.CriterionModifierIsNull || criterion.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if criterion.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addWhere(utils.StrFormat("{table}.{column} IS {not} NULL", utils.StrFormatMap{
					"table":  m.primaryTable,
					"column": m.foreignFK,
					"not":    notClause,
				}))
				return
			}

			if len(criterion.Value) == 0 {
				return
			}

			valuesClause := getHierarchicalValues(ctx, m.tx, criterion.Value, m.foreignTable, m.relationsTable, m.parentFK, criterion.Depth)

			f.addLeftJoin("(SELECT column1 AS root_id, column2 AS item_id FROM ("+valuesClause+"))", m.derivedTable, fmt.Sprintf("%s.item_id = %s.%s", m.derivedTable, m.primaryTable, m.foreignFK))

			addHierarchicalConditionClauses(f, criterion, m.derivedTable, "root_id")
		}
	}
}

type joinedHierarchicalMultiCriterionHandlerBuilder struct {
	tx dbi

	primaryTable string
	foreignTable string
	foreignFK    string

	parentFK       string
	relationsTable string

	joinAs    string
	joinTable string
	primaryFK string
}

func (m *joinedHierarchicalMultiCriterionHandlerBuilder) handler(criterion *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			joinAlias := m.joinAs

			if criterion.Modifier == models.CriterionModifierIsNull || criterion.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if criterion.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin(m.joinTable, joinAlias, fmt.Sprintf("%s.%s = %s.id", joinAlias, m.primaryFK, m.primaryTable))

				f.addWhere(utils.StrFormat("{table}.{column} IS {not} NULL", utils.StrFormatMap{
					"table":  joinAlias,
					"column": m.foreignFK,
					"not":    notClause,
				}))
				return
			}

			if len(criterion.Value) == 0 {
				return
			}

			valuesClause := getHierarchicalValues(ctx, m.tx, criterion.Value, m.foreignTable, m.relationsTable, m.parentFK, criterion.Depth)

			joinTable := utils.StrFormat(`(
	SELECT j.*, d.column1 AS root_id, d.column2 AS item_id FROM {joinTable} AS j
	INNER JOIN ({valuesClause}) AS d ON j.{foreignFK} = d.column2
)
`, utils.StrFormatMap{
				"joinTable":    m.joinTable,
				"foreignFK":    m.foreignFK,
				"valuesClause": valuesClause,
			})

			f.addLeftJoin(joinTable, joinAlias, fmt.Sprintf("%s.%s = %s.id", joinAlias, m.primaryFK, m.primaryTable))

			addHierarchicalConditionClauses(f, criterion, joinAlias, "root_id")
		}
	}
}
