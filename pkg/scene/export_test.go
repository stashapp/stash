package scene

import (
	"errors"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/utils"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

const (
	sceneID    = 1
	noImageID  = 2
	errImageID = 3

	studioID        = 4
	missingStudioID = 5
	errStudioID     = 6

	// noGalleryID  = 7
	// errGalleryID = 8

	noTagsID  = 11
	errTagsID = 12

	noMoviesID     = 13
	errFindMovieID = 15

	noMarkersID         = 16
	errMarkersID        = 17
	errFindPrimaryTagID = 18
	errFindByMarkerID   = 19
)

var (
	url                = "url"
	checksum           = "checksum"
	oshash             = "oshash"
	title              = "title"
	phash        int64 = -3846826108889195
	date               = "2001-01-01"
	dateObj            = models.NewDate(date)
	rating             = 5
	ocounter           = 2
	organized          = true
	details            = "details"
	size               = "size"
	duration           = 1.23
	durationStr        = "1.23"
	videoCodec         = "videoCodec"
	audioCodec         = "audioCodec"
	format             = "format"
	width              = 100
	height             = 100
	framerate          = 3.21
	framerateStr       = "3.21"
	bitrate      int64 = 1
)

var (
	studioName = "studioName"
	// galleryChecksum = "galleryChecksum"

	validMovie1  = 1
	validMovie2  = 2
	invalidMovie = 3

	movie1Name = "movie1Name"
	movie2Name = "movie2Name"

	movie1Scene = 1
	movie2Scene = 2
)

var names = []string{
	"name1",
	"name2",
}

var imageBytes = []byte("imageBytes")

var stashID = models.StashID{
	StashID:  "StashID",
	Endpoint: "Endpoint",
}

const imageBase64 = "aW1hZ2VCeXRlcw=="

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

func createFullScene(id int) models.Scene {
	return models.Scene{
		ID:         id,
		Title:      title,
		AudioCodec: &audioCodec,
		Bitrate:    &bitrate,
		Checksum:   &checksum,
		Date:       &dateObj,
		Details:    details,
		Duration:   &duration,
		Format:     &format,
		Framerate:  &framerate,
		Height:     &height,
		OCounter:   ocounter,
		OSHash:     &oshash,
		Phash:      &phash,
		Rating:     &rating,
		Organized:  organized,
		Size:       &size,
		VideoCodec: &videoCodec,
		Width:      &width,
		URL:        url,
		StashIDs: []models.StashID{
			stashID,
		},
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createEmptyScene(id int) models.Scene {
	return models.Scene{
		ID:        id,
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createFullJSONScene(image string) *jsonschema.Scene {
	return &jsonschema.Scene{
		Title:     title,
		Checksum:  checksum,
		Date:      date,
		Details:   details,
		OCounter:  ocounter,
		OSHash:    oshash,
		Phash:     utils.PhashToString(phash),
		Rating:    rating,
		Organized: organized,
		URL:       url,
		File: &jsonschema.SceneFile{
			AudioCodec: audioCodec,
			Bitrate:    int(bitrate),
			Duration:   durationStr,
			Format:     format,
			Framerate:  framerateStr,
			Height:     height,
			Size:       size,
			VideoCodec: videoCodec,
			Width:      width,
		},
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
		Cover: image,
		StashIDs: []models.StashID{
			stashID,
		},
	}
}

func createEmptyJSONScene() *jsonschema.Scene {
	return &jsonschema.Scene{
		File: &jsonschema.SceneFile{},
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
	}
}

type basicTestScenario struct {
	input    models.Scene
	expected *jsonschema.Scene
	err      bool
}

var scenarios = []basicTestScenario{
	{
		createFullScene(sceneID),
		createFullJSONScene(imageBase64),
		false,
	},
	{
		createEmptyScene(noImageID),
		createEmptyJSONScene(),
		false,
	},
	{
		createFullScene(errImageID),
		nil,
		true,
	},
}

func TestToJSON(t *testing.T) {
	mockSceneReader := &mocks.SceneReaderWriter{}

	imageErr := errors.New("error getting image")

	mockSceneReader.On("GetCover", testCtx, sceneID).Return(imageBytes, nil).Once()
	mockSceneReader.On("GetCover", testCtx, noImageID).Return(nil, nil).Once()
	mockSceneReader.On("GetCover", testCtx, errImageID).Return(nil, imageErr).Once()

	for i, s := range scenarios {
		scene := s.input
		json, err := ToBasicJSON(testCtx, mockSceneReader, &scene)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockSceneReader.AssertExpectations(t)
}

func createStudioScene(studioID int) models.Scene {
	return models.Scene{
		StudioID: &studioID,
	}
}

type stringTestScenario struct {
	input    models.Scene
	expected string
	err      bool
}

var getStudioScenarios = []stringTestScenario{
	{
		createStudioScene(studioID),
		studioName,
		false,
	},
	{
		createStudioScene(missingStudioID),
		"",
		false,
	},
	{
		createStudioScene(errStudioID),
		"",
		true,
	},
}

func TestGetStudioName(t *testing.T) {
	mockStudioReader := &mocks.StudioReaderWriter{}

	studioErr := errors.New("error getting image")

	mockStudioReader.On("Find", testCtx, studioID).Return(&models.Studio{
		Name: models.NullString(studioName),
	}, nil).Once()
	mockStudioReader.On("Find", testCtx, missingStudioID).Return(nil, nil).Once()
	mockStudioReader.On("Find", testCtx, errStudioID).Return(nil, studioErr).Once()

	for i, s := range getStudioScenarios {
		scene := s.input
		json, err := GetStudioName(testCtx, mockStudioReader, &scene)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockStudioReader.AssertExpectations(t)
}

type stringSliceTestScenario struct {
	input    models.Scene
	expected []string
	err      bool
}

var getTagNamesScenarios = []stringSliceTestScenario{
	{
		createEmptyScene(sceneID),
		names,
		false,
	},
	{
		createEmptyScene(noTagsID),
		nil,
		false,
	},
	{
		createEmptyScene(errTagsID),
		nil,
		true,
	},
}

func getTags(names []string) []*models.Tag {
	var ret []*models.Tag
	for _, n := range names {
		ret = append(ret, &models.Tag{
			Name: n,
		})
	}

	return ret
}

func TestGetTagNames(t *testing.T) {
	mockTagReader := &mocks.TagReaderWriter{}

	tagErr := errors.New("error getting tag")

	mockTagReader.On("FindBySceneID", testCtx, sceneID).Return(getTags(names), nil).Once()
	mockTagReader.On("FindBySceneID", testCtx, noTagsID).Return(nil, nil).Once()
	mockTagReader.On("FindBySceneID", testCtx, errTagsID).Return(nil, tagErr).Once()

	for i, s := range getTagNamesScenarios {
		scene := s.input
		json, err := GetTagNames(testCtx, mockTagReader, &scene)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockTagReader.AssertExpectations(t)
}

type sceneMoviesTestScenario struct {
	input    models.Scene
	expected []jsonschema.SceneMovie
	err      bool
}

var validMovies = []models.MoviesScenes{
	{
		MovieID:    validMovie1,
		SceneIndex: &movie1Scene,
	},
	{
		MovieID:    validMovie2,
		SceneIndex: &movie2Scene,
	},
}

var invalidMovies = []models.MoviesScenes{
	{
		MovieID:    invalidMovie,
		SceneIndex: &movie1Scene,
	},
}

var getSceneMoviesJSONScenarios = []sceneMoviesTestScenario{
	{
		models.Scene{
			ID:     sceneID,
			Movies: validMovies,
		},
		[]jsonschema.SceneMovie{
			{
				MovieName:  movie1Name,
				SceneIndex: movie1Scene,
			},
			{
				MovieName:  movie2Name,
				SceneIndex: movie2Scene,
			},
		},
		false,
	},
	{
		models.Scene{
			ID: noMoviesID,
		},
		nil,
		false,
	},
	{
		models.Scene{
			ID:     errFindMovieID,
			Movies: invalidMovies,
		},
		nil,
		true,
	},
}

func TestGetSceneMoviesJSON(t *testing.T) {
	mockMovieReader := &mocks.MovieReaderWriter{}
	movieErr := errors.New("error getting movie")

	mockMovieReader.On("Find", testCtx, validMovie1).Return(&models.Movie{
		Name: models.NullString(movie1Name),
	}, nil).Once()
	mockMovieReader.On("Find", testCtx, validMovie2).Return(&models.Movie{
		Name: models.NullString(movie2Name),
	}, nil).Once()
	mockMovieReader.On("Find", testCtx, invalidMovie).Return(nil, movieErr).Once()

	for i, s := range getSceneMoviesJSONScenarios {
		scene := s.input
		json, err := GetSceneMoviesJSON(testCtx, mockMovieReader, &scene)

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
}

const (
	validMarkerID1 = 1
	validMarkerID2 = 2

	invalidMarkerID1 = 3
	invalidMarkerID2 = 4

	validTagID1 = 1
	validTagID2 = 2

	validTagName1 = "validTagName1"
	validTagName2 = "validTagName2"

	invalidTagID = 3

	markerTitle1 = "markerTitle1"
	markerTitle2 = "markerTitle2"

	markerSeconds1 = 1.0
	markerSeconds2 = 2.3

	markerSeconds1Str = "1.0"
	markerSeconds2Str = "2.3"
)

type sceneMarkersTestScenario struct {
	input    models.Scene
	expected []jsonschema.SceneMarker
	err      bool
}

var getSceneMarkersJSONScenarios = []sceneMarkersTestScenario{
	{
		createEmptyScene(sceneID),
		[]jsonschema.SceneMarker{
			{
				Title:      markerTitle1,
				PrimaryTag: validTagName1,
				Seconds:    markerSeconds1Str,
				Tags: []string{
					validTagName1,
					validTagName2,
				},
				CreatedAt: json.JSONTime{
					Time: createTime,
				},
				UpdatedAt: json.JSONTime{
					Time: updateTime,
				},
			},
			{
				Title:      markerTitle2,
				PrimaryTag: validTagName2,
				Seconds:    markerSeconds2Str,
				Tags: []string{
					validTagName2,
				},
				CreatedAt: json.JSONTime{
					Time: createTime,
				},
				UpdatedAt: json.JSONTime{
					Time: updateTime,
				},
			},
		},
		false,
	},
	{
		createEmptyScene(noMarkersID),
		nil,
		false,
	},
	{
		createEmptyScene(errMarkersID),
		nil,
		true,
	},
	{
		createEmptyScene(errFindPrimaryTagID),
		nil,
		true,
	},
	{
		createEmptyScene(errFindByMarkerID),
		nil,
		true,
	},
}

var validMarkers = []*models.SceneMarker{
	{
		ID:           validMarkerID1,
		Title:        markerTitle1,
		PrimaryTagID: validTagID1,
		Seconds:      markerSeconds1,
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	},
	{
		ID:           validMarkerID2,
		Title:        markerTitle2,
		PrimaryTagID: validTagID2,
		Seconds:      markerSeconds2,
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	},
}

var invalidMarkers1 = []*models.SceneMarker{
	{
		ID:           invalidMarkerID1,
		PrimaryTagID: invalidTagID,
	},
}

var invalidMarkers2 = []*models.SceneMarker{
	{
		ID:           invalidMarkerID2,
		PrimaryTagID: validTagID1,
	},
}

func TestGetSceneMarkersJSON(t *testing.T) {
	mockTagReader := &mocks.TagReaderWriter{}
	mockMarkerReader := &mocks.SceneMarkerReaderWriter{}

	markersErr := errors.New("error getting scene markers")
	tagErr := errors.New("error getting tags")

	mockMarkerReader.On("FindBySceneID", testCtx, sceneID).Return(validMarkers, nil).Once()
	mockMarkerReader.On("FindBySceneID", testCtx, noMarkersID).Return(nil, nil).Once()
	mockMarkerReader.On("FindBySceneID", testCtx, errMarkersID).Return(nil, markersErr).Once()
	mockMarkerReader.On("FindBySceneID", testCtx, errFindPrimaryTagID).Return(invalidMarkers1, nil).Once()
	mockMarkerReader.On("FindBySceneID", testCtx, errFindByMarkerID).Return(invalidMarkers2, nil).Once()

	mockTagReader.On("Find", testCtx, validTagID1).Return(&models.Tag{
		Name: validTagName1,
	}, nil)
	mockTagReader.On("Find", testCtx, validTagID2).Return(&models.Tag{
		Name: validTagName2,
	}, nil)
	mockTagReader.On("Find", testCtx, invalidTagID).Return(nil, tagErr)

	mockTagReader.On("FindBySceneMarkerID", testCtx, validMarkerID1).Return([]*models.Tag{
		{
			Name: validTagName1,
		},
		{
			Name: validTagName2,
		},
	}, nil)
	mockTagReader.On("FindBySceneMarkerID", testCtx, validMarkerID2).Return([]*models.Tag{
		{
			Name: validTagName2,
		},
	}, nil)
	mockTagReader.On("FindBySceneMarkerID", testCtx, invalidMarkerID2).Return(nil, tagErr).Once()

	for i, s := range getSceneMarkersJSONScenarios {
		scene := s.input
		json, err := GetSceneMarkersJSON(testCtx, mockMarkerReader, mockTagReader, &scene)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockTagReader.AssertExpectations(t)
}
