package scene

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

const invalidImage = "aW1hZ2VCeXRlcw&&"

const (
	path = "path"

	sceneNameErr = "sceneNameErr"
	// existingSceneName = "existingSceneName"

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
	errGalleriesID  = 201

	missingChecksum = "missingChecksum"
	missingOSHash   = "missingOSHash"
	errChecksum     = "errChecksum"
	errOSHash       = "errOSHash"
)

var testCtx = context.Background()

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

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.Input.Cover = imageBase64

	err = i.PreImport(testCtx)
	assert.Nil(t, err)
}

func TestImporterPreImportWithStudio(t *testing.T) {
	studioReaderWriter := &mocks.StudioReaderWriter{}
	testCtx := context.Background()

	i := Importer{
		StudioWriter: studioReaderWriter,
		Path:         path,
		Input: jsonschema.Scene{
			Studio: existingStudioName,
		},
	}

	studioReaderWriter.On("FindByName", testCtx, existingStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	studioReaderWriter.On("FindByName", testCtx, existingStudioErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, int64(existingStudioID), i.scene.StudioID.Int64)

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
		Input: jsonschema.Scene{
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

	studioReaderWriter.On("FindByName", testCtx, missingStudioName, false).Return(nil, nil).Once()
	studioReaderWriter.On("Create", testCtx, mock.AnythingOfType("models.Studio")).Return(nil, errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)
}

func TestImporterPreImportWithGallery(t *testing.T) {
	galleryReaderWriter := &mocks.GalleryReaderWriter{}

	i := Importer{
		GalleryWriter:       galleryReaderWriter,
		Path:                path,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Scene{
			Galleries: []string{
				existingGalleryChecksum,
			},
		},
	}

	galleryReaderWriter.On("FindByChecksums", testCtx, []string{existingGalleryChecksum}).Return([]*models.Gallery{
		{
			ID:       existingGalleryID,
			Checksum: existingGalleryChecksum,
		},
	}, nil).Once()

	galleryReaderWriter.On("FindByChecksums", testCtx, []string{existingGalleryErr}).Return(nil, errors.New("FindByChecksums error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingGalleryID, i.galleries[0].ID)

	i.Input.Galleries = []string{existingGalleryErr}
	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	galleryReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingGallery(t *testing.T) {
	galleryReaderWriter := &mocks.GalleryReaderWriter{}

	i := Importer{
		Path:          path,
		GalleryWriter: galleryReaderWriter,
		Input: jsonschema.Scene{
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

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

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
		Input: jsonschema.Scene{
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
		Input: jsonschema.Scene{
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

func TestImporterPreImportWithMovie(t *testing.T) {
	movieReaderWriter := &mocks.MovieReaderWriter{}
	testCtx := context.Background()

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

	movieReaderWriter.On("FindByName", testCtx, existingMovieName, false).Return(&models.Movie{
		ID:   existingMovieID,
		Name: models.NullString(existingMovieName),
	}, nil).Once()
	movieReaderWriter.On("FindByName", testCtx, existingMovieErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingMovieID, i.movies[0].MovieID)

	i.Input.Movies[0].MovieName = existingMovieErr
	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	movieReaderWriter.AssertExpectations(t)
}

func TestImporterPreImportWithMissingMovie(t *testing.T) {
	movieReaderWriter := &mocks.MovieReaderWriter{}
	testCtx := context.Background()

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

	movieReaderWriter.On("FindByName", testCtx, missingMovieName, false).Return(nil, nil).Times(3)
	movieReaderWriter.On("Create", testCtx, mock.AnythingOfType("models.Movie")).Return(&models.Movie{
		ID: existingMovieID,
	}, nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
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

	movieReaderWriter.On("FindByName", testCtx, missingMovieName, false).Return(nil, nil).Once()
	movieReaderWriter.On("Create", testCtx, mock.AnythingOfType("models.Movie")).Return(nil, errors.New("Create error"))

	err := i.PreImport(testCtx)
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
		Input: jsonschema.Scene{
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
		Input: jsonschema.Scene{
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

func TestImporterPostImport(t *testing.T) {
	readerWriter := &mocks.SceneReaderWriter{}

	i := Importer{
		ReaderWriter:   readerWriter,
		coverImageData: imageBytes,
	}

	updateSceneImageErr := errors.New("UpdateCover error")

	readerWriter.On("UpdateCover", testCtx, sceneID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateCover", testCtx, errImageID, imageBytes).Return(updateSceneImageErr).Once()

	err := i.PostImport(testCtx, sceneID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterPostImportUpdateGalleries(t *testing.T) {
	sceneReaderWriter := &mocks.SceneReaderWriter{}

	i := Importer{
		ReaderWriter: sceneReaderWriter,
		galleries: []*models.Gallery{
			{
				ID: existingGalleryID,
			},
		},
	}

	updateErr := errors.New("UpdateGalleries error")

	sceneReaderWriter.On("UpdateGalleries", testCtx, sceneID, []int{existingGalleryID}).Return(nil).Once()
	sceneReaderWriter.On("UpdateGalleries", testCtx, errGalleriesID, mock.AnythingOfType("[]int")).Return(updateErr).Once()

	err := i.PostImport(testCtx, sceneID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errGalleriesID)
	assert.NotNil(t, err)

	sceneReaderWriter.AssertExpectations(t)
}

func TestImporterPostImportUpdatePerformers(t *testing.T) {
	sceneReaderWriter := &mocks.SceneReaderWriter{}

	i := Importer{
		ReaderWriter: sceneReaderWriter,
		performers: []*models.Performer{
			{
				ID: existingPerformerID,
			},
		},
	}

	updateErr := errors.New("UpdatePerformers error")

	sceneReaderWriter.On("UpdatePerformers", testCtx, sceneID, []int{existingPerformerID}).Return(nil).Once()
	sceneReaderWriter.On("UpdatePerformers", testCtx, errPerformersID, mock.AnythingOfType("[]int")).Return(updateErr).Once()

	err := i.PostImport(testCtx, sceneID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errPerformersID)
	assert.NotNil(t, err)

	sceneReaderWriter.AssertExpectations(t)
}

func TestImporterPostImportUpdateMovies(t *testing.T) {
	sceneReaderWriter := &mocks.SceneReaderWriter{}

	i := Importer{
		ReaderWriter: sceneReaderWriter,
		movies: []models.MoviesScenes{
			{
				MovieID: existingMovieID,
			},
		},
	}

	updateErr := errors.New("UpdateMovies error")

	sceneReaderWriter.On("UpdateMovies", testCtx, sceneID, []models.MoviesScenes{
		{
			MovieID: existingMovieID,
			SceneID: sceneID,
		},
	}).Return(nil).Once()
	sceneReaderWriter.On("UpdateMovies", testCtx, errMoviesID, mock.AnythingOfType("[]models.MoviesScenes")).Return(updateErr).Once()

	err := i.PostImport(testCtx, sceneID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errMoviesID)
	assert.NotNil(t, err)

	sceneReaderWriter.AssertExpectations(t)
}

func TestImporterPostImportUpdateTags(t *testing.T) {
	sceneReaderWriter := &mocks.SceneReaderWriter{}

	i := Importer{
		ReaderWriter: sceneReaderWriter,
		tags: []*models.Tag{
			{
				ID: existingTagID,
			},
		},
	}

	updateErr := errors.New("UpdateTags error")

	sceneReaderWriter.On("UpdateTags", testCtx, sceneID, []int{existingTagID}).Return(nil).Once()
	sceneReaderWriter.On("UpdateTags", testCtx, errTagsID, mock.AnythingOfType("[]int")).Return(updateErr).Once()

	err := i.PostImport(testCtx, sceneID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errTagsID)
	assert.NotNil(t, err)

	sceneReaderWriter.AssertExpectations(t)
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
	readerWriter.On("FindByChecksum", testCtx, missingChecksum).Return(nil, nil).Once()
	readerWriter.On("FindByChecksum", testCtx, checksum).Return(&models.Scene{
		ID: existingSceneID,
	}, nil).Once()
	readerWriter.On("FindByChecksum", testCtx, errChecksum).Return(nil, expectedErr).Once()

	readerWriter.On("FindByOSHash", testCtx, missingOSHash).Return(nil, nil).Once()
	readerWriter.On("FindByOSHash", testCtx, oshash).Return(&models.Scene{
		ID: existingSceneID,
	}, nil).Once()
	readerWriter.On("FindByOSHash", testCtx, errOSHash).Return(nil, expectedErr).Once()

	id, err := i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Checksum = checksum
	id, err = i.FindExistingID(testCtx)
	assert.Equal(t, existingSceneID, *id)
	assert.Nil(t, err)

	i.Input.Checksum = errChecksum
	id, err = i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	i.FileNamingAlgorithm = models.HashAlgorithmOshash
	id, err = i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.OSHash = oshash
	id, err = i.FindExistingID(testCtx)
	assert.Equal(t, existingSceneID, *id)
	assert.Nil(t, err)

	i.Input.OSHash = errOSHash
	id, err = i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.SceneReaderWriter{}

	scene := models.Scene{
		Title: models.NullString(title),
	}

	sceneErr := models.Scene{
		Title: models.NullString(sceneNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		scene:        scene,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", testCtx, scene).Return(&models.Scene{
		ID: sceneID,
	}, nil).Once()
	readerWriter.On("Create", testCtx, sceneErr).Return(nil, errCreate).Once()

	id, err := i.Create(testCtx)
	assert.Equal(t, sceneID, *id)
	assert.Nil(t, err)
	assert.Equal(t, sceneID, i.ID)

	i.scene = sceneErr
	id, err = i.Create(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.SceneReaderWriter{}

	scene := models.Scene{
		Title: models.NullString(title),
	}

	sceneErr := models.Scene{
		Title: models.NullString(sceneNameErr),
	}

	i := Importer{
		ReaderWriter: readerWriter,
		scene:        scene,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	scene.ID = sceneID
	readerWriter.On("UpdateFull", testCtx, scene).Return(nil, nil).Once()

	err := i.Update(testCtx, sceneID)
	assert.Nil(t, err)
	assert.Equal(t, sceneID, i.ID)

	i.scene = sceneErr

	// need to set id separately
	sceneErr.ID = errImageID
	readerWriter.On("UpdateFull", testCtx, sceneErr).Return(nil, errUpdate).Once()

	err = i.Update(testCtx, errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
