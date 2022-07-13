```
  __ _  ___   __ _ _   _
 / _` |/ _ \ / _` | | | |
| (_| | (_) | (_| | |_| |
 \__, |\___/ \__, |\__,_|
 |___/          |_|
```
[![GitHub tag](https://img.shields.io/github/tag/doug-martin/goqu.svg?style=flat)](https://github.com/doug-martin/goqu/releases)
[![Test](https://github.com/doug-martin/goqu/workflows/Test/badge.svg?branch=master&event=push)](https://github.com/doug-martin/goqu/actions?query=workflow%3ATest+and+branch%3Amaster+)
[![Go Reference](https://pkg.go.dev/badge/github.com/doug-martin/goqu/v9.svg)](https://pkg.go.dev/github.com/doug-martin/goqu/v9)
[![codecov](https://codecov.io/gh/doug-martin/goqu/branch/master/graph/badge.svg)](https://codecov.io/gh/doug-martin/goqu)
[![Go Report Card](https://goreportcard.com/badge/github.com/doug-martin/goqu/v9)](https://goreportcard.com/report/github.com/doug-martin/goqu/v9)

`goqu` is an expressive SQL builder and executor
    
If you are upgrading from an older version please read the [Migrating Between Versions](./docs/version_migration.md) docs.


## Installation

If using go modules.

```sh
go get -u github.com/doug-martin/goqu/v9
```

If you are not using go modules...

**NOTE** You should still be able to use this package if you are using go version `>v1.10` but, you will need to drop the version from the package. `import "github.com/doug-martin/goqu/v9` -> `import "github.com/doug-martin/goqu"`

```sh
go get -u github.com/doug-martin/goqu
```

### [Migrating Between Versions](./docs/version_migration.md)

## Features

`goqu` comes with many features but here are a few of the more notable ones

* Query Builder
* Parameter interpolation (e.g `SELECT * FROM "items" WHERE "id" = ?` -> `SELECT * FROM "items" WHERE "id" = 1`)
* Built from the ground up with multiple dialects in mind
* Insert, Multi Insert, Update, and Delete support
* Scanning of rows to struct[s] or primitive value[s]

While goqu may support the scanning of rows into structs it is not intended to be used as an ORM if you are looking for common ORM features like associations,
or hooks I would recommend looking at some of the great ORM libraries such as:

* [gorm](https://github.com/jinzhu/gorm)
* [hood](https://github.com/eaigner/hood)

## Why?

We tried a few other sql builders but each was a thin wrapper around sql fragments that we found error prone. `goqu` was built with the following goals in mind:

* Make the generation of SQL easy and enjoyable
* Create an expressive DSL that would find common errors with SQL at compile time.
* Provide a DSL that accounts for the common SQL expressions, NOT every nuance for each database.
* Provide developers the ability to:
  * Use SQL when desired
  * Easily scan results into primitive values and structs
  * Use the native sql.Db methods when desired

## Docs

* [Dialect](./docs/dialect.md) - Introduction to different dialects (`mysql`, `postgres`, `sqlite3`, `sqlserver` etc) 
* [Expressions](./docs/expressions.md) - Introduction to `goqu` expressions and common examples.
* [Select Dataset](./docs/selecting.md) - Docs and examples about creating and executing SELECT sql statements.
* [Insert Dataset](./docs/inserting.md) - Docs and examples about creating and executing INSERT sql statements.
* [Update Dataset](./docs/updating.md) - Docs and examples about creating and executing UPDATE sql statements.
* [Delete Dataset](./docs/deleting.md) - Docs and examples about creating and executing DELETE sql statements.
* [Prepared Statements](./docs/interpolation.md) - Docs about interpolation and prepared statements in `goqu`.
* [Database](./docs/database.md) - Docs and examples of using a Database to execute queries in `goqu`
* [Working with time.Time](./docs/time.md) - Docs on how to use alternate time locations.

## Quick Examples

### Select

See the [select dataset](./docs/selecting.md) docs for more in depth examples

```go
sql, _, _ := goqu.From("test").ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT * FROM "test"
```

```go
sql, _, _ := goqu.From("test").Where(goqu.Ex{
	"d": []string{"a", "b", "c"},
}).ToSQL()
fmt.Println(sql)
```

Output:

```
SELECT * FROM "test" WHERE ("d" IN ('a', 'b', 'c'))
```

### Insert

See the [insert dataset](./docs/inserting.md) docs for more in depth examples

```go
ds := goqu.Insert("user").
	Cols("first_name", "last_name").
	Vals(
		goqu.Vals{"Greg", "Farley"},
		goqu.Vals{"Jimmy", "Stewart"},
		goqu.Vals{"Jeff", "Jeffers"},
	)
insertSQL, args, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```

Output: 
```sql
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```

```go
ds := goqu.Insert("user").Rows(
	goqu.Record{"first_name": "Greg", "last_name": "Farley"},
	goqu.Record{"first_name": "Jimmy", "last_name": "Stewart"},
	goqu.Record{"first_name": "Jeff", "last_name": "Jeffers"},
)
insertSQL, args, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```


```go
type User struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
ds := goqu.Insert("user").Rows(
	User{FirstName: "Greg", LastName: "Farley"},
	User{FirstName: "Jimmy", LastName: "Stewart"},
	User{FirstName: "Jeff", LastName: "Jeffers"},
)
insertSQL, args, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```

```go
ds := goqu.Insert("user").Prepared(true).
	FromQuery(goqu.From("other_table"))
insertSQL, args, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" SELECT * FROM "other_table" []
```

```go
ds := goqu.Insert("user").Prepared(true).
	Cols("first_name", "last_name").
	FromQuery(goqu.From("other_table").Select("fn", "ln"))
insertSQL, args, _ := ds.ToSQL()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") SELECT "fn", "ln" FROM "other_table" []
```

### Update

See the [update dataset](./docs/updating.md) docs for more in depth examples

```go
sql, args, _ := goqu.Update("items").Set(
	goqu.Record{"name": "Test", "address": "111 Test Addr"},
).ToSQL()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
```

```go
type item struct {
	Address string `db:"address"`
	Name    string `db:"name" goqu:"skipupdate"`
}
sql, args, _ := goqu.Update("items").Set(
	item{Name: "Test", Address: "111 Test Addr"},
).ToSQL()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr' []
```

```go
sql, _, _ := goqu.Update("test").
	Set(goqu.Record{"foo": "bar"}).
	Where(goqu.Ex{
		"a": goqu.Op{"gt": 10}
	}).ToSQL()
fmt.Println(sql)
```

Output:
```
UPDATE "test" SET "foo"='bar' WHERE ("a" > 10)
```

### Delete

See the [delete dataset](./docs/deleting.md) docs for more in depth examples

```go
ds := goqu.Delete("items")

sql, args, _ := ds.ToSQL()
fmt.Println(sql, args)
```

```go
sql, _, _ := goqu.Delete("test").Where(goqu.Ex{
		"c": nil
	}).ToSQL()
fmt.Println(sql)
```

Output:
```
DELETE FROM "test" WHERE ("c" IS NULL)
```

<a name="contributions"></a>
## Contributions

I am always welcoming contributions of any type. Please open an issue or create a PR if you find an issue with any of the following.

* An issue with Documentation
* You found the documentation lacking in some way

If you have an issue with the package please include the following

* The dialect you are using
* A description of the problem
* A short example of how to reproduce (if applicable)

Without those basics it can be difficult to reproduce your issue locally. You may be asked for more information but that is a good starting point.

### New Features

New features and/or enhancements are great and I encourage you to either submit a PR or create an issue. In both cases include the following as the need/requirement may not be readily apparent.

1. The use case
2. A short example

If you are issuing a PR also include the following

1. Tests - otherwise the PR will not be merged
2. Documentation - otherwise the PR will not be merged
3. Examples - [If applicable] see example_test.go for examples

If you find an issue you want to work on please comment on it letting other people know you are looking at it and I will assign the issue to you.

If want to work on an issue but dont know where to start just leave a comment and I'll be more than happy to point you in the right direction.

### Running tests
The test suite requires a postgres, mysql and sqlserver databases. You can override the connection strings with the [`MYSQL_URI`, `PG_URI`, `SQLSERVER_URI` environment variables](https://github.com/doug-martin/goqu/blob/2fe3349/docker-compose.yml#L26)*

```sh
go test -v -race ./...
```

You can also run the tests in a container using [docker-compose](https://docs.docker.com/compose/).

```sh
MYSQL_VERSION=8 POSTGRES_VERSION=13.4 SQLSERVER_VERSION=2017-CU8-ubuntu GO_VERSION=latest docker-compose run goqu
```

## License

`goqu` is released under the [MIT License](http://www.opensource.org/licenses/MIT).





