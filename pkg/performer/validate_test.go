package performer

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func nameFilter(n string) *models.PerformerFilterType {
	return &models.PerformerFilterType{
		Name: &models.StringCriterionInput{
			Value:    n,
			Modifier: models.CriterionModifierEquals,
		},
		Disambiguation: &models.StringCriterionInput{
			Modifier: models.CriterionModifierIsNull,
		},
	}
}

func disambigFilter(n string, d string) *models.PerformerFilterType {
	return &models.PerformerFilterType{
		Name: &models.StringCriterionInput{
			Value:    n,
			Modifier: models.CriterionModifierEquals,
		},
		Disambiguation: &models.StringCriterionInput{
			Value:    d,
			Modifier: models.CriterionModifierEquals,
		},
	}
}

func TestValidateName(t *testing.T) {
	db := mocks.NewDatabase()

	const (
		name1       = "name 1"
		name2       = "name 2"
		disambig    = "disambiguation"
		newName     = "new name"
		newDisambig = "new disambiguation"
	)

	pp := 1
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	db.Performer.On("QueryCount", testCtx, nameFilter(name1), findFilter).Return(1, nil)
	db.Performer.On("QueryCount", testCtx, nameFilter(name2), findFilter).Return(1, nil)
	db.Performer.On("QueryCount", testCtx, disambigFilter(name2, disambig), findFilter).Return(1, nil)
	db.Performer.On("QueryCount", testCtx, mock.Anything, findFilter).Return(0, nil)

	tests := []struct {
		tName    string
		name     string
		disambig string
		want     error
	}{
		{"missing name", "", newDisambig, ErrNameMissing},
		{"new name", newName, "", nil},
		{"new name new disambig", newName, newDisambig, nil},
		{"new name existing disambig", newName, disambig, nil},
		{"existing name", name1, "", &NameExistsError{name1, ""}},
		{"existing name new disambig", name1, newDisambig, nil},
		{"existing name existing disambig", name1, disambig, nil},
		{"existing name and disambig", name2, disambig, &NameExistsError{name2, disambig}},
	}

	for _, tt := range tests {
		t.Run(tt.tName, func(t *testing.T) {
			got := ValidateName(testCtx, tt.name, tt.disambig, db.Performer)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateUpdateName(t *testing.T) {
	db := mocks.NewDatabase()

	const (
		name1       = "name 1"
		name2       = "name 2"
		disambig1   = "disambiguation 1"
		disambig2   = "disambiguation 2"
		newName     = "new name"
		newDisambig = "new disambiguation"
	)

	osUnset := models.OptionalString{}
	osNull := models.OptionalString{Set: true, Null: true}
	osName1 := models.NewOptionalString(name1)
	osName2 := models.NewOptionalString(name2)
	osDisambig1 := models.NewOptionalString(disambig1)
	osDisambig2 := models.NewOptionalString(disambig2)
	osNewName := models.NewOptionalString(newName)
	osNewDisambig := models.NewOptionalString(newDisambig)

	existing1 := models.Performer{
		ID:   1,
		Name: name1,
	}
	existing2 := models.Performer{
		ID:             2,
		Name:           name2,
		Disambiguation: disambig1,
	}
	existing3 := models.Performer{
		ID:             3,
		Name:           name2,
		Disambiguation: disambig2,
	}

	pp := 2
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	db.Performer.On("Query", testCtx, nameFilter(name1), findFilter).Return([]*models.Performer{&existing1}, 1, nil)
	db.Performer.On("Query", testCtx, nameFilter(name2), findFilter).Return([]*models.Performer{&existing2, &existing3}, 2, nil)
	db.Performer.On("Query", testCtx, disambigFilter(name2, disambig1), findFilter).Return([]*models.Performer{&existing2}, 1, nil)
	db.Performer.On("Query", testCtx, disambigFilter(name2, disambig2), findFilter).Return([]*models.Performer{&existing3}, 1, nil)
	db.Performer.On("Query", testCtx, mock.Anything, findFilter).Return(nil, 0, nil)

	tests := []struct {
		tName     string
		performer models.Performer
		name      models.OptionalString
		disambig  models.OptionalString
		want      error
	}{
		{"missing name", existing1, osNull, osUnset, ErrNameMissing},
		{"same name", existing3, osName2, osUnset, nil},
		{"same disambig", existing2, osUnset, osDisambig1, nil},
		{"same name same disambig", existing2, osName2, osDisambig1, nil},
		{"new name", existing1, osNewName, osUnset, nil},
		{"new disambig", existing1, osUnset, osNewDisambig, nil},
		{"new name new disambig", existing1, osNewName, osNewDisambig, nil},
		{"remove disambig", existing3, osUnset, osNull, &NameExistsError{name2, ""}},
		{"existing name keep disambig", existing3, osName1, osUnset, nil},
		{"existing name remove disambig", existing3, osName1, osNull, &NameExistsError{name1, ""}},
		{"existing disambig", existing2, osUnset, osDisambig2, &NameExistsError{name2, disambig2}},
		{"existing name and disambig", existing1, osName2, osDisambig1, &NameExistsError{name2, disambig1}},
	}

	for _, tt := range tests {
		t.Run(tt.tName, func(t *testing.T) {
			got := ValidateUpdateName(testCtx, tt.performer, tt.name, tt.disambig, db.Performer)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateAliases(t *testing.T) {
	const (
		name1  = "name 1"
		name1U = "NAME 1"
		name2  = "name 2"
		name3  = "name 3"
		name4  = "name 4"
	)

	tests := []struct {
		tName   string
		name    string
		aliases []string
		want    error
	}{
		{"no aliases", name1, nil, nil},
		{"valid aliases", name2, []string{name3, name4}, nil},
		{"duplicate alias", name1, []string{name2, name3, name2}, &DuplicateAliasError{name2}},
		{"duplicate name", name4, []string{name4, name3}, &DuplicateAliasError{name4}},
		{"duplicate alias caps", name2, []string{name1, name1U}, &DuplicateAliasError{name1U}},
		{"duplicate name caps", name1U, []string{name1}, &DuplicateAliasError{name1}},
	}

	for _, tt := range tests {
		t.Run(tt.tName, func(t *testing.T) {
			got := ValidateAliases(tt.name, models.NewRelatedStrings(tt.aliases))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateUpdateAliases(t *testing.T) {
	const (
		name1  = "name 1"
		name1U = "NAME 1"
		name2  = "name 2"
		name3  = "name 3"
		name4  = "name 4"
	)

	existing := models.Performer{
		Name:    name1,
		Aliases: models.NewRelatedStrings([]string{name2}),
	}

	osUnset := models.OptionalString{}
	os1 := models.NewOptionalString(name1)
	os2 := models.NewOptionalString(name2)
	os3 := models.NewOptionalString(name3)
	os4 := models.NewOptionalString(name4)

	tests := []struct {
		tName   string
		name    models.OptionalString
		aliases []string
		want    error
	}{
		{"both unset", osUnset, nil, nil},
		{"invalid name set", os2, nil, &DuplicateAliasError{name2}},
		{"valid name set", os3, nil, nil},
		{"valid aliases empty", os1, []string{}, nil},
		{"invalid aliases set", osUnset, []string{name1U}, &DuplicateAliasError{name1U}},
		{"valid aliases set", osUnset, []string{name3, name2}, nil},
		{"invalid both set", os4, []string{name4}, &DuplicateAliasError{name4}},
		{"valid both set", os2, []string{name1}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.tName, func(t *testing.T) {
			var aliases *models.UpdateStrings
			if tt.aliases != nil {
				aliases = &models.UpdateStrings{
					Values: tt.aliases,
					Mode:   models.RelationshipUpdateModeSet,
				}
			}
			got := ValidateUpdateAliases(existing, tt.name, aliases)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateDeathDate(t *testing.T) {
	date1, _ := models.ParseDate("2001-01-01")
	date2, _ := models.ParseDate("2002-01-01")
	date3, _ := models.ParseDate("2003-01-01")
	date4, _ := models.ParseDate("2004-01-01")

	tests := []struct {
		name      string
		birthdate *models.Date
		deathdate *models.Date
		want      error
	}{
		{"both nil", nil, nil, nil},
		{"birthdate nil", nil, &date1, nil},
		{"deathdate nil", nil, &date2, nil},
		{"valid", &date3, &date4, nil},
		{"invalid", &date3, &date2, &DeathDateError{date3, date2}},
		{"same date", &date1, &date1, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateDeathDate(tt.birthdate, tt.deathdate)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateUpdateDeathDate(t *testing.T) {
	date1, _ := models.ParseDate("2001-01-01")
	date2, _ := models.ParseDate("2002-01-01")
	date3, _ := models.ParseDate("2003-01-01")
	date4, _ := models.ParseDate("2004-01-01")

	existing := models.Performer{
		Birthdate: &date2,
		DeathDate: &date3,
	}

	odUnset := models.OptionalDate{}
	odNull := models.OptionalDate{Set: true, Null: true}
	od1 := models.NewOptionalDate(date1)
	od2 := models.NewOptionalDate(date2)
	od3 := models.NewOptionalDate(date3)
	od4 := models.NewOptionalDate(date4)

	tests := []struct {
		name      string
		birthdate models.OptionalDate
		deathdate models.OptionalDate
		want      error
	}{
		{"both unset", odUnset, odUnset, nil},
		{"invalid birthdate set", od4, odUnset, &DeathDateError{date4, date3}},
		{"valid birthdate set", od1, odUnset, nil},
		{"valid birthdate set null", odNull, odUnset, nil},
		{"invalid deathdate set", odUnset, od1, &DeathDateError{date2, date1}},
		{"valid deathdate set", odUnset, od4, nil},
		{"valid deathdate set null", odUnset, odNull, nil},
		{"invalid both set", od3, od2, &DeathDateError{date3, date2}},
		{"valid both set", od2, od3, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateUpdateDeathDate(existing, tt.birthdate, tt.deathdate)
			assert.Equal(t, tt.want, got)
		})
	}
}
