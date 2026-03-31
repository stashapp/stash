//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestFingerprintSubmissionCreate(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		submission := &models.FingerprintSubmission{
			Endpoint:  "https://stashdb.org/graphql",
			StashID:   "test-stash-id-1",
			SceneID:   sceneIDs[sceneIdxWithGallery],
			Vote:      models.FingerprintVoteInvalid,
			CreatedAt: time.Now(),
		}

		err := db.FingerprintSubmission.Create(ctx, submission)
		assert.NoError(t, err)

		// Verify it was created
		found, err := db.FingerprintSubmission.FindByEndpoint(ctx, submission.Endpoint)
		assert.NoError(t, err)
		assert.Len(t, found, 1)
		assert.Equal(t, submission.StashID, found[0].StashID)
		assert.Equal(t, submission.SceneID, found[0].SceneID)
		assert.Equal(t, submission.Vote, found[0].Vote)

		return nil
	})
}

func TestFingerprintSubmissionCreateDuplicate(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		submission := &models.FingerprintSubmission{
			Endpoint:  "https://stashdb.org/graphql",
			StashID:   "test-stash-id-dup",
			SceneID:   sceneIDs[sceneIdxWithGallery],
			Vote:      models.FingerprintVoteValid,
			CreatedAt: time.Now(),
		}

		err := db.FingerprintSubmission.Create(ctx, submission)
		assert.NoError(t, err)

		// Creating again with same endpoint+stash_id should not error (ON CONFLICT DO NOTHING)
		submission2 := &models.FingerprintSubmission{
			Endpoint:  "https://stashdb.org/graphql",
			StashID:   "test-stash-id-dup",
			SceneID:   sceneIDs[sceneIdxWithPerformer],
			Vote:      models.FingerprintVoteInvalid,
			CreatedAt: time.Now(),
		}

		err = db.FingerprintSubmission.Create(ctx, submission2)
		assert.NoError(t, err)

		// Original should still exist unchanged
		found, err := db.FingerprintSubmission.FindByEndpoint(ctx, submission.Endpoint)
		assert.NoError(t, err)
		assert.Len(t, found, 1)
		assert.Equal(t, models.FingerprintVoteValid, found[0].Vote)

		return nil
	})
}

func TestFingerprintSubmissionFindByEndpoint(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		endpoint := "https://test-endpoint.org/graphql"

		// Create multiple submissions for the same endpoint
		for i := 0; i < 3; i++ {
			submission := &models.FingerprintSubmission{
				Endpoint:  endpoint,
				StashID:   "stash-id-" + string(rune('a'+i)),
				SceneID:   sceneIDs[sceneIdxWithGallery],
				Vote:      models.FingerprintVoteInvalid,
				CreatedAt: time.Now(),
			}
			err := db.FingerprintSubmission.Create(ctx, submission)
			assert.NoError(t, err)
		}

		// Create one for a different endpoint
		otherSubmission := &models.FingerprintSubmission{
			Endpoint:  "https://other-endpoint.org/graphql",
			StashID:   "other-stash-id",
			SceneID:   sceneIDs[sceneIdxWithGallery],
			Vote:      models.FingerprintVoteValid,
			CreatedAt: time.Now(),
		}
		err := db.FingerprintSubmission.Create(ctx, otherSubmission)
		assert.NoError(t, err)

		// Find by endpoint should return only the 3
		found, err := db.FingerprintSubmission.FindByEndpoint(ctx, endpoint)
		assert.NoError(t, err)
		assert.Len(t, found, 3)

		return nil
	})
}

func TestFingerprintSubmissionDelete(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		submission := &models.FingerprintSubmission{
			Endpoint:  "https://delete-test.org/graphql",
			StashID:   "delete-test-stash-id",
			SceneID:   sceneIDs[sceneIdxWithGallery],
			Vote:      models.FingerprintVoteInvalid,
			CreatedAt: time.Now(),
		}

		err := db.FingerprintSubmission.Create(ctx, submission)
		assert.NoError(t, err)

		// Delete it
		err = db.FingerprintSubmission.Delete(ctx, submission.Endpoint, submission.StashID)
		assert.NoError(t, err)

		// Verify it's gone
		found, err := db.FingerprintSubmission.FindByEndpoint(ctx, submission.Endpoint)
		assert.NoError(t, err)
		assert.Len(t, found, 0)

		return nil
	})
}

func TestFingerprintSubmissionDeleteByEndpoint(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		endpoint := "https://delete-all-test.org/graphql"

		// Create multiple submissions
		for i := 0; i < 3; i++ {
			submission := &models.FingerprintSubmission{
				Endpoint:  endpoint,
				StashID:   "delete-all-stash-id-" + string(rune('a'+i)),
				SceneID:   sceneIDs[sceneIdxWithGallery],
				Vote:      models.FingerprintVoteInvalid,
				CreatedAt: time.Now(),
			}
			err := db.FingerprintSubmission.Create(ctx, submission)
			assert.NoError(t, err)
		}

		// Delete all by endpoint
		err := db.FingerprintSubmission.DeleteByEndpoint(ctx, endpoint)
		assert.NoError(t, err)

		// Verify all are gone
		found, err := db.FingerprintSubmission.FindByEndpoint(ctx, endpoint)
		assert.NoError(t, err)
		assert.Len(t, found, 0)

		return nil
	})
}
