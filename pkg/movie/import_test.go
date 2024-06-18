package movie

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
	movieNameErr      = "movieNameErr"
	existingMovieName = "existingMovieName"

	existingMovieID  = 100
	existingStudioID = 101

	existingStudioName = "existingStudioName"
	existingStudioErr  = "existingStudioErr"
	missingStudioName  = "existingStudioName"

	errImageID = 3

	existingTagID = 105
	errTagsID     = 106

	existingTagName = "existingTagName"
	existingTagErr  = "existingTagErr"
	missingTagName  = "missingTagName"
)

var testCtx = context.Background()

func TestImporterName(t *testing.T) {
	i := Importer{
		Input: jsonschema.Movie{
			Name: movieName,
		},
	}

	assert.Equal(t, movieName, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Input: jsonschema.Movie{
			Name:       movieName,
			FrontImage: invalidImage,
		},
	}

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.Input.FrontImage = frontImage
	i.Input.BackImage = invalidImage

	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.Input.BackImage = ""

	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.Input.BackImage = backImage

	err = i.PreImport(testCtx)
	assert.Nil(t, err)
}

func TestImporterPreImportWithStudio(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Movie,
		StudioWriter: db.Studio,
		Input: jsonschema.Movie{
			Name:       movieName,
			FrontImage: frontImage,
			Studio:     existingStudioName,
			Rating:     5,
			Duration:   10,
		},
	}

	db.Studio.On("FindByName", testCtx, existingStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	db.Studio.On("FindByName", testCtx, existingStudioErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingStudioID, *i.movie.StudioID)

	i.Input.Studio = existingStudioErr
	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingStudio(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Movie,
		StudioWriter: db.Studio,
		Input: jsonschema.Movie{
			Name:       movieName,
			FrontImage: frontImage,
			Studio:     missingStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	db.Studio.On("FindByName", testCtx, missingStudioName, false).Return(nil, nil).Times(3)
	db.Studio.On("Create", testCtx, mock.AnythingOfType("*models.Studio")).Run(func(args mock.Arguments) {
		s := args.Get(1).(*models.Studio)
		s.ID = existingStudioID
	}).Return(nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingStudioID, *i.movie.StudioID)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingStudioCreateErr(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Movie,
		StudioWriter: db.Studio,
		Input: jsonschema.Movie{
			Name:       movieName,
			FrontImage: frontImage,
			Studio:     missingStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	db.Studio.On("FindByName", testCtx, missingStudioName, false).Return(nil, nil).Once()
	db.Studio.On("Create", testCtx, mock.AnythingOfType("*models.Studio")).Return(errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithTag(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter:        db.Movie,
		TagWriter:           db.Tag,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Movie{
			Tags: []string{
				existingTagName,
			},
		},
	}

	db.Tag.On("FindByNames", testCtx, []string{existingTagName}, false).Return([]*models.Tag{
		{
			ID:   existingTagID,
			Name: existingTagName,
		},
	}, nil).Once()
	db.Tag.On("FindByNames", testCtx, []string{existingTagErr}, false).Return(nil, errors.New("FindByNames error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingTagID, i.movie.TagIDs.List()[0])

	i.Input.Tags = []string{existingTagErr}
	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingTag(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Movie,
		TagWriter:    db.Tag,
		Input: jsonschema.Movie{
			Tags: []string{
				missingTagName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	db.Tag.On("FindByNames", testCtx, []string{missingTagName}, false).Return(nil, nil).Times(3)
	db.Tag.On("Create", testCtx, mock.AnythingOfType("*models.Tag")).Run(func(args mock.Arguments) {
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
	assert.Equal(t, existingTagID, i.movie.TagIDs.List()[0])

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingTagCreateErr(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Movie,
		TagWriter:    db.Tag,
		Input: jsonschema.Movie{
			Tags: []string{
				missingTagName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	db.Tag.On("FindByNames", testCtx, []string{missingTagName}, false).Return(nil, nil).Once()
	db.Tag.On("Create", testCtx, mock.AnythingOfType("*models.Tag")).Return(errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPostImport(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter:   db.Movie,
		StudioWriter:   db.Studio,
		frontImageData: frontImageBytes,
		backImageData:  backImageBytes,
	}

	updateMovieImageErr := errors.New("UpdateImages error")

	db.Movie.On("UpdateFrontImage", testCtx, movieID, frontImageBytes).Return(nil).Once()
	db.Movie.On("UpdateBackImage", testCtx, movieID, backImageBytes).Return(nil).Once()
	db.Movie.On("UpdateFrontImage", testCtx, errImageID, frontImageBytes).Return(updateMovieImageErr).Once()

	err := i.PostImport(testCtx, movieID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errImageID)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Movie,
		StudioWriter: db.Studio,
		Input: jsonschema.Movie{
			Name: movieName,
		},
	}

	errFindByName := errors.New("FindByName error")
	db.Movie.On("FindByName", testCtx, movieName, false).Return(nil, nil).Once()
	db.Movie.On("FindByName", testCtx, existingMovieName, false).Return(&models.Movie{
		ID: existingMovieID,
	}, nil).Once()
	db.Movie.On("FindByName", testCtx, movieNameErr, false).Return(nil, errFindByName).Once()

	id, err := i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Name = existingMovieName
	id, err = i.FindExistingID(testCtx)
	assert.Equal(t, existingMovieID, *id)
	assert.Nil(t, err)

	i.Input.Name = movieNameErr
	id, err = i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	db := mocks.NewDatabase()

	movie := models.Movie{
		Name: movieName,
	}

	movieErr := models.Movie{
		Name: movieNameErr,
	}

	i := Importer{
		ReaderWriter: db.Movie,
		StudioWriter: db.Studio,
		movie:        movie,
	}

	errCreate := errors.New("Create error")
	db.Movie.On("Create", testCtx, &movie).Run(func(args mock.Arguments) {
		m := args.Get(1).(*models.Movie)
		m.ID = movieID
	}).Return(nil).Once()
	db.Movie.On("Create", testCtx, &movieErr).Return(errCreate).Once()

	id, err := i.Create(testCtx)
	assert.Equal(t, movieID, *id)
	assert.Nil(t, err)

	i.movie = movieErr
	id, err = i.Create(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	db := mocks.NewDatabase()

	movie := models.Movie{
		Name: movieName,
	}

	movieErr := models.Movie{
		Name: movieNameErr,
	}

	i := Importer{
		ReaderWriter: db.Movie,
		StudioWriter: db.Studio,
		movie:        movie,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	movie.ID = movieID
	db.Movie.On("Update", testCtx, &movie).Return(nil).Once()

	err := i.Update(testCtx, movieID)
	assert.Nil(t, err)

	i.movie = movieErr

	// need to set id separately
	movieErr.ID = errImageID
	db.Movie.On("Update", testCtx, &movieErr).Return(errUpdate).Once()

	err = i.Update(testCtx, errImageID)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}
