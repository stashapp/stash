package tag

import (
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

const image = "aW1hZ2VCeXRlcw=="
const invalidImage = "aW1hZ2VCeXRlcw&&"

var imageBytes = []byte("imageBytes")

const (
	tagNameErr      = "tagNameErr"
	existingTagName = "existingTagName"

	existingTagID = 100
)

func TestImporterName(t *testing.T) {
	i := Importer{
		Input: jsonschema.Tag{
			Name: tagName,
		},
	}

	assert.Equal(t, tagName, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Input: jsonschema.Tag{
			Name:  tagName,
			Image: invalidImage,
		},
	}

	err := i.PreImport()

	assert.NotNil(t, err)

	i.Input.Image = image

	err = i.PreImport()

	assert.Nil(t, err)
}

func TestImporterPostImport(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		imageData:    imageBytes,
	}

	updateTagImageErr := errors.New("UpdateTagImage error")

	readerWriter.On("UpdateTagImage", tagID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateTagImage", errImageID, imageBytes).Return(updateTagImageErr).Once()

	err := i.PostImport(tagID)
	assert.Nil(t, err)

	err = i.PostImport(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Tag{
			Name: tagName,
		},
	}

	errFindByName := errors.New("FindByName error")
	readerWriter.On("FindByName", tagName, false).Return(nil, nil).Once()
	readerWriter.On("FindByName", existingTagName, false).Return(&models.Tag{
		ID: existingTagID,
	}, nil).Once()
	readerWriter.On("FindByName", tagNameErr, false).Return(nil, errFindByName).Once()

	id, err := i.FindExistingID()
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Name = existingTagName
	id, err = i.FindExistingID()
	assert.Equal(t, existingTagID, *id)
	assert.Nil(t, err)

	i.Input.Name = tagNameErr
	id, err = i.FindExistingID()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}

	tag := models.Tag{
		Name: tagName,
	}

	tagErr := models.Tag{
		Name: tagNameErr,
	}

	i := Importer{
		ReaderWriter: readerWriter,
		tag:          tag,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", tag).Return(&models.Tag{
		ID: tagID,
	}, nil).Once()
	readerWriter.On("Create", tagErr).Return(nil, errCreate).Once()

	id, err := i.Create()
	assert.Equal(t, tagID, *id)
	assert.Nil(t, err)

	i.tag = tagErr
	id, err = i.Create()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}

	tag := models.Tag{
		Name: tagName,
	}

	tagErr := models.Tag{
		Name: tagNameErr,
	}

	i := Importer{
		ReaderWriter: readerWriter,
		tag:          tag,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	tag.ID = tagID
	readerWriter.On("Update", tag).Return(nil, nil).Once()

	err := i.Update(tagID)
	assert.Nil(t, err)

	i.tag = tagErr

	// need to set id separately
	tagErr.ID = errImageID
	readerWriter.On("Update", tagErr).Return(nil, errUpdate).Once()

	err = i.Update(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
