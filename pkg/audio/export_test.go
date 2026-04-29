package audio

import (
	"errors"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testing"
	"time"
)

const (
	audioID = 1

	studioID        = 4
	missingStudioID = 5
	errStudioID     = 6
	customFieldsID  = 7

	noTagsID  = 11
	errTagsID = 12

	noGroupsID     = 13
	errFindGroupID = 15

	errCustomFieldsID = 20
)

var (
	url        = "url"
	title      = "title"
	date       = "2001-01-01"
	dateObj, _ = models.ParseDate(date)
	rating     = 5
	organized  = true
	details    = "details"
)

var (
	studioName = "studioName"
	// galleryChecksum = "galleryChecksum"

	validGroup1  = 1
	validGroup2  = 2
	invalidGroup = 3

	group1Name = "group1Name"
	group2Name = "group2Name"

	group1Audio = 1
	group2Audio = 2
)

var names = []string{
	"name1",
	"name2",
}

const (
	path = "path"
)

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

var (
	emptyCustomFields = make(map[string]interface{})
	customFields      = map[string]interface{}{
		"customField1": "customValue1",
	}
)

func createFullAudio(id int) models.Audio {
	return models.Audio{
		ID:        id,
		Title:     title,
		Date:      &dateObj,
		Details:   details,
		Rating:    &rating,
		Organized: organized,
		URLs:      models.NewRelatedStrings([]string{url}),
		Files: models.NewRelatedAudioFiles([]*models.AudioFile{
			{
				BaseFile: &models.BaseFile{
					Path: path,
				},
			},
		}),
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createEmptyAudio(id int) models.Audio {
	return models.Audio{
		ID: id,
		Files: models.NewRelatedAudioFiles([]*models.AudioFile{
			{
				BaseFile: &models.BaseFile{
					Path: path,
				},
			},
		}),
		URLs:      models.NewRelatedStrings([]string{}),
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createFullJSONAudio(customFields map[string]interface{}) *jsonschema.Audio {
	return &jsonschema.Audio{
		Title:     title,
		Files:     []string{path},
		Date:      date,
		Details:   details,
		Rating:    rating,
		Organized: organized,
		URLs:      []string{url},
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
		CustomFields: customFields,
	}
}

func createEmptyJSONAudio() *jsonschema.Audio {
	return &jsonschema.Audio{
		URLs:  []string{},
		Files: []string{path},
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
		CustomFields: emptyCustomFields,
	}
}

type basicTestScenario struct {
	input        models.Audio
	customFields map[string]interface{}
	expected     *jsonschema.Audio
	err          bool
}

var scenarios = []basicTestScenario{
	{
		createFullAudio(audioID),
		emptyCustomFields,
		createFullJSONAudio(emptyCustomFields),
		false,
	},
	{
		createFullAudio(customFieldsID),
		customFields,
		createFullJSONAudio(customFields),
		false,
	},
	{
		createFullAudio(errCustomFieldsID),
		customFields,
		createFullJSONAudio(customFields),
		true,
	},
	{
		createEmptyAudio(audioID),
		emptyCustomFields,
		createEmptyJSONAudio(),
		false,
	},
}

func TestToJSON(t *testing.T) {
	db := mocks.NewDatabase()

	db.Audio.On("GetViewDates", testCtx, mock.Anything).Return(nil, nil)
	db.Audio.On("GetODates", testCtx, mock.Anything).Return(nil, nil)
	db.Audio.On("GetCustomFields", testCtx, customFieldsID).Return(customFields, nil).Once()
	db.Audio.On("GetCustomFields", testCtx, errCustomFieldsID).Return(nil, errors.New("error getting custom fields")).Once()
	db.Audio.On("GetCustomFields", testCtx, mock.Anything).Return(emptyCustomFields, nil)

	for i, s := range scenarios {
		audio := s.input
		json, err := ToBasicJSON(testCtx, db.Audio, &audio)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		case err != nil:
			// error case already handled, no need for assertion
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	db.AssertExpectations(t)
}

func createStudioAudio(studioID int) models.Audio {
	return models.Audio{
		StudioID: &studioID,
	}
}

type stringTestScenario struct {
	input    models.Audio
	expected string
	err      bool
}

var getStudioScenarios = []stringTestScenario{
	{
		createStudioAudio(studioID),
		studioName,
		false,
	},
	{
		createStudioAudio(missingStudioID),
		"",
		false,
	},
	{
		createStudioAudio(errStudioID),
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
		audio := s.input
		json, err := GetStudioName(testCtx, db.Studio, &audio)

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

type stringSliceTestScenario struct {
	input    models.Audio
	expected []string
	err      bool
}

var getTagNamesScenarios = []stringSliceTestScenario{
	{
		createEmptyAudio(audioID),
		names,
		false,
	},
	{
		createEmptyAudio(noTagsID),
		nil,
		false,
	},
	{
		createEmptyAudio(errTagsID),
		nil,
		true,
	},
}

func getTags(names []string) []*models.Tag {
	var ret []*models.Tag
	for _, n := range names {
		ret = append(ret, &models.Tag{
			Name: n,
		})
	}

	return ret
}

func TestGetTagNames(t *testing.T) {
	db := mocks.NewDatabase()

	tagErr := errors.New("error getting tag")

	db.Tag.On("FindByAudioID", testCtx, audioID).Return(getTags(names), nil).Once()
	db.Tag.On("FindByAudioID", testCtx, noTagsID).Return(nil, nil).Once()
	db.Tag.On("FindByAudioID", testCtx, errTagsID).Return(nil, tagErr).Once()

	for i, s := range getTagNamesScenarios {
		audio := s.input
		json, err := GetTagNames(testCtx, db.Tag, &audio)

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

type audioGroupsTestScenario struct {
	input    models.Audio
	expected []jsonschema.AudioGroup
	err      bool
}

var validGroups = models.NewRelatedGroupsAudio([]models.GroupsAudios{
	{
		GroupID:    validGroup1,
		AudioIndex: &group1Audio,
	},
	{
		GroupID:    validGroup2,
		AudioIndex: &group2Audio,
	},
})

var invalidGroups = models.NewRelatedGroupsAudio([]models.GroupsAudios{
	{
		GroupID:    invalidGroup,
		AudioIndex: &group1Audio,
	},
})

var getAudioGroupsJSONScenarios = []audioGroupsTestScenario{
	{
		models.Audio{
			ID:     audioID,
			Groups: validGroups,
		},
		[]jsonschema.AudioGroup{
			{
				GroupName:  group1Name,
				AudioIndex: group1Audio,
			},
			{
				GroupName:  group2Name,
				AudioIndex: group2Audio,
			},
		},
		false,
	},
	{
		models.Audio{
			ID:     noGroupsID,
			Groups: models.NewRelatedGroupsAudio([]models.GroupsAudios{}),
		},
		nil,
		false,
	},
	{
		models.Audio{
			ID:     errFindGroupID,
			Groups: invalidGroups,
		},
		nil,
		true,
	},
}

func TestGetAudioGroupsJSON(t *testing.T) {
	db := mocks.NewDatabase()

	groupErr := errors.New("error getting group")

	db.Group.On("Find", testCtx, validGroup1).Return(&models.Group{
		Name: group1Name,
	}, nil).Once()
	db.Group.On("Find", testCtx, validGroup2).Return(&models.Group{
		Name: group2Name,
	}, nil).Once()
	db.Group.On("Find", testCtx, invalidGroup).Return(nil, groupErr).Once()

	for i, s := range getAudioGroupsJSONScenarios {
		audio := s.input
		json, err := GetAudioGroupsJSON(testCtx, db.Group, &audio)

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
