package gallery

import (
	"errors"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/models/modelstest"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

const (
	galleryID = 1

	studioID        = 4
	missingStudioID = 5
	errStudioID     = 6

	noTagsID  = 11
	errTagsID = 12
)

const (
	path     = "path"
	zip      = true
	url      = "url"
	checksum = "checksum"
	title    = "title"
	date     = "2001-01-01"
	rating   = 5
	details  = "details"
)

const (
	studioName = "studioName"
)

var names = []string{
	"name1",
	"name2",
}

var createTime time.Time = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
var updateTime time.Time = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)

func createFullGallery(id int) models.Gallery {
	return models.Gallery{
		ID:       id,
		Path:     modelstest.NullString(path),
		Zip:      zip,
		Title:    modelstest.NullString(title),
		Checksum: checksum,
		Date: models.SQLiteDate{
			String: date,
			Valid:  true,
		},
		Details: modelstest.NullString(details),
		Rating:  modelstest.NullInt64(rating),
		URL:     modelstest.NullString(url),
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	}
}

func createEmptyGallery(id int) models.Gallery {
	return models.Gallery{
		ID: id,
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	}
}

func createFullJSONGallery() *jsonschema.Gallery {
	return &jsonschema.Gallery{
		Title:    title,
		Path:     path,
		Zip:      zip,
		Checksum: checksum,
		Date:     date,
		Details:  details,
		Rating:   rating,
		URL:      url,
		CreatedAt: models.JSONTime{
			Time: createTime,
		},
		UpdatedAt: models.JSONTime{
			Time: updateTime,
		},
	}
}

func createEmptyJSONGallery() *jsonschema.Gallery {
	return &jsonschema.Gallery{
		CreatedAt: models.JSONTime{
			Time: createTime,
		},
		UpdatedAt: models.JSONTime{
			Time: updateTime,
		},
	}
}

type basicTestScenario struct {
	input    models.Gallery
	expected *jsonschema.Gallery
	err      bool
}

var scenarios = []basicTestScenario{
	{
		createFullGallery(galleryID),
		createFullJSONGallery(),
		false,
	},
}

func TestToJSON(t *testing.T) {
	for i, s := range scenarios {
		gallery := s.input
		json, err := ToBasicJSON(&gallery)

		if !s.err && err != nil {
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		} else if s.err && err == nil {
			t.Errorf("[%d] expected error not returned", i)
		} else {
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}
}

func createStudioGallery(studioID int) models.Gallery {
	return models.Gallery{
		StudioID: modelstest.NullInt64(int64(studioID)),
	}
}

type stringTestScenario struct {
	input    models.Gallery
	expected string
	err      bool
}

var getStudioScenarios = []stringTestScenario{
	{
		createStudioGallery(studioID),
		studioName,
		false,
	},
	{
		createStudioGallery(missingStudioID),
		"",
		false,
	},
	{
		createStudioGallery(errStudioID),
		"",
		true,
	},
}

func TestGetStudioName(t *testing.T) {
	mockStudioReader := &mocks.StudioReaderWriter{}

	studioErr := errors.New("error getting image")

	mockStudioReader.On("Find", studioID).Return(&models.Studio{
		Name: modelstest.NullString(studioName),
	}, nil).Once()
	mockStudioReader.On("Find", missingStudioID).Return(nil, nil).Once()
	mockStudioReader.On("Find", errStudioID).Return(nil, studioErr).Once()

	for i, s := range getStudioScenarios {
		gallery := s.input
		json, err := GetStudioName(mockStudioReader, &gallery)

		if !s.err && err != nil {
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		} else if s.err && err == nil {
			t.Errorf("[%d] expected error not returned", i)
		} else {
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockStudioReader.AssertExpectations(t)
}
