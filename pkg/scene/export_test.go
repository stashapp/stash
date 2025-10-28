package scene

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
	sceneID    = 1
	noImageID  = 2
	errImageID = 3

	studioID        = 4
	missingStudioID = 5
	errStudioID     = 6

	noTagsID  = 11
	errTagsID = 12

	noGroupsID     = 13
	errFindGroupID = 15

	noMarkersID         = 16
	errMarkersID        = 17
	errFindPrimaryTagID = 18
	errFindByMarkerID   = 19
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

	group1Scene = 1
	group2Scene = 2
)

var names = []string{
	"name1",
	"name2",
}

var imageBytes = []byte("imageBytes")

var stashID = models.StashID{
	StashID:  "StashID",
	Endpoint: "Endpoint",
}

const (
	path        = "path"
	imageBase64 = "aW1hZ2VCeXRlcw=="
)

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

func createFullScene(id int) models.Scene {
	return models.Scene{
		ID:        id,
		Title:     title,
		Date:      &dateObj,
		Details:   details,
		Rating:    &rating,
		Organized: organized,
		URLs:      models.NewRelatedStrings([]string{url}),
		Files: models.NewRelatedVideoFiles([]*models.VideoFile{
			{
				BaseFile: &models.BaseFile{
					Path: path,
				},
			},
		}),
		StashIDs: models.NewRelatedStashIDs([]models.StashID{
			stashID,
		}),
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createEmptyScene(id int) models.Scene {
	return models.Scene{
		ID: id,
		Files: models.NewRelatedVideoFiles([]*models.VideoFile{
			{
				BaseFile: &models.BaseFile{
					Path: path,
				},
			},
		}),
		URLs:      models.NewRelatedStrings([]string{}),
		StashIDs:  models.NewRelatedStashIDs([]models.StashID{}),
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createFullJSONScene(image string) *jsonschema.Scene {
	return &jsonschema.Scene{
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
		Cover: image,
		StashIDs: []models.StashID{
			stashID,
		},
	}
}

func createEmptyJSONScene() *jsonschema.Scene {
	return &jsonschema.Scene{
		URLs:  []string{},
		Files: []string{path},
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
	}
}

type basicTestScenario struct {
	input    models.Scene
	expected *jsonschema.Scene
	err      bool
}

var scenarios = []basicTestScenario{
	{
		createFullScene(sceneID),
		createFullJSONScene(imageBase64),
		false,
	},
	{
		createEmptyScene(noImageID),
		createEmptyJSONScene(),
		false,
	},
	{
		createFullScene(errImageID),
		createFullJSONScene(""),
		// failure to get image should not cause an error
		false,
	},
}

func TestToJSON(t *testing.T) {
	db := mocks.NewDatabase()

	imageErr := errors.New("error getting image")

	db.Scene.On("GetCover", testCtx, sceneID).Return(imageBytes, nil).Once()
	db.Scene.On("GetCover", testCtx, noImageID).Return(nil, nil).Once()
	db.Scene.On("GetCover", testCtx, errImageID).Return(nil, imageErr).Once()
	db.Scene.On("GetViewDates", testCtx, mock.Anything).Return(nil, nil)
	db.Scene.On("GetODates", testCtx, mock.Anything).Return(nil, nil)

	for i, s := range scenarios {
		scene := s.input
		json, err := ToBasicJSON(testCtx, db.Scene, &scene)

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

func createStudioScene(studioID int) models.Scene {
	return models.Scene{
		StudioID: &studioID,
	}
}

type stringTestScenario struct {
	input    models.Scene
	expected string
	err      bool
}

var getStudioScenarios = []stringTestScenario{
	{
		createStudioScene(studioID),
		studioName,
		false,
	},
	{
		createStudioScene(missingStudioID),
		"",
		false,
	},
	{
		createStudioScene(errStudioID),
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
		scene := s.input
		json, err := GetStudioName(testCtx, db.Studio, &scene)

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
	input    models.Scene
	expected []string
	err      bool
}

var getTagNamesScenarios = []stringSliceTestScenario{
	{
		createEmptyScene(sceneID),
		names,
		false,
	},
	{
		createEmptyScene(noTagsID),
		nil,
		false,
	},
	{
		createEmptyScene(errTagsID),
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

	db.Tag.On("FindBySceneID", testCtx, sceneID).Return(getTags(names), nil).Once()
	db.Tag.On("FindBySceneID", testCtx, noTagsID).Return(nil, nil).Once()
	db.Tag.On("FindBySceneID", testCtx, errTagsID).Return(nil, tagErr).Once()

	for i, s := range getTagNamesScenarios {
		scene := s.input
		json, err := GetTagNames(testCtx, db.Tag, &scene)

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

type sceneGroupsTestScenario struct {
	input    models.Scene
	expected []jsonschema.SceneGroup
	err      bool
}

var validGroups = models.NewRelatedGroups([]models.GroupsScenes{
	{
		GroupID:    validGroup1,
		SceneIndex: &group1Scene,
	},
	{
		GroupID:    validGroup2,
		SceneIndex: &group2Scene,
	},
})

var invalidGroups = models.NewRelatedGroups([]models.GroupsScenes{
	{
		GroupID:    invalidGroup,
		SceneIndex: &group1Scene,
	},
})

var getSceneGroupsJSONScenarios = []sceneGroupsTestScenario{
	{
		models.Scene{
			ID:     sceneID,
			Groups: validGroups,
		},
		[]jsonschema.SceneGroup{
			{
				GroupName:  group1Name,
				SceneIndex: group1Scene,
			},
			{
				GroupName:  group2Name,
				SceneIndex: group2Scene,
			},
		},
		false,
	},
	{
		models.Scene{
			ID:     noGroupsID,
			Groups: models.NewRelatedGroups([]models.GroupsScenes{}),
		},
		nil,
		false,
	},
	{
		models.Scene{
			ID:     errFindGroupID,
			Groups: invalidGroups,
		},
		nil,
		true,
	},
}

func TestGetSceneGroupsJSON(t *testing.T) {
	db := mocks.NewDatabase()

	groupErr := errors.New("error getting group")

	db.Group.On("Find", testCtx, validGroup1).Return(&models.Group{
		Name: group1Name,
	}, nil).Once()
	db.Group.On("Find", testCtx, validGroup2).Return(&models.Group{
		Name: group2Name,
	}, nil).Once()
	db.Group.On("Find", testCtx, invalidGroup).Return(nil, groupErr).Once()

	for i, s := range getSceneGroupsJSONScenarios {
		scene := s.input
		json, err := GetSceneGroupsJSON(testCtx, db.Group, &scene)

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

const (
	validMarkerID1 = 1
	validMarkerID2 = 2

	invalidMarkerID1 = 3
	invalidMarkerID2 = 4

	validTagID1 = 1
	validTagID2 = 2

	validTagName1 = "validTagName1"
	validTagName2 = "validTagName2"

	invalidTagID = 3

	markerTitle1 = "markerTitle1"
	markerTitle2 = "markerTitle2"

	markerSeconds1 = 1.0
	markerSeconds2 = 2.3

	markerSeconds1Str = "1.0"
	markerSeconds2Str = "2.3"
)

type sceneMarkersTestScenario struct {
	input    models.Scene
	expected []jsonschema.SceneMarker
	err      bool
}

var getSceneMarkersJSONScenarios = []sceneMarkersTestScenario{
	{
		createEmptyScene(sceneID),
		[]jsonschema.SceneMarker{
			{
				Title:      markerTitle1,
				PrimaryTag: validTagName1,
				Seconds:    markerSeconds1Str,
				Tags: []string{
					validTagName1,
					validTagName2,
				},
				CreatedAt: json.JSONTime{
					Time: createTime,
				},
				UpdatedAt: json.JSONTime{
					Time: updateTime,
				},
			},
			{
				Title:      markerTitle2,
				PrimaryTag: validTagName2,
				Seconds:    markerSeconds2Str,
				Tags: []string{
					validTagName2,
				},
				CreatedAt: json.JSONTime{
					Time: createTime,
				},
				UpdatedAt: json.JSONTime{
					Time: updateTime,
				},
			},
		},
		false,
	},
	{
		createEmptyScene(noMarkersID),
		nil,
		false,
	},
	{
		createEmptyScene(errMarkersID),
		nil,
		true,
	},
	{
		createEmptyScene(errFindPrimaryTagID),
		nil,
		true,
	},
	{
		createEmptyScene(errFindByMarkerID),
		nil,
		true,
	},
}

var validMarkers = []*models.SceneMarker{
	{
		ID:           validMarkerID1,
		Title:        markerTitle1,
		PrimaryTagID: validTagID1,
		Seconds:      markerSeconds1,
		CreatedAt:    createTime,
		UpdatedAt:    updateTime,
	},
	{
		ID:           validMarkerID2,
		Title:        markerTitle2,
		PrimaryTagID: validTagID2,
		Seconds:      markerSeconds2,
		CreatedAt:    createTime,
		UpdatedAt:    updateTime,
	},
}

var invalidMarkers1 = []*models.SceneMarker{
	{
		ID:           invalidMarkerID1,
		PrimaryTagID: invalidTagID,
	},
}

var invalidMarkers2 = []*models.SceneMarker{
	{
		ID:           invalidMarkerID2,
		PrimaryTagID: validTagID1,
	},
}

func TestGetSceneMarkersJSON(t *testing.T) {
	db := mocks.NewDatabase()

	markersErr := errors.New("error getting scene markers")
	tagErr := errors.New("error getting tags")

	db.SceneMarker.On("FindBySceneID", testCtx, sceneID).Return(validMarkers, nil).Once()
	db.SceneMarker.On("FindBySceneID", testCtx, noMarkersID).Return(nil, nil).Once()
	db.SceneMarker.On("FindBySceneID", testCtx, errMarkersID).Return(nil, markersErr).Once()
	db.SceneMarker.On("FindBySceneID", testCtx, errFindPrimaryTagID).Return(invalidMarkers1, nil).Once()
	db.SceneMarker.On("FindBySceneID", testCtx, errFindByMarkerID).Return(invalidMarkers2, nil).Once()

	db.Tag.On("Find", testCtx, validTagID1).Return(&models.Tag{
		Name: validTagName1,
	}, nil)
	db.Tag.On("Find", testCtx, validTagID2).Return(&models.Tag{
		Name: validTagName2,
	}, nil)
	db.Tag.On("Find", testCtx, invalidTagID).Return(nil, tagErr)

	db.Tag.On("FindBySceneMarkerID", testCtx, validMarkerID1).Return([]*models.Tag{
		{
			Name: validTagName1,
		},
		{
			Name: validTagName2,
		},
	}, nil)
	db.Tag.On("FindBySceneMarkerID", testCtx, validMarkerID2).Return([]*models.Tag{
		{
			Name: validTagName2,
		},
	}, nil)
	db.Tag.On("FindBySceneMarkerID", testCtx, invalidMarkerID2).Return(nil, tagErr).Once()

	for i, s := range getSceneMarkersJSONScenarios {
		scene := s.input
		json, err := GetSceneMarkersJSON(testCtx, db.SceneMarker, db.Tag, &scene)

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
