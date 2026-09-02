package scene

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestInvalidateGeneratedFiles(t *testing.T) {
	const hash = "abc123"

	tmpDir := t.TempDir()
	p := paths.NewPaths(tmpDir, filepath.Join(tmpDir, "blobs"))

	seedFile := func(path string) {
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
		require.NoError(t, os.WriteFile(path, []byte("stale"), 0644))
	}

	videoPreview := p.Scene.GetVideoPreviewPath(hash)
	webpPreview := p.Scene.GetWebpPreviewPath(hash)
	transcode := p.Scene.GetTranscodePath(hash)
	spriteVtt := p.Scene.GetSpriteVttFilePath(hash)
	spriteImage := p.Scene.GetSpriteImageFilePath(hash)
	heatmap := p.Scene.GetInteractiveHeatmapPath(hash)
	markersFolder := filepath.Join(p.Generated.Markers, hash)

	seedFile(videoPreview)
	seedFile(webpPreview)
	seedFile(transcode)
	seedFile(spriteVtt)
	seedFile(spriteImage)
	seedFile(heatmap)
	seedFile(filepath.Join(markersFolder, "1.mp4"))

	// no files exist for this hash - should be a no-op
	assert.NotPanics(t, func() {
		InvalidateGeneratedFiles(&p, "no-such-hash")
	})

	InvalidateGeneratedFiles(&p, hash)

	assert.NoFileExists(t, videoPreview)
	assert.NoFileExists(t, webpPreview)
	assert.NoFileExists(t, transcode)
	assert.NoFileExists(t, spriteVtt)
	assert.NoFileExists(t, spriteImage)
	assert.NoFileExists(t, heatmap)
	assert.NoDirExists(t, markersFolder)
}

func TestHandle_ContentChangedAtSamePath(t *testing.T) {
	const testSceneID = 1
	const testFileID = 100

	const oldHash = "oldhash0000000000000000000000000"
	const newHash = "newhash0000000000000000000000000"

	tests := []struct {
		name               string
		newFingerprint     string
		expectInvalidation bool
	}{
		{
			name:               "clears stale generated files and cover when content changed at the same path",
			newFingerprint:     newHash,
			expectInvalidation: true,
		},
		{
			name:               "leaves generated files and cover alone when the hash is unchanged",
			newFingerprint:     oldHash,
			expectInvalidation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFile := &models.VideoFile{
				BaseFile: &models.BaseFile{
					ID:   models.FileID(testFileID),
					Path: "test.mp4",
					Fingerprints: models.Fingerprints{
						{Type: models.FingerprintTypeOshash, Fingerprint: oldHash},
					},
				},
			}
			newFile := &models.VideoFile{
				BaseFile: &models.BaseFile{
					ID:   models.FileID(testFileID),
					Path: "test.mp4",
					Fingerprints: models.Fingerprints{
						{Type: models.FingerprintTypeOshash, Fingerprint: tt.newFingerprint},
					},
				},
			}

			scene := &models.Scene{
				ID:    testSceneID,
				Files: models.NewRelatedVideoFiles([]*models.VideoFile{oldFile}),
			}

			tmpDir := t.TempDir()
			p := paths.NewPaths(tmpDir, filepath.Join(tmpDir, "blobs"))

			staleSprite := p.Scene.GetSpriteImageFilePath(oldHash)
			require.NoError(t, os.MkdirAll(filepath.Dir(staleSprite), 0755))
			require.NoError(t, os.WriteFile(staleSprite, []byte("stale"), 0644))

			db := mocks.NewDatabase()
			db.File.On("GetCaptions", mock.Anything, mock.Anything).Return(nil, nil)
			db.Scene.On("FindByFileID", mock.Anything, models.FileID(testFileID)).Return([]*models.Scene{scene}, nil)
			db.Scene.On("GetFiles", mock.Anything, testSceneID).Return([]*models.VideoFile{oldFile}, nil)
			db.Scene.On("UpdatePartial", mock.Anything, testSceneID, mock.Anything).
				Return(&models.Scene{ID: testSceneID}, nil)
			db.Scene.On("UpdateCover", mock.Anything, testSceneID, mock.Anything).Return(nil)
			db.Gallery.On("FindByPath", mock.Anything, mock.Anything).Return(nil, nil)

			h := &ScanHandler{
				CreatorUpdater:       db.Scene,
				GalleryFinderUpdater: db.Gallery,
				CaptionUpdater:       db.File,
				ScanGenerator:        noopScanGenerator{},
				PluginCache:          &plugin.Cache{},
				FileNamingAlgorithm:  models.HashAlgorithmOshash,
				Paths:                &p,
			}

			db.WithTxnCtx(func(ctx context.Context) {
				err := h.Handle(ctx, newFile, oldFile)
				assert.NoError(t, err)
			})

			if tt.expectInvalidation {
				assert.NoFileExists(t, staleSprite)
				db.Scene.AssertCalled(t, "UpdateCover", mock.Anything, testSceneID, mock.Anything)
			} else {
				assert.FileExists(t, staleSprite)
				db.Scene.AssertNotCalled(t, "UpdateCover", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

type noopScanGenerator struct{}

func (noopScanGenerator) Generate(ctx context.Context, s *models.Scene, f *models.VideoFile) error {
	return nil
}
