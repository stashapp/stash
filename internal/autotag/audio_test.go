package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const audioExt = "mp3"

// asserts that got == expected
// ignores expected.UpdatedAt, but ensures that got.UpdatedAt is set and not null
func audioPartialsEqual(got, expected models.AudioPartial) bool {
	// updated at should be set and not null
	if !got.UpdatedAt.Set || got.UpdatedAt.Null {
		return false
	}
	// else ignore the exact value
	got.UpdatedAt = models.OptionalTime{}

	return assert.ObjectsAreEqual(got, expected)
}

func TestAudioPerformers(t *testing.T) {
	t.Parallel()

	const audioID = 1
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

	testTables := generateTestTable(performerName, audioExt)

	assert := assert.New(t)

	for _, test := range testTables {
		db := mocks.NewDatabase()

		db.Performer.On("Query", testCtx, mock.Anything, mock.Anything).Return(nil, 0, nil)
		db.Performer.On("QueryForAutoTag", testCtx, mock.Anything).Return([]*models.Performer{&performer, &reversedPerformer}, nil).Once()

		audio := models.Audio{
			ID:           audioID,
			Path:         test.Path,
			PerformerIDs: models.NewRelatedIDs([]int{}),
		}

		if test.Matches {
			matchPartial := mock.MatchedBy(func(got models.AudioPartial) bool {
				expected := models.AudioPartial{
					PerformerIDs: &models.UpdateIDs{
						IDs:  []int{performerID},
						Mode: models.RelationshipUpdateModeAdd,
					},
				}

				return audioPartialsEqual(got, expected)
			})
			db.Audio.On("UpdatePartial", testCtx, audioID, matchPartial).Return(nil, nil).Once()
		}

		err := AudioPerformers(testCtx, &audio, db.Audio, db.Performer, nil)

		assert.Nil(err)
		db.AssertExpectations(t)
	}
}

func TestAudioStudios(t *testing.T) {
	t.Parallel()

	var (
		audioID    = 1
		studioName = "studio name"
		studioID   = 2
	)
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

	testTables := generateTestTable(studioName, audioExt)

	assert := assert.New(t)

	doTest := func(db *mocks.Database, test pathTestTable) {
		if test.Matches {
			matchPartial := mock.MatchedBy(func(got models.AudioPartial) bool {
				expected := models.AudioPartial{
					StudioID: models.NewOptionalInt(studioID),
				}

				return audioPartialsEqual(got, expected)
			})
			db.Audio.On("UpdatePartial", testCtx, audioID, matchPartial).Return(nil, nil).Once()
		}

		audio := models.Audio{
			ID:   audioID,
			Path: test.Path,
		}
		err := AudioStudios(testCtx, &audio, db.Audio, db.Studio, nil)

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

	const unmatchedName = "unmatched"
	studio.Name = unmatchedName

	// test against aliases
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

func TestAudioTags(t *testing.T) {
	t.Parallel()

	const audioID = 1
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

	testTables := generateTestTable(tagName, audioExt)

	assert := assert.New(t)

	doTest := func(db *mocks.Database, test pathTestTable) {
		if test.Matches {
			matchPartial := mock.MatchedBy(func(got models.AudioPartial) bool {
				expected := models.AudioPartial{
					TagIDs: &models.UpdateIDs{
						IDs:  []int{tagID},
						Mode: models.RelationshipUpdateModeAdd,
					},
				}

				return audioPartialsEqual(got, expected)
			})
			db.Audio.On("UpdatePartial", testCtx, audioID, matchPartial).Return(nil, nil).Once()
		}

		audio := models.Audio{
			ID:     audioID,
			Path:   test.Path,
			TagIDs: models.NewRelatedIDs([]int{}),
		}
		err := AudioTags(testCtx, &audio, db.Audio, db.Tag, nil)

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

	// test against aliases
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
