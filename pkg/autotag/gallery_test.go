package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const galleryExt = "zip"

func TestGalleryPerformers(t *testing.T) {
	const galleryID = 1
	const performerName = "performer name"
	const performerID = 2
	performer := models.Performer{
		ID:   performerID,
		Name: models.NullString(performerName),
	}

	const reversedPerformerName = "name performer"
	const reversedPerformerID = 3
	reversedPerformer := models.Performer{
		ID:   reversedPerformerID,
		Name: models.NullString(reversedPerformerName),
	}

	testTables := generateTestTable(performerName, galleryExt)

	assert := assert.New(t)

	for _, test := range testTables {
		mockPerformerReader := &mocks.PerformerReaderWriter{}
		mockGalleryReader := &mocks.GalleryReaderWriter{}

		mockPerformerReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Performer{&performer, &reversedPerformer}, nil).Once()

		if test.Matches {
			mockGalleryReader.On("GetPerformerIDs", galleryID).Return(nil, nil).Once()
			mockGalleryReader.On("UpdatePerformers", galleryID, []int{performerID}).Return(nil).Once()
		}

		gallery := models.Gallery{
			ID:   galleryID,
			Path: models.NullString(test.Path),
		}
		err := GalleryPerformers(&gallery, mockGalleryReader, mockPerformerReader)

		assert.Nil(err)
		mockPerformerReader.AssertExpectations(t)
		mockGalleryReader.AssertExpectations(t)
	}
}

func TestGalleryStudios(t *testing.T) {
	const galleryID = 1
	const studioName = "studio name"
	const studioID = 2
	studio := models.Studio{
		ID:   studioID,
		Name: models.NullString(studioName),
	}

	const reversedStudioName = "name studio"
	const reversedStudioID = 3
	reversedStudio := models.Studio{
		ID:   reversedStudioID,
		Name: models.NullString(reversedStudioName),
	}

	testTables := generateTestTable(studioName, galleryExt)

	assert := assert.New(t)

	for _, test := range testTables {
		mockStudioReader := &mocks.StudioReaderWriter{}
		mockGalleryReader := &mocks.GalleryReaderWriter{}

		mockStudioReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Studio{&studio, &reversedStudio}, nil).Once()
		mockStudioReader.On("GetAliases", mock.Anything).Return([]string{}, nil).Maybe()

		if test.Matches {
			mockGalleryReader.On("Find", galleryID).Return(&models.Gallery{}, nil).Once()
			expectedStudioID := models.NullInt64(studioID)
			mockGalleryReader.On("UpdatePartial", models.GalleryPartial{
				ID:       galleryID,
				StudioID: &expectedStudioID,
			}).Return(nil, nil).Once()
		}

		gallery := models.Gallery{
			ID:   galleryID,
			Path: models.NullString(test.Path),
		}
		err := GalleryStudios(&gallery, mockGalleryReader, mockStudioReader)

		assert.Nil(err)
		mockStudioReader.AssertExpectations(t)
		mockGalleryReader.AssertExpectations(t)
	}
}

func TestGalleryTags(t *testing.T) {
	const galleryID = 1
	const tagName = "tag name"
	const tagID = 2
	tag := models.Tag{
		ID:   tagID,
		Name: tagName,
	}

	const reversedTagName = "name tag"
	const reversedTagID = 3
	reversedTag := models.Tag{
		ID:   reversedTagID,
		Name: reversedTagName,
	}

	testTables := generateTestTable(tagName, galleryExt)

	assert := assert.New(t)

	for _, test := range testTables {
		mockTagReader := &mocks.TagReaderWriter{}
		mockGalleryReader := &mocks.GalleryReaderWriter{}

		mockTagReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Tag{&tag, &reversedTag}, nil).Once()

		if test.Matches {
			mockGalleryReader.On("GetTagIDs", galleryID).Return(nil, nil).Once()
			mockGalleryReader.On("UpdateTags", galleryID, []int{tagID}).Return(nil).Once()
		}

		gallery := models.Gallery{
			ID:   galleryID,
			Path: models.NullString(test.Path),
		}
		err := GalleryTags(&gallery, mockGalleryReader, mockTagReader)

		assert.Nil(err)
		mockTagReader.AssertExpectations(t)
		mockGalleryReader.AssertExpectations(t)
	}
}
