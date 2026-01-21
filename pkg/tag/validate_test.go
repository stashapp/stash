package tag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

func nameFilter(n string) *models.TagFilterType {
	return &models.TagFilterType{
		Name: &models.StringCriterionInput{
			Value:    n,
			Modifier: models.CriterionModifierEquals,
		},
	}
}

func aliasFilter(n string) *models.TagFilterType {
	return &models.TagFilterType{
		Aliases: &models.StringCriterionInput{
			Value:    n,
			Modifier: models.CriterionModifierEquals,
		},
	}
}

func TestEnsureAliasesUnique(t *testing.T) {
	db := mocks.NewDatabase()

	const (
		name1    = "name 1"
		name2    = "name 2"
		alias1   = "alias 1"
		newAlias = "new alias"
	)

	existing2 := models.Tag{
		ID:   2,
		Name: name2,
	}

	pp := 1
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	// name1 matches existing1 name - ok
	// EnsureAliasesUnique calls EnsureTagNameUnique.
	// EnsureTagNameUnique calls ByName then ByAlias.

	// Case 1: valid alias
	// ByName "alias 1" -> nil
	// ByAlias "alias 1" -> nil
	db.Tag.On("Query", testCtx, nameFilter(alias1), findFilter).Return(nil, 0, nil)
	db.Tag.On("Query", testCtx, aliasFilter(alias1), findFilter).Return(nil, 0, nil)

	// Case 2: alias duplicates existing2 name
	// ByName "name 2" -> existing2
	db.Tag.On("Query", testCtx, nameFilter(name2), findFilter).Return([]*models.Tag{&existing2}, 1, nil)

	// Case 3: alias duplicates existing2 alias
	// ByName "new alias" -> nil
	// ByAlias "new alias" -> existing2
	db.Tag.On("Query", testCtx, nameFilter(newAlias), findFilter).Return(nil, 0, nil)
	db.Tag.On("Query", testCtx, aliasFilter(newAlias), findFilter).Return([]*models.Tag{&existing2}, 1, nil)

	tests := []struct {
		tName   string
		id      int
		aliases []string
		want    error
	}{
		{"valid alias", 1, []string{alias1}, nil},
		{"alias duplicates other name", 1, []string{name2}, &NameExistsError{name2}},
		{"alias duplicates other alias", 1, []string{newAlias}, &NameUsedByAliasError{newAlias, existing2.Name}},
	}

	for _, tt := range tests {
		t.Run(tt.tName, func(t *testing.T) {
			got := EnsureAliasesUnique(testCtx, tt.id, tt.aliases, db.Tag)
			assert.Equal(t, tt.want, got)
		})
	}
}
