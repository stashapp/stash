package exp

import (
	"strings"
)

type (
	identifier struct {
		schema string
		table  string
		col    interface{}
	}
)

var (
	tableAndColumnParts                 = 2
	schemaTableAndColumnIdentifierParts = 3
)

func ParseIdentifier(ident string) IdentifierExpression {
	parts := strings.Split(ident, ".")
	switch len(parts) {
	case tableAndColumnParts:
		return NewIdentifierExpression("", parts[0], parts[1])
	case schemaTableAndColumnIdentifierParts:
		return NewIdentifierExpression(parts[0], parts[1], parts[2])
	}
	return NewIdentifierExpression("", "", ident)
}

func NewIdentifierExpression(schema, table string, col interface{}) IdentifierExpression {
	return identifier{}.Schema(schema).Table(table).Col(col)
}

func (i identifier) clone() identifier {
	return identifier{schema: i.schema, table: i.table, col: i.col}
}

func (i identifier) Clone() Expression {
	return i.clone()
}

func (i identifier) IsQualified() bool {
	schema, table, col := i.schema, i.table, i.col
	switch c := col.(type) {
	case string:
		if c != "" {
			return len(table) > 0 || len(schema) > 0
		}
	default:
		if c != nil {
			return len(table) > 0 || len(schema) > 0
		}
	}
	if len(table) > 0 {
		return len(schema) > 0
	}
	return false
}

// Sets the table on the current identifier
//  I("col").Table("table") -> "table"."col" //postgres
//  I("col").Table("table") -> `table`.`col` //mysql
//  I("col").Table("table") -> `table`.`col` //sqlite3
func (i identifier) Table(table string) IdentifierExpression {
	i.table = table
	return i
}

func (i identifier) GetTable() string {
	return i.table
}

// Sets the table on the current identifier
//  I("table").Schema("schema") -> "schema"."table" //postgres
//  I("col").Schema("table") -> `schema`.`table` //mysql
//  I("col").Schema("table") -> `schema`.`table` //sqlite3
func (i identifier) Schema(schema string) IdentifierExpression {
	i.schema = schema
	return i
}

func (i identifier) GetSchema() string {
	return i.schema
}

// Sets the table on the current identifier
//  I("table").Col("col") -> "table"."col" //postgres
//  I("table").Schema("col") -> `table`.`col` //mysql
//  I("table").Schema("col") -> `table`.`col` //sqlite3
func (i identifier) Col(col interface{}) IdentifierExpression {
	if col == "*" {
		i.col = Star()
	} else {
		i.col = col
	}
	return i
}

func (i identifier) Expression() Expression { return i }

// Qualifies the epression with a * literal (e.g. "table".*)
func (i identifier) All() IdentifierExpression { return i.Col("*") }

func (i identifier) IsEmpty() bool {
	isEmpty := i.schema == "" && i.table == ""
	if isEmpty {
		switch t := i.col.(type) {
		case nil:
			return true
		case string:
			return t == ""
		default:
			return false
		}
	}
	return isEmpty
}

// Gets the column identifier
func (i identifier) GetCol() interface{} { return i.col }

// Used within updates to set a column value
func (i identifier) Set(val interface{}) UpdateExpression { return set(i, val) }

// Alias an identifier (e.g "my_col" AS "other_col")
func (i identifier) As(val interface{}) AliasedExpression {
	if v, ok := val.(string); ok {
		ident := ParseIdentifier(v)
		if i.col != nil && i.col != "" {
			return NewAliasExpression(i, ident)
		}
		aliasCol := ident.GetCol()
		if i.table != "" {
			return NewAliasExpression(i, NewIdentifierExpression("", aliasCol.(string), nil))
		} else if i.schema != "" {
			return NewAliasExpression(i, NewIdentifierExpression(aliasCol.(string), "", nil))
		}
	}
	return NewAliasExpression(i, val)
}

// Returns a BooleanExpression for equality (e.g "my_col" = 1)
func (i identifier) Eq(val interface{}) BooleanExpression { return eq(i, val) }

// Returns a BooleanExpression for in equality (e.g "my_col" != 1)
func (i identifier) Neq(val interface{}) BooleanExpression { return neq(i, val) }

// Returns a BooleanExpression for checking that a identifier is greater than another value (e.g "my_col" > 1)
func (i identifier) Gt(val interface{}) BooleanExpression { return gt(i, val) }

// Returns a BooleanExpression for checking that a identifier is greater than or equal to another value
// (e.g "my_col" >= 1)
func (i identifier) Gte(val interface{}) BooleanExpression { return gte(i, val) }

// Returns a BooleanExpression for checking that a identifier is less than another value (e.g "my_col" < 1)
func (i identifier) Lt(val interface{}) BooleanExpression { return lt(i, val) }

// Returns a BooleanExpression for checking that a identifier is less than or equal to another value
// (e.g "my_col" <= 1)
func (i identifier) Lte(val interface{}) BooleanExpression { return lte(i, val) }

// Returns a BooleanExpression for bit inversion (e.g ~ "my_col")
func (i identifier) BitwiseInversion() BitwiseExpression { return bitwiseInversion(i) }

// Returns a BooleanExpression for bit OR (e.g "my_col" | 1)
func (i identifier) BitwiseOr(val interface{}) BitwiseExpression { return bitwiseOr(i, val) }

// Returns a BooleanExpression for bit AND (e.g "my_col" & 1)
func (i identifier) BitwiseAnd(val interface{}) BitwiseExpression { return bitwiseAnd(i, val) }

// Returns a BooleanExpression for bit XOR (e.g "my_col" ^ 1)
func (i identifier) BitwiseXor(val interface{}) BitwiseExpression { return bitwiseXor(i, val) }

// Returns a BooleanExpression for bit LEFT shift (e.g "my_col" << 1)
func (i identifier) BitwiseLeftShift(val interface{}) BitwiseExpression {
	return bitwiseLeftShift(i, val)
}

// Returns a BooleanExpression for bit RIGHT shift (e.g "my_col" >> 1)
func (i identifier) BitwiseRightShift(val interface{}) BitwiseExpression {
	return bitwiseRightShift(i, val)
}

// Returns a BooleanExpression for checking that a identifier is in a list of values or  (e.g "my_col" > 1)
func (i identifier) In(vals ...interface{}) BooleanExpression         { return in(i, vals...) }
func (i identifier) NotIn(vals ...interface{}) BooleanExpression      { return notIn(i, vals...) }
func (i identifier) Like(val interface{}) BooleanExpression           { return like(i, val) }
func (i identifier) NotLike(val interface{}) BooleanExpression        { return notLike(i, val) }
func (i identifier) ILike(val interface{}) BooleanExpression          { return iLike(i, val) }
func (i identifier) NotILike(val interface{}) BooleanExpression       { return notILike(i, val) }
func (i identifier) RegexpLike(val interface{}) BooleanExpression     { return regexpLike(i, val) }
func (i identifier) RegexpNotLike(val interface{}) BooleanExpression  { return regexpNotLike(i, val) }
func (i identifier) RegexpILike(val interface{}) BooleanExpression    { return regexpILike(i, val) }
func (i identifier) RegexpNotILike(val interface{}) BooleanExpression { return regexpNotILike(i, val) }
func (i identifier) Is(val interface{}) BooleanExpression             { return is(i, val) }
func (i identifier) IsNot(val interface{}) BooleanExpression          { return isNot(i, val) }
func (i identifier) IsNull() BooleanExpression                        { return is(i, nil) }
func (i identifier) IsNotNull() BooleanExpression                     { return isNot(i, nil) }
func (i identifier) IsTrue() BooleanExpression                        { return is(i, true) }
func (i identifier) IsNotTrue() BooleanExpression                     { return isNot(i, true) }
func (i identifier) IsFalse() BooleanExpression                       { return is(i, false) }
func (i identifier) IsNotFalse() BooleanExpression                    { return isNot(i, false) }
func (i identifier) Asc() OrderedExpression                           { return asc(i) }
func (i identifier) Desc() OrderedExpression                          { return desc(i) }
func (i identifier) Distinct() SQLFunctionExpression                  { return NewSQLFunctionExpression("DISTINCT", i) }
func (i identifier) Cast(t string) CastExpression                     { return NewCastExpression(i, t) }

// Returns a RangeExpression for checking that a identifier is between two values (e.g "my_col" BETWEEN 1 AND 10)
func (i identifier) Between(val RangeVal) RangeExpression { return between(i, val) }

// Returns a RangeExpression for checking that a identifier is between two values (e.g "my_col" BETWEEN 1 AND 10)
func (i identifier) NotBetween(val RangeVal) RangeExpression { return notBetween(i, val) }
