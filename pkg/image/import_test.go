package image

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
	path = "path"

	imageNameErr      = "imageNameErr"
	existingImageName = "existingImageName"

	existingImageID     = 100
	existingStudioID    = 101
	existingGalleryID   = 102
	existingPerformerID = 103
	existingMovieID     = 104
	existingTagID       = 105

	existingStudioName = "existingStudioName"
	existingStudioErr  = "existingStudioErr"
	missingStudioName  = "missingStudioName"

	existingGalleryChecksum = "existingGalleryChecksum"
	existingGalleryErr      = "existingGalleryErr"
	missingGalleryChecksum  = "missingGalleryChecksum"

	existingPerformerName = "existingPerformerName"
	existingPerformerErr  = "existingPerformerErr"
	missingPerformerName  = "missingPerformerName"

	existingTagName = "existingTagName"
	existingTagErr  = "existingTagErr"
	missingTagName  = "missingTagName"

	errPerformersID = 200
	errGalleriesID  = 201

	missingChecksum = "missingChecksum"
	errChecksum     = "errChecksum"
)

func TestImporterName(t *testing.T) {
	i := Importer{
		Path:  path,
		Input: jsonschema.Image{},
	}

	assert.Equal(t, path, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Path: path,
	}

	err := i.PreImport()
	assert.Nil(t, err)
}

func TestImporterPreImportWithStudio(t *testing.T) {
	studioReaderWriter := &mocks.StudioReaderWriter{}

	i := Importer{
		StudioWriter: studioReaderWriter,
		Path:         path,
		Input: jsonschema.Image{
			Studio: existingStudioName,
		},
	}

	studioReaderWriter.On("FindByName", existingStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	studioReaderWriter.On("FindByName", existingStudioErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, int64(existingStudioID), i.image.StudioID.Int64)

	i.Input.Studio = existingStudioErr
	err = i.PreImport()
	assert.NotNil(t, err)

	studioReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingStudio(t *testing.T) {
	studioReaderWriter := &mocks.StudioReaderWriter{}

	i := Importer{
		Path:         path,
		StudioWriter: studioReaderWriter,
		Input: jsonschema.Image{
			Studio: missingStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	studioReaderWriter.On("FindByName", missingStudioName, false).Return(nil, nil).Times(3)
	studioReaderWriter.On("Create", mock.AnythingOfType("models.Studio")).Return(&models.Studio{
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
	assert.Equal(t, int64(existingStudioID), i.image.StudioID.Int64)

	studioReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingStudioCreateErr(t *testing.T) {
	studioReaderWriter := &mocks.StudioReaderWriter{}

	i := Importer{
		StudioWriter: studioReaderWriter,
		Path:         path,
		Input: jsonschema.Image{
			Studio: missingStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	studioReaderWriter.On("FindByName", missingStudioName, false).Return(nil, nil).Once()
	studioReaderWriter.On("Create", mock.AnythingOfType("models.Studio")).Return(nil, errors.New("Create error"))

	err := i.PreImport()
	assert.NotNil(t, err)
}

func TestImporterPreImportWithGallery(t *testing.T) {
	galleryReaderWriter := &mocks.GalleryReaderWriter{}

	i := Importer{
		GalleryWriter: galleryReaderWriter,
		Path:          path,
		Input: jsonschema.Image{
			Galleries: []string{
				existingGalleryChecksum,
			},
		},
	}

	galleryReaderWriter.On("FindByChecksum", existingGalleryChecksum).Return(&models.Gallery{
		ID: existingGalleryID,
	}, nil).Once()
	galleryReaderWriter.On("FindByChecksum", existingGalleryErr).Return(nil, errors.New("FindByChecksum error")).Once()

	err := i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, existingGalleryID, i.galleries[0].ID)

	i.Input.Galleries = []string{
		existingGalleryErr,
	}

	err = i.PreImport()
	assert.NotNil(t, err)

	galleryReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingGallery(t *testing.T) {
	galleryReaderWriter := &mocks.GalleryReaderWriter{}

	i := Importer{
		Path:          path,
		GalleryWriter: galleryReaderWriter,
		Input: jsonschema.Image{
			Galleries: []string{
				missingGalleryChecksum,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	galleryReaderWriter.On("FindByChecksum", missingGalleryChecksum).Return(nil, nil).Times(3)

	err := i.PreImport()
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport()
	assert.Nil(t, err)
	assert.Nil(t, i.galleries)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport()
	assert.Nil(t, err)
	assert.Nil(t, i.galleries)

	galleryReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithPerformer(t *testing.T) {
	performerReaderWriter := &mocks.PerformerReaderWriter{}

	i := Importer{
		PerformerWriter:     performerReaderWriter,
		Path:                path,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Image{
			Performers: []string{
				existingPerformerName,
			},
		},
	}

	performerReaderWriter.On("FindByNames", []string{existingPerformerName}, false).Return([]*models.Performer{
		{
			ID:   existingPerformerID,
			Name: modelstest.NullString(existingPerformerName),
		},
	}, nil).Once()
	performerReaderWriter.On("FindByNames", []string{existingPerformerErr}, false).Return(nil, errors.New("FindByNames error")).Once()

	err := i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, existingPerformerID, i.performers[0].ID)

	i.Input.Performers = []string{existingPerformerErr}
	err = i.PreImport()
	assert.NotNil(t, err)

	performerReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingPerformer(t *testing.T) {
	performerReaderWriter := &mocks.PerformerReaderWriter{}

	i := Importer{
		Path:            path,
		PerformerWriter: performerReaderWriter,
		Input: jsonschema.Image{
			Performers: []string{
				missingPerformerName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	performerReaderWriter.On("FindByNames", []string{missingPerformerName}, false).Return(nil, nil).Times(3)
	performerReaderWriter.On("Create", mock.AnythingOfType("models.Performer")).Return(&models.Performer{
		ID: existingPerformerID,
	}, nil)

	err := i.PreImport()
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport()
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, existingPerformerID, i.performers[0].ID)

	performerReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingPerformerCreateErr(t *testing.T) {
	performerReaderWriter := &mocks.PerformerReaderWriter{}

	i := Importer{
		PerformerWriter: performerReaderWriter,
		Path:            path,
		Input: jsonschema.Image{
			Performers: []string{
				missingPerformerName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	performerReaderWriter.On("FindByNames", []string{missingPerformerName}, false).Return(nil, nil).Once()
	performerReaderWriter.On("Create", mock.AnythingOfType("models.Performer")).Return(nil, errors.New("Create error"))

	err := i.PreImport()
	assert.NotNil(t, err)
}

func TestImporterPreImportWithTag(t *testing.T) {
	tagReaderWriter := &mocks.TagReaderWriter{}

	i := Importer{
		TagWriter:           tagReaderWriter,
		Path:                path,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Image{
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
		Path:      path,
		TagWriter: tagReaderWriter,
		Input: jsonschema.Image{
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
		Path:      path,
		Input: jsonschema.Image{
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

func TestImporterPostImportUpdateGallery(t *testing.T) {
	joinReaderWriter := &mocks.JoinReaderWriter{}

	i := Importer{
		JoinWriter: joinReaderWriter,
		galleries: []*models.Gallery{
			{
				ID: existingGalleryID,
			},
		},
	}

	updateErr := errors.New("UpdateGalleriesImages error")

	joinReaderWriter.On("UpdateGalleriesImages", imageID, []models.GalleriesImages{
		{
			GalleryID: existingGalleryID,
			ImageID:   imageID,
		},
	}).Return(nil).Once()
	joinReaderWriter.On("UpdateGalleriesImages", errGalleriesID, mock.AnythingOfType("[]models.GalleriesImages")).Return(updateErr).Once()

	err := i.PostImport(imageID)
	assert.Nil(t, err)

	err = i.PostImport(errGalleriesID)
	assert.NotNil(t, err)

	joinReaderWriter.AssertExpectations(t)
}

func TestImporterPostImportUpdatePerformers(t *testing.T) {
	joinReaderWriter := &mocks.JoinReaderWriter{}

	i := Importer{
		JoinWriter: joinReaderWriter,
		performers: []*models.Performer{
			{
				ID: existingPerformerID,
			},
		},
	}

	updateErr := errors.New("UpdatePerformersImages error")

	joinReaderWriter.On("UpdatePerformersImages", imageID, []models.PerformersImages{
		{
			PerformerID: existingPerformerID,
			ImageID:     imageID,
		},
	}).Return(nil).Once()
	joinReaderWriter.On("UpdatePerformersImages", errPerformersID, mock.AnythingOfType("[]models.PerformersImages")).Return(updateErr).Once()

	err := i.PostImport(imageID)
	assert.Nil(t, err)

	err = i.PostImport(errPerformersID)
	assert.NotNil(t, err)

	joinReaderWriter.AssertExpectations(t)
}

func TestImporterPostImportUpdateTags(t *testing.T) {
	joinReaderWriter := &mocks.JoinReaderWriter{}

	i := Importer{
		JoinWriter: joinReaderWriter,
		tags: []*models.Tag{
			{
				ID: existingTagID,
			},
		},
	}

	updateErr := errors.New("UpdateImagesTags error")

	joinReaderWriter.On("UpdateImagesTags", imageID, []models.ImagesTags{
		{
			TagID:   existingTagID,
			ImageID: imageID,
		},
	}).Return(nil).Once()
	joinReaderWriter.On("UpdateImagesTags", errTagsID, mock.AnythingOfType("[]models.ImagesTags")).Return(updateErr).Once()

	err := i.PostImport(imageID)
	assert.Nil(t, err)

	err = i.PostImport(errTagsID)
	assert.NotNil(t, err)

	joinReaderWriter.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	readerWriter := &mocks.ImageReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Path:         path,
		Input: jsonschema.Image{
			Checksum: missingChecksum,
		},
	}

	expectedErr := errors.New("FindBy* error")
	readerWriter.On("FindByChecksum", missingChecksum).Return(nil, nil).Once()
	readerWriter.On("FindByChecksum", checksum).Return(&models.Image{
		ID: existingImageID,
	}, nil).Once()
	readerWriter.On("FindByChecksum", errChecksum).Return(nil, expectedErr).Once()

	id, err := i.FindExistingID()
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Checksum = checksum
	id, err = i.FindExistingID()
	assert.Equal(t, existingImageID, *id)
	assert.Nil(t, err)

	i.Input.Checksum = errChecksum
	id, err = i.FindExistingID()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.ImageReaderWriter{}

	image := models.Image{
		Title: modelstest.NullString(title),
	}

	imageErr := models.Image{
		Title: modelstest.NullString(imageNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		image:        image,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", image).Return(&models.Image{
		ID: imageID,
	}, nil).Once()
	readerWriter.On("Create", imageErr).Return(nil, errCreate).Once()

	id, err := i.Create()
	assert.Equal(t, imageID, *id)
	assert.Nil(t, err)
	assert.Equal(t, imageID, i.ID)

	i.image = imageErr
	id, err = i.Create()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.ImageReaderWriter{}

	image := models.Image{
		Title: modelstest.NullString(title),
	}

	imageErr := models.Image{
		Title: modelstest.NullString(imageNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		image:        image,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	image.ID = imageID
	readerWriter.On("UpdateFull", image).Return(nil, nil).Once()

	err := i.Update(imageID)
	assert.Nil(t, err)
	assert.Equal(t, imageID, i.ID)

	i.image = imageErr

	// need to set id separately
	imageErr.ID = errImageID
	readerWriter.On("UpdateFull", imageErr).Return(nil, errUpdate).Once()

	err = i.Update(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
