package tag

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

type tagNameFinderMock struct {
	existingTags []*models.Tag
}

func (m tagNameFinderMock) FindByName(ctx context.Context, name string, nocase bool) (*models.Tag, error) {
	for _, n := range m.existingTags {
		if n.Name == name {
			return n, nil
		}
	}

	return nil, nil
}

func (m tagNameFinderMock) FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Tag, error) {
	panic("not implemented")
}

func (m tagNameFinderMock) FindByAlias(ctx context.Context, alias string, nocase bool) (*models.Tag, error) {
	for _, n := range m.existingTags {
		for _, a := range n.Aliases.List() {
			if a == alias {
				return n, nil
			}
		}
	}

	return nil, nil
}

func TestEnsureAliasesUnique(t *testing.T) {
	const (
		name1    = "name 1"
		name2    = "name 2"
		name3    = "name 3"
		alias1   = "alias 1"
		newAlias = "new alias"
	)

	tagMock := tagNameFinderMock{
		existingTags: []*models.Tag{
			{Name: name1, Aliases: models.NewRelatedStrings([]string{})},
			{Name: name2, Aliases: models.NewRelatedStrings([]string{})},
			{Name: name3, Aliases: models.NewRelatedStrings([]string{newAlias})},
		},
	}

	tests := []struct {
		tName   string
		id      int
		aliases []string
		want    error
	}{
		{"valid alias", 1, []string{alias1}, nil},
		{"alias duplicates other name", 1, []string{name2}, &NameExistsError{name2}},
		{"alias duplicates other alias", 1, []string{newAlias}, &NameUsedByAliasError{newAlias, tagMock.existingTags[2].Name}},
	}

	for _, tt := range tests {
		t.Run(tt.tName, func(t *testing.T) {
			got := EnsureAliasesUnique(testCtx, tt.id, tt.aliases, tagMock)
			assert.Equal(t, tt.want, got)
		})
	}
}
