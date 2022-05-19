---
linkTitle: Configuration
title: How to configure gqlgen using gqlgen.yml
description: How to configure gqlgen using gqlgen.yml
menu: main
weight: -5
---

gqlgen can be configured using a `gqlgen.yml` file, by default it will be loaded from the current directory, or any parent directory.

Example:

```yml
# Where are all the schema files located? globs are supported eg  src/**/*.graphqls
schema:
  - graph/*.graphqls

# Where should the generated server code go?
exec:
  layout: follow-schema
  dir: graph/generated
  package: generated

# Enable Apollo federation support
federation:
  filename: graph/generated/federation.go
  package: generated

# Where should any generated models go?
model:
  filename: graph/model/models_gen.go
  package: model

# Where should the resolver implementations go?
resolver:
  layout: follow-schema
  dir: graph
  package: graph
  filename_template: "{name}.resolvers.go"

# Optional: turn on use ` + "`" + `gqlgen:"fieldName"` + "`" + ` tags in your models
# struct_tag: json

# Optional: turn on to use []Thing instead of []*Thing
# omit_slice_element_pointers: false

# Optional: set to speed up generation time by not performing a final validation pass.
# skip_validation: true

# Optional: set to skip running `go mod tidy` when generating server code
# skip_mod_tidy: true

# gqlgen will search for any type names in the schema in these go packages
# if they match it will use them, otherwise it will generate them.
# autobind:
#   - "github.com/[YOUR_APP_DIR]/graph/model"

# This section declares type mapping between the GraphQL and go type systems
#
# The first line in each type will be used as defaults for resolver arguments and
# modelgen, the others will be allowed when binding to fields. Configure them to
# your liking
models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
```

Everything has defaults, so add things as you need.

## Inline config with directives

gqlgen ships with some builtin directives that make it a little easier to manage wiring.

To start using them you first need to define them:

```graphql
directive @goModel(
	model: String
	models: [String!]
) on OBJECT | INPUT_OBJECT | SCALAR | ENUM | INTERFACE | UNION

directive @goField(
	forceResolver: Boolean
	name: String
) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION

directive @goTag(
	key: String!
	value: String
) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION
```

> Here be dragons
>
> gqlgen doesnt currently support user-configurable directives for SCALAR, ENUM, INTERFACE or UNION. This only works
> for internal directives. You can track the progress [here](https://github.com/99designs/gqlgen/issues/760)

Now you can use these directives when defining types in your schema:

```graphql
type User @goModel(model: "github.com/my/app/models.User") {
	id: ID! @goField(name: "todoId")
	name: String!
		@goField(forceResolver: true)
		@goTag(key: "xorm", value: "-")
		@goTag(key: "yaml")
}
```

The builtin directives `goField`, `goModel` and `goTag` are automatically registered to `skip_runtime`. Any directives registered as `skip_runtime` will not exposed during introspection and are used during code generation only. 

If you have created a new code generation plugin using a directive which does not require runtime execution, the directive will need to be set to `skip_runtime`.

e.g. a custom directive called `constraint` would be set as `skip_runtime` using the following configuration
```yml
# custom directives which are not exposed during introspection. These directives are
# used for code generation only
directives:
  constraint:
    skip_runtime: true
```
