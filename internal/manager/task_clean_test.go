package manager

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFindGalleriesToCleanIncludesNoGalleryFolders(t *testing.T) {
	cleanPath := t.TempDir()
	outsideCleanPath := t.TempDir()

	createGalleryFolder := func(root, name string, markers ...string) string {
		t.Helper()

		path := filepath.Join(root, name)
		require.NoError(t, os.MkdirAll(path, 0755))
		for _, marker := range markers {
			require.NoError(t, os.WriteFile(filepath.Join(path, marker), nil, 0644))
		}

		return path
	}

	// Empty path-backed galleries retain the existing clean behaviour.
	emptyGallery := &models.Gallery{ID: 1, Path: createGalleryFolder(cleanPath, "empty")}

	// A gallery matching both clean rules must only be returned once.
	emptyNoGalleryFolderID := models.FolderID(2)
	emptyNoGallery := &models.Gallery{
		ID:       2,
		Path:     createGalleryFolder(cleanPath, "empty-nogallery", ".nogallery"),
		FolderID: &emptyNoGalleryFolderID,
	}

	// A populated folder gallery is cleaned when its folder opts out.
	noGalleryFolderID := models.FolderID(3)
	noGallery := &models.Gallery{
		ID:       3,
		Path:     createGalleryFolder(cleanPath, "nogallery", ".nogallery"),
		FolderID: &noGalleryFolderID,
	}

	// An ordinary folder gallery has no reason to be cleaned.
	regularFolderID := models.FolderID(4)
	regular := &models.Gallery{
		ID:       4,
		Path:     createGalleryFolder(cleanPath, "regular"),
		FolderID: &regularFolderID,
	}

	// Selective clean must not touch galleries outside the requested paths.
	outsideFolderID := models.FolderID(5)
	outside := &models.Gallery{
		ID:       5,
		Path:     createGalleryFolder(outsideCleanPath, "nogallery", ".nogallery"),
		FolderID: &outsideFolderID,
	}

	// Match Scan semantics: .forcegallery wins when both markers exist.
	forcedFolderID := models.FolderID(7)
	forced := &models.Gallery{
		ID:       7,
		Path:     createGalleryFolder(cleanPath, "forced", ".nogallery", ".forcegallery"),
		FolderID: &forcedFolderID,
	}

	db := mocks.NewDatabase()
	emptyGalleryFilter := mock.MatchedBy(func(filter *models.GalleryFilterType) bool {
		return filter != nil && filter.ImageCount != nil &&
			filter.ImageCount.Value == 0 &&
			filter.ImageCount.Modifier == models.CriterionModifierEquals
	})
	// Only folder galleries need marker checks; filter out user and file galleries in SQL.
	folderGalleryFilter := mock.MatchedBy(func(filter *models.GalleryFilterType) bool {
		return filter != nil && filter.FoldersFilter != nil &&
			filter.FoldersFilter.Path != nil &&
			filter.FoldersFilter.Path.Modifier == models.CriterionModifierNotNull
	})

	db.Gallery.On("Query", mock.Anything, emptyGalleryFilter, mock.Anything).
		Return([]*models.Gallery{emptyGallery, emptyNoGallery}, 2, nil).Once()
	db.Gallery.On("Query", mock.Anything, emptyGalleryFilter, mock.Anything).
		Return([]*models.Gallery{}, 0, nil).Once()
	db.Gallery.On("Query", mock.Anything, folderGalleryFilter, mock.Anything).
		Return([]*models.Gallery{emptyNoGallery, noGallery, regular, outside, forced}, 5, nil).Once()
	db.Gallery.On("Query", mock.Anything, folderGalleryFilter, mock.Anything).
		Return([]*models.Gallery{}, 0, nil).Once()

	j := cleanJob{
		repository: db.Repository(),
		input: CleanMetadataInput{
			Paths: []string{cleanPath},
		},
	}

	got, err := j.findGalleriesToClean(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []int{emptyGallery.ID, emptyNoGallery.ID, noGallery.ID}, got)
	db.Gallery.AssertExpectations(t)
}
