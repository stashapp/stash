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
	imageID    = 1
	errImageID = 3

	studioID        = 4
	missingStudioID = 5
	errStudioID     = 6
)

var (
	checksum        = "checksum"
	title           = "title"
	rating          = 5
	organized       = true
	ocounter        = 2
	size      int64 = 123
	width           = 100
	height          = 100
)

const (
	studioName = "studioName"
)

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

func createFullImage(id int) models.Image {
	return models.Image{
		ID:        id,
		Title:     title,
		Checksum:  checksum,
		Height:    &height,
		OCounter:  ocounter,
		Rating:    &rating,
		Size:      &size,
		Organized: organized,
		Width:     &width,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createFullJSONImage() *jsonschema.Image {
	return &jsonschema.Image{
		Title:     title,
		Checksum:  checksum,
		OCounter:  ocounter,
		Rating:    rating,
		Organized: organized,
		File: &jsonschema.ImageFile{
			Height: height,
			Size:   size,
			Width:  width,
		},
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
	mockStudioReader := &mocks.StudioReaderWriter{}

	studioErr := errors.New("error getting image")

	mockStudioReader.On("Find", testCtx, studioID).Return(&models.Studio{
		Name: models.NullString(studioName),
	}, nil).Once()
	mockStudioReader.On("Find", testCtx, missingStudioID).Return(nil, nil).Once()
	mockStudioReader.On("Find", testCtx, errStudioID).Return(nil, studioErr).Once()

	for i, s := range getStudioScenarios {
		image := s.input
		json, err := GetStudioName(testCtx, mockStudioReader, &image)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockStudioReader.AssertExpectations(t)
}
