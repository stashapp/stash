package tag

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
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
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

	ctx := context.Background()
	mockTagReader := &mocks.TagReaderWriter{}

	imageErr := errors.New("error getting image")
	aliasErr := errors.New("error getting aliases")
	parentsErr := errors.New("error getting parents")

	mockTagReader.On("GetAliases", ctx, tagID).Return([]string{"alias"}, nil).Once()
	mockTagReader.On("GetAliases", ctx, noImageID).Return(nil, nil).Once()
	mockTagReader.On("GetAliases", ctx, errImageID).Return(nil, nil).Once()
	mockTagReader.On("GetAliases", ctx, errAliasID).Return(nil, aliasErr).Once()
	mockTagReader.On("GetAliases", ctx, withParentsID).Return(nil, nil).Once()
	mockTagReader.On("GetAliases", ctx, errParentsID).Return(nil, nil).Once()

	mockTagReader.On("GetImage", ctx, tagID).Return(imageBytes, nil).Once()
	mockTagReader.On("GetImage", ctx, noImageID).Return(nil, nil).Once()
	mockTagReader.On("GetImage", ctx, errImageID).Return(nil, imageErr).Once()
	mockTagReader.On("GetImage", ctx, withParentsID).Return(imageBytes, nil).Once()
	mockTagReader.On("GetImage", ctx, errParentsID).Return(nil, nil).Once()

	mockTagReader.On("FindByChildTagID", ctx, tagID).Return(nil, nil).Once()
	mockTagReader.On("FindByChildTagID", ctx, noImageID).Return(nil, nil).Once()
	mockTagReader.On("FindByChildTagID", ctx, withParentsID).Return([]*models.Tag{{Name: "parent"}}, nil).Once()
	mockTagReader.On("FindByChildTagID", ctx, errParentsID).Return(nil, parentsErr).Once()

	for i, s := range scenarios {
		tag := s.tag
		json, err := ToJSON(ctx, mockTagReader, &tag)

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
