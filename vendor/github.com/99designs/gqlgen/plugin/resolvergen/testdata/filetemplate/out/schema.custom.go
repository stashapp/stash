package customresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	customresolver "github.com/99designs/gqlgen/plugin/resolvergen/testdata/singlefile/out"
)

func (r *queryCustomResolverType) Resolver(ctx context.Context) (*customresolver.Resolver, error) {
	// CustomerResolverType.Resolver implementation
	return nil, nil
}

func (r *resolverCustomResolverType) Name(ctx context.Context, obj *customresolver.Resolver) (string, error) {
	// CustomerResolverType.Name implementation
	return "", nil
}

// Query returns customresolver.QueryResolver implementation.
func (r *CustomResolverType) Query() customresolver.QueryResolver { return &queryCustomResolverType{r} }

// Resolver returns customresolver.ResolverResolver implementation.
func (r *CustomResolverType) Resolver() customresolver.ResolverResolver {
	return &resolverCustomResolverType{r}
}

type queryCustomResolverType struct{ *CustomResolverType }
type resolverCustomResolverType struct{ *CustomResolverType }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func AUserHelperFunction() {
	// AUserHelperFunction implementation
}
