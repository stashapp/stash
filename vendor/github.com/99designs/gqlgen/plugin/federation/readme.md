# Federation plugin

Add support for graphql federation in your graphql Go server!

TODO(miguel): add details.

# Tests
There are several different tests. Some will process the configuration file directly.  You can see those in the `federation_test.go`.  There are also tests for entity resolvers, which will simulate requests from a federation server like Apollo Federation.

Running entity resolver tests.
1. Go to `plugin/federation`
2. Run the command `go generate`
3. Run the tests with `go test ./...`.

# Architecture

TODO(miguel): add details.

# Entity resolvers - GetMany entities

The federation plugin implements `GetMany` semantics in which entity resolvers get the entire list of representations that need to be resolved. This functionality is currently optin tho, and to enable it you need to specify the directive `@entityResolver` in the federated entity you want this feature for.  E.g.

```
directive @entityResolver(multi: Boolean) on OBJECT

type MultiHello @key(fields: "name") @entityResolver(multi: true) {
    name: String!
}
```

That allows the federation plugin to generate `GetMany` resolver function that can take a list of representations to be resolved.

From that entity type, the resolver function would be

```
func (r *entityResolver) FindManyMultiHellosByName(ctx context.Context, reps []*generated.ManyMultiHellosByNameInput) ([]*generated.MultiHello, error) {
  /// <Your code to resolve the list of items>
}
```
