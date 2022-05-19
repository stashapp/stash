package federation

import (
	"testing"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/stretchr/testify/require"
)

func TestWithEntities(t *testing.T) {
	f, cfg := load(t, "testdata/allthethings/gqlgen.yml")

	require.Equal(t, []string{"ExternalExtension", "Hello", "MoreNesting", "NestedKey", "VeryNestedKey", "World"}, cfg.Schema.Types["_Entity"].Types)

	require.Len(t, cfg.Schema.Types["Entity"].Fields, 6)

	require.Equal(t, "findExternalExtensionByUpc", cfg.Schema.Types["Entity"].Fields[0].Name)
	require.Equal(t, "findHelloByName", cfg.Schema.Types["Entity"].Fields[1].Name)
	// missing on purpose: all @external fields:
	// require.Equal(t, "findMoreNestingByID", cfg.Schema.Types["Entity"].Fields[2].Name)
	require.Equal(t, "findNestedKeyByIDAndHelloName", cfg.Schema.Types["Entity"].Fields[2].Name)
	require.Equal(t, "findVeryNestedKeyByIDAndHelloNameAndWorldFooAndWorldBarAndMoreWorldFoo", cfg.Schema.Types["Entity"].Fields[3].Name)
	require.Equal(t, "findWorldByFoo", cfg.Schema.Types["Entity"].Fields[4].Name)
	require.Equal(t, "findWorldByBar", cfg.Schema.Types["Entity"].Fields[5].Name)

	require.NoError(t, f.MutateConfig(cfg))

	require.Len(t, f.Entities, 6)

	require.Equal(t, "ExternalExtension", f.Entities[0].Name)
	require.Len(t, f.Entities[0].Resolvers, 1)
	require.Len(t, f.Entities[0].Resolvers[0].KeyFields, 1)
	require.Equal(t, "upc", f.Entities[0].Resolvers[0].KeyFields[0].Definition.Name)
	require.Equal(t, "String", f.Entities[0].Resolvers[0].KeyFields[0].Definition.Type.Name())

	require.Equal(t, "Hello", f.Entities[1].Name)
	require.Len(t, f.Entities[1].Resolvers, 1)
	require.Len(t, f.Entities[1].Resolvers[0].KeyFields, 1)
	require.Equal(t, "name", f.Entities[1].Resolvers[0].KeyFields[0].Definition.Name)
	require.Equal(t, "String", f.Entities[1].Resolvers[0].KeyFields[0].Definition.Type.Name())

	require.Equal(t, "MoreNesting", f.Entities[2].Name)
	require.Len(t, f.Entities[2].Resolvers, 0)

	require.Equal(t, "NestedKey", f.Entities[3].Name)
	require.Len(t, f.Entities[3].Resolvers, 1)
	require.Len(t, f.Entities[3].Resolvers[0].KeyFields, 2)
	require.Equal(t, "id", f.Entities[3].Resolvers[0].KeyFields[0].Definition.Name)
	require.Equal(t, "String", f.Entities[3].Resolvers[0].KeyFields[0].Definition.Type.Name())
	require.Equal(t, "helloName", f.Entities[3].Resolvers[0].KeyFields[1].Definition.Name)
	require.Equal(t, "String", f.Entities[3].Resolvers[0].KeyFields[1].Definition.Type.Name())

	require.Equal(t, "VeryNestedKey", f.Entities[4].Name)
	require.Len(t, f.Entities[4].Resolvers, 1)
	require.Len(t, f.Entities[4].Resolvers[0].KeyFields, 5)
	require.Equal(t, "id", f.Entities[4].Resolvers[0].KeyFields[0].Definition.Name)
	require.Equal(t, "String", f.Entities[4].Resolvers[0].KeyFields[0].Definition.Type.Name())
	require.Equal(t, "helloName", f.Entities[4].Resolvers[0].KeyFields[1].Definition.Name)
	require.Equal(t, "String", f.Entities[4].Resolvers[0].KeyFields[1].Definition.Type.Name())
	require.Equal(t, "worldFoo", f.Entities[4].Resolvers[0].KeyFields[2].Definition.Name)
	require.Equal(t, "String", f.Entities[4].Resolvers[0].KeyFields[2].Definition.Type.Name())
	require.Equal(t, "worldBar", f.Entities[4].Resolvers[0].KeyFields[3].Definition.Name)
	require.Equal(t, "Int", f.Entities[4].Resolvers[0].KeyFields[3].Definition.Type.Name())
	require.Equal(t, "moreWorldFoo", f.Entities[4].Resolvers[0].KeyFields[4].Definition.Name)
	require.Equal(t, "String", f.Entities[4].Resolvers[0].KeyFields[4].Definition.Type.Name())

	require.Len(t, f.Entities[4].Requires, 2)
	require.Equal(t, f.Entities[4].Requires[0].Name, "id")
	require.Equal(t, f.Entities[4].Requires[1].Name, "helloSecondary")

	require.Equal(t, "World", f.Entities[5].Name)
	require.Len(t, f.Entities[5].Resolvers, 2)
	require.Len(t, f.Entities[5].Resolvers[0].KeyFields, 1)
	require.Equal(t, "foo", f.Entities[5].Resolvers[0].KeyFields[0].Definition.Name)
	require.Equal(t, "String", f.Entities[5].Resolvers[0].KeyFields[0].Definition.Type.Name())
	require.Len(t, f.Entities[5].Resolvers[1].KeyFields, 1)
	require.Equal(t, "bar", f.Entities[5].Resolvers[1].KeyFields[0].Definition.Name)
	require.Equal(t, "Int", f.Entities[5].Resolvers[1].KeyFields[0].Definition.Type.Name())
}

func TestNoEntities(t *testing.T) {
	f, cfg := load(t, "testdata/entities/nokey.yml")

	err := f.MutateConfig(cfg)
	require.NoError(t, err)
	require.Len(t, f.Entities, 0)
}

func TestInterfaceKeyDirective(t *testing.T) {
	f, cfg := load(t, "testdata/interfaces/key.yml")

	err := f.MutateConfig(cfg)
	require.NoError(t, err)
	require.Len(t, f.Entities, 0)
}

func TestInterfaceExtendsDirective(t *testing.T) {
	require.Panics(t, func() {
		load(t, "testdata/interfaces/extends.yml")
	})
}

func TestCodeGeneration(t *testing.T) {
	f, cfg := load(t, "testdata/allthethings/gqlgen.yml")

	require.Len(t, cfg.Schema.Types["_Entity"].Types, 6)
	require.Len(t, f.Entities, 6)

	require.NoError(t, f.MutateConfig(cfg))

	data, err := codegen.BuildData(cfg)
	if err != nil {
		panic(err)
	}
	require.NoError(t, f.GenerateCode(data))
}

func TestInjectSourceLate(t *testing.T) {
	_, cfg := load(t, "testdata/allthethings/gqlgen.yml")
	entityGraphqlGenerated := false
	for _, source := range cfg.Sources {
		if source.Name != "federation/entity.graphql" {
			continue
		}
		entityGraphqlGenerated = true
		require.Contains(t, source.Input, "union _Entity")
		require.Contains(t, source.Input, "type _Service {")
		require.Contains(t, source.Input, "extend type Query {")
		require.Contains(t, source.Input, "_entities(representations: [_Any!]!): [_Entity]!")
		require.Contains(t, source.Input, "_service: _Service!")
	}
	require.True(t, entityGraphqlGenerated)

	_, cfg = load(t, "testdata/entities/nokey.yml")
	entityGraphqlGenerated = false
	for _, source := range cfg.Sources {
		if source.Name != "federation/entity.graphql" {
			continue
		}
		entityGraphqlGenerated = true
		require.NotContains(t, source.Input, "union _Entity")
		require.Contains(t, source.Input, "type _Service {")
		require.Contains(t, source.Input, "extend type Query {")
		require.NotContains(t, source.Input, "_entities(representations: [_Any!]!): [_Entity]!")
		require.Contains(t, source.Input, "_service: _Service!")
	}
	require.True(t, entityGraphqlGenerated)

	_, cfg = load(t, "testdata/schema/customquerytype.yml")
	for _, source := range cfg.Sources {
		if source.Name != "federation/entity.graphql" {
			continue
		}
		require.Contains(t, source.Input, "extend type CustomQuery {")
	}
}

func load(t *testing.T, name string) (*federation, *config.Config) {
	t.Helper()

	cfg, err := config.LoadConfig(name)
	require.NoError(t, err)

	f := &federation{}
	cfg.Sources = append(cfg.Sources, f.InjectSourceEarly())
	require.NoError(t, cfg.LoadSchema())

	if src := f.InjectSourceLate(cfg.Schema); src != nil {
		cfg.Sources = append(cfg.Sources, src)
	}
	require.NoError(t, cfg.LoadSchema())

	require.NoError(t, cfg.Init())
	return f, cfg
}
