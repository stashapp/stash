package entityresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/plugin/federation/testdata/entityresolver/generated"
)

func (r *entityResolver) FindHelloByName(ctx context.Context, name string) (*generated.Hello, error) {
	return &generated.Hello{
		Name: name,
	}, nil
}

func (r *entityResolver) FindHelloMultiSingleKeysByKey1AndKey2(ctx context.Context, key1 string, key2 string) (*generated.HelloMultiSingleKeys, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindHelloWithErrorsByName(ctx context.Context, name string) (*generated.HelloWithErrors, error) {
	if name == "inject error" {
		return nil, generated.ErrResolvingHelloWithErrorsByName
	} else if name == "" {
		return nil, generated.ErrEmptyKeyResolvingHelloWithErrorsByName
	}

	return &generated.HelloWithErrors{
		Name: name,
	}, nil
}

func (r *entityResolver) FindManyMultiHelloByNames(ctx context.Context, reps []*generated.MultiHelloByNamesInput) ([]*generated.MultiHello, error) {
	results := []*generated.MultiHello{}

	for _, item := range reps {
		results = append(results, &generated.MultiHello{
			Name: item.Name + " - from multiget",
		})
	}

	return results, nil
}

func (r *entityResolver) FindManyMultiHelloMultipleRequiresByNames(ctx context.Context, reps []*generated.MultiHelloMultipleRequiresByNamesInput) ([]*generated.MultiHelloMultipleRequires, error) {
	results := make([]*generated.MultiHelloMultipleRequires, len(reps))

	for i := range reps {
		results[i] = &generated.MultiHelloMultipleRequires{
			Name: reps[i].Name,
		}
	}

	return results, nil
}

func (r *entityResolver) FindManyMultiHelloRequiresByNames(ctx context.Context, reps []*generated.MultiHelloRequiresByNamesInput) ([]*generated.MultiHelloRequires, error) {
	results := make([]*generated.MultiHelloRequires, len(reps))

	for i := range reps {
		results[i] = &generated.MultiHelloRequires{
			Name: reps[i].Name,
		}
	}

	return results, nil
}

func (r *entityResolver) FindManyMultiHelloWithErrorByNames(ctx context.Context, reps []*generated.MultiHelloWithErrorByNamesInput) ([]*generated.MultiHelloWithError, error) {
	return nil, fmt.Errorf("error resolving MultiHelloWorldWithError")
}

func (r *entityResolver) FindManyMultiPlanetRequiresNestedByNames(ctx context.Context, reps []*generated.MultiPlanetRequiresNestedByNamesInput) ([]*generated.MultiPlanetRequiresNested, error) {
	worlds := map[string]*generated.World{
		"earth": {
			Foo: "A",
		},
		"mars": {
			Foo: "B",
		},
	}

	results := make([]*generated.MultiPlanetRequiresNested, len(reps))

	for i := range reps {
		name := reps[i].Name
		world, ok := worlds[name]
		if !ok {
			return nil, fmt.Errorf("unknown planet: %s", name)
		}

		results[i] = &generated.MultiPlanetRequiresNested{
			Name:  name,
			World: world,
		}
	}

	return results, nil
}

func (r *entityResolver) FindPlanetMultipleRequiresByName(ctx context.Context, name string) (*generated.PlanetMultipleRequires, error) {
	return &generated.PlanetMultipleRequires{Name: name}, nil
}

func (r *entityResolver) FindPlanetRequiresByName(ctx context.Context, name string) (*generated.PlanetRequires, error) {
	return &generated.PlanetRequires{
		Name: name,
	}, nil
}

func (r *entityResolver) FindPlanetRequiresNestedByName(ctx context.Context, name string) (*generated.PlanetRequiresNested, error) {
	worlds := map[string]*generated.World{
		"earth": {
			Foo: "A",
		},
		"mars": {
			Foo: "B",
		},
	}
	world, ok := worlds[name]
	if !ok {
		return nil, fmt.Errorf("unknown planet: %s", name)
	}

	return &generated.PlanetRequiresNested{
		Name:  name,
		World: world,
	}, nil
}

func (r *entityResolver) FindWorldByHelloNameAndFoo(ctx context.Context, helloName string, foo string) (*generated.World, error) {
	return &generated.World{
		Hello: &generated.Hello{
			Name: helloName,
		},
		Foo: foo,
	}, nil
}

func (r *entityResolver) FindWorldNameByName(ctx context.Context, name string) (*generated.WorldName, error) {
	return &generated.WorldName{
		Name: name,
	}, nil
}

func (r *entityResolver) FindWorldWithMultipleKeysByHelloNameAndFoo(ctx context.Context, helloName string, foo string) (*generated.WorldWithMultipleKeys, error) {
	return &generated.WorldWithMultipleKeys{
		Hello: &generated.Hello{
			Name: helloName,
		},
		Foo: foo,
	}, nil
}

func (r *entityResolver) FindWorldWithMultipleKeysByBar(ctx context.Context, bar int) (*generated.WorldWithMultipleKeys, error) {
	return &generated.WorldWithMultipleKeys{
		Bar: bar,
	}, nil
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
