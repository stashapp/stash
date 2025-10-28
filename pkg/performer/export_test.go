package performer

import (
	"errors"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

const (
	performerID       = 1
	noImageID         = 2
	errImageID        = 3
	customFieldsID    = 4
	errCustomFieldsID = 5
)

const (
	performerName  = "testPerformer"
	disambiguation = "disambiguation"
	url            = "url"
	careerLength   = "careerLength"
	country        = "country"
	ethnicity      = "ethnicity"
	eyeColor       = "eyeColor"
	fakeTits       = "fakeTits"
	instagram      = "instagram"
	measurements   = "measurements"
	piercings      = "piercings"
	tattoos        = "tattoos"
	twitter        = "twitter"
	details        = "details"
	hairColor      = "hairColor"

	autoTagIgnored = true
)

var (
	genderEnum      = models.GenderEnumFemale
	gender          = genderEnum.String()
	aliases         = []string{"alias1", "alias2"}
	rating          = 5
	height          = 123
	weight          = 60
	penisLength     = 1.23
	circumcisedEnum = models.CircumisedEnumCut
	circumcised     = circumcisedEnum.String()

	emptyCustomFields = make(map[string]interface{})
	customFields      = map[string]interface{}{
		"customField1": "customValue1",
	}
)

var imageBytes = []byte("imageBytes")

var stashID = models.StashID{
	StashID:  "StashID",
	Endpoint: "Endpoint",
}
var stashIDs = []models.StashID{
	stashID,
}

const image = "aW1hZ2VCeXRlcw=="

var birthDate, _ = models.ParseDate("2001-01-01")
var deathDate, _ = models.ParseDate("2021-02-02")

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.Local)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.Local)
)

func createFullPerformer(id int, name string) *models.Performer {
	return &models.Performer{
		ID:             id,
		Name:           name,
		Disambiguation: disambiguation,
		URLs:           models.NewRelatedStrings([]string{url, twitter, instagram}),
		Aliases:        models.NewRelatedStrings(aliases),
		Birthdate:      &birthDate,
		CareerLength:   careerLength,
		Country:        country,
		Ethnicity:      ethnicity,
		EyeColor:       eyeColor,
		FakeTits:       fakeTits,
		PenisLength:    &penisLength,
		Circumcised:    &circumcisedEnum,
		Favorite:       true,
		Gender:         &genderEnum,
		Height:         &height,
		Measurements:   measurements,
		Piercings:      piercings,
		Tattoos:        tattoos,
		CreatedAt:      createTime,
		UpdatedAt:      updateTime,
		Rating:         &rating,
		Details:        details,
		DeathDate:      &deathDate,
		HairColor:      hairColor,
		Weight:         &weight,
		IgnoreAutoTag:  autoTagIgnored,
		TagIDs:         models.NewRelatedIDs([]int{}),
		StashIDs:       models.NewRelatedStashIDs(stashIDs),
	}
}

func createEmptyPerformer(id int) models.Performer {
	return models.Performer{
		ID:        id,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
		Aliases:   models.NewRelatedStrings([]string{}),
		URLs:      models.NewRelatedStrings([]string{}),
		TagIDs:    models.NewRelatedIDs([]int{}),
		StashIDs:  models.NewRelatedStashIDs([]models.StashID{}),
	}
}

func createFullJSONPerformer(name string, image string, withCustomFields bool) *jsonschema.Performer {
	ret := &jsonschema.Performer{
		Name:           name,
		Disambiguation: disambiguation,
		URLs:           []string{url, twitter, instagram},
		Aliases:        aliases,
		Birthdate:      birthDate.String(),
		CareerLength:   careerLength,
		Country:        country,
		Ethnicity:      ethnicity,
		EyeColor:       eyeColor,
		FakeTits:       fakeTits,
		PenisLength:    penisLength,
		Circumcised:    circumcised,
		Favorite:       true,
		Gender:         gender,
		Height:         strconv.Itoa(height),
		Measurements:   measurements,
		Piercings:      piercings,
		Tattoos:        tattoos,
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
		Rating:        rating,
		Image:         image,
		Details:       details,
		DeathDate:     deathDate.String(),
		HairColor:     hairColor,
		Weight:        weight,
		StashIDs:      stashIDs,
		IgnoreAutoTag: autoTagIgnored,
		CustomFields:  emptyCustomFields,
	}

	if withCustomFields {
		ret.CustomFields = customFields
	}
	return ret
}

func createEmptyJSONPerformer() *jsonschema.Performer {
	return &jsonschema.Performer{
		Aliases:  []string{},
		URLs:     []string{},
		StashIDs: []models.StashID{},
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
		CustomFields: emptyCustomFields,
	}
}

type testScenario struct {
	input        models.Performer
	customFields map[string]interface{}
	expected     *jsonschema.Performer
	err          bool
}

var scenarios []testScenario

func initTestTable() {
	scenarios = []testScenario{
		{
			*createFullPerformer(performerID, performerName),
			emptyCustomFields,
			createFullJSONPerformer(performerName, image, false),
			false,
		},
		{
			*createFullPerformer(customFieldsID, performerName),
			customFields,
			createFullJSONPerformer(performerName, image, true),
			false,
		},
		{
			createEmptyPerformer(noImageID),
			emptyCustomFields,
			createEmptyJSONPerformer(),
			false,
		},
		{
			*createFullPerformer(errImageID, performerName),
			emptyCustomFields,
			createFullJSONPerformer(performerName, "", false),
			// failure to get image should not cause an error
			false,
		},
		{
			*createFullPerformer(errCustomFieldsID, performerName),
			customFields,
			nil,
			// failure to get custom fields should cause an error
			true,
		},
	}
}

func TestToJSON(t *testing.T) {
	initTestTable()

	db := mocks.NewDatabase()

	imageErr := errors.New("error getting image")
	customFieldsErr := errors.New("error getting custom fields")

	db.Performer.On("GetImage", testCtx, performerID).Return(imageBytes, nil).Once()
	db.Performer.On("GetImage", testCtx, customFieldsID).Return(imageBytes, nil).Once()
	db.Performer.On("GetImage", testCtx, noImageID).Return(nil, nil).Once()
	db.Performer.On("GetImage", testCtx, errImageID).Return(nil, imageErr).Once()

	db.Performer.On("GetCustomFields", testCtx, performerID).Return(emptyCustomFields, nil).Once()
	db.Performer.On("GetCustomFields", testCtx, customFieldsID).Return(customFields, nil).Once()
	db.Performer.On("GetCustomFields", testCtx, noImageID).Return(emptyCustomFields, nil).Once()
	db.Performer.On("GetCustomFields", testCtx, errImageID).Return(emptyCustomFields, nil).Once()
	db.Performer.On("GetCustomFields", testCtx, errCustomFieldsID).Return(nil, customFieldsErr).Once()

	for i, s := range scenarios {
		tag := s.input
		json, err := ToJSON(testCtx, db.Performer, &tag)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	db.AssertExpectations(t)
}
