# v9.18.0
* [FEATURE] Add support for aliasing insert datasets to support upsert alias [#306](https://github.com/doug-martin/goqu/pull/306) - [@XIELongDragon](https://github.com/XIELongDragon)
* [FEATURE] Add support for aliasing BooleanExpressions [#307](https://github.com/doug-martin/goqu/pull/307) - [@XIELongDragon](https://github.com/XIELongDragon)

# v9.17.0
* [FEATURE] Add support bitwise operations [#303](https://github.com/doug-martin/goqu/pull/303) - [@XIELongDragon](https://github.com/XIELongDragon)
* [FEATURE] Add support for specifying tables to be locked in ForUpdate, ForNoKeyUpdate, ForKeyShare, ForShare [#299](https://github.com/doug-martin/goqu/pull/299) - [@jbub](https://github.com/jbub)

# v9.16.0
* [FEATURE] Allow ordering by case expression [#282](https://github.com/doug-martin/goqu/issues/282), [#292](https://github.com/doug-martin/goqu/pull/292)

# v9.15.1
* [FIX] Field indexes in the columnMap getting overridden [290](https://github.com/doug-martin/goqu/issues/290), [#291](https://github.com/doug-martin/goqu/pull/291)

# v9.15.0
* [FEATURE] Add SetDefaultPrepared which controls query value interpolation [#288](https://github.com/doug-martin/goqu/pull/288) - [@Deiz](https://github.com/Deiz)

# v9.14.0
* [FEATURE] GroupByAppend to the SelectDataset and SelectClauses [#276](https://github.com/doug-martin/goqu/pull/276), [#287](https://github.com/doug-martin/goqu/pull/287) - [@ashishkf](https://github.com/ashishkf)
* [FEATURE] Allow untagged fields to be ignored [#285](https://github.com/doug-martin/goqu/pull/285) - [@Deiz](https://github.com/Deiz)
* [FIX] Nil valuer [#277](https://github.com/doug-martin/goqu/pull/277) - [@benzolium](https://github.com/benzolium), [@Diggs](https://github.com/Diggs)
* [FIX] Fix old import URL in doc comments [#286](https://github.com/doug-martin/goqu/pull/286) - [@maito1201](https://github.com/maito1201)

# v9.13.0

* [ADDED] ScanStructs, ScanVals to Scanner interface [#273](https://github.com/doug-martin/goqu/issues/273) - [@vlanse](https://github.com/vlanse)

# v9.12.0

* Update golangci-lint and updates for all associated linters
* Update dependencies
  * github.com/DATA-DOG/go-sqlmock  v1.3.3 -> v1.5.0
  * github.com/denisenkom/go-mssqldb v0.0.0-20200206145737-bbfc9a55622e -> v0.10.0
  * github.com/go-sql-driver/mysql v1.4.1 -> v1.6.0
  * github.com/lib/pq v1.2.0 -> v1.10.1
  * github.com/mattn/go-sqlite3 v1.11.0 -> v1.14.7
  * github.com/stretchr/testify -> v1.4.0 -> v1.7.0
  

# v9.11.1

* [FIXED] Avoid mutation of join slice for separate datasets when joins slice capacity is not yet reached [#261](https://github.com/doug-martin/goqu/pull/261) - [@fhaifler](https://github.com/fhaifler)

# v9.11.0

* [FIXED] Use valid 'IS' operator for sqlserver dialect [#240](https://github.com/doug-martin/goqu/pull/240), [#239](https://github.com/doug-martin/goqu/pull/229)  - [@vlanse](https://github.com/vlanse)
* [ADDED] Implement Orderable interface for SQL Functions [#251](https://github.com/doug-martin/goqu/pull/251)  - [@GlebBeloded](https://github.com/GlebBeloded)
* [ADDED] Support for table hint in multi-table MySQL DELETE queries [#252](https://github.com/doug-martin/goqu/pull/252) - [@vlanse](https://github.com/vlanse)

# v9.10.0

* [FIXED] SELECT inherits dialect from INSERT in INSERT FROM SELECT.  [#229](https://github.com/doug-martin/goqu/pull/229), [#223](https://github.com/doug-martin/goqu/issues/223) - [@vlanse](https://github.com/vlanse)
* [FIXED] SQLServer dialect: support prepared statements with TOP.  [#230](https://github.com/doug-martin/goqu/pull/230), [#225](https://github.com/doug-martin/goqu/issues/225) - [@vlanse](https://github.com/vlanse)
* [ADDED] IsPrepared to SQLExpression interface.  [#231](https://github.com/doug-martin/goqu/pull/231) - [@vlanse](https://github.com/vlanse)

# v9.9.0

* [FIXED] SQLite do not add FOR UPDATE in SELECT. [#218](https://github.com/doug-martin/goqu/pull/218) - [@vlanse](https://github.com/vlanse)
* [ADDED] Support for INSERT ON CONFLICT in SQLite. [#218](https://github.com/doug-martin/goqu/pull/218) - [@vlanse](https://github.com/vlanse)

# v9.8.0

* [ADDED] Support for ANY and ALL operators. [#196](https://github.com/doug-martin/goqu/issues/196)
* [ADDED] Support for CASE statements [#193](https://github.com/doug-martin/goqu/issues/193) 
* [ADDED] Support for getting column identifiers from AliasExpressions. [#203](https://github.com/doug-martin/goqu/issues/203)  

# v9.7.1

* Fix all formatting for golangci-lint
* Move to golangci-lint github action

# v9.7.0

* [ADDED] Support for sqlserver dialect [#197](https://github.com/doug-martin/goqu/issues/197),[#205](https://github.com/doug-martin/goqu/issues/205) - [@vlanse](https://github.com/vlanse)

# v9.6.0

* [ADDED] Support for Lateral queries [#182](https://github.com/doug-martin/goqu/issues/182)

# v9.5.1

* [FIXED] WITH clause without a RETURNING clause will panic [#177](https://github.com/doug-martin/goqu/issues/177)
* [FIXED] SQlite dialect escapes single quotes wrong, leads to SQL syntax error [#178](https://github.com/doug-martin/goqu/issues/178)
* [FIXED] Fix ReturnsColumns() nil pointer panic [#181](https://github.com/doug-martin/goqu/issues/181) - [@yeaha](https://github.com/yeaha)
* [FIXED] SelectDataset From with Error [#183](https://github.com/doug-martin/goqu/issues/183)
* [FIXED] Unable to execute union with order by expression [#185](https://github.com/doug-martin/goqu/issues/185)

# v9.5.0

* [ADDED] Ability to use regexp like, ilike, notlike, and notilike without a regexp [#172](https://github.com/doug-martin/goqu/issues/172)

# v9.4.0

* [ADDED] Ability to scan into struct fields from multiple tables [#160](https://github.com/doug-martin/goqu/issues/160)

# v9.3.0

* [ADDED] Using Update, Insert, or Delete datasets in sub selects and CTEs [#164](https://github.com/doug-martin/goqu/issues/164)

# v9.2.0

* [ADDED] exec.Scanner: New exposed scanner supports iterative scanning [#157](https://github.com/doug-martin/goqu/pull/157) - [@akarl](https://github.com/akarl)

# v9.1.0

* [FIXED] ExampleDoUpdate does't work in postgres [#156](https://github.com/doug-martin/goqu/issues/156)
* [FIXED] Issue with timezone being lost [#163](https://github.com/doug-martin/goqu/issues/163) 

# v9.0.1

* [FIXED] Issue where `NULL`, `TRUE` and `FALSE` are interpolated when using an `IS` clause. [#165](https://github.com/doug-martin/goqu/issues/165)
    
# v9.0.0

* Changed `NULL`, `TRUE`, `FALSE` to not be interpolated when creating prepared statements. [#132](https://github.com/doug-martin/goqu/pull/132), [#158](https://github.com/doug-martin/goqu/pull/158) - [@marshallmcmullen](https://github.com/marshallmcmullen)
* Updated dependencies
    * `github.com/lib/pq v1.1.1 -> v1.2.0`
    * `github.com/mattn/go-sqlite3 v1.10.0 -> v1.11.0`
    * `github.com/stretchr/testify v1.3.0 -> v1.4.0`

## v8.6.0

* [ADDED] `SetError()` and `Error()` to all datasets. [#152](https://github.com/doug-martin/goqu/pull/152) and [#150](https://github.com/doug-martin/goqu/pull/150) - [@marshallmcmullen](https://github.com/marshallmcmullen)

## v8.5.0

* [ADDED] Window Function support [#128](https://github.com/doug-martin/goqu/issues/128) - [@Xuyuanp](https://github.com/Xuyuanp)

## v8.4.1

* [FIXED] Returning func be able to handle nil [#140](https://github.com/doug-martin/goqu/issues/140)

## v8.4.0

* Created new `sqlgen` module to encapsulate sql generation
    * Broke SQLDialect inti new SQL generators for each statement type.
* Test refactor
    * Moved to a test case pattern to allow for quickly adding new test cases.
    
## v8.3.2

* [FIXED] Data race during query factory initialization [#133](https://github.com/doug-martin/goqu/issues/133) and [#136](https://github.com/doug-martin/goqu/issues/136) - [@o1egl](https://github.com/o1egl)    

## v8.3.1

* [FIXED] InsertDataset.WithDialect return old dataset [#126](https://github.com/doug-martin/goqu/issues/126) - [@chen56](https://github.com/chen56)
* Test clean up and more testing pattern consistency
    * Changed to use assertion methods off of suite
    * Updated Equals assertions to have expected output first 
* Increase overall test coverage.

## v8.3.0

* [Added] Support for `DISTINCT ON` clauses [#119](https://github.com/doug-martin/goqu/issues/119)

## v8.2.2

* [FIX] Scanner errors on pointers to primitive values [#122](https://github.com/doug-martin/goqu/issues/122)

## v8.2.1

* [FIX] Return an error when an empty identifier is encountered [#115](https://github.com/doug-martin/goqu/issues/115)

## v8.2.0

* [FIX] Fix reflection errors related to nil pointers and unexported fields [#118](https://github.com/doug-martin/goqu/issues/118)
    * Unexported fields are ignored when creating a columnMap
    * Nil embedded pointers will no longer cause a panic
    * Fields on nil embedded pointers will be ignored when creating update or insert statements.
* [ADDED] You can now ingore embedded structs and their fields by using `db:"-"` tag on the embedded struct.

## v8.1.0

* [ADDED] Support column DEFAULT when inserting/updating via struct [#27](https://github.com/doug-martin/goqu/issues/27)

## v8.0.1

* [ADDED] Multi table update support for `mysql` and `postgres` [#60](https://github.com/doug-martin/goqu/issues/60)
* [ADDED] `goqu.V` so values can be used on the LHS of expressions [#104](https://github.com/doug-martin/goqu/issues/104)

## v8.0.0

A major change the the API was made in `v8` to seperate concerns between the different SQL statement types. 

**Why the change?**

1. There were feature requests that could not be cleanly implemented with everything in a single dataset. 
2. Too much functionality was encapsulated in a single datastructure.
    * It was unclear what methods could be used for each SQL statement type.
    * Changing a feature for one statement type had the possiblity of breaking another statement type.
    * Test coverage was decent but was almost solely concerned about SELECT statements, breaking them up allowed for focused testing on each statement type.
    * Most the SQL generation methods (`ToInsertSQL`, `ToUpdateSQL` etc.) took arguments which lead to an ugly API that was not uniform for each statement type, and proved to be inflexible.

**What Changed**

There are now five dataset types, `SelectDataset`, `InsertDataset`, `UpdateDataset`, `DeleteDataset` and `TruncateDataset`

Each dataset type has its own entry point.

* `goqu.From`, `Database#From`, `DialectWrapper#From` - Create SELECT
* `goqu.Insert`, `Database#Insert`, `DialectWrapper#Insert` - Create INSERT
* `goqu.Update`, `Database#db.Update`, `DialectWrapper#Update` - Create UPDATE
* `goqu.Delete`, `Database#Delete`, `DialectWrapper#Delete` - Create DELETE
* `goqu.Truncate`, `Database#Truncate`, `DialectWrapper#Truncate` - Create TRUNCATE
  
`ToInsertSQL`, `ToUpdateSQL`, `ToDeleteSQL`, and `ToTruncateSQL` (and variations of them) methods have been removed from the `SelectDataset`. Instead use the `ToSQL` methods on each dataset type.

Each dataset type will have an `Executor` and `ToSQL` method so a common interface can be created for each type.


## v7.4.0

* [FIXED] literalTime use t.UTC() , This behavior is different from the original sql.DB [#106](https://github.com/doug-martin/goqu/issues/106) - [chen56](https://github.com/chen56)
* [ADDED] Add new method WithTx for Database [#108](https://github.com/doug-martin/goqu/issues/108) - [Xuyuanp](https://github.com/Xuyuanp)

## v7.3.1

* [ADDED] Exposed `goqu.NewTx` to allow creating a goqu tx directly from a `sql.Tx` instead of using `goqu.Database#Begin` [#95](https://github.com/doug-martin/goqu/issues/95)
* [ADDED] `goqu.Database.BeginTx` [#98](https://github.com/doug-martin/goqu/issues/98)

## v7.3.0

* [ADDED] UPDATE and INSERT should use struct Field name if db tag is not specified [#57](https://github.com/doug-martin/goqu/issues/57)
* [CHANGE] Changed goqu.Database to accept a SQLDatabase interface to allow using goqu.Database with other libraries such as `sqlx` [#95](https://github.com/doug-martin/goqu/issues/95)

## v7.2.0

* [FIXED] Sqlite3 does not accept SELECT * UNION (SELECT *) [#79](https://github.com/doug-martin/goqu/issues/79)
* [FIXED] Where(Ex{}) causes panics [mysql] [#49](https://github.com/doug-martin/goqu/issues/49)
* [ADDED] Support for OrderPrepend [#61](https://github.com/doug-martin/goqu/issues/61)
* [DOCS] Added new section about loading a dialect and using it to build SQL [#44](https://github.com/doug-martin/goqu/issues/44)

## v7.1.0

* [FIXED] Embedded pointers with property names that duplicate parent struct properties. [#23](https://github.com/doug-martin/goqu/issues/23)
* [FIXED] Can't scan values using []byte or []string [#90](https://github.com/doug-martin/goqu/issues/90)
    * When a slice that is `*sql.RawBytes`, `*[]byte` or `sql.Scanner` no errors will be returned. 

## v7.0.1

* Fix issue where structs with pointer fields where not set properly [#86](https://github.com/doug-martin/goqu/pull/86) and [#89](https://github.com/doug-martin/goqu/pull/89) - [@efureev](https://github.com/efureev)

## v7.0.0

**Linting**
* Add linting checks and fixed errors 
    * Renamed all snake_case variables to be camelCase.     
    * Fixed examples to always map to a defined method
* Renamed `adapters` to `dialect` to more closely match their intended purpose.

**API Changes**
* Updated all sql generations methods to from `Sql` to `SQL`
    * `ToSql` -> `ToSQL`
    * `ToInsertSql` -> `ToInsertSQL`
    * `ToUpdateSql` -> `ToUpdateSQL`
    * `ToDeleteSql` -> `ToDeleteSQL`
    * `ToTruncateSql` -> `ToTruncateSQL`
* Abstracted out `dialect_options` from the adapter to make the dialect self contained.
    * This also removed the dataset<->adapter co dependency making the dialect self contained.
* Refactored the `goqu.I` method.
    * Added new `goqu.S`, `goqu.T` and `goqu.C` methods to clarify why type of identifier you are using.
    * `goqu.I` should only be used when you have a qualified identifier (e.g. `goqu.I("my_schema.my_table.my_col")
* Added new `goqu.Dialect` method to make using `goqu` as an SQL builder easier.

**Internal Changes**
* Pulled expressions into their own package
    * Broke up expressions.go into multiple files to make working with and defining them easier.
    * Moved the user facing methods into the main `goqu` to keep the same API as before.
* Added more examples
* Moved non-user facing structs and interfaces to internal modules to clean up API.
* Increased test coverage.
 

## v6.1.0

* Handle nil *time.Time Literal [#73](https://github.com/doug-martin/goqu/pull/73) and [#52](https://github.com/doug-martin/goqu/pull/52) - [@RoarkeRandall](https://github.com/RoarkeRandall) and [@quetz](https://github.com/quetz)
* Add ability to change column rename function [#66](https://github.com/doug-martin/goqu/pull/66) - [@blainehansen](https://github.com/blainehansen)

## v6.0.0

* Updated go support to `1.10`, `1.11` and `1.12`
* Change testify dependency from c2fo/testify back to stretchr/testify.
* Add support for "FOR UPDATE" and "SKIP LOCKED" [#62](https://github.com/doug-martin/goqu/pull/62) - [@btubbs](https://github.com/btubbs)
* Changed to use go modules

## v5.0.0

* Drop go 1.6 support, supported versions are `1.8`, `1.9` and latest
* Add context support [#64](https://github.com/doug-martin/goqu/pull/64) - [@cmoad](https://github.com/cmoad)

## v4.2.0

* Add support for ON CONFLICT when using a dataset [#55](https://github.com/doug-martin/goqu/pull/55) - [@bobrnor](https://github.com/bobrnor)

## v4.1.0

* Support for defining WITH clauses for Common Table Expressions (CTE) [#39](https://github.com/doug-martin/goqu/pull/39) - [@Oscil8](https://github.com/Oscil8)

## v4.0

* Prepared(true) issues when using IS NULL comparisson operation [#33](https://github.com/doug-martin/goqu/pull/33) - [@danielfbm](https://github.com/danielfbm)

## v3.3

* Add `upsert` support via `InsertIgnore` and `InsertConflict` methods - [#25](https://github.com/doug-martin/goqu/pull/28) - [@aheuermann](https://github.com/aheuermann)
* Adding vendor dependencies and updating tests to run in docker containers [#29](https://github.com/doug-martin/goqu/pull/29) - [@aheuermann](https://github.com/aheuermann)

## v3.2

* Add range clauses ([NOT] BETWEEN) support - [#25](https://github.com/doug-martin/goqu/pull/25) - [@denisvm](https://github.com/denisvm)
* Readmefix [#26](https://github.com/doug-martin/goqu/pull/26) - [@tiagopotencia](https://github.com/tiagopotencia)

## v3.1.3

* Bugfix for chained Where() [#20](https://github.com/doug-martin/goqu/pull/20) - [@Emreu](https://github.com/Emreu)


## v3.1.2

* Fixing ScanStruct issue with embedded pointers in crud_exec [#20](https://github.com/doug-martin/goqu/pull/20) - [@ruzz311](https://github.com/ruzz311)

## v3.1.1

* Fixing race condition with struct_map_cache in crud_exec [#18](https://github.com/doug-martin/goqu/pull/18) - [@andymoon](https://github.com/andymoon), [@aheuermann](https://github.com/aheuermann)

## v3.1.0

* Version 3.1 [#14](https://github.com/doug-martin/goqu/pull/14) - [@andymoon](https://github.com/andymoon)
    * Fix an issue with a nil pointer access on the inserts and updates.
    * Allowing ScanStructs to take a struct with an embedded pointer to a struct.
    * Change to check if struct is Anonymous when recursing through an embedded struct.
    * Updated to use the latest version of github.com/DATA-DOG/go-sqlmock.

## v3.0.1

* Add literal bytes and update to c2fo testify [#15](https://github.com/doug-martin/goqu/pull/15) - [@TechnotronicOz](https://github.com/TechnotronicOz)

## v3.0.0

* Added support for embedded structs when inserting or updating. [#13](https://github.com/doug-martin/goqu/pull/13) - [@andymoon](https://github.com/andymoon)

## v2.0.3

* Fixed issue with transient columns and the auto select of columns.

## v2.0.2

* Changed references to "github.com/doug-martin/goqu" to "gopkg.in/doug-martin/goqu.v2"

## v2.0.1

* Fixed issue when `ScanStruct(s)` was used with `SelectDistinct` and caused a panic.

## v2.0.0

* When scanning a struct or slice of structs, the struct(s) will be parsed for the column names to select. [#9](https://github.com/doug-martin/goqu/pull/9) - [@technotronicoz](https://github.com/TechnotronicOz)

## v1.0.0

* You can now passed an IdentiferExpression to `As` [#8](https://github.com/doug-martin/goqu/pull/8) - [@croachrose](https://github.com/croachrose)
* Added info about installation through [gopkg.in](http://labix.org/gopkg.in)

## v0.3.1

* Fixed issue setting Logger when starting a new transaction.

## v0.3.0

* Changed sql generation methods to use a common naming convention. `To(Sql|Insert|Update|Delete)`
   * Also changed to have common return values `string, []interface{}, error)`
* Added `Dataset.Prepared` which allows a user to specify whether or not SQL should be interpolated. [#7](https://github.com/doug-martin/goqu/issues/7)
* Updated Docs
    * More examples
* Increased test coverage.

## v0.2.0

* Changed `CrudExec` to not wrap driver errors in a GoquError [#2](https://github.com/doug-martin/goqu/issues/2)
* Added ability to use a dataset in an `Ex` map or `Eq` expression without having to use `In` [#3](https://github.com/doug-martin/goqu/issues/3)
   * `db.From("test").Where(goqu.Ex{"a": db.From("test").Select("b")})`
* Updated readme with links to [`DefaultAdapter`](https://godoc.org/github.com/doug-martin/goqu#DefaultAdapter)

## v0.1.1

* Added SQLite3 adapter [#1](https://github.com/doug-martin/goqu/pull/1) - [@mattn](https://github.com/mattn)

## v0.1.0

* Added:
    * [`Ex`](https://godoc.org/github.com/doug-martin/goqu#Ex)
    * [`ExOr`](https://godoc.org/github.com/doug-martin/goqu#ExOr)
    * [`Op`](https://godoc.org/github.com/doug-martin/goqu#Op)
* More tests and examples
* Added CONTRIBUTING.md
* Added LICENSE information
* Removed godoc introduction in favor of just maintaining the README.

## v0.0.2

* Fixed issue with goqu.New not returning a pointer to a Database

## v0.0.1

* Initial release
