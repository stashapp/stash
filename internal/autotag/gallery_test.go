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
	t.Parallel()

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

		mockPerformerReader.On("Query", mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockPerformerReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Performer{&performer, &reversedPerformer}, nil).Once()

		if test.Matches {
			mockGalleryReader.On("GetPerformerIDs", galleryID).Return(nil, nil).Once()
			mockGalleryReader.On("UpdatePerformers", galleryID, []int{performerID}).Return(nil).Once()
		}

		gallery := models.Gallery{
			ID:   galleryID,
			Path: models.NullString(test.Path),
		}
		err := GalleryPerformers(&gallery, mockGalleryReader, mockPerformerReader, nil)

		assert.Nil(err)
		mockPerformerReader.AssertExpectations(t)
		mockGalleryReader.AssertExpectations(t)
	}
}

func TestGalleryStudios(t *testing.T) {
	t.Parallel()

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

	doTest := func(mockStudioReader *mocks.StudioReaderWriter, mockGalleryReader *mocks.GalleryReaderWriter, test pathTestTable) {
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
		err := GalleryStudios(&gallery, mockGalleryReader, mockStudioReader, nil)

		assert.Nil(err)
		mockStudioReader.AssertExpectations(t)
		mockGalleryReader.AssertExpectations(t)
	}

	for _, test := range testTables {
		mockStudioReader := &mocks.StudioReaderWriter{}
		mockGalleryReader := &mocks.GalleryReaderWriter{}

		mockStudioReader.On("Query", mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockStudioReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Studio{&studio, &reversedStudio}, nil).Once()
		mockStudioReader.On("GetAliases", mock.Anything).Return([]string{}, nil).Maybe()

		doTest(mockStudioReader, mockGalleryReader, test)
	}

	// test against aliases
	const unmatchedName = "unmatched"
	studio.Name.String = unmatchedName

	for _, test := range testTables {
		mockStudioReader := &mocks.StudioReaderWriter{}
		mockGalleryReader := &mocks.GalleryReaderWriter{}

		mockStudioReader.On("Query", mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockStudioReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Studio{&studio, &reversedStudio}, nil).Once()
		mockStudioReader.On("GetAliases", studioID).Return([]string{
			studioName,
		}, nil).Once()
		mockStudioReader.On("GetAliases", reversedStudioID).Return([]string{}, nil).Once()

		doTest(mockStudioReader, mockGalleryReader, test)
	}
}

func TestGalleryTags(t *testing.T) {
	t.Parallel()

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

	doTest := func(mockTagReader *mocks.TagReaderWriter, mockGalleryReader *mocks.GalleryReaderWriter, test pathTestTable) {
		if test.Matches {
			mockGalleryReader.On("GetTagIDs", galleryID).Return(nil, nil).Once()
			mockGalleryReader.On("UpdateTags", galleryID, []int{tagID}).Return(nil).Once()
		}

		gallery := models.Gallery{
			ID:   galleryID,
			Path: models.NullString(test.Path),
		}
		err := GalleryTags(&gallery, mockGalleryReader, mockTagReader, nil)

		assert.Nil(err)
		mockTagReader.AssertExpectations(t)
		mockGalleryReader.AssertExpectations(t)
	}

	for _, test := range testTables {
		mockTagReader := &mocks.TagReaderWriter{}
		mockGalleryReader := &mocks.GalleryReaderWriter{}

		mockTagReader.On("Query", mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockTagReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Tag{&tag, &reversedTag}, nil).Once()
		mockTagReader.On("GetAliases", mock.Anything).Return([]string{}, nil).Maybe()

		doTest(mockTagReader, mockGalleryReader, test)
	}

	const unmatchedName = "unmatched"
	tag.Name = unmatchedName

	for _, test := range testTables {
		mockTagReader := &mocks.TagReaderWriter{}
		mockGalleryReader := &mocks.GalleryReaderWriter{}

		mockTagReader.On("Query", mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockTagReader.On("QueryForAutoTag", mock.Anything).Return([]*models.Tag{&tag, &reversedTag}, nil).Once()
		mockTagReader.On("GetAliases", tagID).Return([]string{
			tagName,
		}, nil).Once()
		mockTagReader.On("GetAliases", reversedTagID).Return([]string{}, nil).Once()

		doTest(mockTagReader, mockGalleryReader, test)
	}
}
