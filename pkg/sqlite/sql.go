package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

var randomSortFloat = rand.Float64()

func selectAll(tableName string) string {
	idColumn := getColumn(tableName, "*")
	return "SELECT " + idColumn + " FROM " + tableName + " "
}

func distinctIDs(qb *queryBuilder, tableName string) {
	qb.addColumn("DISTINCT " + getColumn(tableName, "id"))
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

func getSortDirection(direction string) string {
	if direction != "ASC" && direction != "DESC" {
		return "ASC"
	} else {
		return direction
	}
}
func getSort(sort string, direction string, tableName string) string {
	direction = getSortDirection(direction)

	const randomSeedPrefix = "random_"

	switch {
	case strings.HasSuffix(sort, "_count"):
		var relationTableName = strings.TrimSuffix(sort, "_count") // TODO: pluralize?
		colName := getColumn(relationTableName, "id")
		return " ORDER BY COUNT(distinct " + colName + ") " + direction
	case strings.Compare(sort, "filesize") == 0:
		colName := getColumn(tableName, "size")
		return " ORDER BY cast(" + colName + " as integer) " + direction
	case strings.Compare(sort, "perceptual_similarity") == 0:
		colName := getColumn(tableName, "phash")
		secondaryColName := getColumn(tableName, "size")
		return " ORDER BY " + colName + " " + direction + ", " + secondaryColName + " DESC"
	case strings.HasPrefix(sort, randomSeedPrefix):
		// seed as a parameter from the UI
		// turn the provided seed into a float
		seedStr := "0." + sort[len(randomSeedPrefix):]
		seed, err := strconv.ParseFloat(seedStr, 32)
		if err != nil {
			// fallback to default seed
			seed = randomSortFloat
		}
		return getRandomSort(tableName, direction, seed)
	case strings.Compare(sort, "random") == 0:
		return getRandomSort(tableName, direction, randomSortFloat)
	default:
		colName := getColumn(tableName, sort)
		var additional string
		if tableName == "scenes" {
			additional = ", bitrate DESC, framerate DESC, scenes.rating DESC, scenes.duration DESC"
		} else if tableName == "scene_markers" {
			additional = ", scene_markers.scene_id ASC, scene_markers.seconds ASC"
		}
		if strings.Compare(sort, "name") == 0 {
			return " ORDER BY " + colName + " COLLATE NOCASE " + direction + additional
		}
		if strings.Compare(sort, "title") == 0 {
			return " ORDER BY " + colName + " COLLATE NATURAL_CS " + direction + additional
		}

		return " ORDER BY " + colName + " " + direction + additional
	}
}

func getRandomSort(tableName string, direction string, seed float64) string {
	// https://stackoverflow.com/a/24511461
	colName := getColumn(tableName, "id")
	randomSortString := strconv.FormatFloat(seed, 'f', 16, 32)
	return " ORDER BY " + "(substr(" + colName + " * " + randomSortString + ", length(" + colName + ") + 2))" + " " + direction
}

func getCountSort(primaryTable, joinTable, primaryFK, direction string) string {
	return fmt.Sprintf(" ORDER BY (SELECT COUNT(*) FROM %s WHERE %s = %s.id) %s", joinTable, primaryFK, primaryTable, getSortDirection(direction))
}

func getSearchBinding(columns []string, q string, not bool) (string, []interface{}) {
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

	return "(" + likes + ")", args
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

	panic("unsupported int modifier type")
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

func getImage(ctx context.Context, tx dbi, query string, args ...interface{}) ([]byte, error) {
	rows, err := tx.Queryx(ctx, query, args...)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	defer rows.Close()

	var ret []byte
	if rows.Next() {
		if err := rows.Scan(&ret); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func coalesce(column string) string {
	return fmt.Sprintf("COALESCE(%s, '')", column)
}

func like(v string) string {
	return "%" + v + "%"
}
