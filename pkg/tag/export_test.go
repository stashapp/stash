package tag

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
	tagID         = 1
	noImageID     = 2
	errImageID    = 3
	errAliasID    = 4
	withParentsID = 5
	errParentsID  = 6
)

const (
	tagName     = "testTag"
	description = "description"
)

var (
	autoTagIgnored = true
	createTime     = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime     = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

func createTag(id int) models.Tag {
	return models.Tag{
		ID:            id,
		Name:          tagName,
		Favorite:      true,
		Description:   description,
		IgnoreAutoTag: autoTagIgnored,
		CreatedAt:     createTime,
		UpdatedAt:     updateTime,
	}
}

func createJSONTag(aliases []string, image string, parents []string) *jsonschema.Tag {
	return &jsonschema.Tag{
		Name:          tagName,
		Favorite:      true,
		Description:   description,
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
			createJSONTag(nil, "", nil),
			// getting the image should not cause an error
			false,
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

	db := mocks.NewDatabase()

	imageErr := errors.New("error getting image")
	aliasErr := errors.New("error getting aliases")
	parentsErr := errors.New("error getting parents")

	db.Tag.On("GetAliases", testCtx, tagID).Return([]string{"alias"}, nil).Once()
	db.Tag.On("GetAliases", testCtx, noImageID).Return(nil, nil).Once()
	db.Tag.On("GetAliases", testCtx, errImageID).Return(nil, nil).Once()
	db.Tag.On("GetAliases", testCtx, errAliasID).Return(nil, aliasErr).Once()
	db.Tag.On("GetAliases", testCtx, withParentsID).Return(nil, nil).Once()
	db.Tag.On("GetAliases", testCtx, errParentsID).Return(nil, nil).Once()

	db.Tag.On("GetImage", testCtx, tagID).Return(imageBytes, nil).Once()
	db.Tag.On("GetImage", testCtx, noImageID).Return(nil, nil).Once()
	db.Tag.On("GetImage", testCtx, errImageID).Return(nil, imageErr).Once()
	db.Tag.On("GetImage", testCtx, withParentsID).Return(imageBytes, nil).Once()
	db.Tag.On("GetImage", testCtx, errParentsID).Return(nil, nil).Once()

	db.Tag.On("FindByChildTagID", testCtx, tagID).Return(nil, nil).Once()
	db.Tag.On("FindByChildTagID", testCtx, noImageID).Return(nil, nil).Once()
	db.Tag.On("FindByChildTagID", testCtx, withParentsID).Return([]*models.Tag{{Name: "parent"}}, nil).Once()
	db.Tag.On("FindByChildTagID", testCtx, errParentsID).Return(nil, parentsErr).Once()
	db.Tag.On("FindByChildTagID", testCtx, errImageID).Return(nil, nil).Once()

	for i, s := range scenarios {
		tag := s.tag
		json, err := ToJSON(testCtx, db.Tag, &tag)

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
