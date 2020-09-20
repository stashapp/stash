package gallery

import (
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

const (
	galleryPath         = "galleryPath"
	galleryPathErr      = "galleryPathErr"
	existingGalleryPath = "existingGalleryPath"

	galleryID         = 1
	idErr             = 2
	existingGalleryID = 100
)

func TestImporterName(t *testing.T) {
	i := Importer{
		Input: jsonschema.PathMapping{
			Path: galleryPath,
		},
	}

	assert.Equal(t, galleryPath, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Input: jsonschema.PathMapping{
			Path: galleryPath,
		},
	}

	err := i.PreImport()
	assert.Nil(t, err)
}

func TestImporterFindExistingID(t *testing.T) {
	readerWriter := &mocks.GalleryReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.PathMapping{
			Path: galleryPath,
		},
	}

	errFindByPath := errors.New("FindByPath error")
	readerWriter.On("FindByPath", galleryPath).Return(nil, nil).Once()
	readerWriter.On("FindByPath", existingGalleryPath).Return(&models.Gallery{
		ID: existingGalleryID,
	}, nil).Once()
	readerWriter.On("FindByPath", galleryPathErr).Return(nil, errFindByPath).Once()

	id, err := i.FindExistingID()
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Path = existingGalleryPath
	id, err = i.FindExistingID()
	assert.Equal(t, existingGalleryID, *id)
	assert.Nil(t, err)

	i.Input.Path = galleryPathErr
	id, err = i.FindExistingID()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.GalleryReaderWriter{}

	gallery := models.Gallery{
		Path: galleryPath,
	}

	galleryErr := models.Gallery{
		Path: galleryPathErr,
	}

	i := Importer{
		ReaderWriter: readerWriter,
		gallery:      gallery,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", gallery).Return(&models.Gallery{
		ID: galleryID,
	}, nil).Once()
	readerWriter.On("Create", galleryErr).Return(nil, errCreate).Once()

	id, err := i.Create()
	assert.Equal(t, galleryID, *id)
	assert.Nil(t, err)

	i.gallery = galleryErr
	id, err = i.Create()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.GalleryReaderWriter{}

	gallery := models.Gallery{
		Path: galleryPath,
	}

	galleryErr := models.Gallery{
		Path: galleryPathErr,
	}

	i := Importer{
		ReaderWriter: readerWriter,
		gallery:      gallery,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	gallery.ID = galleryID
	readerWriter.On("Update", gallery).Return(nil, nil).Once()

	err := i.Update(galleryID)
	assert.Nil(t, err)

	i.gallery = galleryErr

	// need to set id separately
	galleryErr.ID = idErr
	readerWriter.On("Update", galleryErr).Return(nil, errUpdate).Once()

	err = i.Update(idErr)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
