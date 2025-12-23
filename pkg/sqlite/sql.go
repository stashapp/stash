package sqlite

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

func selectAll(tableName string) string {
	idColumn := getColumn(tableName, "*")
	return "SELECT " + idColumn + " FROM " + tableName + " "
}

func distinctIDs(qb *queryBuilder, tableName string) {
	qb.addColumn("DISTINCT " + getColumn(tableName, "id"))
	qb.from = tableName
}

func selectIDs(qb *queryBuilder, tableName string) {
	qb.addColumn(getColumn(tableName, "id"))
	qb.from = tableName
}

func getColumn(tableName string, columnName string) string {
	return tableName + "." + columnName
}

func getPagination(findFilter *models.FindFilterType) string {
	if findFilter == nil {
		panic("nil find filter for pagination")
	}

	if findFilter.IsGetAll() {
		return " "
	}

	return getPaginationSQL(findFilter.GetPage(), findFilter.GetPageSize())
}

func getPaginationSQL(page int, perPage int) string {
	page = (page - 1) * perPage
	return " LIMIT " + strconv.Itoa(perPage) + " OFFSET " + strconv.Itoa(page) + " "
}

const randomSeedPrefix = "random_" // prefix for random sort

type sortOptions []string

func (o sortOptions) validateSort(sort string) error {
	if strings.HasPrefix(sort, randomSeedPrefix) {
		// seed as a parameter from the UI
		seedStr := sort[len(randomSeedPrefix):]
		_, err := strconv.ParseUint(seedStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid random seed: %s", seedStr)
		}
		return nil
	}

	for _, v := range o {
		if v == sort {
			return nil
		}
	}

	return fmt.Errorf("invalid sort: %s", sort)
}

func getSortDirection(direction string) string {
	if direction != "ASC" && direction != "DESC" {
		return "ASC"
	} else {
		return direction
	}
}
func getSort(sort string, direction string, tableName string) string {
	direction = getSortDirection(direction)

	switch {
	case strings.HasSuffix(sort, "_count"):
		var relationTableName = strings.TrimSuffix(sort, "_count") // TODO: pluralize?
		colName := getColumn(relationTableName, "id")
		return " ORDER BY COUNT(distinct " + colName + ") " + direction
	case strings.Compare(sort, "filesize") == 0:
		colName := getColumn(tableName, "size")
		return " ORDER BY " + colName + " " + direction
	case strings.HasPrefix(sort, randomSeedPrefix):
		// seed as a parameter from the UI
		seedStr := sort[len(randomSeedPrefix):]
		seed, err := strconv.ParseUint(seedStr, 10, 64)
		if err != nil {
			// fallback to a random seed
			seed = rand.Uint64()
		}
		return getRandomSort(tableName, direction, seed)
	case strings.Compare(sort, "random") == 0:
		return getRandomSort(tableName, direction, rand.Uint64())
	default:
		colName := getColumn(tableName, sort)
		if strings.Contains(sort, ".") {
			colName = sort
		}
		if strings.Compare(sort, "name") == 0 {
			return " ORDER BY " + colName + " COLLATE NATURAL_CI " + direction
		}
		if strings.Compare(sort, "title") == 0 {
			return " ORDER BY " + colName + " COLLATE NATURAL_CI " + direction
		}

		return " ORDER BY " + colName + " " + direction
	}
}

func getRandomSort(tableName string, direction string, seed uint64) string {
	// cap seed at 10^8
	seed %= 1e8

	colName := getColumn(tableName, "id")

	// https://stackoverflow.com/questions/21949795#comment33255354_21949859
	// p1 := 52959209
	// p2 := 1047483763
	// p3 := 2147483647
	// n := <colName>
	// ORDER BY ((n+seed)*(n+seed)*p1 + (n+seed)*p2) % p3
	// since sqlite converts overflowing numbers to reals, a custom db function that uses uints with overflow should be faster,
	// however in practice the overhead of calling a custom function vastly outweighs the benefits
	return fmt.Sprintf(" ORDER BY mod((%[1]s + %[2]d) * (%[1]s + %[2]d) * 52959209 + (%[1]s + %[2]d) * 1047483763, 2147483647) %[3]s", colName, seed, direction)
}

func getCountSort(primaryTable, joinTable, primaryFK, direction string) string {
	return fmt.Sprintf(" ORDER BY (SELECT COUNT(*) FROM %s AS sort WHERE sort.%s = %s.id) %s", joinTable, primaryFK, primaryTable, getSortDirection(direction))
}

func getStringSearchClause(columns []string, q string, not bool) sqlClause {
	var likeClauses []string
	var args []interface{}

	notStr := ""
	binaryType := " OR "
	if not {
		notStr = " NOT"
		binaryType = " AND "
	}
	q = strings.TrimSpace(q)
	trimmedQuery := strings.Trim(q, "\"")

	if trimmedQuery == q {
		q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")
		queryWords := strings.Split(q, " ")
		// Search for any word
		for _, word := range queryWords {
			for _, column := range columns {
				likeClauses = append(likeClauses, column+notStr+" LIKE ?")
				args = append(args, "%"+word+"%")
			}
		}
	} else {
		// Search the exact query
		for _, column := range columns {
			likeClauses = append(likeClauses, column+notStr+" LIKE ?")
			args = append(args, "%"+trimmedQuery+"%")
		}
	}
	likes := strings.Join(likeClauses, binaryType)

	return makeClause("("+likes+")", args...)
}

func getEnumSearchClause(column string, enumVals []string, not bool) sqlClause {
	var args []interface{}

	notStr := ""
	if not {
		notStr = " NOT"
	}

	clause := fmt.Sprintf("(%s%s IN %s)", column, notStr, getInBinding(len(enumVals)))
	for _, enumVal := range enumVals {
		args = append(args, enumVal)
	}

	return makeClause(clause, args...)
}

func getInBinding(length int) string {
	bindings := strings.Repeat("?, ", length)
	bindings = strings.TrimRight(bindings, ", ")
	return "(" + bindings + ")"
}

func getIntCriterionWhereClause(column string, input models.IntCriterionInput) (string, []interface{}) {
	return getIntWhereClause(column, input.Modifier, input.Value, input.Value2)
}

func getIntWhereClause(column string, modifier models.CriterionModifier, value int, upper *int) (string, []interface{}) {
	if upper == nil {
		u := 0
		upper = &u
	}

	args := []interface{}{value, *upper}
	return getNumericWhereClause(column, modifier, args)
}

func getFloatCriterionWhereClause(column string, input models.FloatCriterionInput) (string, []interface{}) {
	return getFloatWhereClause(column, input.Modifier, input.Value, input.Value2)
}

func getFloatWhereClause(column string, modifier models.CriterionModifier, value float64, upper *float64) (string, []interface{}) {
	if upper == nil {
		u := 0.0
		upper = &u
	}

	args := []interface{}{value, *upper}
	return getNumericWhereClause(column, modifier, args)
}

func getNumericWhereClause(column string, modifier models.CriterionModifier, args []interface{}) (string, []interface{}) {
	singleArgs := args[0:1]

	switch modifier {
	case models.CriterionModifierIsNull:
		return fmt.Sprintf("%s IS NULL", column), nil
	case models.CriterionModifierNotNull:
		return fmt.Sprintf("%s IS NOT NULL", column), nil
	case models.CriterionModifierEquals:
		return fmt.Sprintf("%s = ?", column), singleArgs
	case models.CriterionModifierNotEquals:
		return fmt.Sprintf("%s != ?", column), singleArgs
	case models.CriterionModifierBetween:
		return fmt.Sprintf("%s BETWEEN ? AND ?", column), args
	case models.CriterionModifierNotBetween:
		return fmt.Sprintf("%s NOT BETWEEN ? AND ?", column), args
	case models.CriterionModifierLessThan:
		return fmt.Sprintf("%s < ?", column), singleArgs
	case models.CriterionModifierGreaterThan:
		return fmt.Sprintf("%s > ?", column), singleArgs
	}

	panic("unsupported numeric modifier type " + modifier)
}

func getDateCriterionWhereClause(column string, input models.DateCriterionInput) (string, []interface{}) {
	return getDateWhereClause(column, input.Modifier, input.Value, input.Value2)
}

func getDateWhereClause(column string, modifier models.CriterionModifier, value string, upper *string) (string, []interface{}) {
	if upper == nil {
		u := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
		upper = &u
	}

	args := []interface{}{value}
	betweenArgs := []interface{}{value, *upper}

	switch modifier {
	case models.CriterionModifierIsNull:
		return fmt.Sprintf("(%s IS NULL OR %s = '')", column, column), nil
	case models.CriterionModifierNotNull:
		return fmt.Sprintf("(%s IS NOT NULL AND %s != '')", column, column), nil
	case models.CriterionModifierEquals:
		return fmt.Sprintf("%s = ?", column), args
	case models.CriterionModifierNotEquals:
		return fmt.Sprintf("%s != ?", column), args
	case models.CriterionModifierBetween:
		return fmt.Sprintf("%s BETWEEN ? AND ?", column), betweenArgs
	case models.CriterionModifierNotBetween:
		return fmt.Sprintf("%s NOT BETWEEN ? AND ?", column), betweenArgs
	case models.CriterionModifierLessThan:
		return fmt.Sprintf("%s < ?", column), args
	case models.CriterionModifierGreaterThan:
		return fmt.Sprintf("%s > ?", column), args
	}

	panic("unsupported date modifier type")
}

func getTimestampCriterionWhereClause(column string, input models.TimestampCriterionInput) (string, []interface{}) {
	return getTimestampWhereClause(column, input.Modifier, input.Value, input.Value2)
}

func getTimestampWhereClause(column string, modifier models.CriterionModifier, value string, upper *string) (string, []interface{}) {
	if upper == nil {
		u := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
		upper = &u
	}

	args := []interface{}{value}
	betweenArgs := []interface{}{value, *upper}

	switch modifier {
	case models.CriterionModifierIsNull:
		return fmt.Sprintf("%s IS NULL", column), nil
	case models.CriterionModifierNotNull:
		return fmt.Sprintf("%s IS NOT NULL", column), nil
	case models.CriterionModifierEquals:
		return fmt.Sprintf("%s = ?", column), args
	case models.CriterionModifierNotEquals:
		return fmt.Sprintf("%s != ?", column), args
	case models.CriterionModifierBetween:
		return fmt.Sprintf("%s BETWEEN ? AND ?", column), betweenArgs
	case models.CriterionModifierNotBetween:
		return fmt.Sprintf("%s NOT BETWEEN ? AND ?", column), betweenArgs
	case models.CriterionModifierLessThan:
		return fmt.Sprintf("%s < ?", column), args
	case models.CriterionModifierGreaterThan:
		return fmt.Sprintf("%s > ?", column), args
	}

	panic("unsupported date modifier type")
}

// returns where clause and having clause
func getMultiCriterionClause(primaryTable, foreignTable, joinTable, primaryFK, foreignFK string, criterion *models.MultiCriterionInput) (string, string) {
	whereClause := ""
	havingClause := ""
	switch criterion.Modifier {
	case models.CriterionModifierIncludes:
		// includes any of the provided ids
		if joinTable != "" {
			whereClause = joinTable + "." + foreignFK + " IN " + getInBinding(len(criterion.Value))
		} else {
			whereClause = foreignTable + ".id IN " + getInBinding(len(criterion.Value))
		}
	case models.CriterionModifierIncludesAll:
		// includes all of the provided ids
		if joinTable != "" {
			whereClause = joinTable + "." + foreignFK + " IN " + getInBinding(len(criterion.Value))
			havingClause = "count(distinct " + joinTable + "." + foreignFK + ") IS " + strconv.Itoa(len(criterion.Value))
		} else {
			whereClause = foreignTable + ".id IN " + getInBinding(len(criterion.Value))
			havingClause = "count(distinct " + foreignTable + ".id) IS " + strconv.Itoa(len(criterion.Value))
		}
	case models.CriterionModifierExcludes:
		// excludes all of the provided ids
		if joinTable != "" {
			whereClause = primaryTable + ".id not in (select " + joinTable + "." + primaryFK + " from " + joinTable + " where " + joinTable + "." + foreignFK + " in " + getInBinding(len(criterion.Value)) + ")"
		} else {
			whereClause = "not exists (select s.id from " + primaryTable + " as s where s.id = " + primaryTable + ".id and s." + foreignFK + " in " + getInBinding(len(criterion.Value)) + ")"
		}
	}

	return whereClause, havingClause
}

func getCountCriterionClause(primaryTable, joinTable, primaryFK string, criterion models.IntCriterionInput) (string, []interface{}) {
	lhs := fmt.Sprintf("(SELECT COUNT(*) FROM %s s WHERE s.%s = %s.id)", joinTable, primaryFK, primaryTable)
	return getIntCriterionWhereClause(lhs, criterion)
}

func coalesce(column string) string {
	return fmt.Sprintf("COALESCE(%s, '')", column)
}

func like(v string) string {
	return "%" + v + "%"
}

type sqlTable string

func (t sqlTable) Name() string {
	return string(t)
}

func (t sqlTable) Col(n string) string {
	return fmt.Sprintf("%s.%s", string(t), n)
}
