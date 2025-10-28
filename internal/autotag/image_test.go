package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const imageExt = "jpg"

// returns got == expected
// ignores expected.UpdatedAt, but ensures that got.UpdatedAt is set and not null
func imagePartialsEqual(got, expected models.ImagePartial) bool {
	// updated at should be set and not null
	if !got.UpdatedAt.Set || got.UpdatedAt.Null {
		return false
	}
	// else ignore the exact value
	got.UpdatedAt = models.OptionalTime{}

	return assert.ObjectsAreEqual(got, expected)
}

func TestImagePerformers(t *testing.T) {
	t.Parallel()

	const imageID = 1
	const performerName = "performer name"
	const performerID = 2
	performer := models.Performer{
		ID:      performerID,
		Name:    performerName,
		Aliases: models.NewRelatedStrings([]string{}),
	}

	const reversedPerformerName = "name performer"
	const reversedPerformerID = 3
	reversedPerformer := models.Performer{
		ID:      reversedPerformerID,
		Name:    reversedPerformerName,
		Aliases: models.NewRelatedStrings([]string{}),
	}

	testTables := generateTestTable(performerName, imageExt)

	assert := assert.New(t)

	for _, test := range testTables {
		db := mocks.NewDatabase()

		db.Performer.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		db.Performer.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Performer{&performer, &reversedPerformer}, nil).Once()

		if test.Matches {
			matchPartial := mock.MatchedBy(func(got models.ImagePartial) bool {
				expected := models.ImagePartial{
					PerformerIDs: &models.UpdateIDs{
						IDs:  []int{performerID},
						Mode: models.RelationshipUpdateModeAdd,
					},
				}

				return imagePartialsEqual(got, expected)
			})
			db.Image.On("UpdatePartial", testCtx, imageID, matchPartial).Return(nil, nil).Once()
		}

		image := models.Image{
			ID:           imageID,
			Path:         test.Path,
			PerformerIDs: models.NewRelatedIDs([]int{}),
		}
		err := ImagePerformers(testCtx, &image, db.Image, db.Performer, nil)

		assert.Nil(err)
		db.AssertExpectations(t)
	}
}

func TestImageStudios(t *testing.T) {
	t.Parallel()

	const imageID = 1
	const studioName = "studio name"
	var studioID = 2
	studio := models.Studio{
		ID:   studioID,
		Name: studioName,
	}

	const reversedStudioName = "name studio"
	const reversedStudioID = 3
	reversedStudio := models.Studio{
		ID:   reversedStudioID,
		Name: reversedStudioName,
	}

	testTables := generateTestTable(studioName, imageExt)

	assert := assert.New(t)

	doTest := func(db *mocks.Database, test pathTestTable) {
		if test.Matches {
			matchPartial := mock.MatchedBy(func(got models.ImagePartial) bool {
				expected := models.ImagePartial{
					StudioID: models.NewOptionalInt(studioID),
				}

				return imagePartialsEqual(got, expected)
			})
			db.Image.On("UpdatePartial", testCtx, imageID, matchPartial).Return(nil, nil).Once()
		}

		image := models.Image{
			ID:   imageID,
			Path: test.Path,
		}
		err := ImageStudios(testCtx, &image, db.Image, db.Studio, nil)

		assert.Nil(err)
		db.AssertExpectations(t)
	}

	for _, test := range testTables {
		db := mocks.NewDatabase()

		db.Studio.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		db.Studio.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Studio{&studio, &reversedStudio}, nil).Once()
		db.Studio.On("GetAliases", testCtx, mock.Anything).Return([]string{}, nil).Maybe()

		doTest(db, test)
	}

	// test against aliases
	const unmatchedName = "unmatched"
	studio.Name = unmatchedName

	for _, test := range testTables {
		db := mocks.NewDatabase()

		db.Studio.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		db.Studio.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Studio{&studio, &reversedStudio}, nil).Once()
		db.Studio.On("GetAliases", testCtx, studioID).Return([]string{
			studioName,
		}, nil).Once()
		db.Studio.On("GetAliases", testCtx, reversedStudioID).Return([]string{}, nil).Once()

		doTest(db, test)
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

	doTest := func(db *mocks.Database, test pathTestTable) {
		if test.Matches {
			matchPartial := mock.MatchedBy(func(got models.ImagePartial) bool {
				expected := models.ImagePartial{
					TagIDs: &models.UpdateIDs{
						IDs:  []int{tagID},
						Mode: models.RelationshipUpdateModeAdd,
					},
				}

				return imagePartialsEqual(got, expected)
			})
			db.Image.On("UpdatePartial", testCtx, imageID, matchPartial).Return(nil, nil).Once()
		}

		image := models.Image{
			ID:     imageID,
			Path:   test.Path,
			TagIDs: models.NewRelatedIDs([]int{}),
		}
		err := ImageTags(testCtx, &image, db.Image, db.Tag, nil)

		assert.Nil(err)
		db.AssertExpectations(t)
	}

	for _, test := range testTables {
		db := mocks.NewDatabase()

		db.Tag.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		db.Tag.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Tag{&tag, &reversedTag}, nil).Once()
		db.Tag.On("GetAliases", testCtx, mock.Anything).Return([]string{}, nil).Maybe()

		doTest(db, test)
	}

	// test against aliases
	const unmatchedName = "unmatched"
	tag.Name = unmatchedName

	for _, test := range testTables {
		db := mocks.NewDatabase()

		db.Tag.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		db.Tag.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Tag{&tag, &reversedTag}, nil).Once()
		db.Tag.On("GetAliases", testCtx, tagID).Return([]string{
			tagName,
		}, nil).Once()
		db.Tag.On("GetAliases", testCtx, reversedTagID).Return([]string{}, nil).Once()

		doTest(db, test)
	}
}
