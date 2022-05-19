---
title: Using schema directives to implement permission checks
description: Implementing graphql schema directives in golang for permission checks.
linkTitle: Schema Directives
menu: { main: { parent: 'reference', weight: 10 } }
---

Directives act a bit like annotations, decorators, or HTTP middleware. They give you a way to specify some behaviour based on a field or argument in a generic and reusable way. This can be really useful for cross-cutting concerns like permission checks which can be applied broadly across your API.

**Note**: The current directives implementation is still fairly limited, and is designed to cover the most common "field middleware" case.

## Restricting access based on user role

For example, we might want to restrict which mutations or queries a client can make based on the authenticated user's role:
```graphql
type Mutation {
	deleteUser(userID: ID!): Bool @hasRole(role: ADMIN)
}
```

### Declare it in the schema

Before we can use a directive we must declare it in the schema. Here's how we would define the `@hasRole` directive:

```graphql
directive @hasRole(role: Role!) on FIELD_DEFINITION

enum Role {
    ADMIN
    USER
}
```

Next, run `go generate` and gqlgen will add the directive to the DirectiveRoot:
```go
type DirectiveRoot struct {
	HasRole func(ctx context.Context, obj interface{}, next graphql.Resolver, role Role) (res interface{}, err error)
}
```

The arguments are:
 - *ctx*: the parent context
 - *obj*: the object containing the value this was applied to, e.g.:
    - for field definition directives (`FIELD_DEFINITION`), the object/input object that contains the field
    - for argument directives (`ARGUMENT_DEFINITION`), a map containing all arguments
 - *next*: the next directive in the directive chain, or the field resolver. This should be called to get the
           value of the field/argument/whatever. You can block access to the field by not calling `next(ctx)`
           after checking whether a user has a required permission, for example.
 - *...args*: finally, any args defined in the directive schema definition are passed in

## Implement the directive

Now we must implement the directive. The directive function is assigned to the Config object before registering the GraphQL handler.
```go
package main

func main() {
	c := generated.Config{ Resolvers: &resolvers{} }
	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role Role) (interface{}, error) {
		if !getCurrentUser(ctx).HasRole(role) {
			// block calling the next resolver
			return nil, fmt.Errorf("Access denied")
		}

		// or let it pass through
		return next(ctx)
	}

	http.Handle("/query", handler.NewDefaultServer(generated.NewExecutableSchema(c), ))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
```

That's it! You can now apply the `@hasRole` directive to any mutation or query in your schema.
