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

// createRangeTestFile creates a fresh video file usable for range tests,
// avoiding the scene fixtures which already claim their files exclusively.
func createRangeTestFile(ctx context.Context, t *testing.T) models.FileID {
	f := &models.VideoFile{
		BaseFile: &models.BaseFile{
			Path:           getFilePath(folderIdxWithSceneFiles, "ranges_test.mp4"),
			Basename:       "ranges_test.mp4",
			ParentFolderID: folderIDs[folderIdxWithSceneFiles],
		},
		Duration: 100,
		Height:    100,
		Width:     200,
	}

	require.NoError(t, db.File.Create(ctx, f))
	require.NotZero(t, f.ID)
	return f.ID
}

func TestSceneFileRanges(t *testing.T) {
	qb := db.Scene

	runWithRollbackTxn(t, "set and get range", func(t *testing.T, ctx context.Context) {
		fileID := createRangeTestFile(ctx, t)

		s := &models.Scene{}
		require.NoError(t, qb.Create(ctx, s, []models.FileID{fileID}))
		require.NotZero(t, s.ID)

		var (
			start = 10.5
			end   = 20.25
		)

		// no ranges on a fresh scene
		ranges, err := qb.GetFileRanges(ctx, s.ID)
		require.NoError(t, err)
		assert.Empty(t, ranges)

		// set a range
		require.NoError(t, qb.SetFileRange(ctx, s.ID, fileID, &start, &end))

		ranges, err = qb.GetFileRanges(ctx, s.ID)
		require.NoError(t, err)
		require.Len(t, ranges, 1)
		assert.Equal(t, fileID, ranges[0].FileID)
		require.NotNil(t, ranges[0].StartTime)
		assert.Equal(t, start, *ranges[0].StartTime)
		require.NotNil(t, ranges[0].EndTime)
		assert.Equal(t, end, *ranges[0].EndTime)

		// clear the range
		require.NoError(t, qb.SetFileRange(ctx, s.ID, fileID, nil, nil))
		ranges, err = qb.GetFileRanges(ctx, s.ID)
		require.NoError(t, err)
		assert.Empty(t, ranges)
	})

	runWithRollbackTxn(t, "two scenes share one file via AddFileWithRange", func(t *testing.T, ctx context.Context) {
		fileID := createRangeTestFile(ctx, t)

		s1 := &models.Scene{}
		require.NoError(t, qb.Create(ctx, s1, []models.FileID{fileID}))
		s2 := &models.Scene{}
		require.NoError(t, qb.Create(ctx, s2, nil))

		var (
			start = 30.0
			end   = 40.0
		)

		// s2 joins the same file with a range; must not steal from s1
		require.NoError(t, qb.AddFileWithRange(ctx, s2.ID, fileID, &start, &end))

		// both scenes find the file
		s1Files, err := qb.GetManyFileIDs(ctx, []int{s1.ID})
		require.NoError(t, err)
		assert.Contains(t, s1Files[0], fileID)

		s2Files, err := qb.GetManyFileIDs(ctx, []int{s2.ID})
		require.NoError(t, err)
		assert.Contains(t, s2Files[0], fileID)

		// s2's join is ranged, s1's is not
		s1Ranges, err := qb.GetFileRanges(ctx, s1.ID)
		require.NoError(t, err)
		assert.Empty(t, s1Ranges)

		s2Ranges, err := qb.GetFileRanges(ctx, s2.ID)
		require.NoError(t, err)
		require.Len(t, s2Ranges, 1)
		assert.Equal(t, fileID, s2Ranges[0].FileID)

		// unranged sharing check
		shared, err := qb.FileSharedWithUnrangedScene(ctx, fileID, s2.ID)
		require.NoError(t, err)
		assert.True(t, shared, "s1 holds the file without a range")

		shared, err = qb.FileSharedWithUnrangedScene(ctx, fileID, s1.ID)
		require.NoError(t, err)
		assert.False(t, shared, "s2 holds the file only with a range")

		// FindByFileID returns both scenes
		found, err := qb.FindByFileID(ctx, fileID)
		require.NoError(t, err)
		assert.Len(t, found, 2)
	})

	runWithRollbackTxn(t, "AddFileWithRange marks primary when scene has no files", func(t *testing.T, ctx context.Context) {
		fileID := createRangeTestFile(ctx, t)

		s := &models.Scene{}
		require.NoError(t, qb.Create(ctx, s, nil))

		var (
			start = 1.0
			end   = 2.0
		)

		// first file joined with a range becomes primary
		require.NoError(t, qb.AddFileWithRange(ctx, s.ID, fileID, &start, &end))

		found, err := qb.FindByFileID(ctx, fileID)
		require.NoError(t, err)
		require.Len(t, found, 1)
		assert.Equal(t, s.ID, found[0].ID)
	})

	runWithRollbackTxn(t, "AssignFiles still steals from unranged scenes", func(t *testing.T, ctx context.Context) {
		// regression guard: legacy exclusive assignment behavior unchanged
		fileID := createRangeTestFile(ctx, t)

		s1 := &models.Scene{}
		require.NoError(t, qb.Create(ctx, s1, []models.FileID{fileID}))
		s2 := &models.Scene{}
		require.NoError(t, qb.Create(ctx, s2, nil))

		require.NoError(t, qb.AssignFiles(ctx, s2.ID, []models.FileID{fileID}))

		// s1 no longer holds the file
		s1Files, err := qb.GetManyFileIDs(ctx, []int{s1.ID})
		require.NoError(t, err)
		assert.NotContains(t, s1Files[0], fileID)

		s2Files, err := qb.GetManyFileIDs(ctx, []int{s2.ID})
		require.NoError(t, err)
		assert.Contains(t, s2Files[0], fileID)
	})
}
