---
linkTitle: Getting Started
title: Building GraphQL servers in golang
description: Get started building type-safe GraphQL servers in Golang using gqlgen
menu: main
weight: -7
---

This tutorial will take you through the process of building a GraphQL server with gqlgen that can:

 - Return a list of todos
 - Create new todos
 - Mark off todos as they are completed

You can find the finished code for this tutorial [here](https://github.com/vektah/gqlgen-tutorials/tree/master/gettingstarted)

## Set up Project

Create a directory for your project, and [initialise it as a Go Module](https://golang.org/doc/tutorial/create-module):

```shell
mkdir gqlgen-todos
cd gqlgen-todos
go mod init github.com/[username]/gqlgen-todos
```

Next, create a `tools.go` file and add gqlgen as a [tool dependency for your module](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module).

```go
//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
)
```

To automatically add the dependency to your `go.mod` run
```shell
go mod tidy
```

By default you'll be using the latest version of gqlgen, but if you want to specify a particular version you can use `go get` (replacing `VERSION` with the particular version desired)
```shell
go get -d github.com/99designs/gqlgen@VERSION
```



## Building the server

### Create the project skeleton

```shell
go run github.com/99designs/gqlgen init
```

This will create our suggested package layout. You can modify these paths in gqlgen.yml if you need to.
```
├── go.mod
├── go.sum
├── gqlgen.yml               - The gqlgen config file, knobs for controlling the generated code.
├── graph
│   ├── generated            - A package that only contains the generated runtime
│   │   └── generated.go
│   ├── model                - A package for all your graph models, generated or otherwise
│   │   └── models_gen.go
│   ├── resolver.go          - The root graph resolver type. This file wont get regenerated
│   ├── schema.graphqls      - Some schema. You can split the schema into as many graphql files as you like
│   └── schema.resolvers.go  - the resolver implementation for schema.graphql
└── server.go                - The entry point to your app. Customize it however you see fit
```

### Define your schema

gqlgen is a schema-first library — before writing code, you describe your API using the GraphQL
[Schema Definition Language](http://graphql.org/learn/schema/). By default this goes into a file called
`schema.graphqls` but you can break it up into as many different files as you want.

The schema that was generated for us was:
```graphql
type Todo {
  id: ID!
  text: String!
  done: Boolean!
  user: User!
}

type User {
  id: ID!
  name: String!
}

type Query {
  todos: [Todo!]!
}

input NewTodo {
  text: String!
  userId: String!
}

type Mutation {
  createTodo(input: NewTodo!): Todo!
}
```

### Implement the resolvers

When executed, gqlgen's `generate` command compares the schema file (`graph/schema.graphqls`) with the models `graph/model/*`, and, wherever it
can, it will bind directly to the model.  That was done already when `init` was run.  We'll edit the schema later in the tutorial, but for now, let's look at what was generated already.

If we take a look in `graph/schema.resolvers.go` we will see all the times that gqlgen couldn't match them up. For us
it was twice:

```go
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}
```

We just need to implement these two methods to get our server working:

First we need somewhere to track our state, lets put it in `graph/resolver.go`. The `graph/resolver.go` file is where we declare our app's dependencies, like our database. It gets initialized once in `server.go` when we create the graph.

```go
type Resolver struct{
	todos []*model.Todo
}
```

Returning to `graph/schema.resolvers.go`, let's implement the bodies of those automatically generated resolver functions.  For `CreateTodo`, we'll use `math.rand` to simply return a todo with a randomly generated ID and store that in the in-memory todos list --- in a real app, you're likely to use a database or some other backend service.

```go
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	todo := &model.Todo{
		Text:   input.Text,
		ID:     fmt.Sprintf("T%d", rand.Int()),
		User: &model.User{ID: input.UserID, Name: "user " + input.UserID},
	}
	r.todos = append(r.todos, todo)
	return todo, nil
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	return r.todos, nil
}
```

### Run the server

We now have a working server, to start it:
```bash
go run server.go
```

Open http://localhost:8080 in a browser. Here are some queries to try, starting with creating a todo:
```graphql
mutation createTodo {
  createTodo(input: { text: "todo", userId: "1" }) {
    user {
      id
    }
    text
    done
  }
}
```

And then querying for it:

```graphql
query findTodos {
  todos {
    text
    done
    user {
      name
    }
  }
}
```

### Don't eagerly fetch the user

This example is great, but in the real world fetching most objects is expensive. We dont want to load the User on the
todo unless the user actually asked for it. So lets replace the generated `Todo` model with something slightly more
realistic.

Create a new file called `graph/model/todo.go`
```go
package model

type Todo struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Done   bool   `json:"done"`
	User   *User  `json:"user"`
}
```

> Note
>
> By default gqlgen will use any models in the model directory that match on name, this can be configured in `gqlgen.yml`.

And run `go run github.com/99designs/gqlgen generate`.

>
> If you run into this error `package github.com/99designs/gqlgen: no Go files` while executing the `generate` command above, follow the instructions in [this](https://github.com/99designs/gqlgen/issues/800#issuecomment-888908950) comment for a possible solution.

Now if we look in `graph/schema.resolvers.go` we can see a new resolver, lets implement it and fix `CreateTodo`.
```go
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	todo := &model.Todo{
		Text:   input.Text,
		ID:     fmt.Sprintf("T%d", rand.Int()),
		User:   &model.User{ID: input.UserID, Name: "user " + input.UserID},
	}
	r.todos = append(r.todos, todo)
	return todo, nil
}

func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	return &model.User{ID: obj.UserID, Name: "user " + obj.UserID}, nil
}
```

## Finishing touches

At the top of our `resolver.go`, between `package` and `import`, add the following line:

```go
//go:generate go run github.com/99designs/gqlgen generate
```

This magic comment tells `go generate` what command to run when we want to regenerate our code. To run go generate recursively over your entire project, use this command:

```go
go generate ./...
```
