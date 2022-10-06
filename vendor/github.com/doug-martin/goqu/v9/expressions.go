package goqu

import (
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	Expression = exp.Expression
	Ex         = exp.Ex
	ExOr       = exp.ExOr
	Op         = exp.Op
	Record     = exp.Record
	Vals       = exp.Vals
	// Options to use when generating a TRUNCATE statement
	TruncateOptions = exp.TruncateOptions
)

// emptyWindow is an empty WINDOW clause without name
var emptyWindow = exp.NewWindowExpression(nil, nil, nil, nil)

const (
	Wait       = exp.Wait
	NoWait     = exp.NoWait
	SkipLocked = exp.SkipLocked
)

// Creates a new Casted expression
//  Cast(I("a"), "NUMERIC") -> CAST("a" AS NUMERIC)
func Cast(e exp.Expression, t string) exp.CastExpression {
	return exp.NewCastExpression(e, t)
}

// Creates a conflict struct to be passed to InsertConflict to ignore constraint errors
//  InsertConflict(DoNothing(),...) -> INSERT INTO ... ON CONFLICT DO NOTHING
func DoNothing() exp.ConflictExpression {
	return exp.NewDoNothingConflictExpression()
}

// Creates a ConflictUpdate struct to be passed to InsertConflict
// Represents a ON CONFLICT DO UPDATE portion of an INSERT statement (ON DUPLICATE KEY UPDATE for mysql)
//
//  InsertConflict(DoUpdate("target_column", update),...) ->
//  	INSERT INTO ... ON CONFLICT DO UPDATE SET a=b
//  InsertConflict(DoUpdate("target_column", update).Where(Ex{"a": 1},...) ->
//  	INSERT INTO ... ON CONFLICT DO UPDATE SET a=b WHERE a=1
func DoUpdate(target string, update interface{}) exp.ConflictUpdateExpression {
	return exp.NewDoUpdateConflictExpression(target, update)
}

// A list of expressions that should be ORed together
//    Or(I("a").Eq(10), I("b").Eq(11)) //(("a" = 10) OR ("b" = 11))
func Or(expressions ...exp.Expression) exp.ExpressionList {
	return exp.NewExpressionList(exp.OrType, expressions...)
}

// A list of expressions that should be ANDed together
//    And(I("a").Eq(10), I("b").Eq(11)) //(("a" = 10) AND ("b" = 11))
func And(expressions ...exp.Expression) exp.ExpressionList {
	return exp.NewExpressionList(exp.AndType, expressions...)
}

// Creates a new SQLFunctionExpression with the given name and arguments
func Func(name string, args ...interface{}) exp.SQLFunctionExpression {
	return exp.NewSQLFunctionExpression(name, args...)
}

// used internally to normalize the column name if passed in as a string it should be turned into an identifier
func newIdentifierFunc(name string, col interface{}) exp.SQLFunctionExpression {
	if s, ok := col.(string); ok {
		col = I(s)
	}
	return Func(name, col)
}

// Creates a new DISTINCT sql function
//   DISTINCT("a") -> DISTINCT("a")
//   DISTINCT(I("a")) -> DISTINCT("a")
func DISTINCT(col interface{}) exp.SQLFunctionExpression { return newIdentifierFunc("DISTINCT", col) }

// Creates a new COUNT sql function
//   COUNT("a") -> COUNT("a")
//   COUNT("*") -> COUNT("*")
//   COUNT(I("a")) -> COUNT("a")
func COUNT(col interface{}) exp.SQLFunctionExpression { return newIdentifierFunc("COUNT", col) }

// Creates a new MIN sql function
//   MIN("a") -> MIN("a")
//   MIN(I("a")) -> MIN("a")
func MIN(col interface{}) exp.SQLFunctionExpression { return newIdentifierFunc("MIN", col) }

// Creates a new MAX sql function
//   MAX("a") -> MAX("a")
//   MAX(I("a")) -> MAX("a")
func MAX(col interface{}) exp.SQLFunctionExpression { return newIdentifierFunc("MAX", col) }

// Creates a new AVG sql function
//   AVG("a") -> AVG("a")
//   AVG(I("a")) -> AVG("a")
func AVG(col interface{}) exp.SQLFunctionExpression { return newIdentifierFunc("AVG", col) }

// Creates a new FIRST sql function
//   FIRST("a") -> FIRST("a")
//   FIRST(I("a")) -> FIRST("a")
func FIRST(col interface{}) exp.SQLFunctionExpression { return newIdentifierFunc("FIRST", col) }

// Creates a new LAST sql function
//   LAST("a") -> LAST("a")
//   LAST(I("a")) -> LAST("a")
func LAST(col interface{}) exp.SQLFunctionExpression { return newIdentifierFunc("LAST", col) }

// Creates a new SUM sql function
//   SUM("a") -> SUM("a")
//   SUM(I("a")) -> SUM("a")
func SUM(col interface{}) exp.SQLFunctionExpression { return newIdentifierFunc("SUM", col) }

// Creates a new COALESCE sql function
//   COALESCE(I("a"), "a") -> COALESCE("a", 'a')
//   COALESCE(I("a"), I("b"), nil) -> COALESCE("a", "b", NULL)
func COALESCE(vals ...interface{}) exp.SQLFunctionExpression {
	return Func("COALESCE", vals...)
}

//nolint:stylecheck,golint // sql function name
func ROW_NUMBER() exp.SQLFunctionExpression {
	return Func("ROW_NUMBER")
}

func RANK() exp.SQLFunctionExpression {
	return Func("RANK")
}

//nolint:stylecheck,golint // sql function name
func DENSE_RANK() exp.SQLFunctionExpression {
	return Func("DENSE_RANK")
}

//nolint:stylecheck,golint // sql function name
func PERCENT_RANK() exp.SQLFunctionExpression {
	return Func("PERCENT_RANK")
}

//nolint:stylecheck,golint //sql function name
func CUME_DIST() exp.SQLFunctionExpression {
	return Func("CUME_DIST")
}

func NTILE(n int) exp.SQLFunctionExpression {
	return Func("NTILE", n)
}

//nolint:stylecheck,golint //sql function name
func FIRST_VALUE(val interface{}) exp.SQLFunctionExpression {
	return newIdentifierFunc("FIRST_VALUE", val)
}

//nolint:stylecheck,golint //sql function name
func LAST_VALUE(val interface{}) exp.SQLFunctionExpression {
	return newIdentifierFunc("LAST_VALUE", val)
}

//nolint:stylecheck,golint //sql function name
func NTH_VALUE(val interface{}, nth int) exp.SQLFunctionExpression {
	if s, ok := val.(string); ok {
		val = I(s)
	}
	return Func("NTH_VALUE", val, nth)
}

// Creates a new Identifier, the generated sql will use adapter specific quoting or '"' by default, this ensures case
// sensitivity and in certain databases allows for special characters, (e.g. "curr-table", "my table").
//
// The identifier will be split by '.'
//
// Table and Column example
//    I("table.column") -> "table"."column" //A Column and table
// Schema table and column
//    I("schema.table.column") -> "schema"."table"."column"
// Table with star
//    I("table.*") -> "table".*
func I(ident string) exp.IdentifierExpression {
	return exp.ParseIdentifier(ident)
}

// Creates a new Column Identifier, the generated sql will use adapter specific quoting or '"' by default, this ensures case
// sensitivity and in certain databases allows for special characters, (e.g. "curr-table", "my table").
// An Identifier can represent a one or a combination of schema, table, and/or column.
//    C("column") -> "column" //A Column
//    C("column").Table("table") -> "table"."column" //A Column and table
//    C("column").Table("table").Schema("schema") //Schema table and column
//    C("*") //Also handles the * operator
func C(col string) exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", "", col)
}

// Creates a new Schema Identifier, the generated sql will use adapter specific quoting or '"' by default, this ensures case
// sensitivity and in certain databases allows for special characters, (e.g. "curr-schema", "my schema").
//    S("schema") -> "schema" //A Schema
//    S("schema").Table("table") -> "schema"."table" //A Schema and table
//    S("schema").Table("table").Col("col") //Schema table and column
//    S("schema").Table("table").Col("*") //Schema table and all columns
func S(schema string) exp.IdentifierExpression {
	return exp.NewIdentifierExpression(schema, "", "")
}

// Creates a new Table Identifier, the generated sql will use adapter specific quoting or '"' by default, this ensures case
// sensitivity and in certain databases allows for special characters, (e.g. "curr-table", "my table").
//    T("table") -> "table" //A Column
//    T("table").Col("col") -> "table"."column" //A Column and table
//    T("table").Schema("schema").Col("col) -> "schema"."table"."column"  //Schema table and column
//    T("table").Schema("schema").Col("*") -> "schema"."table".*  //Also handles the * operator
func T(table string) exp.IdentifierExpression {
	return exp.NewIdentifierExpression("", table, "")
}

// Create a new WINDOW clause
// 	W() -> ()
// 	W().PartitionBy("a") -> (PARTITION BY "a")
// 	W().PartitionBy("a").OrderBy("b") -> (PARTITION BY "a" ORDER BY "b")
// 	W().PartitionBy("a").OrderBy("b").Inherit("w1") -> ("w1" PARTITION BY "a" ORDER BY "b")
// 	W().PartitionBy("a").OrderBy(I("b").Desc()).Inherit("w1") -> ("w1" PARTITION BY "a" ORDER BY "b" DESC)
// 	W("w") -> "w" AS ()
// 	W("w", "w1") -> "w" AS ("w1")
// 	W("w").Inherit("w1") -> "w" AS ("w1")
// 	W("w").PartitionBy("a") -> "w" AS (PARTITION BY "a")
// 	W("w", "w1").PartitionBy("a") -> "w" AS ("w1" PARTITION BY "a")
// 	W("w", "w1").PartitionBy("a").OrderBy("b") -> "w" AS ("w1" PARTITION BY "a" ORDER BY "b")
func W(ws ...string) exp.WindowExpression {
	switch len(ws) {
	case 0:
		return emptyWindow
	case 1:
		return exp.NewWindowExpression(I(ws[0]), nil, nil, nil)
	default:
		return exp.NewWindowExpression(I(ws[0]), I(ws[1]), nil, nil)
	}
}

// Creates a new ON clause to be used within a join
//    ds.Join(goqu.T("my_table"), goqu.On(
//       goqu.I("my_table.fkey").Eq(goqu.I("other_table.id")),
//    ))
func On(expressions ...exp.Expression) exp.JoinCondition {
	return exp.NewJoinOnCondition(expressions...)
}

// Creates a new USING clause to be used within a join
//    ds.Join(goqu.T("my_table"), goqu.Using("fkey"))
func Using(columns ...interface{}) exp.JoinCondition {
	return exp.NewJoinUsingCondition(columns...)
}

// Creates a new SQL literal with the provided arguments.
//   L("a = 1") -> a = 1
// You can also you placeholders. All placeholders within a Literal are represented by '?'
//   L("a = ?", "b") -> a = 'b'
// Literals can also contain placeholders for other expressions
//   L("(? AND ?) OR (?)", I("a").Eq(1), I("b").Eq("b"), I("c").In([]string{"a", "b", "c"}))
func L(sql string, args ...interface{}) exp.LiteralExpression {
	return Literal(sql, args...)
}

// Alias for goqu.L
func Literal(sql string, args ...interface{}) exp.LiteralExpression {
	return exp.NewLiteralExpression(sql, args...)
}

// Create a new SQL value ( alias for goqu.L("?", val) ). The prrimary use case for this would be in selects.
// See examples.
func V(val interface{}) exp.LiteralExpression {
	return exp.NewLiteralExpression("?", val)
}

// Creates a new Range to be used with a Between expression
//    exp.C("col").Between(exp.Range(1, 10))
func Range(start, end interface{}) exp.RangeVal {
	return exp.NewRangeVal(start, end)
}

// Creates a literal *
func Star() exp.LiteralExpression { return exp.Star() }

// Returns a literal for DEFAULT sql keyword
func Default() exp.LiteralExpression {
	return exp.Default()
}

func Lateral(table exp.AppendableExpression) exp.LateralExpression {
	return exp.NewLateralExpression(table)
}

// Create a new ANY comparison
func Any(val interface{}) exp.SQLFunctionExpression {
	return Func("ANY ", val)
}

// Create a new ALL comparison
func All(val interface{}) exp.SQLFunctionExpression {
	return Func("ALL ", val)
}

func Case() exp.CaseExpression {
	return exp.NewCaseExpression()
}
