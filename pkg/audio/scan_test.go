// TODO(audio): update this file

package audio

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
		testAudioID = 1
		testFileID  = 100
	)

	existingFile := &models.VideoFile{
		BaseFile: &models.BaseFile{ID: models.FileID(testFileID), Path: "test.mp4"},
	}

	makeAudio := func() *models.Audio {
		return &models.Audio{
			ID:    testAudioID,
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
			db.Audio.On("GetFiles", mock.Anything, testAudioID).Return([]*models.VideoFile{existingFile}, nil)

			if tt.expectUpdate {
				db.Audio.On("UpdatePartial", mock.Anything, testAudioID, mock.Anything).
					Return(&models.Audio{ID: testAudioID}, nil)
			}

			h := &ScanHandler{
				CreatorUpdater: db.Audio,
				PluginCache:    &plugin.Cache{},
			}

			db.WithTxnCtx(func(ctx context.Context) {
				err := h.associateExisting(ctx, []*models.Audio{makeAudio()}, existingFile, tt.updateExisting)
				assert.NoError(t, err)
			})

			if tt.expectUpdate {
				db.Audio.AssertCalled(t, "UpdatePartial", mock.Anything, testAudioID, mock.Anything)
			} else {
				db.Audio.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestAssociateExisting_UpdatePartialOnNewFile(t *testing.T) {
	const (
		testAudioID = 1
		existFileID = 100
		newFileID   = 200
	)

	existingFile := &models.VideoFile{
		BaseFile: &models.BaseFile{ID: models.FileID(existFileID), Path: "existing.mp4"},
	}
	newFile := &models.VideoFile{
		BaseFile: &models.BaseFile{ID: models.FileID(newFileID), Path: "new.mp4"},
	}

	audio := &models.Audio{
		ID:    testAudioID,
		Files: models.NewRelatedVideoFiles([]*models.VideoFile{existingFile}),
	}

	db := mocks.NewDatabase()
	db.Audio.On("GetFiles", mock.Anything, testAudioID).Return([]*models.VideoFile{existingFile}, nil)
	db.Audio.On("AddFileID", mock.Anything, testAudioID, models.FileID(newFileID)).Return(nil)
	db.Audio.On("UpdatePartial", mock.Anything, testAudioID, mock.Anything).
		Return(&models.Audio{ID: testAudioID}, nil)

	h := &ScanHandler{
		CreatorUpdater: db.Audio,
		PluginCache:    &plugin.Cache{},
	}

	db.WithTxnCtx(func(ctx context.Context) {
		err := h.associateExisting(ctx, []*models.Audio{audio}, newFile, false)
		assert.NoError(t, err)
	})

	db.Audio.AssertCalled(t, "AddFileID", mock.Anything, testAudioID, models.FileID(newFileID))
	db.Audio.AssertCalled(t, "UpdatePartial", mock.Anything, testAudioID, mock.Anything)
}
