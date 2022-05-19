package sqlite

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

var testCtx = context.Background()

func TestJoinsAddJoin(t *testing.T) {
	var joins joins

	// add a single join
	joins.add(join{table: "test"})

	assert := assert.New(t)

	// ensure join was added
	assert.Len(joins, 1)

	// add the same join and another
	joins.add([]join{
		{
			table: "test",
		},
		{
			table: "foo",
		},
	}...)

	// should have added a single join
	assert.Len(joins, 2)
}

func TestFilterBuilderAnd(t *testing.T) {
	assert := assert.New(t)

	f := &filterBuilder{}
	other := &filterBuilder{}
	newBuilder := &filterBuilder{}

	// and should set the subFilter
	f.and(other)
	assert.Equal(other, f.subFilter)
	assert.Nil(f.getError())

	// and should set error if and is set
	f.and(newBuilder)
	assert.Equal(other, f.subFilter)
	assert.Equal(errSubFilterAlreadySet, f.getError())

	// and should set error if or is set
	// and should not set subFilter if or is set
	f = &filterBuilder{}
	f.or(other)
	f.and(newBuilder)
	assert.Equal(other, f.subFilter)
	assert.Equal(errSubFilterAlreadySet, f.getError())

	// and should set error if not is set
	// and should not set subFilter if not is set
	f = &filterBuilder{}
	f.not(other)
	f.and(newBuilder)
	assert.Equal(other, f.subFilter)
	assert.Equal(errSubFilterAlreadySet, f.getError())
}

func TestFilterBuilderOr(t *testing.T) {
	assert := assert.New(t)

	f := &filterBuilder{}
	other := &filterBuilder{}
	newBuilder := &filterBuilder{}

	// or should set the orFilter
	f.or(other)
	assert.Equal(other, f.subFilter)
	assert.Nil(f.getError())

	// or should set error if or is set
	f.or(newBuilder)
	assert.Equal(newBuilder, f.subFilter)
	assert.Equal(errSubFilterAlreadySet, f.getError())

	// or should set error if and is set
	// or should not set subFilter if and is set
	f = &filterBuilder{}
	f.and(other)
	f.or(newBuilder)
	assert.Equal(other, f.subFilter)
	assert.Equal(errSubFilterAlreadySet, f.getError())

	// or should set error if not is set
	// or should not set subFilter if not is set
	f = &filterBuilder{}
	f.not(other)
	f.or(newBuilder)
	assert.Equal(other, f.subFilter)
	assert.Equal(errSubFilterAlreadySet, f.getError())
}

func TestFilterBuilderNot(t *testing.T) {
	assert := assert.New(t)

	f := &filterBuilder{}
	other := &filterBuilder{}
	newBuilder := &filterBuilder{}

	// not should set the subFilter
	f.not(other)
	// ensure and filter is set
	assert.Equal(other, f.subFilter)
	assert.Nil(f.getError())

	// not should set error if not is set
	f.not(newBuilder)
	assert.Equal(newBuilder, f.subFilter)
	assert.Equal(errSubFilterAlreadySet, f.getError())

	// not should set error if and is set
	// not should not set subFilter if and is set
	f = &filterBuilder{}
	f.and(other)
	f.not(newBuilder)
	assert.Equal(other, f.subFilter)
	assert.Equal(errSubFilterAlreadySet, f.getError())

	// not should set error if or is set
	// not should not set subFilter if or is set
	f = &filterBuilder{}
	f.or(other)
	f.not(newBuilder)
	assert.Equal(other, f.subFilter)
	assert.Equal(errSubFilterAlreadySet, f.getError())
}

func TestAddJoin(t *testing.T) {
	assert := assert.New(t)

	f := &filterBuilder{}

	const (
		table1Name = "table1Name"
		table2Name = "table2Name"

		as1Name = "as1"
		as2Name = "as2"

		onClause = "onClause1"
	)

	f.addLeftJoin(table1Name, as1Name, onClause)

	// ensure join is added
	assert.Len(f.joins, 1)
	assert.Equal(fmt.Sprintf("LEFT JOIN %s AS %s ON %s", table1Name, as1Name, onClause), f.joins[0].toSQL())

	// ensure join with same as is not added
	f.addLeftJoin(table2Name, as1Name, onClause)
	assert.Len(f.joins, 1)

	// ensure same table with different alias can be added
	f.addLeftJoin(table1Name, as2Name, onClause)
	assert.Len(f.joins, 2)
	assert.Equal(fmt.Sprintf("LEFT JOIN %s AS %s ON %s", table1Name, as2Name, onClause), f.joins[1].toSQL())

	// ensure table without alias can be added if tableName != existing alias/tableName
	f.addLeftJoin(table1Name, "", onClause)
	assert.Len(f.joins, 3)
	assert.Equal(fmt.Sprintf("LEFT JOIN %s ON %s", table1Name, onClause), f.joins[2].toSQL())

	// ensure table with alias == table name of a join without alias is not added
	f.addLeftJoin(table2Name, table1Name, onClause)
	assert.Len(f.joins, 3)

	// ensure table without alias cannot be added if tableName == existing alias
	f.addLeftJoin(as2Name, "", onClause)
	assert.Len(f.joins, 3)

	// ensure AS is not used if same as table name
	f.addLeftJoin(table2Name, table2Name, onClause)
	assert.Len(f.joins, 4)
	assert.Equal(fmt.Sprintf("LEFT JOIN %s ON %s", table2Name, onClause), f.joins[3].toSQL())
}

func TestAddWhere(t *testing.T) {
	assert := assert.New(t)

	f := &filterBuilder{}

	// ensure empty sql adds nothing
	f.addWhere("")
	assert.Len(f.whereClauses, 0)

	const whereClause = "a = b"
	var args = []interface{}{"1", "2"}

	// ensure addWhere sets where clause and args
	f.addWhere(whereClause, args...)
	assert.Len(f.whereClauses, 1)
	assert.Equal(whereClause, f.whereClauses[0].sql)
	assert.Equal(args, f.whereClauses[0].args)

	// ensure addWhere without args sets where clause
	f.addWhere(whereClause)
	assert.Len(f.whereClauses, 2)
	assert.Equal(whereClause, f.whereClauses[1].sql)
	assert.Len(f.whereClauses[1].args, 0)
}

func TestAddHaving(t *testing.T) {
	assert := assert.New(t)

	f := &filterBuilder{}

	// ensure empty sql adds nothing
	f.addHaving("")
	assert.Len(f.havingClauses, 0)

	const havingClause = "a = b"
	var args = []interface{}{"1", "2"}

	// ensure addWhere sets where clause and args
	f.addHaving(havingClause, args...)
	assert.Len(f.havingClauses, 1)
	assert.Equal(havingClause, f.havingClauses[0].sql)
	assert.Equal(args, f.havingClauses[0].args)

	// ensure addWhere without args sets where clause
	f.addHaving(havingClause)
	assert.Len(f.havingClauses, 2)
	assert.Equal(havingClause, f.havingClauses[1].sql)
	assert.Len(f.havingClauses[1].args, 0)
}

func TestGenerateWhereClauses(t *testing.T) {
	assert := assert.New(t)

	f := &filterBuilder{}

	const clause1 = "a = 1"
	const clause2 = "b = 2"
	const clause3 = "c = 3"

	const arg1 = "1"
	const arg2 = "2"
	const arg3 = "3"

	// ensure single where clause is generated correctly
	f.addWhere(clause1)
	r, rArgs := f.generateWhereClauses()
	assert.Equal(fmt.Sprintf("(%s)", clause1), r)
	assert.Len(rArgs, 0)

	// ensure multiple where clauses are surrounded with parenthesis and
	// ANDed together
	f.addWhere(clause2, arg1, arg2)
	r, rArgs = f.generateWhereClauses()
	assert.Equal(fmt.Sprintf("((%s) AND (%s))", clause1, clause2), r)
	assert.Len(rArgs, 2)

	// ensure empty subfilter is not added to generated where clause
	sf := &filterBuilder{}
	f.and(sf)

	r, rArgs = f.generateWhereClauses()
	assert.Equal(fmt.Sprintf("((%s) AND (%s))", clause1, clause2), r)
	assert.Len(rArgs, 2)

	// ensure sub-filter is generated correctly
	sf.addWhere(clause3, arg3)
	r, rArgs = f.generateWhereClauses()
	assert.Equal(fmt.Sprintf("((%s) AND (%s)) AND ((%s))", clause1, clause2, clause3), r)
	assert.Len(rArgs, 3)

	// ensure OR sub-filter is generated correctly
	f = &filterBuilder{}
	f.addWhere(clause1)
	f.addWhere(clause2, arg1, arg2)
	f.or(sf)

	r, rArgs = f.generateWhereClauses()
	assert.Equal(fmt.Sprintf("((%s) AND (%s)) OR ((%s))", clause1, clause2, clause3), r)
	assert.Len(rArgs, 3)

	// ensure NOT sub-filter is generated correctly
	f = &filterBuilder{}
	f.addWhere(clause1)
	f.addWhere(clause2, arg1, arg2)
	f.not(sf)

	r, rArgs = f.generateWhereClauses()
	assert.Equal(fmt.Sprintf("((%s) AND (%s)) AND NOT ((%s))", clause1, clause2, clause3), r)
	assert.Len(rArgs, 3)

	// ensure empty filter with ANDed sub-filter does not include AND
	f = &filterBuilder{}
	f.and(sf)

	r, rArgs = f.generateWhereClauses()
	assert.Equal(fmt.Sprintf("((%s))", clause3), r)
	assert.Len(rArgs, 1)

	// ensure empty filter with ORed sub-filter does not include OR
	f = &filterBuilder{}
	f.or(sf)

	r, rArgs = f.generateWhereClauses()
	assert.Equal(fmt.Sprintf("((%s))", clause3), r)
	assert.Len(rArgs, 1)

	// ensure empty filter with NOTed sub-filter does not include AND
	f = &filterBuilder{}
	f.not(sf)

	r, rArgs = f.generateWhereClauses()
	assert.Equal(fmt.Sprintf("NOT ((%s))", clause3), r)
	assert.Len(rArgs, 1)

	// (clause1) AND ((clause2) OR (clause3))
	f = &filterBuilder{}
	f.addWhere(clause1)
	sf2 := &filterBuilder{}
	sf2.addWhere(clause2, arg1, arg2)
	f.and(sf2)
	sf2.or(sf)
	r, rArgs = f.generateWhereClauses()
	assert.Equal(fmt.Sprintf("(%s) AND ((%s) OR ((%s)))", clause1, clause2, clause3), r)
	assert.Len(rArgs, 3)
}

func TestGenerateHavingClauses(t *testing.T) {
	assert := assert.New(t)

	f := &filterBuilder{}

	const clause1 = "a = 1"
	const clause2 = "b = 2"
	const clause3 = "c = 3"

	const arg1 = "1"
	const arg2 = "2"
	const arg3 = "3"

	// ensure single Having clause is generated correctly
	f.addHaving(clause1)
	r, rArgs := f.generateHavingClauses()
	assert.Equal(fmt.Sprintf("(%s)", clause1), r)
	assert.Len(rArgs, 0)

	// ensure multiple Having clauses are surrounded with parenthesis and
	// ANDed together
	f.addHaving(clause2, arg1, arg2)
	r, rArgs = f.generateHavingClauses()
	assert.Equal("(("+clause1+") AND ("+clause2+"))", r)
	assert.Len(rArgs, 2)

	// ensure empty subfilter is not added to generated Having clause
	sf := &filterBuilder{}
	f.and(sf)

	r, rArgs = f.generateHavingClauses()
	assert.Equal("(("+clause1+") AND ("+clause2+"))", r)
	assert.Len(rArgs, 2)

	// ensure sub-filter is generated correctly
	sf.addHaving(clause3, arg3)
	r, rArgs = f.generateHavingClauses()
	assert.Equal("(("+clause1+") AND ("+clause2+")) AND (("+clause3+"))", r)
	assert.Len(rArgs, 3)

	// ensure OR sub-filter is generated correctly
	f = &filterBuilder{}
	f.addHaving(clause1)
	f.addHaving(clause2, arg1, arg2)
	f.or(sf)

	r, rArgs = f.generateHavingClauses()
	assert.Equal("(("+clause1+") AND ("+clause2+")) OR (("+clause3+"))", r)
	assert.Len(rArgs, 3)

	// ensure NOT sub-filter is generated correctly
	f = &filterBuilder{}
	f.addHaving(clause1)
	f.addHaving(clause2, arg1, arg2)
	f.not(sf)

	r, rArgs = f.generateHavingClauses()
	assert.Equal("(("+clause1+") AND ("+clause2+")) AND NOT (("+clause3+"))", r)
	assert.Len(rArgs, 3)
}

func TestGetAllJoins(t *testing.T) {
	assert := assert.New(t)
	f := &filterBuilder{}

	const (
		table1Name = "table1Name"
		table2Name = "table2Name"

		as1Name = "as1"
		as2Name = "as2"

		onClause = "onClause1"
	)

	f.addLeftJoin(table1Name, as1Name, onClause)

	// ensure join is returned
	joins := f.getAllJoins()
	assert.Len(joins, 1)
	assert.Equal(fmt.Sprintf("LEFT JOIN %s AS %s ON %s", table1Name, as1Name, onClause), joins[0].toSQL())

	// ensure joins in sub-filter are returned
	subFilter := &filterBuilder{}
	f.and(subFilter)
	subFilter.addLeftJoin(table2Name, as2Name, onClause)

	joins = f.getAllJoins()
	assert.Len(joins, 2)
	assert.Equal(fmt.Sprintf("LEFT JOIN %s AS %s ON %s", table2Name, as2Name, onClause), joins[1].toSQL())

	// ensure redundant joins are not returned
	subFilter.addLeftJoin(as1Name, "", onClause)
	joins = f.getAllJoins()
	assert.Len(joins, 2)
}

func TestGetError(t *testing.T) {
	assert := assert.New(t)
	f := &filterBuilder{}
	subFilter := &filterBuilder{}

	f.and(subFilter)

	expectedErr := errors.New("test error")
	expectedErr2 := errors.New("test error2")
	f.err = expectedErr
	subFilter.err = expectedErr2

	// ensure getError returns the top-level error state
	assert.Equal(expectedErr, f.getError())

	// ensure getError returns sub-filter error state if top-level error
	// is nil
	f.err = nil
	assert.Equal(expectedErr2, f.getError())

	// ensure getError returns nil if all error states are nil
	subFilter.err = nil
	assert.Nil(f.getError())
}

func TestStringCriterionHandlerIncludes(t *testing.T) {
	assert := assert.New(t)

	const column = "column"
	const value1 = "two words"
	const quotedValue = `"two words"`

	f := &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierIncludes,
		Value:    value1,
	}, column))

	assert.Len(f.whereClauses, 1)
	assert.Equal(fmt.Sprintf("(%[1]s LIKE ? OR %[1]s LIKE ?)", column), f.whereClauses[0].sql)
	assert.Len(f.whereClauses[0].args, 2)
	assert.Equal("%two%", f.whereClauses[0].args[0])
	assert.Equal("%words%", f.whereClauses[0].args[1])

	f = &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierIncludes,
		Value:    quotedValue,
	}, column))

	assert.Len(f.whereClauses, 1)
	assert.Equal(fmt.Sprintf("(%[1]s LIKE ?)", column), f.whereClauses[0].sql)
	assert.Len(f.whereClauses[0].args, 1)
	assert.Equal("%two words%", f.whereClauses[0].args[0])
}

func TestStringCriterionHandlerExcludes(t *testing.T) {
	assert := assert.New(t)

	const column = "column"
	const value1 = "two words"
	const quotedValue = `"two words"`

	f := &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierExcludes,
		Value:    value1,
	}, column))

	assert.Len(f.whereClauses, 1)
	assert.Equal(fmt.Sprintf("(%[1]s NOT LIKE ? AND %[1]s NOT LIKE ?)", column), f.whereClauses[0].sql)
	assert.Len(f.whereClauses[0].args, 2)
	assert.Equal("%two%", f.whereClauses[0].args[0])
	assert.Equal("%words%", f.whereClauses[0].args[1])

	f = &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierExcludes,
		Value:    quotedValue,
	}, column))

	assert.Len(f.whereClauses, 1)
	assert.Equal(fmt.Sprintf("(%[1]s NOT LIKE ?)", column), f.whereClauses[0].sql)
	assert.Len(f.whereClauses[0].args, 1)
	assert.Equal("%two words%", f.whereClauses[0].args[0])
}

func TestStringCriterionHandlerEquals(t *testing.T) {
	assert := assert.New(t)

	const column = "column"
	const value1 = "two words"

	f := &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierEquals,
		Value:    value1,
	}, column))

	assert.Len(f.whereClauses, 1)
	assert.Equal(fmt.Sprintf("%[1]s LIKE ?", column), f.whereClauses[0].sql)
	assert.Len(f.whereClauses[0].args, 1)
	assert.Equal(value1, f.whereClauses[0].args[0])
}

func TestStringCriterionHandlerNotEquals(t *testing.T) {
	assert := assert.New(t)

	const column = "column"
	const value1 = "two words"

	f := &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierNotEquals,
		Value:    value1,
	}, column))

	assert.Len(f.whereClauses, 1)
	assert.Equal(fmt.Sprintf("%[1]s NOT LIKE ?", column), f.whereClauses[0].sql)
	assert.Len(f.whereClauses[0].args, 1)
	assert.Equal(value1, f.whereClauses[0].args[0])
}

func TestStringCriterionHandlerMatchesRegex(t *testing.T) {
	assert := assert.New(t)

	const column = "column"
	const validValue = "two words"
	const invalidValue = "*two words"

	f := &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierMatchesRegex,
		Value:    validValue,
	}, column))

	assert.Len(f.whereClauses, 1)
	assert.Equal(fmt.Sprintf("(%s IS NOT NULL AND %[1]s regexp ?)", column), f.whereClauses[0].sql)
	assert.Len(f.whereClauses[0].args, 1)
	assert.Equal(validValue, f.whereClauses[0].args[0])

	// ensure invalid regex sets error state
	f = &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierMatchesRegex,
		Value:    invalidValue,
	}, column))

	assert.NotNil(f.getError())
}

func TestStringCriterionHandlerNotMatchesRegex(t *testing.T) {
	assert := assert.New(t)

	const column = "column"
	const validValue = "two words"
	const invalidValue = "*two words"

	f := &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierNotMatchesRegex,
		Value:    validValue,
	}, column))

	assert.Len(f.whereClauses, 1)
	assert.Equal(fmt.Sprintf("(%s IS NULL OR %[1]s NOT regexp ?)", column), f.whereClauses[0].sql)
	assert.Len(f.whereClauses[0].args, 1)
	assert.Equal(validValue, f.whereClauses[0].args[0])

	// ensure invalid regex sets error state
	f = &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierNotMatchesRegex,
		Value:    invalidValue,
	}, column))

	assert.NotNil(f.getError())
}

func TestStringCriterionHandlerIsNull(t *testing.T) {
	assert := assert.New(t)

	const column = "column"

	f := &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierIsNull,
	}, column))

	assert.Len(f.whereClauses, 1)
	assert.Equal(fmt.Sprintf("(%[1]s IS NULL OR TRIM(%[1]s) = '')", column), f.whereClauses[0].sql)
	assert.Len(f.whereClauses[0].args, 0)
}

func TestStringCriterionHandlerNotNull(t *testing.T) {
	assert := assert.New(t)

	const column = "column"

	f := &filterBuilder{}
	f.handleCriterion(testCtx, stringCriterionHandler(&models.StringCriterionInput{
		Modifier: models.CriterionModifierNotNull,
	}, column))

	assert.Len(f.whereClauses, 1)
	assert.Equal(fmt.Sprintf("(%[1]s IS NOT NULL AND TRIM(%[1]s) != '')", column), f.whereClauses[0].sql)
	assert.Len(f.whereClauses[0].args, 0)
}
