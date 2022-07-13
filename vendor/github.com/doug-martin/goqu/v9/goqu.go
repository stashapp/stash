/*
goqu an idiomatch SQL builder, and query package.

      __ _  ___   __ _ _   _
     / _` |/ _ \ / _` | | | |
    | (_| | (_) | (_| | |_| |
     \__, |\___/ \__, |\__,_|
     |___/          |_|


Please see https://github.com/doug-martin/goqu for an introduction to goqu.
*/
package goqu

import (
	"time"

	"github.com/doug-martin/goqu/v9/internal/util"
	"github.com/doug-martin/goqu/v9/sqlgen"
)

type DialectWrapper struct {
	dialect string
}

// Creates a new DialectWrapper to create goqu.Datasets or goqu.Databases with the specified dialect.
func Dialect(dialect string) DialectWrapper {
	return DialectWrapper{dialect: dialect}
}

// Create a new dataset for creating SELECT sql statements
func (dw DialectWrapper) From(table ...interface{}) *SelectDataset {
	return From(table...).WithDialect(dw.dialect)
}

// Create a new dataset for creating SELECT sql statements
func (dw DialectWrapper) Select(cols ...interface{}) *SelectDataset {
	return newDataset(dw.dialect, nil).Select(cols...)
}

// Create a new dataset for creating UPDATE sql statements
func (dw DialectWrapper) Update(table interface{}) *UpdateDataset {
	return Update(table).WithDialect(dw.dialect)
}

// Create a new dataset for creating INSERT sql statements
func (dw DialectWrapper) Insert(table interface{}) *InsertDataset {
	return Insert(table).WithDialect(dw.dialect)
}

// Create a new dataset for creating DELETE sql statements
func (dw DialectWrapper) Delete(table interface{}) *DeleteDataset {
	return Delete(table).WithDialect(dw.dialect)
}

// Create a new dataset for creating TRUNCATE sql statements
func (dw DialectWrapper) Truncate(table ...interface{}) *TruncateDataset {
	return Truncate(table...).WithDialect(dw.dialect)
}

func (dw DialectWrapper) DB(db SQLDatabase) *Database {
	return newDatabase(dw.dialect, db)
}

func New(dialect string, db SQLDatabase) *Database {
	return newDatabase(dialect, db)
}

// Set the behavior when encountering struct fields that do not have a db tag.
// By default this is false; if set to true any field without a db tag will not
// be targeted by Select or Scan operations.
func SetIgnoreUntaggedFields(ignore bool) {
	util.SetIgnoreUntaggedFields(ignore)
}

// Set the column rename function. This is used for struct fields that do not have a db tag to specify the column name
// By default all struct fields that do not have a db tag will be converted lowercase
func SetColumnRenameFunction(renameFunc func(string) string) {
	util.SetColumnRenameFunction(renameFunc)
}

// Set the location to use when interpolating time.Time instances. See https://golang.org/pkg/time/#LoadLocation
// NOTE: This has no effect when using prepared statements.
func SetTimeLocation(loc *time.Location) {
	sqlgen.SetTimeLocation(loc)
}
