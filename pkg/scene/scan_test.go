package scene

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
		testSceneID = 1
		testFileID  = 100
	)

	existingFile := &models.VideoFile{
		BaseFile: &models.BaseFile{ID: models.FileID(testFileID), Path: "test.mp4"},
	}

	makeScene := func() *models.Scene {
		return &models.Scene{
			ID:    testSceneID,
			Files: models.NewRelatedVideoFiles([]*models.VideoFile{existingFile}),
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
			db.Scene.On("GetFiles", mock.Anything, testSceneID).Return([]*models.VideoFile{existingFile}, nil)

			if tt.expectUpdate {
				db.Scene.On("UpdatePartial", mock.Anything, testSceneID, mock.Anything).
					Return(&models.Scene{ID: testSceneID}, nil)
			}

			h := &ScanHandler{
				CreatorUpdater: db.Scene,
				PluginCache:    &plugin.Cache{},
			}

			db.WithTxnCtx(func(ctx context.Context) {
				err := h.associateExisting(ctx, []*models.Scene{makeScene()}, existingFile, tt.updateExisting)
				assert.NoError(t, err)
			})

			if tt.expectUpdate {
				db.Scene.AssertCalled(t, "UpdatePartial", mock.Anything, testSceneID, mock.Anything)
			} else {
				db.Scene.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestAssociateExisting_UpdatePartialOnNewFile(t *testing.T) {
	const (
		testSceneID = 1
		existFileID = 100
		newFileID   = 200
	)

	existingFile := &models.VideoFile{
		BaseFile: &models.BaseFile{ID: models.FileID(existFileID), Path: "existing.mp4"},
	}
	newFile := &models.VideoFile{
		BaseFile: &models.BaseFile{ID: models.FileID(newFileID), Path: "new.mp4"},
	}

	scene := &models.Scene{
		ID:    testSceneID,
		Files: models.NewRelatedVideoFiles([]*models.VideoFile{existingFile}),
	}

	db := mocks.NewDatabase()
	db.Scene.On("GetFiles", mock.Anything, testSceneID).Return([]*models.VideoFile{existingFile}, nil)
	db.Scene.On("AddFileID", mock.Anything, testSceneID, models.FileID(newFileID)).Return(nil)
	db.Scene.On("UpdatePartial", mock.Anything, testSceneID, mock.Anything).
		Return(&models.Scene{ID: testSceneID}, nil)

	h := &ScanHandler{
		CreatorUpdater: db.Scene,
		PluginCache:    &plugin.Cache{},
	}

	db.WithTxnCtx(func(ctx context.Context) {
		err := h.associateExisting(ctx, []*models.Scene{scene}, newFile, false)
		assert.NoError(t, err)
	})

	db.Scene.AssertCalled(t, "AddFileID", mock.Anything, testSceneID, models.FileID(newFileID))
	db.Scene.AssertCalled(t, "UpdatePartial", mock.Anything, testSceneID, mock.Anything)
}
