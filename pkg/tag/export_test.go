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
	tagID = iota + 1
	customFieldsID
	noImageID
	errImageID
	errAliasID
	withParentsID
	errParentsID
	errCustomFieldsID
)

const (
	tagName     = "testTag"
	sortName    = "sortName"
	description = "description"
)

var (
	autoTagIgnored = true
	createTime     = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime     = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)

	emptyCustomFields = make(map[string]interface{})
	customFields      = map[string]interface{}{
		"customField1": "customValue1",
	}
)

func createTag(id int) models.Tag {
	return models.Tag{
		ID:            id,
		Name:          tagName,
		SortName:      sortName,
		Favorite:      true,
		Description:   description,
		IgnoreAutoTag: autoTagIgnored,
		CreatedAt:     createTime,
		UpdatedAt:     updateTime,
	}
}

func createJSONTag(aliases []string, image string, parents []string, withCustomFields bool) *jsonschema.Tag {
	ret := &jsonschema.Tag{
		Name:          tagName,
		SortName:      sortName,
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
		Image:        image,
		Parents:      parents,
		CustomFields: emptyCustomFields,
	}

	if withCustomFields {
		ret.CustomFields = customFields
	}

	return ret
}

type testScenario struct {
	tag          models.Tag
	customFields map[string]interface{}
	expected     *jsonschema.Tag
	err          bool
}

var scenarios []testScenario

func initTestTable() {
	scenarios = []testScenario{
		{
			createTag(tagID),
			emptyCustomFields,
			createJSONTag([]string{"alias"}, image, nil, false),
			false,
		},
		{
			createTag(customFieldsID),
			customFields,
			createJSONTag([]string{"alias"}, image, nil, true),
			false,
		},
		{
			createTag(noImageID),
			emptyCustomFields,
			createJSONTag(nil, "", nil, false),
			false,
		},
		{
			createTag(errImageID),
			emptyCustomFields,
			createJSONTag(nil, "", nil, false),
			// getting the image should not cause an error
			false,
		},
		{
			createTag(errAliasID),
			emptyCustomFields,
			nil,
			true,
		},
		{
			createTag(withParentsID),
			emptyCustomFields,
			createJSONTag(nil, image, []string{"parent"}, false),
			false,
		},
		{
			createTag(errParentsID),
			emptyCustomFields,
			nil,
			true,
		},
		{
			createTag(errCustomFieldsID),
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
	aliasErr := errors.New("error getting aliases")
	parentsErr := errors.New("error getting parents")
	customFieldsErr := errors.New("error getting custom fields")

	db.Tag.On("GetAliases", testCtx, tagID).Return([]string{"alias"}, nil).Once()
	db.Tag.On("GetAliases", testCtx, customFieldsID).Return([]string{"alias"}, nil).Once()
	db.Tag.On("GetAliases", testCtx, noImageID).Return(nil, nil).Once()
	db.Tag.On("GetAliases", testCtx, errImageID).Return(nil, nil).Once()
	db.Tag.On("GetAliases", testCtx, errAliasID).Return(nil, aliasErr).Once()
	db.Tag.On("GetAliases", testCtx, withParentsID).Return(nil, nil).Once()
	db.Tag.On("GetAliases", testCtx, errParentsID).Return(nil, nil).Once()
	db.Tag.On("GetAliases", testCtx, errCustomFieldsID).Return(nil, nil).Once()

	db.Tag.On("GetStashIDs", testCtx, tagID).Return(nil, nil).Once()
	db.Tag.On("GetStashIDs", testCtx, customFieldsID).Return(nil, nil).Once()
	db.Tag.On("GetStashIDs", testCtx, noImageID).Return(nil, nil).Once()
	db.Tag.On("GetStashIDs", testCtx, errImageID).Return(nil, nil).Once()
	// errAliasID test fails before GetStashIDs is called, so no mock needed
	db.Tag.On("GetStashIDs", testCtx, withParentsID).Return(nil, nil).Once()
	db.Tag.On("GetStashIDs", testCtx, errParentsID).Return(nil, nil).Once()
	db.Tag.On("GetStashIDs", testCtx, errCustomFieldsID).Return(nil, nil).Once()

	db.Tag.On("GetImage", testCtx, tagID).Return(imageBytes, nil).Once()
	db.Tag.On("GetImage", testCtx, customFieldsID).Return(imageBytes, nil).Once()
	db.Tag.On("GetImage", testCtx, noImageID).Return(nil, nil).Once()
	db.Tag.On("GetImage", testCtx, errImageID).Return(nil, imageErr).Once()
	db.Tag.On("GetImage", testCtx, withParentsID).Return(imageBytes, nil).Once()
	db.Tag.On("GetImage", testCtx, errParentsID).Return(nil, nil).Once()
	db.Tag.On("GetImage", testCtx, errCustomFieldsID).Return(nil, nil).Once()

	db.Tag.On("FindByChildTagID", testCtx, tagID).Return(nil, nil).Once()
	db.Tag.On("FindByChildTagID", testCtx, customFieldsID).Return(nil, nil).Once()
	db.Tag.On("FindByChildTagID", testCtx, noImageID).Return(nil, nil).Once()
	db.Tag.On("FindByChildTagID", testCtx, withParentsID).Return([]*models.Tag{{Name: "parent"}}, nil).Once()
	db.Tag.On("FindByChildTagID", testCtx, errParentsID).Return(nil, parentsErr).Once()
	db.Tag.On("FindByChildTagID", testCtx, errImageID).Return(nil, nil).Once()
	db.Tag.On("FindByChildTagID", testCtx, errCustomFieldsID).Return(nil, nil).Once()

	db.Tag.On("GetCustomFields", testCtx, tagID).Return(emptyCustomFields, nil).Once()
	db.Tag.On("GetCustomFields", testCtx, customFieldsID).Return(customFields, nil).Once()
	db.Tag.On("GetCustomFields", testCtx, noImageID).Return(emptyCustomFields, nil).Once()
	db.Tag.On("GetCustomFields", testCtx, errImageID).Return(emptyCustomFields, nil).Once()
	db.Tag.On("GetCustomFields", testCtx, withParentsID).Return(emptyCustomFields, nil).Once()
	db.Tag.On("GetCustomFields", testCtx, errCustomFieldsID).Return(nil, customFieldsErr).Once()

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
