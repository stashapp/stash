package group

import (
	"errors"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testing"
	"time"
)

const (
	movieID = iota + 1
	emptyID
	errFrontImageID
	errBackImageID
	errStudioMovieID
	missingStudioMovieID
	errCustomFieldsID
)

const (
	studioID = iota + 1
	missingStudioID
	errStudioID
)

const movieName = "testMovie"
const movieAliases = "aliases"

var (
	date       = "2001-01-01"
	dateObj, _ = models.ParseDate(date)
	rating     = 5
	duration   = 100
	director   = "director"
	synopsis   = "synopsis"
	url        = "url"
)

const studioName = "studio"

const (
	frontImage = "ZnJvbnRJbWFnZUJ5dGVz"
	backImage  = "YmFja0ltYWdlQnl0ZXM="
)

var (
	frontImageBytes = []byte("frontImageBytes")
	backImageBytes  = []byte("backImageBytes")

	emptyCustomFields = make(map[string]interface{})
	customFields      = map[string]interface{}{
		"customField1": "customValue1",
	}
)

var movieStudio models.Studio = models.Studio{
	Name: studioName,
}

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

func createFullMovie(id int, studioID int) models.Group {
	return models.Group{
		ID:        id,
		Name:      movieName,
		Aliases:   movieAliases,
		Date:      &dateObj,
		Rating:    &rating,
		Duration:  &duration,
		Director:  director,
		Synopsis:  synopsis,
		URLs:      models.NewRelatedStrings([]string{url}),
		StudioID:  &studioID,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createEmptyMovie(id int) models.Group {
	return models.Group{
		ID:        id,
		URLs:      models.NewRelatedStrings([]string{}),
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createFullJSONMovie(studio, frontImage, backImage string, customFields map[string]interface{}) *jsonschema.Group {
	return &jsonschema.Group{
		Name:       movieName,
		Aliases:    movieAliases,
		Date:       date,
		Rating:     rating,
		Duration:   duration,
		Director:   director,
		Synopsis:   synopsis,
		URLs:       []string{url},
		Studio:     studio,
		FrontImage: frontImage,
		BackImage:  backImage,
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
		CustomFields: customFields,
	}
}

func createEmptyJSONMovie() *jsonschema.Group {
	return &jsonschema.Group{
		URLs: []string{},
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
	movie        models.Group
	customFields map[string]interface{}
	expected     *jsonschema.Group
	err          bool
}

var scenarios []testScenario

func initTestTable() {
	scenarios = []testScenario{
		{
			createFullMovie(movieID, studioID),
			customFields,
			createFullJSONMovie(studioName, frontImage, backImage, customFields),
			false,
		},
		{
			createEmptyMovie(emptyID),
			emptyCustomFields,
			createEmptyJSONMovie(),
			false,
		},
		{
			createFullMovie(errFrontImageID, studioID),
			emptyCustomFields,
			createFullJSONMovie(studioName, "", backImage, emptyCustomFields),
			// failure to get front image should not cause error
			false,
		},
		{
			createFullMovie(errBackImageID, studioID),
			emptyCustomFields,
			createFullJSONMovie(studioName, frontImage, "", emptyCustomFields),
			// failure to get back image should not cause error
			false,
		},
		{
			createFullMovie(errStudioMovieID, errStudioID),
			emptyCustomFields,
			nil,
			true,
		},
		{
			createFullMovie(missingStudioMovieID, missingStudioID),
			emptyCustomFields,
			createFullJSONMovie("", frontImage, backImage, emptyCustomFields),
			false,
		},
		{
			createFullMovie(errCustomFieldsID, studioID),
			customFields,
			nil,
			true,
		},
	}
}

func TestToJSON(t *testing.T) {
	initTestTable()

	db := mocks.NewDatabase()

	imageErr := errors.New("error getting image")

	db.Group.On("GetFrontImage", testCtx, movieID).Return(frontImageBytes, nil).Once()
	db.Group.On("GetFrontImage", testCtx, missingStudioMovieID).Return(frontImageBytes, nil).Once()
	db.Group.On("GetFrontImage", testCtx, emptyID).Return(nil, nil).Once().Maybe()
	db.Group.On("GetFrontImage", testCtx, errFrontImageID).Return(nil, imageErr).Once()
	db.Group.On("GetFrontImage", testCtx, errBackImageID).Return(frontImageBytes, nil).Once()
	db.Group.On("GetFrontImage", testCtx, errCustomFieldsID).Return(nil, nil).Once()

	db.Group.On("GetBackImage", testCtx, movieID).Return(backImageBytes, nil).Once()
	db.Group.On("GetBackImage", testCtx, missingStudioMovieID).Return(backImageBytes, nil).Once()
	db.Group.On("GetBackImage", testCtx, emptyID).Return(nil, nil).Once()
	db.Group.On("GetBackImage", testCtx, errBackImageID).Return(nil, imageErr).Once()
	db.Group.On("GetBackImage", testCtx, errFrontImageID).Return(backImageBytes, nil).Maybe()
	db.Group.On("GetBackImage", testCtx, errStudioMovieID).Return(backImageBytes, nil).Maybe()
	db.Group.On("GetBackImage", testCtx, errCustomFieldsID).Return(nil, nil).Once()

	db.Group.On("GetCustomFields", testCtx, movieID).Return(customFields, nil).Once()
	db.Group.On("GetCustomFields", testCtx, errCustomFieldsID).Return(nil, errors.New("error getting custom fields")).Once()
	db.Group.On("GetCustomFields", testCtx, mock.Anything).Return(emptyCustomFields, nil).Times(4)

	studioErr := errors.New("error getting studio")

	db.Studio.On("Find", testCtx, studioID).Return(&movieStudio, nil)
	db.Studio.On("Find", testCtx, missingStudioID).Return(nil, nil)
	db.Studio.On("Find", testCtx, errStudioID).Return(nil, studioErr)

	for i, s := range scenarios {
		movie := s.movie
		json, err := ToJSON(testCtx, db.Group, db.Studio, &movie)

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
