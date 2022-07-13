package scene

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/stashapp/stash/pkg/models"
// 	"github.com/stashapp/stash/pkg/models/jsonschema"
// 	"github.com/stashapp/stash/pkg/models/mocks"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// const (
// 	seconds      = "5"
// 	secondsFloat = 5.0
// 	errSceneID   = 999
// )

// func TestMarkerImporterName(t *testing.T) {
// 	i := MarkerImporter{
// 		Input: jsonschema.SceneMarker{
// 			Title:   title,
// 			Seconds: seconds,
// 		},
// 	}

// 	assert.Equal(t, title+" (5)", i.Name())
// }

// func TestMarkerImporterPreImportWithTag(t *testing.T) {
// 	tagReaderWriter := &mocks.TagReaderWriter{}
// 	ctx := context.Background()

// 	i := MarkerImporter{
// 		TagWriter:           tagReaderWriter,
// 		MissingRefBehaviour: models.ImportMissingRefEnumFail,
// 		Input: jsonschema.SceneMarker{
// 			PrimaryTag: existingTagName,
// 		},
// 	}

// 	tagReaderWriter.On("FindByNames", ctx, []string{existingTagName}, false).Return([]*models.Tag{
// 		{
// 			ID:   existingTagID,
// 			Name: existingTagName,
// 		},
// 	}, nil).Times(4)
// 	tagReaderWriter.On("FindByNames", ctx, []string{existingTagErr}, false).Return(nil, errors.New("FindByNames error")).Times(2)

// 	err := i.PreImport(ctx)
// 	assert.Nil(t, err)
// 	assert.Equal(t, existingTagID, i.marker.PrimaryTagID)

// 	i.Input.PrimaryTag = existingTagErr
// 	err = i.PreImport(ctx)
// 	assert.NotNil(t, err)

// 	i.Input.PrimaryTag = existingTagName
// 	i.Input.Tags = []string{
// 		existingTagName,
// 	}
// 	err = i.PreImport(ctx)
// 	assert.Nil(t, err)
// 	assert.Equal(t, existingTagID, i.tags[0].ID)

// 	i.Input.Tags[0] = existingTagErr
// 	err = i.PreImport(ctx)
// 	assert.NotNil(t, err)

// 	tagReaderWriter.AssertExpectations(t)
// }

// func TestMarkerImporterPostImportUpdateTags(t *testing.T) {
// 	sceneMarkerReaderWriter := &mocks.SceneMarkerReaderWriter{}
// 	ctx := context.Background()

// 	i := MarkerImporter{
// 		ReaderWriter: sceneMarkerReaderWriter,
// 		tags: []*models.Tag{
// 			{
// 				ID: existingTagID,
// 			},
// 		},
// 	}

// 	updateErr := errors.New("UpdateTags error")

// 	sceneMarkerReaderWriter.On("UpdateTags", ctx, sceneID, []int{existingTagID}).Return(nil).Once()
// 	sceneMarkerReaderWriter.On("UpdateTags", ctx, errTagsID, mock.AnythingOfType("[]int")).Return(updateErr).Once()

// 	err := i.PostImport(ctx, sceneID)
// 	assert.Nil(t, err)

// 	err = i.PostImport(ctx, errTagsID)
// 	assert.NotNil(t, err)

// 	sceneMarkerReaderWriter.AssertExpectations(t)
// }

// func TestMarkerImporterFindExistingID(t *testing.T) {
// 	readerWriter := &mocks.SceneMarkerReaderWriter{}
// 	ctx := context.Background()

// 	i := MarkerImporter{
// 		ReaderWriter: readerWriter,
// 		SceneID:      sceneID,
// 		marker: models.SceneMarker{
// 			Seconds: secondsFloat,
// 		},
// 	}

// 	expectedErr := errors.New("FindBy* error")
// 	readerWriter.On("FindBySceneID", ctx, sceneID).Return([]*models.SceneMarker{
// 		{
// 			ID:      existingSceneID,
// 			Seconds: secondsFloat,
// 		},
// 	}, nil).Times(2)
// 	readerWriter.On("FindBySceneID", ctx, errSceneID).Return(nil, expectedErr).Once()

// 	id, err := i.FindExistingID(ctx)
// 	assert.Equal(t, existingSceneID, *id)
// 	assert.Nil(t, err)

// 	i.marker.Seconds++
// 	id, err = i.FindExistingID(ctx)
// 	assert.Nil(t, id)
// 	assert.Nil(t, err)

// 	i.SceneID = errSceneID
// 	id, err = i.FindExistingID(ctx)
// 	assert.Nil(t, id)
// 	assert.NotNil(t, err)

// 	readerWriter.AssertExpectations(t)
// }

// func TestMarkerImporterCreate(t *testing.T) {
// 	readerWriter := &mocks.SceneMarkerReaderWriter{}
// 	ctx := context.Background()

// 	scene := models.SceneMarker{
// 		Title: title,
// 	}

// 	sceneErr := models.SceneMarker{
// 		Title: sceneNameErr,
// 	}

// 	i := MarkerImporter{
// 		ReaderWriter: readerWriter,
// 		marker:       scene,
// 	}

// 	errCreate := errors.New("Create error")
// 	readerWriter.On("Create", ctx, scene).Return(&models.SceneMarker{
// 		ID: sceneID,
// 	}, nil).Once()
// 	readerWriter.On("Create", ctx, sceneErr).Return(nil, errCreate).Once()

// 	id, err := i.Create(ctx)
// 	assert.Equal(t, sceneID, *id)
// 	assert.Nil(t, err)

// 	i.marker = sceneErr
// 	id, err = i.Create(ctx)
// 	assert.Nil(t, id)
// 	assert.NotNil(t, err)

// 	readerWriter.AssertExpectations(t)
// }

// func TestMarkerImporterUpdate(t *testing.T) {
// 	readerWriter := &mocks.SceneMarkerReaderWriter{}
// 	ctx := context.Background()

// 	scene := models.SceneMarker{
// 		Title: title,
// 	}

// 	sceneErr := models.SceneMarker{
// 		Title: sceneNameErr,
// 	}

// 	i := MarkerImporter{
// 		ReaderWriter: readerWriter,
// 		marker:       scene,
// 	}

// 	errUpdate := errors.New("Update error")

// 	// id needs to be set for the mock input
// 	scene.ID = sceneID
// 	readerWriter.On("Update", ctx, scene).Return(nil, nil).Once()

// 	err := i.Update(ctx, sceneID)
// 	assert.Nil(t, err)

// 	i.marker = sceneErr

// 	// need to set id separately
// 	sceneErr.ID = errImageID
// 	readerWriter.On("Update", ctx, sceneErr).Return(nil, errUpdate).Once()

// 	err = i.Update(ctx, errImageID)
// 	assert.NotNil(t, err)

// 	readerWriter.AssertExpectations(t)
// }
