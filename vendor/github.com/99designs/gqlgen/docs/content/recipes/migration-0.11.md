---
title: "Migrating to 0.11"
description: Changes in gqlgen 0.11
linkTitle: Migrating to 0.11
menu: { main: { parent: 'recipes' } }
---

## Updated gqlparser

gqlparser had a breaking change, if you have any references to it in your project your going to need to update
them from `github.com/vektah/gqlparser` to `github.com/vektah/gqlparser/v2`.

```bash
sed -i 's/github.com\/vektah\/gqlparser/github.com\/vektah\/gqlparser\/v2/' $(find -name '*.go')
```

## Handler Refactor

The handler package has grown organically for a long time, 0.11 is a large cleanup of the handler package to make it
more modular and easier to maintain once we get to 1.0.


### Transports

Transports are the first thing that run, they handle decoding the incoming http request, and encoding the graphql
response. Supported transports are:

 - GET
 - JSON POST
 - Multipart form
 - Websockets

new usage looks like this
```go
srv := New(es)

srv.AddTransport(transport.Websocket{
	KeepAlivePingInterval: 10 * time.Second,
})
srv.AddTransport(transport.Options{})
srv.AddTransport(transport.GET{})
srv.AddTransport(transport.POST{})
srv.AddTransport(transport.MultipartForm{})
```

### New handler extension API

The core of this changes the handler package to be a set of composable extensions. The extensions implement a set of optional interfaces:

 - **OperationParameterMutator** runs before creating a OperationContext (formerly RequestContext). allows manipulating the raw query before parsing.
 - **OperationContextMutator** runs after creating the OperationContext, but before executing the root resolver.
 - **OperationInterceptor** runs for each incoming query after parsing and validation, for basic requests the writer will be invoked once, for subscriptions it will be invoked multiple times.
 - **ResponseInterceptor** runs around each graphql operation response. This can be called many times for a single operation the case of subscriptions.
 - **FieldInterceptor** runs around each field

![Anatomy of a request@2x](https://user-images.githubusercontent.com/2247982/68181515-c8a27c00-ffeb-11e9-86f6-1673e7179ecb.png)

Users of an extension should not need to know which extension points are being used by a given extension, they are added to the server simply by calling `Use(extension)`.

There are a few convenience methods for defining middleware inline, instead of creating an extension

```go
srv := handler.New(es)
srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	// this function will be called around every field. next() will evaluate the field and return
	// its computed value.
	return next(ctx)
})
srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	// This function will be called around every operation, next() will return a function that when
	// called will evaluate one response. Eventually next will return nil, signalling there are no
	// more results to be returned by the server.
	return next(ctx)
})
srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	// This function will be called around each response in the operation. next() will evaluate
	// and return a single response.
	return next(ctx)
})
```


Some of the features supported out of the box by handler extensions:

 - APQ
 - Query Complexity
 - Error Presenters and Recover Func
 - Introspection Query support
 - Query AST cache
 - Tracing API

They can be `Use`'d like this:

```go
srv := handler.New(es)
srv.Use(extension.Introspection{})
srv.Use(extension.AutomaticPersistedQuery{
	Cache: lru.New(100),
})
srv.Use(apollotracing.Tracer{})
```

### Default server

We provide a set of default extensions and transports if you aren't ready to customize them yet. Simply:
```go
handler.NewDefaultServer(es)
```

### More consistent naming

As part of cleaning up the names the RequestContext has been renamed to OperationContext, as there can be multiple created during the lifecycle of a request. A new ResponseContext has also been created and error handling has been moved here. This allows each response in a subscription to have its own errors. I'm not sure what bugs this might have been causing before...

### Removal of tracing

Many of the old interfaces collapse down into just a few extension points:

![Anatomy of a request](/request_anatomy.png)

The tracing interface has also been removed, tracing stats are now measured in core (eg time to parse query) and made available on the operation/response contexts. Much of the old interface was designed so that users of a tracer dont need to know which extension points it was listening to, the new handler extensions have the same goal.

### Backward compatibility

There is a backwards compatibility layer that keeps most of the original interface in place. There are a few places where BC is known to be broken:

 - ResponseMiddleware: The signature used to be `func(ctx context.Context, next func(ctx context.Context) []byte) []byte` and is now `func(ctx context.Context) *Response`. We could maintain BC by marshalling to json before and after, but the change is pretty easy to make and is likely to cause less issues.
 - The Tracer interface has been removed, any tracers will need to be reimplemented against the new extension interface.

## New resolver layout

0.11 also added a new way to generate and layout resolvers on disk. We used to only generate resolver implementations
whenever the file didnt exist. This behaviour is still there for those that are already used to it, However there is a
new mode you can turn on in config:

```yaml
resolver:
  layout: follow-schema
  dir: graph
```

This tells gqlgen to generate resolvers next to the schema file that declared the graphql field, which looks like this:

![follow-schema layout](/schema_layout.png)
