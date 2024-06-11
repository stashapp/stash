package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type criterionHandler interface {
	handle(ctx context.Context, f *filterBuilder)
}

type criterionHandlerFunc func(ctx context.Context, f *filterBuilder)

func (h criterionHandlerFunc) handle(ctx context.Context, f *filterBuilder) {
	h(ctx, f)
}

type compoundHandler []criterionHandler

func (h compoundHandler) handle(ctx context.Context, f *filterBuilder) {
	for _, h := range h {
		h.handle(ctx, f)
	}
}

// shared criterion handlers go here

func stringCriterionHandler(c *models.StringCriterionInput, column string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			if modifier := c.Modifier; c.Modifier.IsValid() {
				switch modifier {
				case models.CriterionModifierIncludes:
					f.whereClauses = append(f.whereClauses, getStringSearchClause([]string{column}, c.Value, false))
				case models.CriterionModifierExcludes:
					f.whereClauses = append(f.whereClauses, getStringSearchClause([]string{column}, c.Value, true))
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

func enumCriterionHandler(modifier models.CriterionModifier, values []string, column string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if modifier.IsValid() {
			switch modifier {
			case models.CriterionModifierIncludes, models.CriterionModifierEquals:
				if len(values) > 0 {
					f.whereClauses = append(f.whereClauses, getEnumSearchClause(column, values, false))
				}
			case models.CriterionModifierExcludes, models.CriterionModifierNotEquals:
				if len(values) > 0 {
					f.whereClauses = append(f.whereClauses, getEnumSearchClause(column, values, true))
				}
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

func pathCriterionHandler(c *models.StringCriterionInput, pathColumn string, basenameColumn string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}
			addWildcards := true
			not := false

			if modifier := c.Modifier; c.Modifier.IsValid() {
				switch modifier {
				case models.CriterionModifierIncludes:
					f.whereClauses = append(f.whereClauses, getPathSearchClauseMany(pathColumn, basenameColumn, c.Value, addWildcards, not))
				case models.CriterionModifierExcludes:
					not = true
					f.whereClauses = append(f.whereClauses, getPathSearchClauseMany(pathColumn, basenameColumn, c.Value, addWildcards, not))
				case models.CriterionModifierEquals:
					addWildcards = false
					f.whereClauses = append(f.whereClauses, getPathSearchClause(pathColumn, basenameColumn, c.Value, addWildcards, not))
				case models.CriterionModifierNotEquals:
					addWildcards = false
					not = true
					f.whereClauses = append(f.whereClauses, getPathSearchClause(pathColumn, basenameColumn, c.Value, addWildcards, not))
				case models.CriterionModifierMatchesRegex:
					if _, err := regexp.Compile(c.Value); err != nil {
						f.setError(err)
						return
					}
					filepathColumn := fmt.Sprintf("%s || '%s' || %s", pathColumn, string(filepath.Separator), basenameColumn)
					f.addWhere(fmt.Sprintf("%s IS NOT NULL AND %s IS NOT NULL AND %s regexp ?", pathColumn, basenameColumn, filepathColumn), c.Value)
				case models.CriterionModifierNotMatchesRegex:
					if _, err := regexp.Compile(c.Value); err != nil {
						f.setError(err)
						return
					}
					filepathColumn := fmt.Sprintf("%s || '%s' || %s", pathColumn, string(filepath.Separator), basenameColumn)
					f.addWhere(fmt.Sprintf("%s IS NULL OR %s IS NULL OR %s NOT regexp ?", pathColumn, basenameColumn, filepathColumn), c.Value)
				case models.CriterionModifierIsNull:
					f.addWhere(fmt.Sprintf("%s IS NULL OR TRIM(%[1]s) = '' OR %s IS NULL OR TRIM(%[2]s) = ''", pathColumn, basenameColumn))
				case models.CriterionModifierNotNull:
					f.addWhere(fmt.Sprintf("%s IS NOT NULL AND TRIM(%[1]s) != '' AND %s IS NOT NULL AND TRIM(%[2]s) != ''", pathColumn, basenameColumn))
				default:
					panic("unsupported string filter modifier")
				}
			}
		}
	}
}

func getPathSearchClause(pathColumn, basenameColumn, p string, addWildcards, not bool) sqlClause {
	if addWildcards {
		p = "%" + p + "%"
	}

	filepathColumn := fmt.Sprintf("%s || '%s' || %s", pathColumn, string(filepath.Separator), basenameColumn)
	ret := makeClause(fmt.Sprintf("%s LIKE ?", filepathColumn), p)

	if not {
		ret = ret.not()
	}

	return ret
}

// getPathSearchClauseMany splits the query string p on whitespace
// Used for backwards compatibility for the includes/excludes modifiers
func getPathSearchClauseMany(pathColumn, basenameColumn, p string, addWildcards, not bool) sqlClause {
	q := strings.TrimSpace(p)
	trimmedQuery := strings.Trim(q, "\"")

	if trimmedQuery == q {
		q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")
		queryWords := strings.Split(q, " ")

		var ret []sqlClause
		// Search for any word
		for _, word := range queryWords {
			ret = append(ret, getPathSearchClause(pathColumn, basenameColumn, word, addWildcards, not))
		}

		if !not {
			return orClauses(ret...)
		}

		return andClauses(ret...)
	}

	return getPathSearchClause(pathColumn, basenameColumn, trimmedQuery, addWildcards, not)
}

func intCriterionHandler(c *models.IntCriterionInput, column string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}
			clause, args := getIntCriterionWhereClause(column, *c)
			f.addWhere(clause, args...)
		}
	}
}

func floatCriterionHandler(c *models.FloatCriterionInput, column string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}
			clause, args := getFloatCriterionWhereClause(column, *c)
			f.addWhere(clause, args...)
		}
	}
}

func floatIntCriterionHandler(durationFilter *models.IntCriterionInput, column string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if durationFilter != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}
			clause, args := getIntCriterionWhereClause("cast("+column+" as int)", *durationFilter)
			f.addWhere(clause, args...)
		}
	}
}

func boolCriterionHandler(c *bool, column string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}
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

type dateCriterionHandler struct {
	c      *models.DateCriterionInput
	column string
	joinFn func(f *filterBuilder)
}

func (h *dateCriterionHandler) handle(ctx context.Context, f *filterBuilder) {
	if h.c != nil {
		if h.joinFn != nil {
			h.joinFn(f)
		}
		clause, args := getDateCriterionWhereClause(h.column, *h.c)
		f.addWhere(clause, args...)
	}
}

type timestampCriterionHandler struct {
	c      *models.TimestampCriterionInput
	column string
	joinFn func(f *filterBuilder)
}

func (h *timestampCriterionHandler) handle(ctx context.Context, f *filterBuilder) {
	if h.c != nil {
		if h.joinFn != nil {
			h.joinFn(f)
		}
		clause, args := getTimestampCriterionWhereClause(h.column, *h.c)
		f.addWhere(clause, args...)
	}
}

func yearFilterCriterionHandler(year *models.IntCriterionInput, col string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if year != nil && year.Modifier.IsValid() {
			clause, args := getIntCriterionWhereClause("cast(strftime('%Y', "+col+") as int)", *year)
			f.addWhere(clause, args...)
		}
	}
}

func resolutionCriterionHandler(resolution *models.ResolutionCriterionInput, heightColumn string, widthColumn string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if resolution != nil && resolution.Value.IsValid() {
			if addJoinFn != nil {
				addJoinFn(f)
			}

			min := resolution.Value.GetMinResolution()
			max := resolution.Value.GetMaxResolution()

			widthHeight := fmt.Sprintf("MIN(%s, %s)", widthColumn, heightColumn)

			switch resolution.Modifier {
			case models.CriterionModifierEquals:
				f.addWhere(fmt.Sprintf("%s BETWEEN %d AND %d", widthHeight, min, max))
			case models.CriterionModifierNotEquals:
				f.addWhere(fmt.Sprintf("%s NOT BETWEEN %d AND %d", widthHeight, min, max))
			case models.CriterionModifierLessThan:
				f.addWhere(fmt.Sprintf("%s < %d", widthHeight, min))
			case models.CriterionModifierGreaterThan:
				f.addWhere(fmt.Sprintf("%s > %d", widthHeight, max))
			}
		}
	}
}

func orientationCriterionHandler(orientation *models.OrientationCriterionInput, heightColumn string, widthColumn string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if orientation != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}

			var clauses []sqlClause

			for _, v := range orientation.Value {
				// width mod height
				mod := ""
				switch v {
				case models.OrientationPortrait:
					mod = "<"
				case models.OrientationLandscape:
					mod = ">"
				case models.OrientationSquare:
					mod = "="
				}

				if mod != "" {
					clauses = append(clauses, makeClause(fmt.Sprintf("%s %s %s", widthColumn, mod, heightColumn)))
				}
			}

			if len(clauses) > 0 {
				f.whereClauses = append(f.whereClauses, orClauses(clauses...))
			}
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

func (m *joinedMultiCriterionHandlerBuilder) handler(c *models.MultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			// make local copy so we can modify it
			criterion := *c

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

			if len(criterion.Value) == 0 && len(criterion.Excludes) == 0 {
				return
			}

			// combine excludes if excludes modifier is selected
			if criterion.Modifier == models.CriterionModifierExcludes {
				criterion.Modifier = models.CriterionModifierIncludesAll
				criterion.Excludes = append(criterion.Excludes, criterion.Value...)
				criterion.Value = nil
			}

			if len(criterion.Value) > 0 {
				whereClause := ""
				havingClause := ""

				var args []interface{}
				for _, tagID := range criterion.Value {
					args = append(args, tagID)
				}

				switch criterion.Modifier {
				case models.CriterionModifierIncludes:
					// includes any of the provided ids
					m.addJoinTable(f)
					whereClause = fmt.Sprintf("%s.%s IN %s", joinAlias, m.foreignFK, getInBinding(len(criterion.Value)))
				case models.CriterionModifierEquals:
					// includes only the provided ids
					m.addJoinTable(f)
					whereClause = utils.StrFormat("{joinAlias}.{foreignFK} IN {inBinding} AND (SELECT COUNT(*) FROM {joinTable} s WHERE s.{primaryFK} = {primaryTable}.id) = ?", utils.StrFormatMap{
						"joinAlias":    joinAlias,
						"foreignFK":    m.foreignFK,
						"inBinding":    getInBinding(len(criterion.Value)),
						"joinTable":    m.joinTable,
						"primaryFK":    m.primaryFK,
						"primaryTable": m.primaryTable,
					})
					havingClause = fmt.Sprintf("count(distinct %s.%s) IS %d", joinAlias, m.foreignFK, len(criterion.Value))
					args = append(args, len(criterion.Value))
				case models.CriterionModifierNotEquals:
					f.setError(fmt.Errorf("not equals modifier is not supported for multi criterion input"))
				case models.CriterionModifierIncludesAll:
					// includes all of the provided ids
					m.addJoinTable(f)
					whereClause = fmt.Sprintf("%s.%s IN %s", joinAlias, m.foreignFK, getInBinding(len(criterion.Value)))
					havingClause = fmt.Sprintf("count(distinct %s.%s) IS %d", joinAlias, m.foreignFK, len(criterion.Value))
				}

				f.addWhere(whereClause, args...)
				f.addHaving(havingClause)
			}

			if len(criterion.Excludes) > 0 {
				var args []interface{}
				for _, tagID := range criterion.Excludes {
					args = append(args, tagID)
				}

				// excludes all of the provided ids
				// need to use actual join table name for this
				// <primaryTable>.id NOT IN (select <joinTable>.<primaryFK> from <joinTable> where <joinTable>.<foreignFK> in <values>)
				whereClause := fmt.Sprintf("%[1]s.id NOT IN (SELECT %[3]s.%[2]s from %[3]s where %[3]s.%[4]s in %[5]s)", m.primaryTable, m.primaryFK, m.joinTable, m.foreignFK, getInBinding(len(criterion.Excludes)))

				f.addWhere(whereClause, args...)
			}
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
	primaryTable string
	// foreign key of the primary object on the join table
	primaryFK string
	// table joining primary and foreign objects
	joinTable string
	// string field on the join table
	stringColumn string

	addJoinTable   func(f *filterBuilder)
	excludeHandler func(f *filterBuilder, criterion *models.StringCriterionInput)
}

func (m *stringListCriterionHandlerBuilder) handler(criterion *models.StringCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			if criterion.Modifier == models.CriterionModifierExcludes {
				// special handling for excludes
				if m.excludeHandler != nil {
					m.excludeHandler(f, criterion)
					return
				}

				// excludes all of the provided values
				// need to use actual join table name for this
				// <primaryTable>.id NOT IN (select <joinTable>.<primaryFK> from <joinTable> where <joinTable>.<foreignFK> in <values>)
				whereClause := utils.StrFormat("{primaryTable}.id NOT IN (SELECT {joinTable}.{primaryFK} from {joinTable} where {joinTable}.{stringColumn} LIKE ?)",
					utils.StrFormatMap{
						"primaryTable": m.primaryTable,
						"joinTable":    m.joinTable,
						"primaryFK":    m.primaryFK,
						"stringColumn": m.stringColumn,
					},
				)

				f.addWhere(whereClause, "%"+criterion.Value+"%")

				// TODO - should we also exclude null values?
				// m.addJoinTable(f)
				// stringCriterionHandler(&models.StringCriterionInput{
				// 	Modifier: models.CriterionModifierNotNull,
				// }, m.joinTable+"."+m.stringColumn)(ctx, f)
			} else {
				m.addJoinTable(f)
				stringCriterionHandler(criterion, m.joinTable+"."+m.stringColumn)(ctx, f)
			}
		}
	}
}

func studioCriterionHandler(primaryTable string, studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if studios == nil {
			return
		}

		studiosCopy := *studios
		switch studiosCopy.Modifier {
		case models.CriterionModifierEquals:
			studiosCopy.Modifier = models.CriterionModifierIncludesAll
		case models.CriterionModifierNotEquals:
			studiosCopy.Modifier = models.CriterionModifierExcludes
		}

		hh := hierarchicalMultiCriterionHandlerBuilder{
			primaryTable: primaryTable,
			foreignTable: studioTable,
			foreignFK:    studioIDColumn,
			parentFK:     "parent_id",
		}

		hh.handler(&studiosCopy)(ctx, f)
	}
}

type hierarchicalMultiCriterionHandlerBuilder struct {
	primaryTable string
	foreignTable string
	foreignFK    string

	parentFK       string
	childFK        string
	relationsTable string
}

func getHierarchicalValues(ctx context.Context, values []string, table, relationsTable, parentFK string, childFK string, depth *int) (string, error) {
	var args []interface{}

	if parentFK == "" {
		parentFK = "parent_id"
	}
	if childFK == "" {
		childFK = "child_id"
	}

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
			return "VALUES" + strings.Join(valuesClauses, ","), nil
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
		"childFK":         childFK,
		"depthCondition":  depthCondition,
		"unionClause":     "",
	}

	if relationsTable != "" {
		withClauseMap["recursiveSelect"] = utils.StrFormat(`SELECT p.root_id, c.{childFK}, depth + 1 FROM {relationsTable} AS c
INNER JOIN items as p ON c.{parentFK} = p.item_id
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

	var valuesClause sql.NullString
	err := dbWrapper.Get(ctx, &valuesClause, query, args...)
	if err != nil {
		return "", fmt.Errorf("failed to get hierarchical values: %w", err)
	}

	// if no values are found, just return a values string with the values only
	if !valuesClause.Valid {
		for i, value := range values {
			values[i] = fmt.Sprintf("(%s, %s)", value, value)
		}
		valuesClause.String = "VALUES" + strings.Join(values, ",")
	}

	return valuesClause.String, nil
}

func addHierarchicalConditionClauses(f *filterBuilder, criterion models.HierarchicalMultiCriterionInput, table, idColumn string) {
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

func (m *hierarchicalMultiCriterionHandlerBuilder) handler(c *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			// make a copy so we don't modify the original
			criterion := *c

			// don't support equals/not equals
			if criterion.Modifier == models.CriterionModifierEquals || criterion.Modifier == models.CriterionModifierNotEquals {
				f.setError(fmt.Errorf("modifier %s is not supported for hierarchical multi criterion", criterion.Modifier))
				return
			}

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

			if len(criterion.Value) == 0 && len(criterion.Excludes) == 0 {
				return
			}

			// combine excludes if excludes modifier is selected
			if criterion.Modifier == models.CriterionModifierExcludes {
				criterion.Modifier = models.CriterionModifierIncludesAll
				criterion.Excludes = append(criterion.Excludes, criterion.Value...)
				criterion.Value = nil
			}

			if len(criterion.Value) > 0 {
				valuesClause, err := getHierarchicalValues(ctx, criterion.Value, m.foreignTable, m.relationsTable, m.parentFK, m.childFK, criterion.Depth)
				if err != nil {
					f.setError(err)
					return
				}

				switch criterion.Modifier {
				case models.CriterionModifierIncludes:
					f.addWhere(fmt.Sprintf("%s.%s IN (SELECT column2 FROM (%s))", m.primaryTable, m.foreignFK, valuesClause))
				case models.CriterionModifierIncludesAll:
					f.addWhere(fmt.Sprintf("%s.%s IN (SELECT column2 FROM (%s))", m.primaryTable, m.foreignFK, valuesClause))
					f.addHaving(fmt.Sprintf("count(distinct %s.%s) IS %d", m.primaryTable, m.foreignFK, len(criterion.Value)))
				}
			}

			if len(criterion.Excludes) > 0 {
				valuesClause, err := getHierarchicalValues(ctx, criterion.Excludes, m.foreignTable, m.relationsTable, m.parentFK, m.childFK, criterion.Depth)
				if err != nil {
					f.setError(err)
					return
				}

				f.addWhere(fmt.Sprintf("%s.%s NOT IN (SELECT column2 FROM (%s)) OR %[1]s.%[2]s IS NULL", m.primaryTable, m.foreignFK, valuesClause))
			}
		}
	}
}

type joinedHierarchicalMultiCriterionHandlerBuilder struct {
	primaryTable string
	primaryKey   string
	foreignTable string
	foreignFK    string

	parentFK       string
	childFK        string
	relationsTable string

	joinAs    string
	joinTable string
	primaryFK string
}

func (m *joinedHierarchicalMultiCriterionHandlerBuilder) addHierarchicalConditionClauses(f *filterBuilder, criterion models.HierarchicalMultiCriterionInput, table, idColumn string) {
	primaryKey := m.primaryKey
	if primaryKey == "" {
		primaryKey = "id"
	}

	switch criterion.Modifier {
	case models.CriterionModifierEquals:
		// includes only the provided ids
		f.addWhere(fmt.Sprintf("%s.%s IS NOT NULL", table, idColumn))
		f.addHaving(fmt.Sprintf("count(distinct %s.%s) IS %d", table, idColumn, len(criterion.Value)))
		f.addWhere(utils.StrFormat("(SELECT COUNT(*) FROM {joinTable} s WHERE s.{primaryFK} = {primaryTable}.{primaryKey}) = ?", utils.StrFormatMap{
			"joinTable":    m.joinTable,
			"primaryFK":    m.primaryFK,
			"primaryTable": m.primaryTable,
			"primaryKey":   primaryKey,
		}), len(criterion.Value))
	case models.CriterionModifierNotEquals:
		f.setError(fmt.Errorf("not equals modifier is not supported for hierarchical multi criterion input"))
	default:
		addHierarchicalConditionClauses(f, criterion, table, idColumn)
	}
}

func (m *joinedHierarchicalMultiCriterionHandlerBuilder) handler(c *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			// make a copy so we don't modify the original
			criterion := *c
			joinAlias := m.joinAs
			primaryKey := m.primaryKey
			if primaryKey == "" {
				primaryKey = "id"
			}

			if criterion.Modifier == models.CriterionModifierEquals && criterion.Depth != nil && *criterion.Depth != 0 {
				f.setError(fmt.Errorf("depth is not supported for equals modifier in hierarchical multi criterion input"))
				return
			}

			if criterion.Modifier == models.CriterionModifierIsNull || criterion.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if criterion.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin(m.joinTable, joinAlias, fmt.Sprintf("%s.%s = %s.%s", joinAlias, m.primaryFK, m.primaryTable, primaryKey))

				f.addWhere(utils.StrFormat("{table}.{column} IS {not} NULL", utils.StrFormatMap{
					"table":  joinAlias,
					"column": m.foreignFK,
					"not":    notClause,
				}))
				return
			}

			// combine excludes if excludes modifier is selected
			if criterion.Modifier == models.CriterionModifierExcludes {
				criterion.Modifier = models.CriterionModifierIncludesAll
				criterion.Excludes = append(criterion.Excludes, criterion.Value...)
				criterion.Value = nil
			}

			if len(criterion.Value) == 0 && len(criterion.Excludes) == 0 {
				return
			}

			if len(criterion.Value) > 0 {
				valuesClause, err := getHierarchicalValues(ctx, criterion.Value, m.foreignTable, m.relationsTable, m.parentFK, m.childFK, criterion.Depth)
				if err != nil {
					f.setError(err)
					return
				}

				joinTable := utils.StrFormat(`(
		SELECT j.*, d.column1 AS root_id, d.column2 AS item_id FROM {joinTable} AS j
		INNER JOIN ({valuesClause}) AS d ON j.{foreignFK} = d.column2
	)
	`, utils.StrFormatMap{
					"joinTable":    m.joinTable,
					"foreignFK":    m.foreignFK,
					"valuesClause": valuesClause,
				})

				f.addLeftJoin(joinTable, joinAlias, fmt.Sprintf("%s.%s = %s.%s", joinAlias, m.primaryFK, m.primaryTable, primaryKey))

				m.addHierarchicalConditionClauses(f, criterion, joinAlias, "root_id")
			}

			if len(criterion.Excludes) > 0 {
				valuesClause, err := getHierarchicalValues(ctx, criterion.Excludes, m.foreignTable, m.relationsTable, m.parentFK, m.childFK, criterion.Depth)
				if err != nil {
					f.setError(err)
					return
				}

				joinTable := utils.StrFormat(`(
		SELECT j2.*, e.column1 AS root_id, e.column2 AS item_id FROM {joinTable} AS j2
		INNER JOIN ({valuesClause}) AS e ON j2.{foreignFK} = e.column2
	)
	`, utils.StrFormatMap{
					"joinTable":    m.joinTable,
					"foreignFK":    m.foreignFK,
					"valuesClause": valuesClause,
				})

				joinAlias2 := joinAlias + "2"

				f.addLeftJoin(joinTable, joinAlias2, fmt.Sprintf("%s.%s = %s.%s", joinAlias2, m.primaryFK, m.primaryTable, primaryKey))

				// modify for exclusion
				criterionCopy := criterion
				criterionCopy.Modifier = models.CriterionModifierExcludes
				criterionCopy.Value = c.Excludes

				m.addHierarchicalConditionClauses(f, criterionCopy, joinAlias2, "root_id")
			}
		}
	}
}

type joinedPerformerTagsHandler struct {
	criterion *models.HierarchicalMultiCriterionInput

	primaryTable   string // eg scenes
	joinTable      string // eg performers_scenes
	joinPrimaryKey string // eg scene_id
}

func (h *joinedPerformerTagsHandler) handle(ctx context.Context, f *filterBuilder) {
	tags := h.criterion

	if tags != nil {
		criterion := tags.CombineExcludes()

		// validate the modifier
		switch criterion.Modifier {
		case models.CriterionModifierIncludesAll, models.CriterionModifierIncludes, models.CriterionModifierExcludes, models.CriterionModifierIsNull, models.CriterionModifierNotNull:
			// valid
		default:
			f.setError(fmt.Errorf("invalid modifier %s for performer tags", criterion.Modifier))
		}

		strFormatMap := utils.StrFormatMap{
			"primaryTable":   h.primaryTable,
			"joinTable":      h.joinTable,
			"joinPrimaryKey": h.joinPrimaryKey,
			"inBinding":      getInBinding(len(criterion.Value)),
		}

		if criterion.Modifier == models.CriterionModifierIsNull || criterion.Modifier == models.CriterionModifierNotNull {
			var notClause string
			if criterion.Modifier == models.CriterionModifierNotNull {
				notClause = "NOT"
			}

			f.addLeftJoin(h.joinTable, "", utils.StrFormat("{primaryTable}.id = {joinTable}.{joinPrimaryKey}", strFormatMap))
			f.addLeftJoin("performers_tags", "", utils.StrFormat("{joinTable}.performer_id = performers_tags.performer_id", strFormatMap))

			f.addWhere(fmt.Sprintf("performers_tags.tag_id IS %s NULL", notClause))
			return
		}

		if len(criterion.Value) == 0 && len(criterion.Excludes) == 0 {
			return
		}

		if len(criterion.Value) > 0 {
			valuesClause, err := getHierarchicalValues(ctx, criterion.Value, tagTable, "tags_relations", "", "", criterion.Depth)
			if err != nil {
				f.setError(err)
				return
			}

			f.addWith(utils.StrFormat(`performer_tags AS (
SELECT ps.{joinPrimaryKey} as primaryID, t.column1 AS root_tag_id FROM {joinTable} ps
INNER JOIN performers_tags pt ON pt.performer_id = ps.performer_id
INNER JOIN (`+valuesClause+`) t ON t.column2 = pt.tag_id
)`, strFormatMap))

			f.addLeftJoin("performer_tags", "", utils.StrFormat("performer_tags.primaryID = {primaryTable}.id", strFormatMap))

			addHierarchicalConditionClauses(f, criterion, "performer_tags", "root_tag_id")
		}

		if len(criterion.Excludes) > 0 {
			valuesClause, err := getHierarchicalValues(ctx, criterion.Excludes, tagTable, "tags_relations", "", "", criterion.Depth)
			if err != nil {
				f.setError(err)
				return
			}

			clause := utils.StrFormat("{primaryTable}.id NOT IN (SELECT {joinTable}.{joinPrimaryKey} FROM {joinTable} INNER JOIN performers_tags ON {joinTable}.performer_id = performers_tags.performer_id WHERE performers_tags.tag_id IN (SELECT column2 FROM (%s)))", strFormatMap)
			f.addWhere(fmt.Sprintf(clause, valuesClause))
		}
	}
}

type stashIDCriterionHandler struct {
	c                 *models.StashIDCriterionInput
	stashIDRepository *stashIDRepository
	stashIDTableAs    string
	parentIDCol       string
}

func (h *stashIDCriterionHandler) handle(ctx context.Context, f *filterBuilder) {
	if h.c == nil {
		return
	}

	stashIDRepo := h.stashIDRepository
	t := stashIDRepo.tableName
	if h.stashIDTableAs != "" {
		t = h.stashIDTableAs
	}

	joinClause := fmt.Sprintf("%s.%s = %s", t, stashIDRepo.idColumn, h.parentIDCol)
	if h.c.Endpoint != nil && *h.c.Endpoint != "" {
		joinClause += fmt.Sprintf(" AND %s.endpoint = '%s'", t, *h.c.Endpoint)
	}

	f.addLeftJoin(stashIDRepo.tableName, h.stashIDTableAs, joinClause)

	v := ""
	if h.c.StashID != nil {
		v = *h.c.StashID
	}

	stringCriterionHandler(&models.StringCriterionInput{
		Value:    v,
		Modifier: h.c.Modifier,
	}, t+".stash_id")(ctx, f)
}

type relatedFilterHandler struct {
	relatedIDCol   string
	relatedRepo    repository
	relatedHandler criterionHandler
	joinFn         func(f *filterBuilder)
}

func (h *relatedFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	ff := filterBuilderFromHandler(ctx, h.relatedHandler)
	if ff.err != nil {
		f.setError(ff.err)
		return
	}

	if ff.empty() {
		return
	}

	subQuery := h.relatedRepo.newQuery()
	selectIDs(&subQuery, subQuery.repository.tableName)
	if err := subQuery.addFilter(ff); err != nil {
		f.setError(err)
		return
	}

	if h.joinFn != nil {
		h.joinFn(f)
	}

	f.addWhere(fmt.Sprintf("%s IN ("+subQuery.toSQL(false)+")", h.relatedIDCol), subQuery.args...)
}
