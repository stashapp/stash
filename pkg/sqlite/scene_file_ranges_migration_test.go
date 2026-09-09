//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stashapp/stash/pkg/models"
)

// TestSceneFileRangesMigration verifies that the test database - which is
// created through the full migration chain - ends at the new schema version
// and that scenes created without ranges keep whole-file semantics.
func TestSceneFileRangesMigration(t *testing.T) {
	// the standard test harness runs the full migration chain up to
	// appSchemaVersion; reaching this point proves 86->87 migrated cleanly.
	assert.Equal(t, db.Version(), uint(87))

	require.NoError(t, withRollbackTxn(func(ctx context.Context) error {
		// legacy-style scene: no ranges involved
		s := &models.Scene{}
		if err := db.Scene.Create(ctx, s, []models.FileID{fileIDs[fileIdxStartVideoFiles]}); err != nil {
			return err
		}

		ranges, err := db.Scene.GetFileRanges(ctx, s.ID)
		if err != nil {
			return err
		}
		assert.Empty(t, ranges, "scene created without ranges must have none")

		// existing fixture scenes (created by the populate step, which does
		// not use ranges) must likewise expose no ranges
		fixtureRanges, err := db.Scene.GetFileRanges(ctx, sceneIDs[sceneIdxWithGallery])
		if err != nil {
			return err
		}
		assert.Empty(t, fixtureRanges, "pre-existing scenes must keep whole-file semantics")

		return nil
	}))
}
