package scene

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const invalidImage = "aW1hZ2VCeXRlcw&&"

var (
	existingStudioID    = 101
	existingPerformerID = 103
	existingGroupID     = 104
	existingTagID       = 105

	existingStudioName = "existingStudioName"
	existingStudioErr  = "existingStudioErr"
	missingStudioName  = "missingStudioName"

	existingPerformerName = "existingPerformerName"
	existingPerformerErr  = "existingPerformerErr"
	missingPerformerName  = "missingPerformerName"

	existingGroupName = "existingGroupName"
	existingGroupErr  = "existingGroupErr"
	missingGroupName  = "missingGroupName"

	existingTagName = "existingTagName"
	existingTagErr  = "existingTagErr"
	missingTagName  = "missingTagName"
)

var testCtx = context.Background()

func TestImporterPreImport(t *testing.T) {
	var (
		title     = "title"
		code      = "code"
		details   = "details"
		director  = "director"
		endpoint1 = "endpoint1"
		stashID1  = "stashID1"
		endpoint2 = "endpoint2"
		stashID2  = "stashID2"
		url1      = "url1"
		url2      = "url2"
		rating    = 3
		organized = true

		createdAt = time.Now().Add(-time.Hour)
		updatedAt = time.Now().Add(-time.Minute)

		resumeTime   = 1.234
		playDuration = 2.345
	)
	tests := []struct {
		name   string
		input  jsonschema.Scene
		output models.Scene
	}{
		{
			"basic",
			jsonschema.Scene{
				Title:    title,
				Code:     code,
				Details:  details,
				Director: director,
				StashIDs: []models.StashID{
					{Endpoint: endpoint1, StashID: stashID1},
					{Endpoint: endpoint2, StashID: stashID2},
				},
				URLs:         []string{url1, url2},
				Rating:       rating,
				Organized:    organized,
				CreatedAt:    json.JSONTime{Time: createdAt},
				UpdatedAt:    json.JSONTime{Time: updatedAt},
				ResumeTime:   resumeTime,
				PlayDuration: playDuration,
			},
			models.Scene{
				Title:    title,
				Code:     code,
				Details:  details,
				Director: director,
				StashIDs: models.NewRelatedStashIDs([]models.StashID{
					{Endpoint: endpoint1, StashID: stashID1},
					{Endpoint: endpoint2, StashID: stashID2},
				}),
				URLs:         models.NewRelatedStrings([]string{url1, url2}),
				Rating:       &rating,
				Organized:    organized,
				CreatedAt:    createdAt.Truncate(0),
				UpdatedAt:    updatedAt.Truncate(0),
				ResumeTime:   resumeTime,
				PlayDuration: playDuration,

				Files:        models.NewRelatedVideoFiles([]*models.VideoFile{}),
				GalleryIDs:   models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				Groups:       models.NewRelatedGroups([]models.GroupsScenes{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Importer{
				Input: tt.input,
			}

			if err := i.PreImport(testCtx); err != nil {
				t.Errorf("PreImport() error = %v", err)
				return
			}

			assert.Equal(t, tt.output, i.scene)
		})
	}
}

func truncateTimes(t []time.Time) []time.Time {
	return sliceutil.Map(t, func(t time.Time) time.Time { return t.Truncate(0) })
}

func TestImporterPreImportHistory(t *testing.T) {
	var (
		playTime1 = time.Now().Add(-time.Hour * 2)
		playTime2 = time.Now().Add(-time.Minute * 2)
		oTime1    = time.Now().Add(-time.Hour * 3)
		oTime2    = time.Now().Add(-time.Minute * 3)
	)
	tests := []struct {
		name                string
		input               jsonschema.Scene
		expectedPlayHistory []time.Time
		expectedOHistory    []time.Time
	}{
		{
			"basic",
			jsonschema.Scene{
				PlayHistory: []json.JSONTime{
					{Time: playTime1},
					{Time: playTime2},
				},
				OHistory: []json.JSONTime{
					{Time: oTime1},
					{Time: oTime2},
				},
			},
			[]time.Time{playTime1, playTime2},
			[]time.Time{oTime1, oTime2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Importer{
				Input: tt.input,
			}

			if err := i.PreImport(testCtx); err != nil {
				t.Errorf("PreImport() error = %v", err)
				return
			}

			// convert histories to unix timestamps for comparison
			eph := truncateTimes(tt.expectedPlayHistory)
			vh := truncateTimes(i.viewHistory)

			eoh := truncateTimes(tt.expectedOHistory)
			oh := truncateTimes(i.oHistory)

			assert.Equal(t, eph, vh, "view history mismatch")
			assert.Equal(t, eoh, oh, "o history mismatch")
		})
	}
}

func TestImporterPreImportCoverImage(t *testing.T) {
	i := Importer{
		Input: jsonschema.Scene{
			Cover: invalidImage,
		},
	}

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.Input.Cover = imageBase64

	err = i.PreImport(testCtx)
	assert.Nil(t, err)
}

func TestImporterPreImportWithStudio(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		StudioWriter: db.Studio,
		Input: jsonschema.Scene{
			Studio: existingStudioName,
		},
	}

	db.Studio.On("FindByName", testCtx, existingStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	db.Studio.On("FindByName", testCtx, existingStudioErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingStudioID, *i.scene.StudioID)

	i.Input.Studio = existingStudioErr
	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingStudio(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		StudioWriter: db.Studio,
		Input: jsonschema.Scene{
			Studio: missingStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	db.Studio.On("FindByName", testCtx, missingStudioName, false).Return(nil, nil).Times(3)
	db.Studio.On("Create", testCtx, mock.AnythingOfType("*models.Studio")).Run(func(args mock.Arguments) {
		s := args.Get(1).(*models.Studio)
		s.ID = existingStudioID
	}).Return(nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingStudioID, *i.scene.StudioID)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingStudioCreateErr(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		StudioWriter: db.Studio,
		Input: jsonschema.Scene{
			Studio: missingStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	db.Studio.On("FindByName", testCtx, missingStudioName, false).Return(nil, nil).Once()
	db.Studio.On("Create", testCtx, mock.AnythingOfType("*models.Studio")).Return(errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithPerformer(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		PerformerWriter:     db.Performer,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Scene{
			Performers: []string{
				existingPerformerName,
			},
		},
	}

	db.Performer.On("FindByNames", testCtx, []string{existingPerformerName}, false).Return([]*models.Performer{
		{
			ID:   existingPerformerID,
			Name: existingPerformerName,
		},
	}, nil).Once()
	db.Performer.On("FindByNames", testCtx, []string{existingPerformerErr}, false).Return(nil, errors.New("FindByNames error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, []int{existingPerformerID}, i.scene.PerformerIDs.List())

	i.Input.Performers = []string{existingPerformerErr}
	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingPerformer(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		PerformerWriter: db.Performer,
		Input: jsonschema.Scene{
			Performers: []string{
				missingPerformerName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	db.Performer.On("FindByNames", testCtx, []string{missingPerformerName}, false).Return(nil, nil).Times(3)
	db.Performer.On("Create", testCtx, mock.AnythingOfType("*models.CreatePerformerInput")).Run(func(args mock.Arguments) {
		p := args.Get(1).(*models.CreatePerformerInput)
		p.ID = existingPerformerID
	}).Return(nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, []int{existingPerformerID}, i.scene.PerformerIDs.List())

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingPerformerCreateErr(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		PerformerWriter: db.Performer,
		Input: jsonschema.Scene{
			Performers: []string{
				missingPerformerName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	db.Performer.On("FindByNames", testCtx, []string{missingPerformerName}, false).Return(nil, nil).Once()
	db.Performer.On("Create", testCtx, mock.AnythingOfType("*models.CreatePerformerInput")).Return(errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithGroup(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		GroupWriter:         db.Group,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Scene{
			Groups: []jsonschema.SceneGroup{
				{
					GroupName:  existingGroupName,
					SceneIndex: 1,
				},
			},
		},
	}

	db.Group.On("FindByName", testCtx, existingGroupName, false).Return(&models.Group{
		ID:   existingGroupID,
		Name: existingGroupName,
	}, nil).Once()
	db.Group.On("FindByName", testCtx, existingGroupErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingGroupID, i.scene.Groups.List()[0].GroupID)

	i.Input.Groups[0].GroupName = existingGroupErr
	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingGroup(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		GroupWriter: db.Group,
		Input: jsonschema.Scene{
			Groups: []jsonschema.SceneGroup{
				{
					GroupName: missingGroupName,
				},
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	db.Group.On("FindByName", testCtx, missingGroupName, false).Return(nil, nil).Times(3)
	db.Group.On("Create", testCtx, mock.AnythingOfType("*models.Group")).Run(func(args mock.Arguments) {
		m := args.Get(1).(*models.Group)
		m.ID = existingGroupID
	}).Return(nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingGroupID, i.scene.Groups.List()[0].GroupID)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingGroupCreateErr(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		GroupWriter: db.Group,
		Input: jsonschema.Scene{
			Groups: []jsonschema.SceneGroup{
				{
					GroupName: missingGroupName,
				},
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	db.Group.On("FindByName", testCtx, missingGroupName, false).Return(nil, nil).Once()
	db.Group.On("Create", testCtx, mock.AnythingOfType("*models.Group")).Return(errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithTag(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		TagWriter:           db.Tag,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Scene{
			Tags: []string{
				existingTagName,
			},
		},
	}

	db.Tag.On("FindByNames", testCtx, []string{existingTagName}, false).Return([]*models.Tag{
		{
			ID:   existingTagID,
			Name: existingTagName,
		},
	}, nil).Once()
	db.Tag.On("FindByNames", testCtx, []string{existingTagErr}, false).Return(nil, errors.New("FindByNames error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, []int{existingTagID}, i.scene.TagIDs.List())

	i.Input.Tags = []string{existingTagErr}
	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingTag(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		TagWriter: db.Tag,
		Input: jsonschema.Scene{
			Tags: []string{
				missingTagName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	db.Tag.On("FindByNames", testCtx, []string{missingTagName}, false).Return(nil, nil).Times(3)
	db.Tag.On("Create", testCtx, mock.AnythingOfType("*models.Tag")).Run(func(args mock.Arguments) {
		t := args.Get(1).(*models.Tag)
		t.ID = existingTagID
	}).Return(nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, []int{existingTagID}, i.scene.TagIDs.List())

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingTagCreateErr(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		TagWriter: db.Tag,
		Input: jsonschema.Scene{
			Tags: []string{
				missingTagName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	db.Tag.On("FindByNames", testCtx, []string{missingTagName}, false).Return(nil, nil).Once()
	db.Tag.On("Create", testCtx, mock.AnythingOfType("*models.Tag")).Return(errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}
