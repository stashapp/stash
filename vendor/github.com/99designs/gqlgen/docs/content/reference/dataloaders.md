---
title: "Optimizing N+1 database queries using Dataloaders"
description: Speeding up your GraphQL requests by reducing the number of round trips to the database.
linkTitle: Dataloaders
menu: { main: { parent: 'reference', weight: 10 } }
---

Dataloaders consolidate the retrieval of information into fewer, batched calls. This example demonstrates the value of dataloaders by consolidating many SQL queries into a single bulk query.

## The Problem

Imagine your graph has query that lists todos...

```graphql
query { todos { user { name } } }
```

and the `todo.user` resolver reads the `User` from a database...
```go
func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	res := db.LogAndQuery(
		r.Conn,
		"SELECT id, name FROM users WHERE id = ?",
		obj.UserID,
	)
	defer res.Close()

	if !res.Next() {
		return nil, nil
	}
	var user model.User
	if err := res.Scan(&user.ID, &user.Name); err != nil {
		panic(err)
	}
	return &user, nil
}
```

The query executor will call the `Query.Todos` resolver which does a `select * from todo` and returns `N` todos. If the nested `User` is selected, the above `UserRaw` resolver will run a separate query for each user, resulting in `N+1` database queries.

eg:
```sql
SELECT id, todo, user_id FROM todo
SELECT id, name FROM users WHERE id = ?
SELECT id, name FROM users WHERE id = ?
SELECT id, name FROM users WHERE id = ?
SELECT id, name FROM users WHERE id = ?
SELECT id, name FROM users WHERE id = ?
SELECT id, name FROM users WHERE id = ?
```

Whats even worse? most of those todos are all owned by the same user! We can do better than this.

## Dataloader

Dataloaders allow us to consolidate the fetching of `todo.user` across all resolvers for a given GraphQL request into a single database query and even cache the results for subsequent requests.

We're going to use [graph-gophers/dataloader](https://github.com/graph-gophers/dataloader) to implement a dataloader for bulk-fetching users.

```bash
go get -u github.com/graph-gophers/dataloader
```

Next, we implement a data loader and a middleware for injecting the data loader on a request context.

```go
package storage

// import graph gophers with your other imports
import (
	"github.com/graph-gophers/dataloader"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// UserReader reads Users from a database
type UserReader struct {
	conn *sql.DB
}

// GetUsers implements a batch function that can retrieve many users by ID,
// for use in a dataloader
func (u *UserReader) GetUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// read all requested users in a single query
	userIDs := make([]string, len(keys))
	for ix, key := range keys {
		userIDs[ix] = key.String()
	}
	res := u.db.Exec(
		r.Conn,
		"SELECT id, name
		FROM users
		WHERE id IN (?" + strings.Repeat(",?", len(userIDs-1)) + ")",
		userIDs...,
	)
	defer res.Close()
	// return User records into a map by ID
	userById := map[string]*model.User{}
	for res.Next() {
		user := model.User{}
		if err := res.Scan(&user.ID, &user.Name); err != nil {
			panic(err)
		}
		userById[user.ID] = &user
	}
	// return users in the same order requested
	output := make([]*dataloader.Result, len(keys))
	for index, userKey := range keys {
		user, ok := userById[userKey.String()]
		if ok {
			output[index] = &dataloader.Result{Data: record, Error: nil}
		} else {
			err := fmt.Errorf("user not found %s", userKey.String())
			output[index] = &dataloader.Result{Data: nil, Error: err}
		}
	}
	return output
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	UserLoader *dataloader.Loader
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *sql.DB) *Loaders {
	// define the data loader
	userReader := &UserReader{conn: conn}
	loaders := &Loaders{
		UserLoader: dataloader.NewBatchedLoader(u.GetUsers),
	}
	return loaders
}

// Middleware injects data loaders into the context
func Middleware(loaders *Loaders, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loaders)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// GetUser wraps the User dataloader for efficient retrieval by user ID
func GetUser(ctx context.Context, userID string) (*model.User, error) {
	loaders := For(ctx)
	thunk := loaders.UserLoader.Load(ctx, dataloader.StringKey(userID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.User), nil
}

```

Now lets update our resolver to call the dataloader:
```go
func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	return storage.GetUser(ctx, obj.UserID)
}
```

The end result? Just 2 queries!
```sql
SELECT id, todo, user_id FROM todo
SELECT id, name from user WHERE id IN (?,?,?,?,?)
```

You can see an end-to-end example [here](https://github.com/zenyui/gqlgen-dataloader).
