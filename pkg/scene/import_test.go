package scene

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

	sceneNameErr      = "sceneNameErr"
	existingSceneName = "existingSceneName"

	existingSceneID     = 100
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

	existingMovieName = "existingMovieName"
	existingMovieErr  = "existingMovieErr"
	missingMovieName  = "missingMovieName"

	existingTagName = "existingTagName"
	existingTagErr  = "existingTagErr"
	missingTagName  = "missingTagName"

	errPerformersID = 200

	missingChecksum = "missingChecksum"
	missingOSHash   = "missingOSHash"
	errChecksum     = "errChecksum"
	errOSHash       = "errOSHash"
)

func TestImporterName(t *testing.T) {
	i := Importer{
		Path:  path,
		Input: jsonschema.Scene{},
	}

	assert.Equal(t, path, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Path: path,
		Input: jsonschema.Scene{
			Cover: invalidImage,
		},
	}

	err := i.PreImport()
	assert.NotNil(t, err)

	i.Input.Cover = image

	err = i.PreImport()
	assert.Nil(t, err)
}

func TestImporterPreImportWithStudio(t *testing.T) {
	studioReaderWriter := &mocks.StudioReaderWriter{}

	i := Importer{
		StudioWriter: studioReaderWriter,
		Path:         path,
		Input: jsonschema.Scene{
			Studio: existingStudioName,
		},
	}

	studioReaderWriter.On("FindByName", existingStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	studioReaderWriter.On("FindByName", existingStudioErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, int64(existingStudioID), i.scene.StudioID.Int64)

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
		Input: jsonschema.Scene{
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
	assert.Equal(t, int64(existingStudioID), i.scene.StudioID.Int64)

	studioReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingStudioCreateErr(t *testing.T) {
	studioReaderWriter := &mocks.StudioReaderWriter{}

	i := Importer{
		StudioWriter: studioReaderWriter,
		Path:         path,
		Input: jsonschema.Scene{
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
		Input: jsonschema.Scene{
			Gallery: existingGalleryChecksum,
		},
	}

	galleryReaderWriter.On("FindByChecksum", existingGalleryChecksum).Return(&models.Gallery{
		ID: existingGalleryID,
	}, nil).Once()
	galleryReaderWriter.On("FindByChecksum", existingGalleryErr).Return(nil, errors.New("FindByChecksum error")).Once()

	err := i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, existingGalleryID, i.gallery.ID)

	i.Input.Gallery = existingGalleryErr
	err = i.PreImport()
	assert.NotNil(t, err)

	galleryReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingGallery(t *testing.T) {
	galleryReaderWriter := &mocks.GalleryReaderWriter{}

	i := Importer{
		Path:          path,
		GalleryWriter: galleryReaderWriter,
		Input: jsonschema.Scene{
			Gallery: missingGalleryChecksum,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	galleryReaderWriter.On("FindByChecksum", missingGalleryChecksum).Return(nil, nil).Times(3)

	err := i.PreImport()
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport()
	assert.Nil(t, err)
	assert.Nil(t, i.gallery)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport()
	assert.Nil(t, err)
	assert.Nil(t, i.gallery)

	galleryReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithPerformer(t *testing.T) {
	performerReaderWriter := &mocks.PerformerReaderWriter{}

	i := Importer{
		PerformerWriter:     performerReaderWriter,
		Path:                path,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Scene{
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
		Input: jsonschema.Scene{
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
		Input: jsonschema.Scene{
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

func TestImporterPreImportWithMovie(t *testing.T) {
	movieReaderWriter := &mocks.MovieReaderWriter{}

	i := Importer{
		MovieWriter:         movieReaderWriter,
		Path:                path,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Scene{
			Movies: []jsonschema.SceneMovie{
				{
					MovieName:  existingMovieName,
					SceneIndex: 1,
				},
			},
		},
	}

	movieReaderWriter.On("FindByName", existingMovieName, false).Return(&models.Movie{
		ID:   existingMovieID,
		Name: modelstest.NullString(existingMovieName),
	}, nil).Once()
	movieReaderWriter.On("FindByName", existingMovieErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, existingMovieID, i.movies[0].MovieID)

	i.Input.Movies[0].MovieName = existingMovieErr
	err = i.PreImport()
	assert.NotNil(t, err)

	movieReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingMovie(t *testing.T) {
	movieReaderWriter := &mocks.MovieReaderWriter{}

	i := Importer{
		Path:        path,
		MovieWriter: movieReaderWriter,
		Input: jsonschema.Scene{
			Movies: []jsonschema.SceneMovie{
				{
					MovieName: missingMovieName,
				},
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	movieReaderWriter.On("FindByName", missingMovieName, false).Return(nil, nil).Times(3)
	movieReaderWriter.On("Create", mock.AnythingOfType("models.Movie")).Return(&models.Movie{
		ID: existingMovieID,
	}, nil)

	err := i.PreImport()
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport()
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport()
	assert.Nil(t, err)
	assert.Equal(t, existingMovieID, i.movies[0].MovieID)

	movieReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingMovieCreateErr(t *testing.T) {
	movieReaderWriter := &mocks.MovieReaderWriter{}

	i := Importer{
		MovieWriter: movieReaderWriter,
		Path:        path,
		Input: jsonschema.Scene{
			Movies: []jsonschema.SceneMovie{
				{
					MovieName: missingMovieName,
				},
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	movieReaderWriter.On("FindByName", missingMovieName, false).Return(nil, nil).Once()
	movieReaderWriter.On("Create", mock.AnythingOfType("models.Movie")).Return(nil, errors.New("Create error"))

	err := i.PreImport()
	assert.NotNil(t, err)
}

func TestImporterPreImportWithTag(t *testing.T) {
	tagReaderWriter := &mocks.TagReaderWriter{}

	i := Importer{
		TagWriter:           tagReaderWriter,
		Path:                path,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Scene{
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
		Input: jsonschema.Scene{
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
		Input: jsonschema.Scene{
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
	readerWriter := &mocks.SceneReaderWriter{}

	i := Importer{
		ReaderWriter:   readerWriter,
		coverImageData: imageBytes,
	}

	updateSceneImageErr := errors.New("UpdateSceneCover error")

	readerWriter.On("UpdateSceneCover", sceneID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateSceneCover", errImageID, imageBytes).Return(updateSceneImageErr).Once()

	err := i.PostImport(sceneID)
	assert.Nil(t, err)

	err = i.PostImport(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterPostImportUpdateGallery(t *testing.T) {
	galleryReaderWriter := &mocks.GalleryReaderWriter{}

	i := Importer{
		GalleryWriter: galleryReaderWriter,
		gallery: &models.Gallery{
			ID: existingGalleryID,
		},
	}

	updateErr := errors.New("Update error")

	updateArg := *i.gallery
	updateArg.SceneID = modelstest.NullInt64(sceneID)

	galleryReaderWriter.On("Update", updateArg).Return(nil, nil).Once()

	updateArg.SceneID = modelstest.NullInt64(errGalleryID)
	galleryReaderWriter.On("Update", updateArg).Return(nil, updateErr).Once()

	err := i.PostImport(sceneID)
	assert.Nil(t, err)

	err = i.PostImport(errGalleryID)
	assert.NotNil(t, err)

	galleryReaderWriter.AssertExpectations(t)
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

	updateErr := errors.New("UpdatePerformersScenes error")

	joinReaderWriter.On("UpdatePerformersScenes", sceneID, []models.PerformersScenes{
		{
			PerformerID: existingPerformerID,
			SceneID:     sceneID,
		},
	}).Return(nil).Once()
	joinReaderWriter.On("UpdatePerformersScenes", errPerformersID, mock.AnythingOfType("[]models.PerformersScenes")).Return(updateErr).Once()

	err := i.PostImport(sceneID)
	assert.Nil(t, err)

	err = i.PostImport(errPerformersID)
	assert.NotNil(t, err)

	joinReaderWriter.AssertExpectations(t)
}

func TestImporterPostImportUpdateMovies(t *testing.T) {
	joinReaderWriter := &mocks.JoinReaderWriter{}

	i := Importer{
		JoinWriter: joinReaderWriter,
		movies: []models.MoviesScenes{
			{
				MovieID: existingMovieID,
			},
		},
	}

	updateErr := errors.New("UpdateMoviesScenes error")

	joinReaderWriter.On("UpdateMoviesScenes", sceneID, []models.MoviesScenes{
		{
			MovieID: existingMovieID,
			SceneID: sceneID,
		},
	}).Return(nil).Once()
	joinReaderWriter.On("UpdateMoviesScenes", errMoviesID, mock.AnythingOfType("[]models.MoviesScenes")).Return(updateErr).Once()

	err := i.PostImport(sceneID)
	assert.Nil(t, err)

	err = i.PostImport(errMoviesID)
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

	updateErr := errors.New("UpdateScenesTags error")

	joinReaderWriter.On("UpdateScenesTags", sceneID, []models.ScenesTags{
		{
			TagID:   existingTagID,
			SceneID: sceneID,
		},
	}).Return(nil).Once()
	joinReaderWriter.On("UpdateScenesTags", errTagsID, mock.AnythingOfType("[]models.ScenesTags")).Return(updateErr).Once()

	err := i.PostImport(sceneID)
	assert.Nil(t, err)

	err = i.PostImport(errTagsID)
	assert.NotNil(t, err)

	joinReaderWriter.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	readerWriter := &mocks.SceneReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Path:         path,
		Input: jsonschema.Scene{
			Checksum: missingChecksum,
			OSHash:   missingOSHash,
		},
		FileNamingAlgorithm: models.HashAlgorithmMd5,
	}

	expectedErr := errors.New("FindBy* error")
	readerWriter.On("FindByChecksum", missingChecksum).Return(nil, nil).Once()
	readerWriter.On("FindByChecksum", checksum).Return(&models.Scene{
		ID: existingSceneID,
	}, nil).Once()
	readerWriter.On("FindByChecksum", errChecksum).Return(nil, expectedErr).Once()

	readerWriter.On("FindByOSHash", missingOSHash).Return(nil, nil).Once()
	readerWriter.On("FindByOSHash", oshash).Return(&models.Scene{
		ID: existingSceneID,
	}, nil).Once()
	readerWriter.On("FindByOSHash", errOSHash).Return(nil, expectedErr).Once()

	id, err := i.FindExistingID()
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Checksum = checksum
	id, err = i.FindExistingID()
	assert.Equal(t, existingSceneID, *id)
	assert.Nil(t, err)

	i.Input.Checksum = errChecksum
	id, err = i.FindExistingID()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	i.FileNamingAlgorithm = models.HashAlgorithmOshash
	id, err = i.FindExistingID()
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.OSHash = oshash
	id, err = i.FindExistingID()
	assert.Equal(t, existingSceneID, *id)
	assert.Nil(t, err)

	i.Input.OSHash = errOSHash
	id, err = i.FindExistingID()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.SceneReaderWriter{}

	scene := models.Scene{
		Title: modelstest.NullString(title),
	}

	sceneErr := models.Scene{
		Title: modelstest.NullString(sceneNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		scene:        scene,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", scene).Return(&models.Scene{
		ID: sceneID,
	}, nil).Once()
	readerWriter.On("Create", sceneErr).Return(nil, errCreate).Once()

	id, err := i.Create()
	assert.Equal(t, sceneID, *id)
	assert.Nil(t, err)
	assert.Equal(t, sceneID, i.ID)

	i.scene = sceneErr
	id, err = i.Create()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.SceneReaderWriter{}

	scene := models.Scene{
		Title: modelstest.NullString(title),
	}

	sceneErr := models.Scene{
		Title: modelstest.NullString(sceneNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		scene:        scene,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	scene.ID = sceneID
	readerWriter.On("UpdateFull", scene).Return(nil, nil).Once()

	err := i.Update(sceneID)
	assert.Nil(t, err)
	assert.Equal(t, sceneID, i.ID)

	i.scene = sceneErr

	// need to set id separately
	sceneErr.ID = errImageID
	readerWriter.On("UpdateFull", sceneErr).Return(nil, errUpdate).Once()

	err = i.Update(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
