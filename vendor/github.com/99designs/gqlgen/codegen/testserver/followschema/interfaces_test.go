package followschema

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

func TestInterfaces(t *testing.T) {
	t.Run("slices of interfaces are not pointers", func(t *testing.T) {
		field, ok := reflect.TypeOf((*QueryResolver)(nil)).Elem().MethodByName("Shapes")
		require.True(t, ok)
		require.Equal(t, "[]followschema.Shape", field.Type.Out(0).String())
	})

	t.Run("models returning interfaces", func(t *testing.T) {
		resolvers := &Stub{}
		resolvers.QueryResolver.Node = func(ctx context.Context) (node Node, err error) {
			return &ConcreteNodeA{
				ID:   "1234",
				Name: "asdf",
				child: &ConcreteNodeA{
					ID:    "5678",
					Name:  "hjkl",
					child: nil,
				},
			}, nil
		}

		srv := handler.NewDefaultServer(
			NewExecutableSchema(Config{
				Resolvers: resolvers,
			}),
		)

		c := client.New(srv)

		var resp struct {
			Node struct {
				ID    string
				Child struct {
					ID string
				}
			}
		}
		c.MustPost(`{ node { id, child { id } } }`, &resp)
		require.Equal(t, "1234", resp.Node.ID)
		require.Equal(t, "5678", resp.Node.Child.ID)
	})

	t.Run("interfaces can be nil", func(t *testing.T) {
		resolvers := &Stub{}
		resolvers.QueryResolver.NoShape = func(ctx context.Context) (shapes Shape, e error) {
			return nil, nil
		}

		srv := handler.NewDefaultServer(
			NewExecutableSchema(Config{
				Resolvers: resolvers,
				Directives: DirectiveRoot{
					MakeNil: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
						return nil, nil
					},
				},
			}),
		)

		c := client.New(srv)

		var resp interface{}
		c.MustPost(`{ noShape { area } }`, &resp)
	})

	t.Run("interfaces can be typed nil", func(t *testing.T) {
		resolvers := &Stub{}
		resolvers.QueryResolver.NoShapeTypedNil = func(ctx context.Context) (shapes Shape, e error) {
			panic("should not be called")
		}

		srv := handler.NewDefaultServer(
			NewExecutableSchema(Config{
				Resolvers: resolvers,
				Directives: DirectiveRoot{
					MakeTypedNil: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
						var circle *Circle
						return circle, nil
					},
				},
			}),
		)

		c := client.New(srv)

		var resp interface{}
		c.MustPost(`{ noShapeTypedNil { area } }`, &resp)
	})

	t.Run("interfaces can be nil (test with code-generated resolver)", func(t *testing.T) {
		resolvers := &Stub{}
		resolvers.QueryResolver.Animal = func(ctx context.Context) (animal Animal, e error) {
			panic("should not be called")
		}

		srv := handler.NewDefaultServer(
			NewExecutableSchema(Config{
				Resolvers: resolvers,
				Directives: DirectiveRoot{
					MakeTypedNil: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
						var dog *Dog // return a typed nil, not just nil
						return dog, nil
					},
				},
			}),
		)

		c := client.New(srv)

		var resp interface{}
		c.MustPost(`{ animal { species } }`, &resp)
	})

	t.Run("can bind to interfaces even when the graphql is not", func(t *testing.T) {
		resolvers := &Stub{}
		resolvers.BackedByInterfaceResolver.ID = func(ctx context.Context, obj BackedByInterface) (s string, err error) {
			return "ID:" + obj.ThisShouldBind(), nil
		}
		resolvers.QueryResolver.NotAnInterface = func(ctx context.Context) (byInterface BackedByInterface, err error) {
			return &BackedByInterfaceImpl{
				Value: "A",
				Error: nil,
			}, nil
		}

		c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

		var resp struct {
			NotAnInterface struct {
				ID                      string
				ThisShouldBind          string
				ThisShouldBindWithError string
			}
		}
		c.MustPost(`{ notAnInterface { id, thisShouldBind, thisShouldBindWithError } }`, &resp)
		require.Equal(t, "ID:A", resp.NotAnInterface.ID)
		require.Equal(t, "A", resp.NotAnInterface.ThisShouldBind)
		require.Equal(t, "A", resp.NotAnInterface.ThisShouldBindWithError)
	})

	t.Run("can return errors from interface funcs", func(t *testing.T) {
		resolvers := &Stub{}
		resolvers.BackedByInterfaceResolver.ID = func(ctx context.Context, obj BackedByInterface) (s string, err error) {
			return "ID:" + obj.ThisShouldBind(), nil
		}
		resolvers.QueryResolver.NotAnInterface = func(ctx context.Context) (byInterface BackedByInterface, err error) {
			return &BackedByInterfaceImpl{
				Value: "A",
				Error: fmt.Errorf("boom"),
			}, nil
		}

		c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

		var resp struct {
			NotAnInterface struct {
				ID                      string
				ThisShouldBind          string
				ThisShouldBindWithError string
			}
		}
		err := c.Post(`{ notAnInterface { id, thisShouldBind, thisShouldBindWithError } }`, &resp)
		require.EqualError(t, err, `[{"message":"boom","path":["notAnInterface","thisShouldBindWithError"]}]`)
	})

	t.Run("interfaces can implement other interfaces", func(t *testing.T) {
		resolvers := &Stub{}
		resolvers.QueryResolver.Node = func(ctx context.Context) (node Node, err error) {
			return ConcreteNodeInterfaceImplementor{}, nil
		}

		c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

		var resp struct {
			Node struct {
				ID    string
				Child struct {
					ID string
				}
			}
		}
		c.MustPost(`{ node { id, child { id } } }`, &resp)
		require.Equal(t, "CNII", resp.Node.ID)
		require.Equal(t, "Child", resp.Node.Child.ID)
	})

	t.Run("interface implementors should return merged base fields", func(t *testing.T) {
		resolvers := &Stub{}
		resolvers.QueryResolver.Shapes = func(ctx context.Context) (shapes []Shape, err error) {
			return []Shape{
				&Rectangle{
					Coordinates: Coordinates{
						X: -1,
						Y: -1,
					},
				},
				&Circle{
					Coordinates: Coordinates{
						X: 1,
						Y: 1,
					},
				},
			}, nil
		}

		c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))
		var resp struct {
			Shapes []struct {
				Coordinates struct {
					X float64
					Y float64
				}
			}
		}

		c.MustPost(`
			{
				shapes {
					coordinates {
						x
					}
					... on Rectangle {
						coordinates {
							x
						}
					}
					... on Circle {
						coordinates {
							y
						}
					}
				}
			}
		`, &resp)

		require.Equal(t, 2, len(resp.Shapes))
		require.Equal(t, float64(-1), resp.Shapes[0].Coordinates.X)
		require.Equal(t, float64(0), resp.Shapes[0].Coordinates.Y)
		require.Equal(t, float64(1), resp.Shapes[1].Coordinates.X)
		require.Equal(t, float64(1), resp.Shapes[1].Coordinates.Y)
	})
}
