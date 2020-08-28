package studio

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
	studioID              = 1
	noImageID             = 2
	errImageID            = 3
	missingParentStudioID = 4
	errStudioID           = 5

	parentStudioID    = 10
	missingStudioID   = 11
	errParentStudioID = 12
)

const studioName = "testStudio"
const url = "url"

const parentStudioName = "parentStudio"

var parentStudio models.Studio = models.Studio{
	Name: modelstest.NullString(parentStudioName),
}

var imageBytes = []byte("imageBytes")

const image = "aW1hZ2VCeXRlcw=="

var createTime time.Time = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
var updateTime time.Time = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)

func createFullStudio(id int, parentID int) models.Studio {
	return models.Studio{
		ID:       id,
		Name:     modelstest.NullString(studioName),
		URL:      modelstest.NullString(url),
		ParentID: modelstest.NullInt64(int64(parentID)),
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	}
}

func createEmptyStudio(id int) models.Studio {
	return models.Studio{
		ID: id,
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	}
}

func createFullJSONStudio(parentStudio, image string) *jsonschema.Studio {
	return &jsonschema.Studio{
		Name: studioName,
		URL:  url,
		CreatedAt: models.JSONTime{
			Time: createTime,
		},
		UpdatedAt: models.JSONTime{
			Time: updateTime,
		},
		ParentStudio: parentStudio,
		Image:        image,
	}
}

func createEmptyJSONStudio() *jsonschema.Studio {
	return &jsonschema.Studio{
		CreatedAt: models.JSONTime{
			Time: createTime,
		},
		UpdatedAt: models.JSONTime{
			Time: updateTime,
		},
	}
}

type testScenario struct {
	input    models.Studio
	expected *jsonschema.Studio
	err      bool
}

var scenarios []testScenario

func initTestTable() {
	scenarios = []testScenario{
		testScenario{
			createFullStudio(studioID, parentStudioID),
			createFullJSONStudio(parentStudioName, image),
			false,
		},
		testScenario{
			createEmptyStudio(noImageID),
			createEmptyJSONStudio(),
			false,
		},
		testScenario{
			createFullStudio(errImageID, parentStudioID),
			nil,
			true,
		},
		testScenario{
			createFullStudio(missingParentStudioID, missingStudioID),
			createFullJSONStudio("", image),
			false,
		},
		testScenario{
			createFullStudio(errStudioID, errParentStudioID),
			nil,
			true,
		},
	}
}

func TestToJSON(t *testing.T) {
	initTestTable()

	mockStudioReader := &mocks.StudioReaderWriter{}

	imageErr := errors.New("error getting image")

	mockStudioReader.On("GetStudioImage", studioID).Return(imageBytes, nil).Once()
	mockStudioReader.On("GetStudioImage", noImageID).Return(nil, nil).Once()
	mockStudioReader.On("GetStudioImage", errImageID).Return(nil, imageErr).Once()
	mockStudioReader.On("GetStudioImage", missingParentStudioID).Return(imageBytes, nil).Maybe()
	mockStudioReader.On("GetStudioImage", errStudioID).Return(imageBytes, nil).Maybe()

	parentStudioErr := errors.New("error getting parent studio")

	mockStudioReader.On("Find", parentStudioID).Return(&parentStudio, nil)
	mockStudioReader.On("Find", missingStudioID).Return(nil, nil)
	mockStudioReader.On("Find", errParentStudioID).Return(nil, parentStudioErr)

	for i, s := range scenarios {
		studio := s.input
		json, err := ToJSON(mockStudioReader, &studio)

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
