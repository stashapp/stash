//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var missingSinceNotNull = &models.TimestampCriterionInput{Modifier: models.CriterionModifierNotNull}

func TestFileStore_MissingSince(t *testing.T) {
	qb := db.File

	runWithRollbackTxn(t, "file missing since", func(t *testing.T, ctx context.Context) {
		now := time.Now()
		create := func(folderIdx int, basename string) *models.BaseFile {
			f := &models.BaseFile{
				Path:           getFilePath(folderIdx, basename),
				ParentFolderID: folderIDs[folderIdx],
				Basename:       basename,
				DirEntry:       models.DirEntry{ModTime: now},
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			require.NoError(t, qb.Create(ctx, f))
			return f
		}

		inFolder := create(folderIdxWithFiles, "missing_a.mp4")
		elsewhere := create(folderIdxWithSceneFiles, "missing_b.mp4")
		present := create(folderIdxWithFiles, "present.mp4")

		require.NoError(t, qb.SetMissing(ctx, inFolder.ID, &now))
		require.NoError(t, qb.SetMissing(ctx, elsewhere.ID, &now))
		assert.Error(t, qb.SetMissing(ctx, models.FileID(-1), &now), "unknown id")

		// path scoping
		scoped := []string{folderPaths[folderIdxWithFiles]}
		got, err := qb.FindMissingInPaths(ctx, scoped, nil, -1, 0)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, inFolder.ID, got[0].Base().ID)
		assert.NotNil(t, got[0].Base().MissingSince)

		n, err := qb.CountMissingInPaths(ctx, scoped, nil)
		require.NoError(t, err)
		assert.Equal(t, 1, n)

		// ordered by id, so paging is stable
		got, err = qb.FindMissingInPaths(ctx, nil, nil, 1, 1)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, elsewhere.ID, got[0].Base().ID)

		// findFiles filter
		res, err := qb.Query(ctx, models.FileQueryOptions{FileFilter: &models.FileFilterType{MissingSince: missingSinceNotNull}})
		require.NoError(t, err)
		assert.ElementsMatch(t, []models.FileID{inFolder.ID, elsewhere.ID}, res.IDs)

		inFolderFiles, err := qb.FindByFolderID(ctx, folderIDs[folderIdxWithFiles])
		require.NoError(t, err)
		var ids []models.FileID
		for _, f := range inFolderFiles {
			ids = append(ids, f.Base().ID)
		}
		assert.Contains(t, ids, inFolder.ID)
		assert.Contains(t, ids, present.ID)
		assert.NotContains(t, ids, elsewhere.ID)

		// clearing
		require.NoError(t, qb.SetMissing(ctx, inFolder.ID, nil))
		n, err = qb.CountMissingInPaths(ctx, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, 1, n)
	})
}

func TestFolderStore_MissingSince(t *testing.T) {
	qb := db.Folder

	runWithRollbackTxn(t, "folder missing since", func(t *testing.T, ctx context.Context) {
		now := time.Now()
		root := filepath.Join(folderPaths[folderIdxRoot], "missing_since_root")
		create := func(path string) *models.Folder {
			f := &models.Folder{
				Path:      path,
				DirEntry:  models.DirEntry{ModTime: now},
				CreatedAt: now,
				UpdatedAt: now,
			}
			require.NoError(t, qb.Create(ctx, f))
			return f
		}

		inRoot := create(filepath.Join(root, "a"))
		elsewhere := create(filepath.Join(folderPaths[folderIdxRoot], "missing_since_other"))

		require.NoError(t, qb.SetMissing(ctx, inRoot.ID, &now))
		require.NoError(t, qb.SetMissing(ctx, elsewhere.ID, &now))
		assert.Error(t, qb.SetMissing(ctx, models.FolderID(-1), &now), "unknown id")

		got, err := qb.FindMissingInPaths(ctx, []string{root}, nil, -1, 0)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, inRoot.ID, got[0].ID)
		assert.NotNil(t, got[0].MissingSince)

		n, err := qb.CountMissingInPaths(ctx, []string{root}, nil)
		require.NoError(t, err)
		assert.Equal(t, 1, n)

		// missingSinceBefore excludes folders marked at or after the cutoff
		before := now.Add(-time.Hour)
		got, err = qb.FindMissingInPaths(ctx, nil, &before, -1, 0)
		require.NoError(t, err)
		assert.Empty(t, got)

		res, err := qb.Query(ctx, models.FolderQueryOptions{FolderFilter: &models.FolderFilterType{MissingSince: missingSinceNotNull}})
		require.NoError(t, err)
		assert.ElementsMatch(t, []models.FolderID{inRoot.ID, elsewhere.ID}, res.IDs)

		require.NoError(t, qb.SetMissing(ctx, inRoot.ID, nil))
		found, err := qb.Find(ctx, inRoot.ID)
		require.NoError(t, err)
		assert.Nil(t, found.MissingSince)
	})
}
