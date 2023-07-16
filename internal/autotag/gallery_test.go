package autotag

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const galleryExt = "zip"

var testCtx = context.Background()

// returns got == expected
// ignores expected.UpdatedAt, but ensures that got.UpdatedAt is set and not null
func galleryPartialsEqual(got, expected models.GalleryPartial) bool {
	// updated at should be set and not null
	if !got.UpdatedAt.Set || got.UpdatedAt.Null {
		return false
	}
	// else ignore the exact value
	got.UpdatedAt = models.OptionalTime{}

	return assert.ObjectsAreEqual(got, expected)
}

func TestGalleryPerformers(t *testing.T) {
	t.Parallel()

	const galleryID = 1
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

	testTables := generateTestTable(performerName, galleryExt)

	assert := assert.New(t)

	for _, test := range testTables {
		db := mocks.NewDatabase()

		db.Performer.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		db.Performer.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Performer{&performer, &reversedPerformer}, nil).Once()

		if test.Matches {
			matchPartial := mock.MatchedBy(func(got models.GalleryPartial) bool {
				expected := models.GalleryPartial{
					PerformerIDs: &models.UpdateIDs{
						IDs:  []int{performerID},
						Mode: models.RelationshipUpdateModeAdd,
					},
				}

				return galleryPartialsEqual(got, expected)
			})
			db.Gallery.On("UpdatePartial", testCtx, galleryID, matchPartial).Return(nil, nil).Once()
		}

		gallery := models.Gallery{
			ID:           galleryID,
			Path:         test.Path,
			PerformerIDs: models.NewRelatedIDs([]int{}),
		}
		err := GalleryPerformers(testCtx, &gallery, db.Gallery, db.Performer, nil)

		assert.Nil(err)
		db.AssertExpectations(t)
	}
}

func TestGalleryStudios(t *testing.T) {
	t.Parallel()

	const galleryID = 1
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

	testTables := generateTestTable(studioName, galleryExt)

	assert := assert.New(t)

	doTest := func(db *mocks.Database, test pathTestTable) {
		if test.Matches {
			matchPartial := mock.MatchedBy(func(got models.GalleryPartial) bool {
				expected := models.GalleryPartial{
					StudioID: models.NewOptionalInt(studioID),
				}

				return galleryPartialsEqual(got, expected)
			})
			db.Gallery.On("UpdatePartial", testCtx, galleryID, matchPartial).Return(nil, nil).Once()
		}

		gallery := models.Gallery{
			ID:   galleryID,
			Path: test.Path,
		}
		err := GalleryStudios(testCtx, &gallery, db.Gallery, db.Studio, nil)

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

	doTest := func(db *mocks.Database, test pathTestTable) {
		if test.Matches {
			matchPartial := mock.MatchedBy(func(got models.GalleryPartial) bool {
				expected := models.GalleryPartial{
					TagIDs: &models.UpdateIDs{
						IDs:  []int{tagID},
						Mode: models.RelationshipUpdateModeAdd,
					},
				}

				return galleryPartialsEqual(got, expected)
			})
			db.Gallery.On("UpdatePartial", testCtx, galleryID, matchPartial).Return(nil, nil).Once()
		}

		gallery := models.Gallery{
			ID:     galleryID,
			Path:   test.Path,
			TagIDs: models.NewRelatedIDs([]int{}),
		}
		err := GalleryTags(testCtx, &gallery, db.Gallery, db.Tag, nil)

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
