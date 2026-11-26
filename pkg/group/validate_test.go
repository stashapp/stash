package group

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateAliases(t *testing.T) {
	const name = "name"

	tests := []struct {
		name    string
		aliases models.RelatedStrings
		want    error
	}{
		{"not loaded", models.RelatedStrings{}, nil},
		{"empty aliases", models.NewRelatedStrings([]string{}), nil},
		{"valid multiple aliases", models.NewRelatedStrings([]string{"alias one", "alias two"}), nil},
		{"empty string alias", models.NewRelatedStrings([]string{""}), ErrEmptyAlias},
		{"whitespace-only alias", models.NewRelatedStrings([]string{"   "}), ErrEmptyAlias},
		{"duplicate aliases", models.NewRelatedStrings([]string{"alias", "ALIAS"}), &DuplicateAliasError{"ALIAS"}},
		{"alias same as name", models.NewRelatedStrings([]string{"Name"}), &DuplicateAliasError{"Name"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ValidateAliases(name, tt.aliases))
		})
	}
}

func TestValidateUpdateAliases(t *testing.T) {
	existing := models.Group{
		Name:    "name",
		Aliases: models.NewRelatedStrings([]string{"alias1", "alias2"}),
	}

	tests := []struct {
		name    string
		newName models.OptionalString
		aliases *models.UpdateStrings
		want    error
	}{
		{"nothing set", models.OptionalString{}, nil, nil},
		{"name same as existing alias", models.NewOptionalString("alias1"), nil, &DuplicateAliasError{"alias1"}},
		{"name not in aliases", models.NewOptionalString("new"), nil, nil},
		{"add name as alias", models.OptionalString{}, &models.UpdateStrings{
			Values: []string{"NAME"},
			Mode:   models.RelationshipUpdateModeAdd,
		}, &DuplicateAliasError{"NAME"}},
		{"set aliases", models.OptionalString{}, &models.UpdateStrings{
			Values: []string{"alias3"},
			Mode:   models.RelationshipUpdateModeSet,
		}, nil},
		{"set alias to new name", models.NewOptionalString("alias3"), &models.UpdateStrings{
			Values: []string{"alias3"},
			Mode:   models.RelationshipUpdateModeSet,
		}, &DuplicateAliasError{"alias3"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ValidateUpdateAliases(existing, tt.newName, tt.aliases))
		})
	}
}
