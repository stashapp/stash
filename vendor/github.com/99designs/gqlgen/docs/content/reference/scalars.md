---
linkTitle: Scalars
title: Mapping GraphQL scalar types to Go types
description: Mapping GraphQL scalar types to Go types
menu: { main: { parent: "reference", weight: 10 } }
---

## Built-in helpers

gqlgen ships with some built-in helpers for common custom scalar use-cases, `Time`, `Any`, `Upload` and `Map`. Adding any of these to a schema will automatically add the marshalling behaviour to Go types.

### Time

```graphql
scalar Time
```

Maps a `Time` GraphQL scalar to a Go `time.Time` struct. This scalar adheres to the [time.RFC3339Nano](https://pkg.go.dev/time#pkg-constants) format.

### Map

```graphql
scalar Map
```

Maps an arbitrary GraphQL value to a `map[string]interface{}` Go type.

### Upload

```graphql
scalar Upload
```

Maps a `Upload` GraphQL scalar to a `graphql.Upload` struct, defined as follows:

```go
type Upload struct {
	File        io.Reader
	Filename    string
	Size        int64
	ContentType string
}
```

### Any

```graphql
scalar Any
```

Maps an arbitrary GraphQL value to a `interface{}` Go type.

## Custom scalars with user defined types

For user defined types you can implement the [graphql.Marshaler](https://pkg.go.dev/github.com/99designs/gqlgen/graphql#Marshaler) and [graphql.Unmarshaler](https://pkg.go.dev/github.com/99designs/gqlgen/graphql#Unmarshaler) or implement the [graphql.ContextMarshaler](https://pkg.go.dev/github.com/99designs/gqlgen/graphql#ContextMarshaler) and [graphql.ContextUnmarshaler](https://pkg.go.dev/github.com/99designs/gqlgen/graphql#ContextUnmarshaler) interfaces and they will be called.

```go
package mypkg

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

//
// Most common scalars
//

type YesNo bool

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (y *YesNo) UnmarshalGQL(v interface{}) error {
	yes, ok := v.(string)
	if !ok {
		return fmt.Errorf("YesNo must be a string")
	}

	if yes == "yes" {
		*y = true
	} else {
		*y = false
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (y YesNo) MarshalGQL(w io.Writer) {
	if y {
		w.Write([]byte(`"yes"`))
	} else {
		w.Write([]byte(`"no"`))
	}
}

//
// Scalars that need access to the request context
//

type Length float64

// UnmarshalGQLContext implements the graphql.ContextUnmarshaler interface
func (l *Length) UnmarshalGQLContext(ctx context.Context, v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("Length must be a string")
	}
	length, err := ParseLength(s)
	if err != nil {
		return err
	}
	*l = length
	return nil
}

// MarshalGQLContext implements the graphql.ContextMarshaler interface
func (l Length) MarshalGQLContext(ctx context.Context, w io.Writer) error {
	s, err := l.FormatContext(ctx)
	if err != nil {
		return err
	}
	w.Write([]byte(strconv.Quote(s)))
	return nil
}

// ParseLength parses a length measurement string with unit on the end (eg: "12.45in")
func ParseLength(string) (Length, error)

// ParseLength formats the string using a value in the context to specify format
func (l Length) FormatContext(ctx context.Context) (string, error)
```

and then wire up the type in .gqlgen.yml or via directives like normal:

```yaml
models:
  YesNo:
    model: github.com/me/mypkg.YesNo
```

## Custom scalars with third party types

Sometimes you are unable to add add methods to a type - perhaps you don't own the type, or it is part of the standard
library (eg string or time.Time). To support this we can build an external marshaler:

```go
package mypkg

import (
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)


func MarshalMyCustomBooleanScalar(b bool) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if b {
			w.Write([]byte("true"))
		} else {
			w.Write([]byte("false"))
		}
	})
}

func UnmarshalMyCustomBooleanScalar(v interface{}) (bool, error) {
	switch v := v.(type) {
	case string:
		return "true" == strings.ToLower(v), nil
	case int:
		return v != 0, nil
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("%T is not a bool", v)
	}
}
```

Then in .gqlgen.yml point to the name without the Marshal|Unmarshal in front:

```yaml
models:
  MyCustomBooleanScalar:
    model: github.com/me/mypkg.MyCustomBooleanScalar
```

**Note:** you also can un/marshal to pointer types via this approach, simply accept a pointer in your
`Marshal...` func and return one in your `Unmarshal...` func.

**Note:** you can also un/marshal with a context by having your custom marshal function return a
`graphql.ContextMarshaler` _and_ your unmarshal function take a `context.Context` as the first argument.

See the [_examples/scalars](https://github.com/99designs/gqlgen/tree/master/_examples/scalars) package for more examples.

## Marshaling/Unmarshaling Errors

The errors that occur as part of custom scalar marshaling/unmarshaling will return a full path to the field.
For example, given the following schema ...

```graphql
extend type Mutation{
    updateUser(userInput: UserInput!): User!
}

input UserInput {
    name: String!
    primaryContactDetails: ContactDetailsInput!
    secondaryContactDetails: ContactDetailsInput!
}

scalar Email
input ContactDetailsInput {
    email: Email!
}
```

... and the following variables:

```json

{
  "userInput": {
    "name": "George",
    "primaryContactDetails": {
      "email": "not-an-email"
    },
    "secondaryContactDetails": {
      "email": "george@gmail.com"
    }
  }
}
```

... and an unmarshal function that returns an error if the email is invalid. The mutation will return an error containing the full path:
```json
{
  "message": "email invalid",
  "path": [
    "updateUser",
    "userInput",
    "primaryContactDetails",
    "email"
  ]
}
```

**Note:** Marshaling errors can only be returned when using the `graphql.ContextMarshaler` style interface.
