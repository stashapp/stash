package studio

import (
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/models/modelstest"
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
			Name:  studioName,
			Image: invalidImage,
		},
	}

	err := i.PreImport()

	assert.NotNil(t, err)

	i.Input.Image = image

	err = i.PreImport()

	assert.Nil(t, err)
}

func TestImporterPreImportWithParent(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Studio{
			Name:         studioName,
			Image:        image,
			ParentStudio: existingParentStudioName,
		},
	}

	readerWriter.On("FindByName", existingParentStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	readerWriter.On("FindByName", existingParentStudioErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, int64(existingStudioID), i.studio.ParentID.Int64)

	i.Input.ParentStudio = existingParentStudioErr
	err = i.PreImport()
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingParent(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Studio{
			Name:         studioName,
			Image:        image,
			ParentStudio: missingParentStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	readerWriter.On("FindByName", missingParentStudioName, false).Return(nil, nil).Times(3)
	readerWriter.On("Create", mock.AnythingOfType("models.Studio")).Return(&models.Studio{
		ID: existingStudioID,
	}, nil)

	err := i.PreImport()
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport()
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, int64(existingStudioID), i.studio.ParentID.Int64)

	readerWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingParentCreateErr(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Studio{
			Name:         studioName,
			Image:        image,
			ParentStudio: missingParentStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	readerWriter.On("FindByName", missingParentStudioName, false).Return(nil, nil).Once()
	readerWriter.On("Create", mock.AnythingOfType("models.Studio")).Return(nil, errors.New("Create error"))

	err := i.PreImport()
	assert.NotNil(t, err)
}

func TestImporterPostImport(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		imageData:    imageBytes,
	}

	updateStudioImageErr := errors.New("UpdateStudioImage error")

	readerWriter.On("UpdateStudioImage", studioID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateStudioImage", errImageID, imageBytes).Return(updateStudioImageErr).Once()

	err := i.PostImport(studioID)
	assert.Nil(t, err)

	err = i.PostImport(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Studio{
			Name: studioName,
		},
	}

	errFindByName := errors.New("FindByName error")
	readerWriter.On("FindByName", studioName, false).Return(nil, nil).Once()
	readerWriter.On("FindByName", existingStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	readerWriter.On("FindByName", studioNameErr, false).Return(nil, errFindByName).Once()

	id, err := i.FindExistingID()
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Name = existingStudioName
	id, err = i.FindExistingID()
	assert.Equal(t, existingStudioID, *id)
	assert.Nil(t, err)

	i.Input.Name = studioNameErr
	id, err = i.FindExistingID()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}

	studio := models.Studio{
		Name: modelstest.NullString(studioName),
	}

	studioErr := models.Studio{
		Name: modelstest.NullString(studioNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		studio:       studio,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", studio).Return(&models.Studio{
		ID: studioID,
	}, nil).Once()
	readerWriter.On("Create", studioErr).Return(nil, errCreate).Once()

	id, err := i.Create()
	assert.Equal(t, studioID, *id)
	assert.Nil(t, err)

	i.studio = studioErr
	id, err = i.Create()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.StudioReaderWriter{}

	studio := models.Studio{
		Name: modelstest.NullString(studioName),
	}

	studioErr := models.Studio{
		Name: modelstest.NullString(studioNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		studio:       studio,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	studio.ID = studioID
	readerWriter.On("UpdateFull", studio).Return(nil, nil).Once()

	err := i.Update(studioID)
	assert.Nil(t, err)

	i.studio = studioErr

	// need to set id separately
	studioErr.ID = errImageID
	readerWriter.On("UpdateFull", studioErr).Return(nil, errUpdate).Once()

	err = i.Update(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
