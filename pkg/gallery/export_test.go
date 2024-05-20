package gallery

import (
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
	galleryID = 1

	studioID        = 4
	missingStudioID = 5
	errStudioID     = 6

	// noTagsID  = 11
	noChaptersID       = 7
	errChaptersID      = 8
	errFindByChapterID = 9
)

var (
	url        = "url"
	title      = "title"
	date       = "2001-01-01"
	dateObj, _ = models.ParseDate(date)
	rating     = 5
	organized  = true
	details    = "details"
)

const (
	studioName = "studioName"
	path       = "path"
)

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
)

func createFullGallery(id int) models.Gallery {
	return models.Gallery{
		ID: id,
		Files: models.NewRelatedFiles([]models.File{
			&models.BaseFile{
				Path: path,
			},
		}),
		Title:     title,
		Date:      &dateObj,
		Details:   details,
		Rating:    &rating,
		Organized: organized,
		URLs:      models.NewRelatedStrings([]string{url}),
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createEmptyGallery(id int) models.Gallery {
	return models.Gallery{
		ID: id,
		Files: models.NewRelatedFiles([]models.File{
			&models.BaseFile{
				Path: path,
			},
		}),
		CreatedAt: createTime,
		UpdatedAt: updateTime,
	}
}

func createFullJSONGallery() *jsonschema.Gallery {
	return &jsonschema.Gallery{
		Title:     title,
		Date:      date,
		Details:   details,
		Rating:    rating,
		Organized: organized,
		URLs:      []string{url},
		ZipFiles:  []string{path},
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
	}
}

type basicTestScenario struct {
	input    models.Gallery
	expected *jsonschema.Gallery
	err      bool
}

var scenarios = []basicTestScenario{
	{
		createFullGallery(galleryID),
		createFullJSONGallery(),
		false,
	},
}

func TestToJSON(t *testing.T) {
	for i, s := range scenarios {
		gallery := s.input
		json, err := ToBasicJSON(&gallery)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}
}

func createStudioGallery(studioID int) models.Gallery {
	return models.Gallery{
		StudioID: &studioID,
	}
}

type stringTestScenario struct {
	input    models.Gallery
	expected string
	err      bool
}

var getStudioScenarios = []stringTestScenario{
	{
		createStudioGallery(studioID),
		studioName,
		false,
	},
	{
		createStudioGallery(missingStudioID),
		"",
		false,
	},
	{
		createStudioGallery(errStudioID),
		"",
		true,
	},
}

func TestGetStudioName(t *testing.T) {
	db := mocks.NewDatabase()

	studioErr := errors.New("error getting image")

	db.Studio.On("Find", testCtx, studioID).Return(&models.Studio{
		Name: studioName,
	}, nil).Once()
	db.Studio.On("Find", testCtx, missingStudioID).Return(nil, nil).Once()
	db.Studio.On("Find", testCtx, errStudioID).Return(nil, studioErr).Once()

	for i, s := range getStudioScenarios {
		gallery := s.input
		json, err := GetStudioName(testCtx, db.Studio, &gallery)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	db.AssertExpectations(t)
}

const (
	validChapterID1 = 1
	validChapterID2 = 2

	chapterTitle1 = "chapterTitle1"
	chapterTitle2 = "chapterTitle2"

	chapterImageIndex1 = 10
	chapterImageIndex2 = 50
)

type galleryChaptersTestScenario struct {
	input    models.Gallery
	expected []jsonschema.GalleryChapter
	err      bool
}

var getGalleryChaptersJSONScenarios = []galleryChaptersTestScenario{
	{
		createEmptyGallery(galleryID),
		[]jsonschema.GalleryChapter{
			{
				Title:      chapterTitle1,
				ImageIndex: chapterImageIndex1,
				CreatedAt: json.JSONTime{
					Time: createTime,
				},
				UpdatedAt: json.JSONTime{
					Time: updateTime,
				},
			},
			{
				Title:      chapterTitle2,
				ImageIndex: chapterImageIndex2,
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
		createEmptyGallery(noChaptersID),
		nil,
		false,
	},
	{
		createEmptyGallery(errChaptersID),
		nil,
		true,
	},
}

var validChapters = []*models.GalleryChapter{
	{
		ID:         validChapterID1,
		Title:      chapterTitle1,
		ImageIndex: chapterImageIndex1,
		CreatedAt:  createTime,
		UpdatedAt:  updateTime,
	},
	{
		ID:         validChapterID2,
		Title:      chapterTitle2,
		ImageIndex: chapterImageIndex2,
		CreatedAt:  createTime,
		UpdatedAt:  updateTime,
	},
}

func TestGetGalleryChaptersJSON(t *testing.T) {
	db := mocks.NewDatabase()

	chaptersErr := errors.New("error getting gallery chapters")

	db.GalleryChapter.On("FindByGalleryID", testCtx, galleryID).Return(validChapters, nil).Once()
	db.GalleryChapter.On("FindByGalleryID", testCtx, noChaptersID).Return(nil, nil).Once()
	db.GalleryChapter.On("FindByGalleryID", testCtx, errChaptersID).Return(nil, chaptersErr).Once()

	for i, s := range getGalleryChaptersJSONScenarios {
		gallery := s.input
		json, err := GetGalleryChaptersJSON(testCtx, db.GalleryChapter, &gallery)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	db.AssertExpectations(t)
}
