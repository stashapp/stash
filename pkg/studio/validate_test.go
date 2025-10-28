package studio

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func nameFilter(n string) *models.StudioFilterType {
	return &models.StudioFilterType{
		Name: &models.StringCriterionInput{
			Value:    n,
			Modifier: models.CriterionModifierEquals,
		},
	}
}

func TestValidateName(t *testing.T) {
	db := mocks.NewDatabase()

	const (
		name1   = "name 1"
		newName = "new name"
	)

	existing1 := models.Studio{
		ID:   1,
		Name: name1,
	}

	pp := 1
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	db.Studio.On("Query", testCtx, nameFilter(name1), findFilter).Return([]*models.Studio{&existing1}, 1, nil)
	db.Studio.On("Query", testCtx, mock.Anything, findFilter).Return(nil, 0, nil)

	tests := []struct {
		tName string
		name  string
		want  error
	}{
		{"missing name", "", ErrNameMissing},
		{"new name", newName, nil},
		{"existing name", name1, &NameExistsError{name1}},
	}

	for _, tt := range tests {
		t.Run(tt.tName, func(t *testing.T) {
			got := validateName(testCtx, 0, tt.name, db.Studio)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateUpdateName(t *testing.T) {
	db := mocks.NewDatabase()

	const (
		name1   = "name 1"
		name2   = "name 2"
		newName = "new name"
	)

	existing1 := models.Studio{
		ID:   1,
		Name: name1,
	}
	existing2 := models.Studio{
		ID:   2,
		Name: name2,
	}

	pp := 1
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	db.Studio.On("Query", testCtx, nameFilter(name1), findFilter).Return([]*models.Studio{&existing1}, 1, nil)
	db.Studio.On("Query", testCtx, nameFilter(name2), findFilter).Return([]*models.Studio{&existing2}, 2, nil)
	db.Studio.On("Query", testCtx, mock.Anything, findFilter).Return(nil, 0, nil)

	tests := []struct {
		tName  string
		studio models.Studio
		name   string
		want   error
	}{
		{"missing name", existing1, "", ErrNameMissing},
		{"same name", existing2, name2, nil},
		{"new name", existing1, newName, nil},
	}

	for _, tt := range tests {
		t.Run(tt.tName, func(t *testing.T) {
			got := validateName(testCtx, tt.studio.ID, tt.name, db.Studio)
			assert.Equal(t, tt.want, got)
		})
	}
}
