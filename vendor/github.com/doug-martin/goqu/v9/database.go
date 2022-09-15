package goqu

import (
	"context"
	"database/sql"
	"sync"

	"github.com/doug-martin/goqu/v9/exec"
)

type (
	Logger interface {
		Printf(format string, v ...interface{})
	}
	// Interface for sql.DB, an interface is used so you can use with other
	// libraries such as sqlx instead of the native sql.DB
	SQLDatabase interface {
		Begin() (*sql.Tx, error)
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}
	// This struct is the wrapper for a Db. The struct delegates most calls to either an Exec instance or to the Db
	// passed into the constructor.
	Database struct {
		logger  Logger
		dialect string
		// nolint: stylecheck // keep for backwards compatibility
		Db     SQLDatabase
		qf     exec.QueryFactory
		qfOnce sync.Once
	}
)

// This is the common entry point into goqu.
//
// dialect: This is the adapter dialect, you should see your database adapter for the string to use. Built in adapters
// can be found at https://github.com/doug-martin/goqu/tree/master/adapters
//
// db: A sql.Db to use for querying the database
//      import (
//          "database/sql"
//          "fmt"
//          "github.com/doug-martin/goqu/v9"
//          _ "github.com/doug-martin/goqu/v9/dialect/postgres"
//          _ "github.com/lib/pq"
//      )
//
//      func main() {
//          sqlDb, err := sql.Open("postgres", "user=postgres dbname=goqupostgres sslmode=disable ")
//          if err != nil {
//              panic(err.Error())
//          }
//          db := goqu.New("postgres", sqlDb)
//      }
// The most commonly used Database method is From, which creates a new Dataset that uses the correct adapter and
// supports queries.
//          var ids []uint32
//          if err := db.From("items").Where(goqu.I("id").Gt(10)).Pluck("id", &ids); err != nil {
//              panic(err.Error())
//          }
//          fmt.Printf("%+v", ids)
func newDatabase(dialect string, db SQLDatabase) *Database {
	return &Database{
		logger:  nil,
		dialect: dialect,
		Db:      db,
		qf:      nil,
		qfOnce:  sync.Once{},
	}
}

// returns this databases dialect
func (d *Database) Dialect() string {
	return d.dialect
}

// Starts a new Transaction.
func (d *Database) Begin() (*TxDatabase, error) {
	sqlTx, err := d.Db.Begin()
	if err != nil {
		return nil, err
	}
	tx := NewTx(d.dialect, sqlTx)
	tx.Logger(d.logger)
	return tx, nil
}

// Starts a new Transaction. See sql.DB#BeginTx for option description
func (d *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*TxDatabase, error) {
	sqlTx, err := d.Db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	tx := NewTx(d.dialect, sqlTx)
	tx.Logger(d.logger)
	return tx, nil
}

// WithTx starts a new transaction and executes it in Wrap method
func (d *Database) WithTx(fn func(*TxDatabase) error) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	return tx.Wrap(func() error { return fn(tx) })
}

// Creates a new Dataset that uses the correct adapter and supports queries.
//          var ids []uint32
//          if err := db.From("items").Where(goqu.I("id").Gt(10)).Pluck("id", &ids); err != nil {
//              panic(err.Error())
//          }
//          fmt.Printf("%+v", ids)
//
// from...: Sources for you dataset, could be table names (strings), a goqu.Literal or another goqu.Dataset
func (d *Database) From(from ...interface{}) *SelectDataset {
	return newDataset(d.dialect, d.queryFactory()).From(from...)
}

func (d *Database) Select(cols ...interface{}) *SelectDataset {
	return newDataset(d.dialect, d.queryFactory()).Select(cols...)
}

func (d *Database) Update(table interface{}) *UpdateDataset {
	return newUpdateDataset(d.dialect, d.queryFactory()).Table(table)
}

func (d *Database) Insert(table interface{}) *InsertDataset {
	return newInsertDataset(d.dialect, d.queryFactory()).Into(table)
}

func (d *Database) Delete(table interface{}) *DeleteDataset {
	return newDeleteDataset(d.dialect, d.queryFactory()).From(table)
}

func (d *Database) Truncate(table ...interface{}) *TruncateDataset {
	return newTruncateDataset(d.dialect, d.queryFactory()).Table(table...)
}

// Sets the logger for to use when logging queries
func (d *Database) Logger(logger Logger) {
	d.logger = logger
}

// Logs a given operation with the specified sql and arguments
func (d *Database) Trace(op, sqlString string, args ...interface{}) {
	if d.logger != nil {
		if sqlString != "" {
			if len(args) != 0 {
				d.logger.Printf("[goqu] %s [query:=`%s` args:=%+v]", op, sqlString, args)
			} else {
				d.logger.Printf("[goqu] %s [query:=`%s`]", op, sqlString)
			}
		} else {
			d.logger.Printf("[goqu] %s", op)
		}
	}
}

// Uses the db to Execute the query with arguments and return the sql.Result
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.ExecContext(context.Background(), query, args...)
}

// Uses the db to Execute the query with arguments and return the sql.Result
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	d.Trace("EXEC", query, args...)
	return d.Db.ExecContext(ctx, query, args...)
}

// Can be used to prepare a query.
//
// You can use this in tandem with a dataset by doing the following.
//    sql, args, err := db.From("items").Where(goqu.I("id").Gt(10)).ToSQL(true)
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    stmt, err := db.Prepare(sql)
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    defer stmt.Close()
//    rows, err := stmt.Query(args)
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    defer rows.Close()
//    for rows.Next(){
//              //scan your rows
//    }
//    if rows.Err() != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//
// query: The SQL statement to prepare.
func (d *Database) Prepare(query string) (*sql.Stmt, error) {
	return d.PrepareContext(context.Background(), query)
}

// Can be used to prepare a query.
//
// You can use this in tandem with a dataset by doing the following.
//    sql, args, err := db.From("items").Where(goqu.I("id").Gt(10)).ToSQL(true)
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    stmt, err := db.Prepare(sql)
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    defer stmt.Close()
//    rows, err := stmt.QueryContext(ctx, args)
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    defer rows.Close()
//    for rows.Next(){
//              //scan your rows
//    }
//    if rows.Err() != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//
// query: The SQL statement to prepare.
func (d *Database) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	d.Trace("PREPARE", query)
	return d.Db.PrepareContext(ctx, query)
}

// Used to query for multiple rows.
//
// You can use this in tandem with a dataset by doing the following.
//    sql, err := db.From("items").Where(goqu.I("id").Gt(10)).ToSQL()
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    rows, err := stmt.Query(args)
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    defer rows.Close()
//    for rows.Next(){
//              //scan your rows
//    }
//    if rows.Err() != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.QueryContext(context.Background(), query, args...)
}

// Used to query for multiple rows.
//
// You can use this in tandem with a dataset by doing the following.
//    sql, err := db.From("items").Where(goqu.I("id").Gt(10)).ToSQL()
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    rows, err := stmt.QueryContext(ctx, args)
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    defer rows.Close()
//    for rows.Next(){
//              //scan your rows
//    }
//    if rows.Err() != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	d.Trace("QUERY", query, args...)
	return d.Db.QueryContext(ctx, query, args...)
}

// Used to query for a single row.
//
// You can use this in tandem with a dataset by doing the following.
//    sql, err := db.From("items").Where(goqu.I("id").Gt(10)).Limit(1).ToSQL()
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    rows, err := stmt.QueryRow(args)
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    //scan your row
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.QueryRowContext(context.Background(), query, args...)
}

// Used to query for a single row.
//
// You can use this in tandem with a dataset by doing the following.
//    sql, err := db.From("items").Where(goqu.I("id").Gt(10)).Limit(1).ToSQL()
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    rows, err := stmt.QueryRowContext(ctx, args)
//    if err != nil{
//        panic(err.Error()) //you could gracefully handle the error also
//    }
//    //scan your row
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	d.Trace("QUERY ROW", query, args...)
	return d.Db.QueryRowContext(ctx, query, args...)
}

func (d *Database) queryFactory() exec.QueryFactory {
	d.qfOnce.Do(func() {
		d.qf = exec.NewQueryFactory(d)
	})
	return d.qf
}

// Queries the database using the supplied query, and args and uses CrudExec.ScanStructs to scan the results into a
// slice of structs
//
// i: A pointer to a slice of structs
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) ScanStructs(i interface{}, query string, args ...interface{}) error {
	return d.ScanStructsContext(context.Background(), i, query, args...)
}

// Queries the database using the supplied context, query, and args and uses CrudExec.ScanStructsContext to scan the
// results into a slice of structs
//
// i: A pointer to a slice of structs
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) ScanStructsContext(ctx context.Context, i interface{}, query string, args ...interface{}) error {
	return d.queryFactory().FromSQL(query, args...).ScanStructsContext(ctx, i)
}

// Queries the database using the supplied query, and args and uses CrudExec.ScanStruct to scan the results into a
// struct
//
// i: A pointer to a struct
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) ScanStruct(i interface{}, query string, args ...interface{}) (bool, error) {
	return d.ScanStructContext(context.Background(), i, query, args...)
}

// Queries the database using the supplied context, query, and args and uses CrudExec.ScanStructContext to scan the
// results into a struct
//
// i: A pointer to a struct
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) ScanStructContext(ctx context.Context, i interface{}, query string, args ...interface{}) (bool, error) {
	return d.queryFactory().FromSQL(query, args...).ScanStructContext(ctx, i)
}

// Queries the database using the supplied query, and args and uses CrudExec.ScanVals to scan the results into a slice
// of primitive values
//
// i: A pointer to a slice of primitive values
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) ScanVals(i interface{}, query string, args ...interface{}) error {
	return d.ScanValsContext(context.Background(), i, query, args...)
}

// Queries the database using the supplied context, query, and args and uses CrudExec.ScanValsContext to scan the
// results into a slice of primitive values
//
// i: A pointer to a slice of primitive values
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) ScanValsContext(ctx context.Context, i interface{}, query string, args ...interface{}) error {
	return d.queryFactory().FromSQL(query, args...).ScanValsContext(ctx, i)
}

// Queries the database using the supplied query, and args and uses CrudExec.ScanVal to scan the results into a
// primitive value
//
// i: A pointer to a primitive value
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) ScanVal(i interface{}, query string, args ...interface{}) (bool, error) {
	return d.ScanValContext(context.Background(), i, query, args...)
}

// Queries the database using the supplied context, query, and args and uses CrudExec.ScanValContext to scan the
// results into a primitive value
//
// i: A pointer to a primitive value
//
// query: The SQL to execute
//
// args...: for any placeholder parameters in the query
func (d *Database) ScanValContext(ctx context.Context, i interface{}, query string, args ...interface{}) (bool, error) {
	return d.queryFactory().FromSQL(query, args...).ScanValContext(ctx, i)
}

// A wrapper around a sql.Tx and works the same way as Database
type (
	// Interface for sql.Tx, an interface is used so you can use with other
	// libraries such as sqlx instead of the native sql.DB
	SQLTx interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		Commit() error
		Rollback() error
	}
	TxDatabase struct {
		logger  Logger
		dialect string
		Tx      SQLTx
		qf      exec.QueryFactory
		qfOnce  sync.Once
	}
)

// Creates a new TxDatabase
func NewTx(dialect string, tx SQLTx) *TxDatabase {
	return &TxDatabase{dialect: dialect, Tx: tx}
}

// returns this databases dialect
func (td *TxDatabase) Dialect() string {
	return td.dialect
}

// Creates a new Dataset for querying a Database.
func (td *TxDatabase) From(cols ...interface{}) *SelectDataset {
	return newDataset(td.dialect, td.queryFactory()).From(cols...)
}

func (td *TxDatabase) Select(cols ...interface{}) *SelectDataset {
	return newDataset(td.dialect, td.queryFactory()).Select(cols...)
}

func (td *TxDatabase) Update(table interface{}) *UpdateDataset {
	return newUpdateDataset(td.dialect, td.queryFactory()).Table(table)
}

func (td *TxDatabase) Insert(table interface{}) *InsertDataset {
	return newInsertDataset(td.dialect, td.queryFactory()).Into(table)
}

func (td *TxDatabase) Delete(table interface{}) *DeleteDataset {
	return newDeleteDataset(td.dialect, td.queryFactory()).From(table)
}

func (td *TxDatabase) Truncate(table ...interface{}) *TruncateDataset {
	return newTruncateDataset(td.dialect, td.queryFactory()).Table(table...)
}

// Sets the logger
func (td *TxDatabase) Logger(logger Logger) {
	td.logger = logger
}

func (td *TxDatabase) Trace(op, sqlString string, args ...interface{}) {
	if td.logger != nil {
		if sqlString != "" {
			if len(args) != 0 {
				td.logger.Printf("[goqu - transaction] %s [query:=`%s` args:=%+v] ", op, sqlString, args)
			} else {
				td.logger.Printf("[goqu - transaction] %s [query:=`%s`] ", op, sqlString)
			}
		} else {
			td.logger.Printf("[goqu - transaction] %s", op)
		}
	}
}

// See Database#Exec
func (td *TxDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return td.ExecContext(context.Background(), query, args...)
}

// See Database#ExecContext
func (td *TxDatabase) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	td.Trace("EXEC", query, args...)
	return td.Tx.ExecContext(ctx, query, args...)
}

// See Database#Prepare
func (td *TxDatabase) Prepare(query string) (*sql.Stmt, error) {
	return td.PrepareContext(context.Background(), query)
}

// See Database#PrepareContext
func (td *TxDatabase) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	td.Trace("PREPARE", query)
	return td.Tx.PrepareContext(ctx, query)
}

// See Database#Query
func (td *TxDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return td.QueryContext(context.Background(), query, args...)
}

// See Database#QueryContext
func (td *TxDatabase) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	td.Trace("QUERY", query, args...)
	return td.Tx.QueryContext(ctx, query, args...)
}

// See Database#QueryRow
func (td *TxDatabase) QueryRow(query string, args ...interface{}) *sql.Row {
	return td.QueryRowContext(context.Background(), query, args...)
}

// See Database#QueryRowContext
func (td *TxDatabase) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	td.Trace("QUERY ROW", query, args...)
	return td.Tx.QueryRowContext(ctx, query, args...)
}

func (td *TxDatabase) queryFactory() exec.QueryFactory {
	td.qfOnce.Do(func() {
		td.qf = exec.NewQueryFactory(td)
	})
	return td.qf
}

// See Database#ScanStructs
func (td *TxDatabase) ScanStructs(i interface{}, query string, args ...interface{}) error {
	return td.ScanStructsContext(context.Background(), i, query, args...)
}

// See Database#ScanStructsContext
func (td *TxDatabase) ScanStructsContext(ctx context.Context, i interface{}, query string, args ...interface{}) error {
	return td.queryFactory().FromSQL(query, args...).ScanStructsContext(ctx, i)
}

// See Database#ScanStruct
func (td *TxDatabase) ScanStruct(i interface{}, query string, args ...interface{}) (bool, error) {
	return td.ScanStructContext(context.Background(), i, query, args...)
}

// See Database#ScanStructContext
func (td *TxDatabase) ScanStructContext(ctx context.Context, i interface{}, query string, args ...interface{}) (bool, error) {
	return td.queryFactory().FromSQL(query, args...).ScanStructContext(ctx, i)
}

// See Database#ScanVals
func (td *TxDatabase) ScanVals(i interface{}, query string, args ...interface{}) error {
	return td.ScanValsContext(context.Background(), i, query, args...)
}

// See Database#ScanValsContext
func (td *TxDatabase) ScanValsContext(ctx context.Context, i interface{}, query string, args ...interface{}) error {
	return td.queryFactory().FromSQL(query, args...).ScanValsContext(ctx, i)
}

// See Database#ScanVal
func (td *TxDatabase) ScanVal(i interface{}, query string, args ...interface{}) (bool, error) {
	return td.ScanValContext(context.Background(), i, query, args...)
}

// See Database#ScanValContext
func (td *TxDatabase) ScanValContext(ctx context.Context, i interface{}, query string, args ...interface{}) (bool, error) {
	return td.queryFactory().FromSQL(query, args...).ScanValContext(ctx, i)
}

// COMMIT the transaction
func (td *TxDatabase) Commit() error {
	td.Trace("COMMIT", "")
	return td.Tx.Commit()
}

// ROLLBACK the transaction
func (td *TxDatabase) Rollback() error {
	td.Trace("ROLLBACK", "")
	return td.Tx.Rollback()
}

// A helper method that will automatically COMMIT or ROLLBACK once the supplied function is done executing
//
//      tx, err := db.Begin()
//      if err != nil{
//           panic(err.Error()) // you could gracefully handle the error also
//      }
//      if err := tx.Wrap(func() error{
//          if _, err := tx.From("test").Insert(Record{"a":1, "b": "b"}).Exec(){
//              // this error will be the return error from the Wrap call
//              return err
//          }
//          return nil
//      }); err != nil{
//           panic(err.Error()) // you could gracefully handle the error also
//      }
func (td *TxDatabase) Wrap(fn func() error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			_ = td.Rollback()
			panic(p)
		}
		if err != nil {
			if rollbackErr := td.Rollback(); rollbackErr != nil {
				err = rollbackErr
			}
		} else {
			if commitErr := td.Commit(); commitErr != nil {
				err = commitErr
			}
		}
	}()
	return fn()
}
