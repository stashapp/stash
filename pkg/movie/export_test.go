package movie

import (
	"database/sql"
	"errors"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

const (
	movieID              = 1
	emptyID              = 2
	errFrontImageID      = 3
	errBackImageID       = 4
	errStudioMovieID     = 5
	missingStudioMovieID = 6
)

const (
	studioID        = 1
	missingStudioID = 2
	errStudioID     = 3
)

const movieName = "testMovie"
const movieAliases = "aliases"

var date = models.SQLiteDate{
	String: "2001-01-01",
	Valid:  true,
}

const rating = 5
const duration = 100
const director = "director"
const synopsis = "synopsis"
const url = "url"

const studioName = "studio"

const (
	frontImage = "ZnJvbnRJbWFnZUJ5dGVz"
	backImage  = "YmFja0ltYWdlQnl0ZXM="
)

var (
	frontImageBytes = []byte("frontImageBytes")
	backImageBytes  = []byte("backImageBytes")
)

var movieStudio models.Studio = models.Studio{
	Name: models.NullString(studioName),
}

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

func createFullMovie(id int, studioID int) models.Movie {
	return models.Movie{
		ID:      id,
		Name:    models.NullString(movieName),
		Aliases: models.NullString(movieAliases),
		Date:    date,
		Rating: sql.NullInt64{
			Int64: rating,
			Valid: true,
		},
		Duration: sql.NullInt64{
			Int64: duration,
			Valid: true,
		},
		Director: models.NullString(director),
		Synopsis: models.NullString(synopsis),
		URL:      models.NullString(url),
		StudioID: sql.NullInt64{
			Int64: int64(studioID),
			Valid: true,
		},
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	}
}

func createEmptyMovie(id int) models.Movie {
	return models.Movie{
		ID: id,
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	}
}

func createFullJSONMovie(studio, frontImage, backImage string) *jsonschema.Movie {
	return &jsonschema.Movie{
		Name:       movieName,
		Aliases:    movieAliases,
		Date:       date.String,
		Rating:     rating,
		Duration:   duration,
		Director:   director,
		Synopsis:   synopsis,
		URL:        url,
		Studio:     studio,
		FrontImage: frontImage,
		BackImage:  backImage,
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
	}
}

func createEmptyJSONMovie() *jsonschema.Movie {
	return &jsonschema.Movie{
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
	}
}

type testScenario struct {
	movie    models.Movie
	expected *jsonschema.Movie
	err      bool
}

var scenarios []testScenario

func initTestTable() {
	scenarios = []testScenario{
		{
			createFullMovie(movieID, studioID),
			createFullJSONMovie(studioName, frontImage, backImage),
			false,
		},
		{
			createEmptyMovie(emptyID),
			createEmptyJSONMovie(),
			false,
		},
		{
			createFullMovie(errFrontImageID, studioID),
			nil,
			true,
		},
		{
			createFullMovie(errBackImageID, studioID),
			nil,
			true,
		},
		{
			createFullMovie(errStudioMovieID, errStudioID),
			nil,
			true,
		},
		{
			createFullMovie(missingStudioMovieID, missingStudioID),
			createFullJSONMovie("", frontImage, backImage),
			false,
		},
	}
}

func TestToJSON(t *testing.T) {
	initTestTable()

	mockMovieReader := &mocks.MovieReaderWriter{}

	imageErr := errors.New("error getting image")

	mockMovieReader.On("GetFrontImage", testCtx, movieID).Return(frontImageBytes, nil).Once()
	mockMovieReader.On("GetFrontImage", testCtx, missingStudioMovieID).Return(frontImageBytes, nil).Once()
	mockMovieReader.On("GetFrontImage", testCtx, emptyID).Return(nil, nil).Once().Maybe()
	mockMovieReader.On("GetFrontImage", testCtx, errFrontImageID).Return(nil, imageErr).Once()
	mockMovieReader.On("GetFrontImage", testCtx, errBackImageID).Return(frontImageBytes, nil).Once()

	mockMovieReader.On("GetBackImage", testCtx, movieID).Return(backImageBytes, nil).Once()
	mockMovieReader.On("GetBackImage", testCtx, missingStudioMovieID).Return(backImageBytes, nil).Once()
	mockMovieReader.On("GetBackImage", testCtx, emptyID).Return(nil, nil).Once()
	mockMovieReader.On("GetBackImage", testCtx, errBackImageID).Return(nil, imageErr).Once()
	mockMovieReader.On("GetBackImage", testCtx, errFrontImageID).Return(backImageBytes, nil).Maybe()
	mockMovieReader.On("GetBackImage", testCtx, errStudioMovieID).Return(backImageBytes, nil).Maybe()

	mockStudioReader := &mocks.StudioReaderWriter{}

	studioErr := errors.New("error getting studio")

	mockStudioReader.On("Find", testCtx, studioID).Return(&movieStudio, nil)
	mockStudioReader.On("Find", testCtx, missingStudioID).Return(nil, nil)
	mockStudioReader.On("Find", testCtx, errStudioID).Return(nil, studioErr)

	for i, s := range scenarios {
		movie := s.movie
		json, err := ToJSON(testCtx, mockMovieReader, mockStudioReader, &movie)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockMovieReader.AssertExpectations(t)
	mockStudioReader.AssertExpectations(t)
}
