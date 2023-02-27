![gqlgen](https://user-images.githubusercontent.com/980499/133180111-d064b38c-6eb9-444b-a60f-7005a6e68222.png)


# gqlgen [![Integration](https://github.com/99designs/gqlgen/actions/workflows/integration.yml/badge.svg)](https://github.com/99designs/gqlgen/actions) [![Coverage Status](https://coveralls.io/repos/github/99designs/gqlgen/badge.svg?branch=master)](https://coveralls.io/github/99designs/gqlgen?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/99designs/gqlgen)](https://goreportcard.com/report/github.com/99designs/gqlgen) [![Go Reference](https://pkg.go.dev/badge/github.com/99designs/gqlgen.svg)](https://pkg.go.dev/github.com/99designs/gqlgen) [![Read the Docs](https://badgen.net/badge/docs/available/green)](http://gqlgen.com/)

## What is gqlgen?

[gqlgen](https://github.com/99designs/gqlgen) is a Go library for building GraphQL servers without any fuss.<br/>

- **gqlgen is based on a Schema first approach** — You get to Define your API using the GraphQL [Schema Definition Language](http://graphql.org/learn/schema/).
- **gqlgen prioritizes Type safety** — You should never see `map[string]interface{}` here.
- **gqlgen enables Codegen** — We generate the boring bits, so you can focus on building your app quickly.

Still not convinced enough to use **gqlgen**? Compare **gqlgen** with other Go graphql [implementations](https://gqlgen.com/feature-comparison/)

## Quick start
1. [Initialise a new go module](https://golang.org/doc/tutorial/create-module)

       mkdir example
       cd example
       go mod init example

2. Add `github.com/99designs/gqlgen` to your [project's tools.go](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)

       printf '// +build tools\npackage tools\nimport (_ "github.com/99designs/gqlgen"\n _ "github.com/99designs/gqlgen/graphql/introspection")' | gofmt > tools.go

       go mod tidy

3. Initialise gqlgen config and generate models

       go run github.com/99designs/gqlgen init

4. Start the graphql server

       go run server.go

More help to get started:
 - [Getting started tutorial](https://gqlgen.com/getting-started/) - a comprehensive guide to help you get started
 - [Real-world examples](https://github.com/99designs/gqlgen/tree/master/_examples) show how to create GraphQL applications
 - [Reference docs](https://pkg.go.dev/github.com/99designs/gqlgen) for the APIs

## Reporting Issues

If you think you've found a bug, or something isn't behaving the way you think it should, please raise an [issue](https://github.com/99designs/gqlgen/issues) on GitHub.

## Contributing

We welcome contributions, Read our [Contribution Guidelines](https://github.com/99designs/gqlgen/blob/master/CONTRIBUTING.md) to learn more about contributing to **gqlgen**
## Frequently asked questions

### How do I prevent fetching child objects that might not be used?

When you have nested or recursive schema like this:

```graphql
type User {
  id: ID!
  name: String!
  friends: [User!]!
}
```

You need to tell gqlgen that it should only fetch friends if the user requested it. There are two ways to do this;

- #### Using Custom Models

Write a custom model that omits the friends field:

```go
type User struct {
  ID int
  Name string
}
```

And reference the model in `gqlgen.yml`:

```yaml
# gqlgen.yml
models:
  User:
    model: github.com/you/pkg/model.User # go import path to the User struct above
```

- #### Using Explicit Resolvers

If you want to Keep using the generated model, mark the field as requiring a resolver explicitly in `gqlgen.yml` like this:

```yaml
# gqlgen.yml
models:
  User:
    fields:
      friends:
        resolver: true # force a resolver to be generated
```

After doing either of the above and running generate we will need to provide a resolver for friends:

```go
func (r *userResolver) Friends(ctx context.Context, obj *User) ([]*User, error) {
  // select * from user where friendid = obj.ID
  return friends,  nil
}
```

You can also use inline config with directives to achieve the same result

```graphql
directive @goModel(model: String, models: [String!]) on OBJECT
    | INPUT_OBJECT
    | SCALAR
    | ENUM
    | INTERFACE
    | UNION

directive @goField(forceResolver: Boolean, name: String) on INPUT_FIELD_DEFINITION
    | FIELD_DEFINITION

type User @goModel(model: "github.com/you/pkg/model.User") {
    id: ID!         @goField(name: "todoId")
    friends: [User!]!   @goField(forceResolver: true)
}
```

### Can I change the type of the ID from type String to Type Int?

Yes! You can by remapping it in config as seen below:

```yaml
models:
  ID: # The GraphQL type ID is backed by
    model:
      - github.com/99designs/gqlgen/graphql.IntID # a go integer
      - github.com/99designs/gqlgen/graphql.ID # or a go string
```

This means gqlgen will be able to automatically bind to strings or ints for models you have written yourself, but the
first model in this list is used as the default type and it will always be used when:

- Generating models based on schema
- As arguments in resolvers

There isn't any way around this, gqlgen has no way to know what you want in a given context.

### Why do my interfaces have getters? Can I disable these?
These were added in v0.17.14 to allow accessing common interface fields without casting to a concrete type.
However, certain fields, like Relay-style Connections, cannot be implemented with simple getters.

If you'd prefer to not have getters generated in your interfaces, you can add the following in your `gqlgen.yml`:
```yaml
# gqlgen.yml
omit_getters: true
```

## Other Resources

- [Christopher Biscardi @ Gophercon UK 2018](https://youtu.be/FdURVezcdcw)
- [Introducing gqlgen: a GraphQL Server Generator for Go](https://99designs.com.au/blog/engineering/gqlgen-a-graphql-server-generator-for-go/)
- [Dive into GraphQL by Iván Corrales Solera](https://medium.com/@ivan.corrales.solera/dive-into-graphql-9bfedf22e1a)
- [Sample Project built on gqlgen with Postgres by Oleg Shalygin](https://github.com/oshalygin/gqlgen-pg-todo-example)
- [Hackernews GraphQL Server with gqlgen by Shayegan Hooshyari](https://www.howtographql.com/graphql-go/0-introduction/)
