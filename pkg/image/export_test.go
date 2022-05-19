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
	// noImageID  = 2
	errImageID = 3

	studioID        = 4
	missingStudioID = 5
	errStudioID     = 6

	// noGalleryID  = 7
	// errGalleryID = 8

	// noTagsID  = 11
	errTagsID = 12

	// noMoviesID     = 13
	// errMoviesID    = 14
	// errFindMovieID = 15

	// noMarkersID         = 16
	// errMarkersID        = 17
	// errFindPrimaryTagID = 18
	// errFindByMarkerID   = 19
)

const (
	checksum  = "checksum"
	title     = "title"
	rating    = 5
	organized = true
	ocounter  = 2
	size      = 123
	width     = 100
	height    = 100
)

const (
	studioName = "studioName"
	// galleryChecksum = "galleryChecksum"
)

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

func createFullImage(id int) models.Image {
	return models.Image{
		ID:        id,
		Title:     models.NullString(title),
		Checksum:  checksum,
		Height:    models.NullInt64(height),
		OCounter:  ocounter,
		Rating:    models.NullInt64(rating),
		Size:      models.NullInt64(int64(size)),
		Organized: organized,
		Width:     models.NullInt64(width),
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
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
		StudioID: models.NullInt64(int64(studioID)),
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

// var getGalleryChecksumScenarios = []stringTestScenario{
// 	{
// 		createEmptyImage(imageID),
// 		galleryChecksum,
// 		false,
// 	},
// 	{
// 		createEmptyImage(noGalleryID),
// 		"",
// 		false,
// 	},
// 	{
// 		createEmptyImage(errGalleryID),
// 		"",
// 		true,
// 	},
// }

// func TestGetGalleryChecksum(t *testing.T) {
// 	mockGalleryReader := &mocks.GalleryReaderWriter{}

// 	galleryErr := errors.New("error getting gallery")

// 	mockGalleryReader.On("FindByImageID", imageID).Return(&models.Gallery{
// 		Checksum: galleryChecksum,
// 	}, nil).Once()
// 	mockGalleryReader.On("FindByImageID", noGalleryID).Return(nil, nil).Once()
// 	mockGalleryReader.On("FindByImageID", errGalleryID).Return(nil, galleryErr).Once()

// 	for i, s := range getGalleryChecksumScenarios {
// 		image := s.input
// 		json, err := GetGalleryChecksum(mockGalleryReader, &image)

// 		if !s.err && err != nil {
// 			t.Errorf("[%d] unexpected error: %s", i, err.Error())
// 		} else if s.err && err == nil {
// 			t.Errorf("[%d] expected error not returned", i)
// 		} else {
// 			assert.Equal(t, s.expected, json, "[%d]", i)
// 		}
// 	}

// 	mockGalleryReader.AssertExpectations(t)
// }
