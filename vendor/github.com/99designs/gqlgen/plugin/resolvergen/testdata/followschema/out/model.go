package customresolver

import "context"

type Resolver struct{}

type QueryResolver interface {
	Resolver(ctx context.Context) (*Resolver, error)
}

type ResolverResolver interface {
	Name(ctx context.Context, obj *Resolver) (string, error)
}
