package image

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
	imageID = 1

	studioID        = 4
	missingStudioID = 5
	errStudioID     = 6
)

var (
	title      = "title"
	rating     = 5
	url        = "http://a.com"
	date       = "2001-01-01"
	dateObj, _ = models.ParseDate(date)
	organized  = true
	ocounter   = 2
)

const (
	studioName = "studioName"
	path       = "path"
)

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

func createFullImage(id int) models.Image {
	return models.Image{
		ID: id,
		Files: models.NewRelatedFiles([]models.File{
			&models.BaseFile{
				Path: path,
			},
		}),
		Title:     title,
		OCounter:  ocounter,
		Rating:    &rating,
		Date:      &dateObj,
		URLs:      models.NewRelatedStrings([]string{url}),
		Organized: organized,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createFullJSONImage() *jsonschema.Image {
	return &jsonschema.Image{
		Title:     title,
		OCounter:  ocounter,
		Rating:    rating,
		Date:      date,
		URLs:      []string{url},
		Organized: organized,
		Files:     []string{path},
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
	}
}

type basicTestScenario struct {
	input    models.Image
	expected *jsonschema.Image
}

var scenarios = []basicTestScenario{
	{
		createFullImage(imageID),
		createFullJSONImage(),
	},
}

func TestToJSON(t *testing.T) {
	for i, s := range scenarios {
		image := s.input
		json := ToBasicJSON(&image)

		assert.Equal(t, s.expected, json, "[%d]", i)
	}
}

func createStudioImage(studioID int) models.Image {
	return models.Image{
		StudioID: &studioID,
	}
}

type stringTestScenario struct {
	input    models.Image
	expected string
	err      bool
}

var getStudioScenarios = []stringTestScenario{
	{
		createStudioImage(studioID),
		studioName,
		false,
	},
	{
		createStudioImage(missingStudioID),
		"",
		false,
	},
	{
		createStudioImage(errStudioID),
		"",
		true,
	},
}

func TestGetStudioName(t *testing.T) {
	db := mocks.NewDatabase()

	studioErr := errors.New("error getting image")

	db.Studio.On("Find", testCtx, studioID).Return(&models.Studio{
		Name: studioName,
	}, nil).Once()
	db.Studio.On("Find", testCtx, missingStudioID).Return(nil, nil).Once()
	db.Studio.On("Find", testCtx, errStudioID).Return(nil, studioErr).Once()

	for i, s := range getStudioScenarios {
		image := s.input
		json, err := GetStudioName(testCtx, db.Studio, &image)

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
