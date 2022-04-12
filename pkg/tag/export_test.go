package tag

import (
	"errors"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

const (
	tagID         = 1
	noImageID     = 2
	errImageID    = 3
	errAliasID    = 4
	withParentsID = 5
	errParentsID  = 6
)

const tagName = "testTag"

var (
	autoTagIgnored = true
	createTime     = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime     = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

func createTag(id int) models.Tag {
	return models.Tag{
		ID:            id,
		Name:          tagName,
		IgnoreAutoTag: autoTagIgnored,
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	}
}

func createJSONTag(aliases []string, image string, parents []string) *jsonschema.Tag {
	return &jsonschema.Tag{
		Name:          tagName,
		Aliases:       aliases,
		IgnoreAutoTag: autoTagIgnored,
		CreatedAt: models.JSONTime{
			Time: createTime,
		},
		UpdatedAt: models.JSONTime{
			Time: updateTime,
		},
		Image:   image,
		Parents: parents,
	}
}

type testScenario struct {
	tag      models.Tag
	expected *jsonschema.Tag
	err      bool
}

var scenarios []testScenario

func initTestTable() {
	scenarios = []testScenario{
		{
			createTag(tagID),
			createJSONTag([]string{"alias"}, image, nil),
			false,
		},
		{
			createTag(noImageID),
			createJSONTag(nil, "", nil),
			false,
		},
		{
			createTag(errImageID),
			nil,
			true,
		},
		{
			createTag(errAliasID),
			nil,
			true,
		},
		{
			createTag(withParentsID),
			createJSONTag(nil, image, []string{"parent"}),
			false,
		},
		{
			createTag(errParentsID),
			nil,
			true,
		},
	}
}

func TestToJSON(t *testing.T) {
	initTestTable()

	mockTagReader := &mocks.TagReaderWriter{}

	imageErr := errors.New("error getting image")
	aliasErr := errors.New("error getting aliases")
	parentsErr := errors.New("error getting parents")

	mockTagReader.On("GetAliases", tagID).Return([]string{"alias"}, nil).Once()
	mockTagReader.On("GetAliases", noImageID).Return(nil, nil).Once()
	mockTagReader.On("GetAliases", errImageID).Return(nil, nil).Once()
	mockTagReader.On("GetAliases", errAliasID).Return(nil, aliasErr).Once()
	mockTagReader.On("GetAliases", withParentsID).Return(nil, nil).Once()
	mockTagReader.On("GetAliases", errParentsID).Return(nil, nil).Once()

	mockTagReader.On("GetImage", tagID).Return(imageBytes, nil).Once()
	mockTagReader.On("GetImage", noImageID).Return(nil, nil).Once()
	mockTagReader.On("GetImage", errImageID).Return(nil, imageErr).Once()
	mockTagReader.On("GetImage", withParentsID).Return(imageBytes, nil).Once()
	mockTagReader.On("GetImage", errParentsID).Return(nil, nil).Once()

	mockTagReader.On("FindByChildTagID", tagID).Return(nil, nil).Once()
	mockTagReader.On("FindByChildTagID", noImageID).Return(nil, nil).Once()
	mockTagReader.On("FindByChildTagID", withParentsID).Return([]*models.Tag{{Name: "parent"}}, nil).Once()
	mockTagReader.On("FindByChildTagID", errParentsID).Return(nil, parentsErr).Once()

	for i, s := range scenarios {
		tag := s.tag
		json, err := ToJSON(mockTagReader, &tag)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockTagReader.AssertExpectations(t)
}
