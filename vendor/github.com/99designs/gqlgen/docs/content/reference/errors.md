---
linkTitle: Handling Errors
title: Sending custom error data in the graphql response
description: Customising graphql error types to send custom error data back to the client using gqlgen.
menu: { main: { parent: 'reference', weight: 10 } }
---

## Returning errors

All resolvers simply return an error to be sent to the user. The assumption is that any error message returned
here is appropriate for end users. If certain messages aren't safe, customise the error presenter.

### Multiple errors

To return multiple errors you can call the `graphql.Error` functions like so:

```go
package foo

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/99designs/gqlgen/graphql"
)

func (r Query) DoThings(ctx context.Context) (bool, error) {
	// Print a formatted string
	graphql.AddErrorf(ctx, "Error %d", 1)

	// Pass an existing error out
	graphql.AddError(ctx, gqlerror.Errorf("zzzzzt"))

	// Or fully customize the error
	graphql.AddError(ctx, &gqlerror.Error{
		Path:       graphql.GetPath(ctx),
		Message:    "A descriptive error message",
		Extensions: map[string]interface{}{
			"code": "10-4",
		},
	})

	// And you can still return an error if you need
	return false, gqlerror.Errorf("BOOM! Headshot")
}
```

They will be returned in the same order in the response, eg:
```json
{
  "data": {
    "todo": null
  },
  "errors": [
    { "message": "Error 1", "path": [ "todo" ] },
    { "message": "zzzzzt", "path": [ "todo" ] },
    { "message": "A descriptive error message", "path": [ "todo" ], "extensions": { "code": "10-4" } },
    { "message": "BOOM! Headshot", "path": [ "todo" ] }
  ]
}
```

## Hooks

### The error presenter

All `errors` returned by resolvers, or from validation, pass through a hook before being displayed to the user.
This hook gives you the ability to customise errors however makes sense in your app.

The default error presenter will capture the resolver path and use the Error() message in the response.

You change this when creating the server:
```go
package bar

import (
	"context"
	"errors"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
)

func main() {
	server := handler.NewDefaultServer(MakeExecutableSchema(resolvers))
	server.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)

		var myErr *MyError
		if errors.As(e, &myErr) {
			err.Message = "Eeek!"
		}

		return err
	})
}

```

This function will be called with the same resolver context that generated it, so you can extract the
current resolver path and whatever other state you might want to notify the client about.


### The panic handler

There is also a panic handler, called whenever a panic happens to gracefully return a message to the user before
stopping parsing. This is a good spot to notify your bug tracker and send a custom message to the user. Any errors
returned from here will also go through the error presenter.

You change this when creating the server:
```go
server := handler.NewDefaultServer(MakeExecutableSchema(resolvers)
server.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
    // notify bug tracker...

		return gqlerror.Errorf("Internal server error!")
})
```

