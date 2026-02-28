package image

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockScanConfig struct{}

func (m *mockScanConfig) GetCreateGalleriesFromFolders() bool { return false }

func TestAssociateExisting_UpdatePartialOnContentChange(t *testing.T) {
	const (
		testImageID = 1
		testFileID  = 100
	)

	existingFile := &models.BaseFile{ID: models.FileID(testFileID), Path: "/images/test.jpg"}

	makeImage := func() *models.Image {
		return &models.Image{
			ID:         testImageID,
			Files:      models.NewRelatedFiles([]models.File{existingFile}),
			GalleryIDs: models.NewRelatedIDs([]int{}),
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
			db.Image.On("GetFiles", mock.Anything, testImageID).Return([]models.File{existingFile}, nil)
			db.Image.On("GetGalleryIDs", mock.Anything, testImageID).Return([]int{}, nil)

			if tt.expectUpdate {
				db.Image.On("UpdatePartial", mock.Anything, testImageID, mock.Anything).
					Return(&models.Image{ID: testImageID}, nil)
			}

			h := &ScanHandler{
				CreatorUpdater: db.Image,
				GalleryFinder:  db.Gallery,
				ScanConfig:     &mockScanConfig{},
				PluginCache:    &plugin.Cache{},
			}

			db.WithTxnCtx(func(ctx context.Context) {
				err := h.associateExisting(ctx, []*models.Image{makeImage()}, existingFile, tt.updateExisting)
				assert.NoError(t, err)
			})

			if tt.expectUpdate {
				db.Image.AssertCalled(t, "UpdatePartial", mock.Anything, testImageID, mock.Anything)
			} else {
				db.Image.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestAssociateExisting_UpdatePartialOnNewFile(t *testing.T) {
	const (
		testImageID = 1
		existFileID = 100
		newFileID   = 200
	)

	existingFile := &models.BaseFile{ID: models.FileID(existFileID), Path: "/images/existing.jpg"}
	newFile := &models.BaseFile{ID: models.FileID(newFileID), Path: "/images/new.jpg"}

	image := &models.Image{
		ID:         testImageID,
		Files:      models.NewRelatedFiles([]models.File{existingFile}),
		GalleryIDs: models.NewRelatedIDs([]int{}),
	}

	db := mocks.NewDatabase()
	db.Image.On("GetFiles", mock.Anything, testImageID).Return([]models.File{existingFile}, nil)
	db.Image.On("GetGalleryIDs", mock.Anything, testImageID).Return([]int{}, nil)
	db.Image.On("AddFileID", mock.Anything, testImageID, models.FileID(newFileID)).Return(nil)
	db.Image.On("UpdatePartial", mock.Anything, testImageID, mock.Anything).
		Return(&models.Image{ID: testImageID}, nil)

	h := &ScanHandler{
		CreatorUpdater: db.Image,
		GalleryFinder:  db.Gallery,
		ScanConfig:     &mockScanConfig{},
		PluginCache:    &plugin.Cache{},
	}

	db.WithTxnCtx(func(ctx context.Context) {
		err := h.associateExisting(ctx, []*models.Image{image}, newFile, false)
		assert.NoError(t, err)
	})

	db.Image.AssertCalled(t, "AddFileID", mock.Anything, testImageID, models.FileID(newFileID))
	db.Image.AssertCalled(t, "UpdatePartial", mock.Anything, testImageID, mock.Anything)
}
