package studio

import (
	"context"
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
	studioID              = 1
	noImageID             = 2
	errImageID            = 3
	missingParentStudioID = 4
	errStudioID           = 5
	errAliasID            = 6

	parentStudioID    = 10
	missingStudioID   = 11
	errParentStudioID = 12
)

const (
	studioName       = "testStudio"
	url              = "url"
	details          = "details"
	rating           = 5
	parentStudioName = "parentStudio"
	autoTagIgnored   = true
)

var parentStudio models.Studio = models.Studio{
	Name: models.NullString(parentStudioName),
}

var imageBytes = []byte("imageBytes")

var stashID = models.StashID{
	StashID:  "StashID",
	Endpoint: "Endpoint",
}
var stashIDs = []*models.StashID{
	&stashID,
}

const image = "aW1hZ2VCeXRlcw=="

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.Local)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.Local)
)

func createFullStudio(id int, parentID int) models.Studio {
	ret := models.Studio{
		ID:      id,
		Name:    models.NullString(studioName),
		URL:     models.NullString(url),
		Details: models.NullString(details),
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
		Rating:        models.NullInt64(rating),
		IgnoreAutoTag: autoTagIgnored,
	}

	if parentID != 0 {
		ret.ParentID = models.NullInt64(int64(parentID))
	}

	return ret
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

func createFullJSONStudio(parentStudio, image string, aliases []string) *jsonschema.Studio {
	return &jsonschema.Studio{
		Name:    studioName,
		URL:     url,
		Details: details,
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
		ParentStudio: parentStudio,
		Image:        image,
		Rating:       rating,
		Aliases:      aliases,
		StashIDs: []models.StashID{
			stashID,
		},
		IgnoreAutoTag: autoTagIgnored,
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
		{
			createFullStudio(studioID, parentStudioID),
			createFullJSONStudio(parentStudioName, image, []string{"alias"}),
			false,
		},
		{
			createEmptyStudio(noImageID),
			createEmptyJSONStudio(),
			false,
		},
		{
			createFullStudio(errImageID, parentStudioID),
			nil,
			true,
		},
		{
			createFullStudio(missingParentStudioID, missingStudioID),
			createFullJSONStudio("", image, nil),
			false,
		},
		{
			createFullStudio(errStudioID, errParentStudioID),
			nil,
			true,
		},
		{
			createFullStudio(errAliasID, parentStudioID),
			nil,
			true,
		},
	}
}

func TestToJSON(t *testing.T) {
	initTestTable()
	ctx := context.Background()

	mockStudioReader := &mocks.StudioReaderWriter{}

	imageErr := errors.New("error getting image")

	mockStudioReader.On("GetImage", ctx, studioID).Return(imageBytes, nil).Once()
	mockStudioReader.On("GetImage", ctx, noImageID).Return(nil, nil).Once()
	mockStudioReader.On("GetImage", ctx, errImageID).Return(nil, imageErr).Once()
	mockStudioReader.On("GetImage", ctx, missingParentStudioID).Return(imageBytes, nil).Maybe()
	mockStudioReader.On("GetImage", ctx, errStudioID).Return(imageBytes, nil).Maybe()
	mockStudioReader.On("GetImage", ctx, errAliasID).Return(imageBytes, nil).Maybe()

	parentStudioErr := errors.New("error getting parent studio")

	mockStudioReader.On("Find", ctx, parentStudioID).Return(&parentStudio, nil)
	mockStudioReader.On("Find", ctx, missingStudioID).Return(nil, nil)
	mockStudioReader.On("Find", ctx, errParentStudioID).Return(nil, parentStudioErr)

	aliasErr := errors.New("error getting aliases")

	mockStudioReader.On("GetAliases", ctx, studioID).Return([]string{"alias"}, nil).Once()
	mockStudioReader.On("GetAliases", ctx, noImageID).Return(nil, nil).Once()
	mockStudioReader.On("GetAliases", ctx, errImageID).Return(nil, nil).Once()
	mockStudioReader.On("GetAliases", ctx, missingParentStudioID).Return(nil, nil).Once()
	mockStudioReader.On("GetAliases", ctx, errAliasID).Return(nil, aliasErr).Once()

	mockStudioReader.On("GetStashIDs", ctx, studioID).Return(stashIDs, nil).Once()
	mockStudioReader.On("GetStashIDs", ctx, noImageID).Return(nil, nil).Once()
	mockStudioReader.On("GetStashIDs", ctx, missingParentStudioID).Return(stashIDs, nil).Once()

	for i, s := range scenarios {
		studio := s.input
		json, err := ToJSON(ctx, mockStudioReader, &studio)

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
