package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const imageExt = "jpg"

func TestImagePerformers(t *testing.T) {
	t.Parallel()

	const imageID = 1
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

	testTables := generateTestTable(performerName, imageExt)

	assert := assert.New(t)

	for _, test := range testTables {
		mockPerformerReader := &mocks.PerformerReaderWriter{}
		mockImageReader := &mocks.ImageReaderWriter{}

		mockPerformerReader.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockPerformerReader.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Performer{&performer, &reversedPerformer}, nil).Once()

		if test.Matches {
			mockImageReader.On("UpdatePartial", testCtx, imageID, models.ImagePartial{
				PerformerIDs: &models.UpdateIDs{
					IDs:  []int{performerID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			}).Return(nil, nil).Once()
		}

		image := models.Image{
			ID:   imageID,
			Path: test.Path,
		}
		err := ImagePerformers(testCtx, &image, mockImageReader, mockPerformerReader, nil)

		assert.Nil(err)
		mockPerformerReader.AssertExpectations(t)
		mockImageReader.AssertExpectations(t)
	}
}

func TestImageStudios(t *testing.T) {
	t.Parallel()

	const imageID = 1
	const studioName = "studio name"
	var studioID = 2
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

	testTables := generateTestTable(studioName, imageExt)

	assert := assert.New(t)

	doTest := func(mockStudioReader *mocks.StudioReaderWriter, mockImageReader *mocks.ImageReaderWriter, test pathTestTable) {
		if test.Matches {
			expectedStudioID := studioID
			mockImageReader.On("UpdatePartial", testCtx, imageID, models.ImagePartial{
				StudioID: models.NewOptionalInt(expectedStudioID),
			}).Return(nil, nil).Once()
		}

		image := models.Image{
			ID:   imageID,
			Path: test.Path,
		}
		err := ImageStudios(testCtx, &image, mockImageReader, mockStudioReader, nil)

		assert.Nil(err)
		mockStudioReader.AssertExpectations(t)
		mockImageReader.AssertExpectations(t)
	}

	for _, test := range testTables {
		mockStudioReader := &mocks.StudioReaderWriter{}
		mockImageReader := &mocks.ImageReaderWriter{}

		mockStudioReader.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockStudioReader.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Studio{&studio, &reversedStudio}, nil).Once()
		mockStudioReader.On("GetAliases", testCtx, mock.Anything).Return([]string{}, nil).Maybe()

		doTest(mockStudioReader, mockImageReader, test)
	}

	// test against aliases
	const unmatchedName = "unmatched"
	studio.Name.String = unmatchedName

	for _, test := range testTables {
		mockStudioReader := &mocks.StudioReaderWriter{}
		mockImageReader := &mocks.ImageReaderWriter{}

		mockStudioReader.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockStudioReader.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Studio{&studio, &reversedStudio}, nil).Once()
		mockStudioReader.On("GetAliases", testCtx, studioID).Return([]string{
			studioName,
		}, nil).Once()
		mockStudioReader.On("GetAliases", testCtx, reversedStudioID).Return([]string{}, nil).Once()

		doTest(mockStudioReader, mockImageReader, test)
	}
}

func TestImageTags(t *testing.T) {
	t.Parallel()

	const imageID = 1
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

	testTables := generateTestTable(tagName, imageExt)

	assert := assert.New(t)

	doTest := func(mockTagReader *mocks.TagReaderWriter, mockImageReader *mocks.ImageReaderWriter, test pathTestTable) {
		if test.Matches {
			mockImageReader.On("UpdatePartial", testCtx, imageID, models.ImagePartial{
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagID},
					Mode: models.RelationshipUpdateModeAdd,
				},
			}).Return(nil, nil).Once()
		}

		image := models.Image{
			ID:   imageID,
			Path: test.Path,
		}
		err := ImageTags(testCtx, &image, mockImageReader, mockTagReader, nil)

		assert.Nil(err)
		mockTagReader.AssertExpectations(t)
		mockImageReader.AssertExpectations(t)
	}

	for _, test := range testTables {
		mockTagReader := &mocks.TagReaderWriter{}
		mockImageReader := &mocks.ImageReaderWriter{}

		mockTagReader.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockTagReader.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Tag{&tag, &reversedTag}, nil).Once()
		mockTagReader.On("GetAliases", testCtx, mock.Anything).Return([]string{}, nil).Maybe()

		doTest(mockTagReader, mockImageReader, test)
	}

	// test against aliases
	const unmatchedName = "unmatched"
	tag.Name = unmatchedName

	for _, test := range testTables {
		mockTagReader := &mocks.TagReaderWriter{}
		mockImageReader := &mocks.ImageReaderWriter{}

		mockTagReader.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		mockTagReader.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Tag{&tag, &reversedTag}, nil).Once()
		mockTagReader.On("GetAliases", testCtx, tagID).Return([]string{
			tagName,
		}, nil).Once()
		mockTagReader.On("GetAliases", testCtx, reversedTagID).Return([]string{}, nil).Once()

		doTest(mockTagReader, mockImageReader, test)
	}
}
