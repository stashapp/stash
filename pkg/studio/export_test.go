package studio

import (
	"errors"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

const (
	noImageID             = 2
	errImageID            = 3
	missingParentStudioID = 4
	errStudioID           = 5
	customFieldsID        = 6

	parentStudioID    = 10
	missingStudioID   = 11
	errParentStudioID = 12
	errCustomFieldsID = 13
)

var (
	studioName        = "testStudio"
	url               = "url"
	details           = "details"
	parentStudioName  = "parentStudio"
	autoTagIgnored    = true
	emptyCustomFields = make(map[string]interface{})
	customFields      = map[string]interface{}{
		"customField1": "customValue1",
	}
)

var studioID = 1
var rating = 5
var parentStudio models.Studio = models.Studio{
	Name: parentStudioName,
}

var imageBytes = []byte("imageBytes")

var aliases = []string{"alias"}
var stashID = models.StashID{
	StashID:  "StashID",
	Endpoint: "Endpoint",
}
var stashIDs = []models.StashID{
	stashID,
}

const image = "aW1hZ2VCeXRlcw=="

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.Local)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.Local)
)

func createFullStudio(id int, parentID int) models.Studio {
	ret := models.Studio{
		ID:            id,
		Name:          studioName,
		URLs:          models.NewRelatedStrings([]string{url}),
		Details:       details,
		Favorite:      true,
		CreatedAt:     createTime,
		UpdatedAt:     updateTime,
		Rating:        &rating,
		IgnoreAutoTag: autoTagIgnored,
		Aliases:       models.NewRelatedStrings(aliases),
		TagIDs:        models.NewRelatedIDs([]int{}),
		StashIDs:      models.NewRelatedStashIDs(stashIDs),
	}

	if parentID != 0 {
		ret.ParentID = &parentID
	}

	return ret
}

func createEmptyStudio(id int) models.Studio {
	return models.Studio{
		ID:        id,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
		URLs:      models.NewRelatedStrings([]string{}),
		Aliases:   models.NewRelatedStrings([]string{}),
		TagIDs:    models.NewRelatedIDs([]int{}),
		StashIDs:  models.NewRelatedStashIDs([]models.StashID{}),
	}
}

func createFullJSONStudio(parentStudio, image string, aliases []string, customFields map[string]interface{}) *jsonschema.Studio {
	return &jsonschema.Studio{
		Name:     studioName,
		URLs:     []string{url},
		Details:  details,
		Favorite: true,
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
		ParentStudio:  parentStudio,
		Image:         image,
		Rating:        rating,
		Aliases:       aliases,
		StashIDs:      stashIDs,
		IgnoreAutoTag: autoTagIgnored,
		CustomFields:  customFields,
	}
}

func createEmptyJSONStudio() *jsonschema.Studio {
	return &jsonschema.Studio{
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
		Aliases:      []string{},
		URLs:         []string{},
		StashIDs:     []models.StashID{},
		CustomFields: emptyCustomFields,
	}
}

type testScenario struct {
	input        models.Studio
	customFields map[string]interface{}
	expected     *jsonschema.Studio
	err          bool
}

var scenarios []testScenario

func initTestTable() {
	scenarios = []testScenario{
		{
			createFullStudio(studioID, parentStudioID),
			emptyCustomFields,
			createFullJSONStudio(parentStudioName, image, []string{"alias"}, emptyCustomFields),
			false,
		},
		{
			createFullStudio(customFieldsID, parentStudioID),
			customFields,
			createFullJSONStudio(parentStudioName, image, []string{"alias"}, customFields),
			false,
		},
		{
			createEmptyStudio(noImageID),
			emptyCustomFields,
			createEmptyJSONStudio(),
			false,
		},
		{
			createFullStudio(errImageID, parentStudioID),
			emptyCustomFields,
			createFullJSONStudio(parentStudioName, "", []string{"alias"}, emptyCustomFields),
			// failure to get image is not an error
			false,
		},
		{
			createFullStudio(missingParentStudioID, missingStudioID),
			emptyCustomFields,
			createFullJSONStudio("", image, []string{"alias"}, emptyCustomFields),
			false,
		},
		{
			createFullStudio(errStudioID, errParentStudioID),
			emptyCustomFields,
			nil,
			true,
		},
		{
			createFullStudio(errCustomFieldsID, parentStudioID),
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

	db.Studio.On("GetImage", testCtx, studioID).Return(imageBytes, nil).Once()
	db.Studio.On("GetImage", testCtx, noImageID).Return(nil, nil).Once()
	db.Studio.On("GetImage", testCtx, errImageID).Return(nil, imageErr).Once()
	db.Studio.On("GetImage", testCtx, missingParentStudioID).Return(imageBytes, nil).Maybe()
	db.Studio.On("GetImage", testCtx, errStudioID).Return(imageBytes, nil).Maybe()
	db.Studio.On("GetImage", testCtx, customFieldsID).Return(imageBytes, nil).Once()

	parentStudioErr := errors.New("error getting parent studio")

	db.Studio.On("Find", testCtx, parentStudioID).Return(&parentStudio, nil)
	db.Studio.On("Find", testCtx, missingStudioID).Return(nil, nil)
	db.Studio.On("Find", testCtx, errParentStudioID).Return(nil, parentStudioErr)

	customFieldsErr := errors.New("error getting custom fields")

	db.Studio.On("GetCustomFields", testCtx, studioID).Return(emptyCustomFields, nil).Once()
	db.Studio.On("GetCustomFields", testCtx, customFieldsID).Return(customFields, nil).Once()
	db.Studio.On("GetCustomFields", testCtx, missingParentStudioID).Return(emptyCustomFields, nil).Once()
	db.Studio.On("GetCustomFields", testCtx, noImageID).Return(emptyCustomFields, nil).Once()
	db.Studio.On("GetCustomFields", testCtx, errImageID).Return(emptyCustomFields, nil).Once()
	db.Studio.On("GetCustomFields", testCtx, errCustomFieldsID).Return(nil, customFieldsErr).Once()

	for i, s := range scenarios {
		studio := s.input
		json, err := ToJSON(testCtx, db.Studio, &studio)

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
