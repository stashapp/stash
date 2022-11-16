package studio

import (
	"context"
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const invalidImage = "aW1hZ2VCeXRlcw&&"

const (
	studioNameErr      = "studioNameErr"
	existingStudioName = "existingTagName"

	existingStudioID = 100

	existingParentStudioName = "existingParentStudioName"
	existingParentStudioErr  = "existingParentStudioErr"
	missingParentStudioName  = "existingParentStudioName"
)

func TestImporterName(t *testing.T) {
	i := Importer{
		Input: jsonschema.Studio{
			Name: studioName,
		},
	}

	assert.Equal(t, studioName, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Input: jsonschema.Studio{
			Name:          studioName,
			Image:         invalidImage,
			IgnoreAutoTag: autoTagIgnored,
		},
	}
	ctx := context.Background()

	err := i.PreImport(ctx)

	assert.NotNil(t, err)

	i.Input.Image = image

	err = i.PreImport(ctx)

	assert.Nil(t, err)

	i.Input = *createFullJSONStudio(studioName, image, []string{"alias"})
	i.Input.ParentStudio = ""

	err = i.PreImport(ctx)

	assert.Nil(t, err)
	expectedStudio := createFullStudio(0, 0)
	expectedStudio.ParentID.Valid = false
	expectedStudio.Checksum = md5.FromString(studioName)
	assert.Equal(t, expectedStudio, i.studio)
}

func TestImporterPreImportWithParent(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}
	ctx := context.Background()

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Studio{
			Name:         studioName,
			Image:        image,
			ParentStudio: existingParentStudioName,
		},
	}

	readerWriter.On("FindByName", ctx, existingParentStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	readerWriter.On("FindByName", ctx, existingParentStudioErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport(ctx)
	assert.Nil(t, err)
	assert.Equal(t, int64(existingStudioID), i.studio.ParentID.Int64)

	i.Input.ParentStudio = existingParentStudioErr
	err = i.PreImport(ctx)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingParent(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}
	ctx := context.Background()

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Studio{
			Name:         studioName,
			Image:        image,
			ParentStudio: missingParentStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	readerWriter.On("FindByName", ctx, missingParentStudioName, false).Return(nil, nil).Times(3)
	readerWriter.On("Create", ctx, mock.AnythingOfType("models.Studio")).Return(&models.Studio{
		ID: existingStudioID,
	}, nil)

	err := i.PreImport(ctx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(ctx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(ctx)
	assert.Nil(t, err)
	assert.Equal(t, int64(existingStudioID), i.studio.ParentID.Int64)

	readerWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingParentCreateErr(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}
	ctx := context.Background()

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Studio{
			Name:         studioName,
			Image:        image,
			ParentStudio: missingParentStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	readerWriter.On("FindByName", ctx, missingParentStudioName, false).Return(nil, nil).Once()
	readerWriter.On("Create", ctx, mock.AnythingOfType("models.Studio")).Return(nil, errors.New("Create error"))

	err := i.PreImport(ctx)
	assert.NotNil(t, err)
}

func TestImporterPostImport(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}
	ctx := context.Background()

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Studio{
			Aliases: []string{"alias"},
		},
		imageData: imageBytes,
	}

	updateStudioImageErr := errors.New("UpdateImage error")
	updateTagAliasErr := errors.New("UpdateAlias error")

	readerWriter.On("UpdateImage", ctx, studioID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateImage", ctx, errImageID, imageBytes).Return(updateStudioImageErr).Once()
	readerWriter.On("UpdateImage", ctx, errAliasID, imageBytes).Return(nil).Once()

	readerWriter.On("UpdateAliases", ctx, studioID, i.Input.Aliases).Return(nil).Once()
	readerWriter.On("UpdateAliases", ctx, errImageID, i.Input.Aliases).Return(nil).Maybe()
	readerWriter.On("UpdateAliases", ctx, errAliasID, i.Input.Aliases).Return(updateTagAliasErr).Once()

	err := i.PostImport(ctx, studioID)
	assert.Nil(t, err)

	err = i.PostImport(ctx, errImageID)
	assert.NotNil(t, err)

	err = i.PostImport(ctx, errAliasID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}
	ctx := context.Background()

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Studio{
			Name: studioName,
		},
	}

	errFindByName := errors.New("FindByName error")
	readerWriter.On("FindByName", ctx, studioName, false).Return(nil, nil).Once()
	readerWriter.On("FindByName", ctx, existingStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	readerWriter.On("FindByName", ctx, studioNameErr, false).Return(nil, errFindByName).Once()

	id, err := i.FindExistingID(ctx)
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Name = existingStudioName
	id, err = i.FindExistingID(ctx)
	assert.Equal(t, existingStudioID, *id)
	assert.Nil(t, err)

	i.Input.Name = studioNameErr
	id, err = i.FindExistingID(ctx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}
	ctx := context.Background()

	studio := models.Studio{
		Name: models.NullString(studioName),
	}

	studioErr := models.Studio{
		Name: models.NullString(studioNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		studio:       studio,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", ctx, studio).Return(&models.Studio{
		ID: studioID,
	}, nil).Once()
	readerWriter.On("Create", ctx, studioErr).Return(nil, errCreate).Once()

	id, err := i.Create(ctx)
	assert.Equal(t, studioID, *id)
	assert.Nil(t, err)

	i.studio = studioErr
	id, err = i.Create(ctx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}
	ctx := context.Background()

	studio := models.Studio{
		Name: models.NullString(studioName),
	}

	studioErr := models.Studio{
		Name: models.NullString(studioNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		studio:       studio,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	studio.ID = studioID
	readerWriter.On("UpdateFull", ctx, studio).Return(nil, nil).Once()

	err := i.Update(ctx, studioID)
	assert.Nil(t, err)

	i.studio = studioErr

	// need to set id separately
	studioErr.ID = errImageID
	readerWriter.On("UpdateFull", ctx, studioErr).Return(nil, errUpdate).Once()

	err = i.Update(ctx, errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
