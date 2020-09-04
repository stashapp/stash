package performer

import (
	"errors"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/models/modelstest"
	"github.com/stashapp/stash/pkg/utils"
	"github.com/stretchr/testify/assert"

	"testing"
)

const invalidImage = "aW1hZ2VCeXRlcw&&"

const (
	existingPerformerID = 100

	existingPerformerName = "existingPerformerName"
	performerNameErr      = "performerNameErr"
)

func TestImporterName(t *testing.T) {
	i := Importer{
		Input: jsonschema.Performer{
			Name: performerName,
		},
	}

	assert.Equal(t, performerName, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Input: jsonschema.Performer{
			Name:  performerName,
			Image: invalidImage,
		},
	}

	err := i.PreImport()

	assert.NotNil(t, err)

	i.Input = *createFullJSONPerformer(performerName, image)

	err = i.PreImport()

	assert.Nil(t, err)
	expectedPerformer := *createFullPerformer(0, performerName)
	expectedPerformer.Checksum = utils.MD5FromString(performerName)
	assert.Equal(t, expectedPerformer, i.performer)
}

func TestImporterPostImport(t *testing.T) {
	readerWriter := &mocks.PerformerReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		imageData:    imageBytes,
	}

	updatePerformerImageErr := errors.New("UpdatePerformerImage error")

	readerWriter.On("UpdatePerformerImage", performerID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdatePerformerImage", errImageID, imageBytes).Return(updatePerformerImageErr).Once()

	err := i.PostImport(performerID)
	assert.Nil(t, err)

	err = i.PostImport(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	readerWriter := &mocks.PerformerReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Performer{
			Name: performerName,
		},
	}

	errFindByNames := errors.New("FindByNames error")
	readerWriter.On("FindByNames", []string{performerName}, false).Return(nil, nil).Once()
	readerWriter.On("FindByNames", []string{existingPerformerName}, false).Return([]*models.Performer{
		{
			ID: existingPerformerID,
		},
	}, nil).Once()
	readerWriter.On("FindByNames", []string{performerNameErr}, false).Return(nil, errFindByNames).Once()

	id, err := i.FindExistingID()
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Name = existingPerformerName
	id, err = i.FindExistingID()
	assert.Equal(t, existingPerformerID, *id)
	assert.Nil(t, err)

	i.Input.Name = performerNameErr
	id, err = i.FindExistingID()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.PerformerReaderWriter{}

	performer := models.Performer{
		Name: modelstest.NullString(performerName),
	}

	performerErr := models.Performer{
		Name: modelstest.NullString(performerNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		performer:    performer,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", performer).Return(&models.Performer{
		ID: performerID,
	}, nil).Once()
	readerWriter.On("Create", performerErr).Return(nil, errCreate).Once()

	id, err := i.Create()
	assert.Equal(t, performerID, *id)
	assert.Nil(t, err)

	i.performer = performerErr
	id, err = i.Create()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.PerformerReaderWriter{}

	performer := models.Performer{
		Name: modelstest.NullString(performerName),
	}

	performerErr := models.Performer{
		Name: modelstest.NullString(performerNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		performer:    performer,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	performer.ID = performerID
	readerWriter.On("Update", performer).Return(nil, nil).Once()

	err := i.Update(performerID)
	assert.Nil(t, err)

	i.performer = performerErr

	// need to set id separately
	performerErr.ID = errImageID
	readerWriter.On("Update", performerErr).Return(nil, errUpdate).Once()

	err = i.Update(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
