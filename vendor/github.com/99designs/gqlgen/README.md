# gqlgen [![Continuous Integration](https://github.com/99designs/gqlgen/workflows/Continuous%20Integration/badge.svg)](https://github.com/99designs/gqlgen/actions) [![Read the Docs](https://badgen.net/badge/docs/available/green)](http://gqlgen.com/)

## What is gqlgen?

[gqlgen](https://github.com/99designs/gqlgen) is a Go library for building GraphQL servers without any fuss. gqlgen is:

- **Schema first** — Define your API using the GraphQL [Schema Definition Language](http://graphql.org/learn/schema/).
- **Type safe** — You should never see `map[string]interface{}` here.
- **Codegen** — Let us generate the boring bits, so you can build your app quickly.

[Feature Comparison](https://gqlgen.com/feature-comparison/)

## Getting Started

First work your way through the [Getting Started](https://gqlgen.com/getting-started/) tutorial.

If you can't find what your looking for, look at our [examples](https://github.com/99designs/gqlgen/tree/master/example) for example usage of gqlgen.

## Reporting Issues

If you think you've found a bug, or something isn't behaving the way you think it should, please raise an [issue](https://github.com/99designs/gqlgen/issues) on GitHub.

## Contributing

Read our [Contribution Guidelines](https://github.com/99designs/gqlgen/blob/master/CONTRIBUTING.md) for information on how you can help out gqlgen.

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

You need to tell gqlgen that we should only fetch friends if the user requested it. There are two ways to do this.

#### Custom Models

Write a custom model that omits the Friends model:

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

#### Explicit Resolvers

If you want to Keep using the generated model: mark the field as requiring a resolver explicitly in `gqlgen.yml`:

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

### IDs are strings but I like ints, why cant I have ints?

You can by remapping it in config:

```yaml
models:
  ID: # The GraphQL type ID is backed by
    model:
      - github.com/99designs/gqlgen/graphql.IntID # An go integer
      - github.com/99designs/gqlgen/graphql.ID # or a go string
```

This means gqlgen will be able to automatically bind to strings or ints for models you have written yourself, but the
first model in this list is used as the default type and it will always be used when:

- Generating models based on schema
- As arguments in resolvers

There isnt any way around this, gqlgen has no way to know what you want in a given context.

## Other Resources

- [Christopher Biscardi @ Gophercon UK 2018](https://youtu.be/FdURVezcdcw)
- [Introducing gqlgen: a GraphQL Server Generator for Go](https://99designs.com.au/blog/engineering/gqlgen-a-graphql-server-generator-for-go/)
- [Dive into GraphQL by Iván Corrales Solera](https://medium.com/@ivan.corrales.solera/dive-into-graphql-9bfedf22e1a)
- [Sample Project built on gqlgen with Postgres by Oleg Shalygin](https://github.com/oshalygin/gqlgen-pg-todo-example)
