package gallery

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAssociateExisting_UpdatePartialOnContentChange(t *testing.T) {
	const (
		testGalleryID = 1
		testFileID    = 100
	)

	existingFile := &models.BaseFile{ID: models.FileID(testFileID), Path: "test.zip"}

	makeGallery := func() *models.Gallery {
		return &models.Gallery{
			ID:    testGalleryID,
			Files: models.NewRelatedFiles([]models.File{existingFile}),
		}
	}

	tests := []struct {
		name           string
		updateExisting bool
		expectUpdate   bool
	}{
		{
			name:           "calls UpdatePartial when file content changed",
			updateExisting: true,
			expectUpdate:   true,
		},
		{
			name:           "skips UpdatePartial when file unchanged and already associated",
			updateExisting: false,
			expectUpdate:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mocks.NewDatabase()
			db.Gallery.On("GetFiles", mock.Anything, testGalleryID).Return([]models.File{existingFile}, nil)

			if tt.expectUpdate {
				db.Gallery.On("UpdatePartial", mock.Anything, testGalleryID, mock.Anything).
					Return(&models.Gallery{ID: testGalleryID}, nil)
			}

			h := &ScanHandler{
				CreatorUpdater: db.Gallery,
				PluginCache:    &plugin.Cache{},
			}

			db.WithTxnCtx(func(ctx context.Context) {
				err := h.associateExisting(ctx, []*models.Gallery{makeGallery()}, existingFile, tt.updateExisting)
				assert.NoError(t, err)
			})

			if tt.expectUpdate {
				db.Gallery.AssertCalled(t, "UpdatePartial", mock.Anything, testGalleryID, mock.Anything)
			} else {
				db.Gallery.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestAssociateExisting_UpdatePartialOnNewFile(t *testing.T) {
	const (
		testGalleryID = 1
		existFileID   = 100
		newFileID     = 200
	)

	existingFile := &models.BaseFile{ID: models.FileID(existFileID), Path: "existing.zip"}
	newFile := &models.BaseFile{ID: models.FileID(newFileID), Path: "new.zip"}

	gallery := &models.Gallery{
		ID:    testGalleryID,
		Files: models.NewRelatedFiles([]models.File{existingFile}),
	}

	db := mocks.NewDatabase()
	db.Gallery.On("GetFiles", mock.Anything, testGalleryID).Return([]models.File{existingFile}, nil)
	db.Gallery.On("AddFileID", mock.Anything, testGalleryID, models.FileID(newFileID)).Return(nil)
	db.Gallery.On("UpdatePartial", mock.Anything, testGalleryID, mock.Anything).
		Return(&models.Gallery{ID: testGalleryID}, nil)

	h := &ScanHandler{
		CreatorUpdater: db.Gallery,
		PluginCache:    &plugin.Cache{},
	}

	db.WithTxnCtx(func(ctx context.Context) {
		err := h.associateExisting(ctx, []*models.Gallery{gallery}, newFile, false)
		assert.NoError(t, err)
	})

	db.Gallery.AssertCalled(t, "AddFileID", mock.Anything, testGalleryID, models.FileID(newFileID))
	db.Gallery.AssertCalled(t, "UpdatePartial", mock.Anything, testGalleryID, mock.Anything)
}
