---
title: 'Disabling introspection'
description: Prevent users from introspecting schemas in production.
linkTitle: Introspection
menu: { main: { parent: 'reference', weight: 10 } }
---

One of the best features of GraphQL is it's powerful discoverability and its is automatically included when using `NewDefaultServer`.

## Disable introspection for the whole server

To opt out of introspection globally you should build your own server with only the features you use. For example a simple server that only does POST, and only has introspection in dev could look like:
```go
srv := handler.New(es)

srv.AddTransport(transport.Options{})
srv.AddTransport(transport.POST{})

if os.GetEnv("ENVIRONMENT") == "development" {
    srv.Use(extension.Introspection{})
}
```

## Disabling introspection based on authentication

Introspection can also be enabled on a per-request context basis. For example, you could modify it in a middleware based on user authentication:

```go
srv := handler.NewDefaultServer(es)
srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
    if !userForContext(ctx).IsAdmin {
        graphql.GetOperationContext(ctx).DisableIntrospection = true
    }

    return next(ctx)
})
```
