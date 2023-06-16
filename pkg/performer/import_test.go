package performer

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
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

var testCtx = context.Background()

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

	err := i.PreImport(testCtx)

	assert.NotNil(t, err)

	i.Input = *createFullJSONPerformer(performerName, image)

	err = i.PreImport(testCtx)

	assert.Nil(t, err)
	expectedPerformer := *createFullPerformer(0, performerName)
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

	tagReaderWriter.On("FindByNames", testCtx, []string{existingTagName}, false).Return([]*models.Tag{
		{
			ID:   existingTagID,
			Name: existingTagName,
		},
	}, nil).Once()
	tagReaderWriter.On("FindByNames", testCtx, []string{existingTagErr}, false).Return(nil, errors.New("FindByNames error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingTagID, i.performer.TagIDs.List()[0])

	i.Input.Tags = []string{existingTagErr}
	err = i.PreImport(testCtx)
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

	tagReaderWriter.On("FindByNames", testCtx, []string{missingTagName}, false).Return(nil, nil).Times(3)
	tagReaderWriter.On("Create", testCtx, mock.AnythingOfType("*models.Tag")).Run(func(args mock.Arguments) {
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
	assert.Equal(t, existingTagID, i.performer.TagIDs.List()[0])

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

	tagReaderWriter.On("FindByNames", testCtx, []string{missingTagName}, false).Return(nil, nil).Once()
	tagReaderWriter.On("Create", testCtx, mock.AnythingOfType("*models.Tag")).Return(errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)
}

func TestImporterPostImport(t *testing.T) {
	readerWriter := &mocks.PerformerReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		imageData:    imageBytes,
	}

	updatePerformerImageErr := errors.New("UpdateImage error")

	readerWriter.On("UpdateImage", testCtx, performerID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateImage", testCtx, errImageID, imageBytes).Return(updatePerformerImageErr).Once()

	err := i.PostImport(testCtx, performerID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errImageID)
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

	pp := 1
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	performerFilter := func(name string) *models.PerformerFilterType {
		return &models.PerformerFilterType{
			Name: &models.StringCriterionInput{
				Value:    name,
				Modifier: models.CriterionModifierEquals,
			},
		}
	}

	errFindByNames := errors.New("FindByNames error")
	readerWriter.On("Query", testCtx, performerFilter(performerName), findFilter).Return(nil, 0, nil).Once()
	readerWriter.On("Query", testCtx, performerFilter(existingPerformerName), findFilter).Return([]*models.Performer{
		{
			ID: existingPerformerID,
		},
	}, 1, nil).Once()
	readerWriter.On("Query", testCtx, performerFilter(performerNameErr), findFilter).Return(nil, 0, errFindByNames).Once()

	id, err := i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Name = existingPerformerName
	id, err = i.FindExistingID(testCtx)
	assert.Equal(t, existingPerformerID, *id)
	assert.Nil(t, err)

	i.Input.Name = performerNameErr
	id, err = i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.PerformerReaderWriter{}

	performer := models.Performer{
		Name: performerName,
	}

	performerErr := models.Performer{
		Name: performerNameErr,
	}

	i := Importer{
		ReaderWriter: readerWriter,
		performer:    performer,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", testCtx, &performer).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.Performer)
		arg.ID = performerID
	}).Return(nil).Once()
	readerWriter.On("Create", testCtx, &performerErr).Return(errCreate).Once()

	id, err := i.Create(testCtx)
	assert.Equal(t, performerID, *id)
	assert.Nil(t, err)

	i.performer = performerErr
	id, err = i.Create(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.PerformerReaderWriter{}

	performer := models.Performer{
		Name: performerName,
	}

	performerErr := models.Performer{
		Name: performerNameErr,
	}

	i := Importer{
		ReaderWriter: readerWriter,
		performer:    performer,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	performer.ID = performerID
	readerWriter.On("Update", testCtx, &performer).Return(nil).Once()

	err := i.Update(testCtx, performerID)
	assert.Nil(t, err)

	i.performer = performerErr

	// need to set id separately
	performerErr.ID = errImageID
	readerWriter.On("Update", testCtx, &performerErr).Return(errUpdate).Once()

	err = i.Update(testCtx, errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
