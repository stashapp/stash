package performer

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/utils"
	"github.com/stretchr/testify/assert"

	"testing"
)

const invalidImage = "aW1hZ2VCeXRlcw&&"

const (
	existingPerformerID = 100
	existingTagID       = 105
	errTagsID           = 106

	existingPerformerName = "existingPerformerName"
	performerNameErr      = "performerNameErr"

	existingTagName = "existingTagName"
	existingTagErr  = "existingTagErr"
	missingTagName  = "missingTagName"
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

func TestImporterPreImportWithTag(t *testing.T) {
	tagReaderWriter := &mocks.TagReaderWriter{}

	i := Importer{
		TagWriter:           tagReaderWriter,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Performer{
			Tags: []string{
				existingTagName,
			},
		},
	}

	tagReaderWriter.On("FindByNames", []string{existingTagName}, false).Return([]*models.Tag{
		{
			ID:   existingTagID,
			Name: existingTagName,
		},
	}, nil).Once()
	tagReaderWriter.On("FindByNames", []string{existingTagErr}, false).Return(nil, errors.New("FindByNames error")).Once()

	err := i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, existingTagID, i.tags[0].ID)

	i.Input.Tags = []string{existingTagErr}
	err = i.PreImport()
	assert.NotNil(t, err)

	tagReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingTag(t *testing.T) {
	tagReaderWriter := &mocks.TagReaderWriter{}

	i := Importer{
		TagWriter: tagReaderWriter,
		Input: jsonschema.Performer{
			Tags: []string{
				missingTagName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	tagReaderWriter.On("FindByNames", []string{missingTagName}, false).Return(nil, nil).Times(3)
	tagReaderWriter.On("Create", mock.AnythingOfType("models.Tag")).Return(&models.Tag{
		ID: existingTagID,
	}, nil)

	err := i.PreImport()
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport()
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, existingTagID, i.tags[0].ID)

	tagReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingTagCreateErr(t *testing.T) {
	tagReaderWriter := &mocks.TagReaderWriter{}

	i := Importer{
		TagWriter: tagReaderWriter,
		Input: jsonschema.Performer{
			Tags: []string{
				missingTagName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	tagReaderWriter.On("FindByNames", []string{missingTagName}, false).Return(nil, nil).Once()
	tagReaderWriter.On("Create", mock.AnythingOfType("models.Tag")).Return(nil, errors.New("Create error"))

	err := i.PreImport()
	assert.NotNil(t, err)
}

func TestImporterPostImport(t *testing.T) {
	readerWriter := &mocks.PerformerReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		imageData:    imageBytes,
	}

	updatePerformerImageErr := errors.New("UpdateImage error")

	readerWriter.On("UpdateImage", performerID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateImage", errImageID, imageBytes).Return(updatePerformerImageErr).Once()

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

func TestImporterPostImportUpdateTags(t *testing.T) {
	readerWriter := &mocks.PerformerReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		tags: []*models.Tag{
			{
				ID: existingTagID,
			},
		},
	}

	updateErr := errors.New("UpdateTags error")

	readerWriter.On("UpdateTags", performerID, []int{existingTagID}).Return(nil).Once()
	readerWriter.On("UpdateTags", errTagsID, mock.AnythingOfType("[]int")).Return(updateErr).Once()

	err := i.PostImport(performerID)
	assert.Nil(t, err)

	err = i.PostImport(errTagsID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.PerformerReaderWriter{}

	performer := models.Performer{
		Name: models.NullString(performerName),
	}

	performerErr := models.Performer{
		Name: models.NullString(performerNameErr),
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
		Name: models.NullString(performerName),
	}

	performerErr := models.Performer{
		Name: models.NullString(performerNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		performer:    performer,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	performer.ID = performerID
	readerWriter.On("UpdateFull", performer).Return(nil, nil).Once()

	err := i.Update(performerID)
	assert.Nil(t, err)

	i.performer = performerErr

	// need to set id separately
	performerErr.ID = errImageID
	readerWriter.On("UpdateFull", performerErr).Return(nil, errUpdate).Once()

	err = i.Update(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
