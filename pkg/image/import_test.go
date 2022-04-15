package image

import (
	"context"
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	path = "path"

	imageNameErr = "imageNameErr"
	// existingImageName = "existingImageName"

	existingImageID     = 100
	existingStudioID    = 101
	existingGalleryID   = 102
	existingPerformerID = 103
	// existingMovieID     = 104
	existingTagID = 105

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

var testCtx = context.Background()

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

	err := i.PreImport(testCtx)
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

	studioReaderWriter.On("FindByName", testCtx, existingStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	studioReaderWriter.On("FindByName", testCtx, existingStudioErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, int64(existingStudioID), i.image.StudioID.Int64)

	i.Input.Studio = existingStudioErr
	err = i.PreImport(testCtx)
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

	studioReaderWriter.On("FindByName", testCtx, missingStudioName, false).Return(nil, nil).Times(3)
	studioReaderWriter.On("Create", testCtx, mock.AnythingOfType("models.Studio")).Return(&models.Studio{
		ID: existingStudioID,
	}, nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
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

	studioReaderWriter.On("FindByName", testCtx, missingStudioName, false).Return(nil, nil).Once()
	studioReaderWriter.On("Create", testCtx, mock.AnythingOfType("models.Studio")).Return(nil, errors.New("Create error"))

	err := i.PreImport(testCtx)
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

	galleryReaderWriter.On("FindByChecksums", testCtx, []string{existingGalleryChecksum}).Return([]*models.Gallery{{
		ID: existingGalleryID,
	}}, nil).Once()
	galleryReaderWriter.On("FindByChecksums", testCtx, []string{existingGalleryErr}).Return(nil, errors.New("FindByChecksum error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingGalleryID, i.galleries[0].ID)

	i.Input.Galleries = []string{
		existingGalleryErr,
	}

	err = i.PreImport(testCtx)
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

	galleryReaderWriter.On("FindByChecksums", testCtx, []string{missingGalleryChecksum}).Return(nil, nil).Times(3)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Nil(t, i.galleries)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
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

	performerReaderWriter.On("FindByNames", testCtx, []string{existingPerformerName}, false).Return([]*models.Performer{
		{
			ID:   existingPerformerID,
			Name: models.NullString(existingPerformerName),
		},
	}, nil).Once()
	performerReaderWriter.On("FindByNames", testCtx, []string{existingPerformerErr}, false).Return(nil, errors.New("FindByNames error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingPerformerID, i.performers[0].ID)

	i.Input.Performers = []string{existingPerformerErr}
	err = i.PreImport(testCtx)
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

	performerReaderWriter.On("FindByNames", testCtx, []string{missingPerformerName}, false).Return(nil, nil).Times(3)
	performerReaderWriter.On("Create", testCtx, mock.AnythingOfType("models.Performer")).Return(&models.Performer{
		ID: existingPerformerID,
	}, nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
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

	performerReaderWriter.On("FindByNames", testCtx, []string{missingPerformerName}, false).Return(nil, nil).Once()
	performerReaderWriter.On("Create", testCtx, mock.AnythingOfType("models.Performer")).Return(nil, errors.New("Create error"))

	err := i.PreImport(testCtx)
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

	tagReaderWriter.On("FindByNames", testCtx, []string{existingTagName}, false).Return([]*models.Tag{
		{
			ID:   existingTagID,
			Name: existingTagName,
		},
	}, nil).Once()
	tagReaderWriter.On("FindByNames", testCtx, []string{existingTagErr}, false).Return(nil, errors.New("FindByNames error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingTagID, i.tags[0].ID)

	i.Input.Tags = []string{existingTagErr}
	err = i.PreImport(testCtx)
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

	tagReaderWriter.On("FindByNames", testCtx, []string{missingTagName}, false).Return(nil, nil).Times(3)
	tagReaderWriter.On("Create", testCtx, mock.AnythingOfType("models.Tag")).Return(&models.Tag{
		ID: existingTagID,
	}, nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
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

	tagReaderWriter.On("FindByNames", testCtx, []string{missingTagName}, false).Return(nil, nil).Once()
	tagReaderWriter.On("Create", testCtx, mock.AnythingOfType("models.Tag")).Return(nil, errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)
}

func TestImporterPostImportUpdateGallery(t *testing.T) {
	readerWriter := &mocks.ImageReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		galleries: []*models.Gallery{
			{
				ID: existingGalleryID,
			},
		},
	}

	updateErr := errors.New("UpdateGalleries error")

	readerWriter.On("UpdateGalleries", testCtx, imageID, []int{existingGalleryID}).Return(nil).Once()
	readerWriter.On("UpdateGalleries", testCtx, errGalleriesID, mock.AnythingOfType("[]int")).Return(updateErr).Once()

	err := i.PostImport(testCtx, imageID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errGalleriesID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterPostImportUpdatePerformers(t *testing.T) {
	readerWriter := &mocks.ImageReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		performers: []*models.Performer{
			{
				ID: existingPerformerID,
			},
		},
	}

	updateErr := errors.New("UpdatePerformers error")

	readerWriter.On("UpdatePerformers", testCtx, imageID, []int{existingPerformerID}).Return(nil).Once()
	readerWriter.On("UpdatePerformers", testCtx, errPerformersID, mock.AnythingOfType("[]int")).Return(updateErr).Once()

	err := i.PostImport(testCtx, imageID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errPerformersID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterPostImportUpdateTags(t *testing.T) {
	readerWriter := &mocks.ImageReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		tags: []*models.Tag{
			{
				ID: existingTagID,
			},
		},
	}

	updateErr := errors.New("UpdateTags error")

	readerWriter.On("UpdateTags", testCtx, imageID, []int{existingTagID}).Return(nil).Once()
	readerWriter.On("UpdateTags", testCtx, errTagsID, mock.AnythingOfType("[]int")).Return(updateErr).Once()

	err := i.PostImport(testCtx, imageID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errTagsID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
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
	readerWriter.On("FindByChecksum", testCtx, missingChecksum).Return(nil, nil).Once()
	readerWriter.On("FindByChecksum", testCtx, checksum).Return(&models.Image{
		ID: existingImageID,
	}, nil).Once()
	readerWriter.On("FindByChecksum", testCtx, errChecksum).Return(nil, expectedErr).Once()

	id, err := i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Checksum = checksum
	id, err = i.FindExistingID(testCtx)
	assert.Equal(t, existingImageID, *id)
	assert.Nil(t, err)

	i.Input.Checksum = errChecksum
	id, err = i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.ImageReaderWriter{}

	image := models.Image{
		Title: models.NullString(title),
	}

	imageErr := models.Image{
		Title: models.NullString(imageNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		image:        image,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", testCtx, image).Return(&models.Image{
		ID: imageID,
	}, nil).Once()
	readerWriter.On("Create", testCtx, imageErr).Return(nil, errCreate).Once()

	id, err := i.Create(testCtx)
	assert.Equal(t, imageID, *id)
	assert.Nil(t, err)
	assert.Equal(t, imageID, i.ID)

	i.image = imageErr
	id, err = i.Create(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.ImageReaderWriter{}

	image := models.Image{
		Title: models.NullString(title),
	}

	imageErr := models.Image{
		Title: models.NullString(imageNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		image:        image,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	image.ID = imageID
	readerWriter.On("UpdateFull", testCtx, image).Return(nil, nil).Once()

	err := i.Update(testCtx, imageID)
	assert.Nil(t, err)
	assert.Equal(t, imageID, i.ID)

	i.image = imageErr

	// need to set id separately
	imageErr.ID = errImageID
	readerWriter.On("UpdateFull", testCtx, imageErr).Return(nil, errUpdate).Once()

	err = i.Update(testCtx, errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
