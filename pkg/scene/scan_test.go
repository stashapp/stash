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
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// withCommittedTxn runs fn in a real (mocked) transaction that commits on
// success, unlike mocks.Database.WithTxnCtx which always rolls back - some
// side effects (e.g. InvalidateGeneratedFiles) are deferred to a
// post-commit hook and only run down this path.
func withCommittedTxn(t *testing.T, db *mocks.Database, fn func(ctx context.Context) error) {
	t.Helper()
	err := txn.WithTxn(context.Background(), db, fn)
	assert.NoError(t, err)
}

// stubServerConfig is a minimal plugin.ServerConfig so plugin.NewCache
// produces a Cache safe to actually execute post-commit hooks against -
// the zero-value plugin.Cache{} has a nil config, which panics as soon as
// a real commit (see withCommittedTxn) runs its post-hook execution path.
type stubServerConfig struct{}

func (stubServerConfig) GetHost() string              { return "" }
func (stubServerConfig) GetPort() int                 { return 0 }
func (stubServerConfig) GetConfigPathAbs() string     { return "" }
func (stubServerConfig) HasTLSConfig() bool           { return false }
func (stubServerConfig) GetPluginsPath() string       { return "" }
func (stubServerConfig) GetDisabledPlugins() []string { return nil }
func (stubServerConfig) GetPythonPath() string        { return "" }

// oldHash/newHash are the before/after hashes shared by every
// TestHandle_ContentChangedAtSamePath* same-path-content-change test below.
const (
	oldHash = "oldhash0000000000000000000000000"
	newHash = "newhash0000000000000000000000000"
)

// newTestScanHandler builds the ScanHandler struct literal shared by every
// TestHandle_ContentChangedAtSamePath* test below.
func newTestScanHandler(db *mocks.Database, p *paths.Paths) *ScanHandler {
	return &ScanHandler{
		CreatorUpdater:       db.Scene,
		GalleryFinderUpdater: db.Gallery,
		CaptionUpdater:       db.File,
		ScanGenerator:        noopScanGenerator{},
		PluginCache:          plugin.NewCache(stubServerConfig{}),
		FileNamingAlgorithm:  models.HashAlgorithmOshash,
		Paths:                p,
	}
}

// seedStaleFile writes a placeholder file at path, creating its parent
// directory as needed, standing in for generated content a test expects
// to be deleted (or to survive, depending on the assertion).
func seedStaleFile(t *testing.T, path string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	require.NoError(t, os.WriteFile(path, []byte("stale"), 0644))
}

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

	videoPreview := p.Scene.GetVideoPreviewPath(hash)
	webpPreview := p.Scene.GetWebpPreviewPath(hash)
	transcode := p.Scene.GetTranscodePath(hash)
	spriteVtt := p.Scene.GetSpriteVttFilePath(hash)
	spriteImage := p.Scene.GetSpriteImageFilePath(hash)
	heatmap := p.Scene.GetInteractiveHeatmapPath(hash)
	markersFolder := filepath.Join(p.Generated.Markers, hash)

	seedStaleFile(t, videoPreview)
	seedStaleFile(t, webpPreview)
	seedStaleFile(t, transcode)
	seedStaleFile(t, spriteVtt)
	seedStaleFile(t, spriteImage)
	seedStaleFile(t, heatmap)
	seedStaleFile(t, filepath.Join(markersFolder, "1.mp4"))

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

			primaryFileID := oldFile.ID
			scene := &models.Scene{
				ID:            testSceneID,
				PrimaryFileID: &primaryFileID,
				Files:         models.NewRelatedVideoFiles([]*models.VideoFile{oldFile}),
			}

			tmpDir := t.TempDir()
			p := paths.NewPaths(tmpDir, filepath.Join(tmpDir, "blobs"))

			staleSprite := p.Scene.GetSpriteImageFilePath(oldHash)
			seedStaleFile(t, staleSprite)

			db := mocks.NewDatabase()
			db.File.On("GetCaptions", mock.Anything, mock.Anything).Return(nil, nil)
			db.Scene.On("FindByFileID", mock.Anything, models.FileID(testFileID)).Return([]*models.Scene{scene}, nil)
			db.Scene.On("GetFiles", mock.Anything, testSceneID).Return([]*models.VideoFile{oldFile}, nil)
			db.Scene.On("UpdatePartial", mock.Anything, testSceneID, mock.Anything).
				Return(&models.Scene{ID: testSceneID}, nil)
			db.Scene.On("UpdateCover", mock.Anything, testSceneID, mock.Anything).Return(nil)
			db.Gallery.On("FindByPath", mock.Anything, mock.Anything).Return(nil, nil)

			h := newTestScanHandler(db, &p)

			withCommittedTxn(t, db, func(ctx context.Context) error {
				return h.Handle(ctx, newFile, oldFile)
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

func TestHandle_ContentChangedAtSamePath_SecondaryFile(t *testing.T) {
	const testSceneID = 1
	const primaryFileID = 100
	const secondaryFileID = 200

	primaryFile := &models.VideoFile{
		BaseFile: &models.BaseFile{
			ID:   models.FileID(primaryFileID),
			Path: "primary.mp4",
		},
	}
	oldFile := &models.VideoFile{
		BaseFile: &models.BaseFile{
			ID:   models.FileID(secondaryFileID),
			Path: "secondary.mp4",
			Fingerprints: models.Fingerprints{
				{Type: models.FingerprintTypeOshash, Fingerprint: oldHash},
			},
		},
	}
	newFile := &models.VideoFile{
		BaseFile: &models.BaseFile{
			ID:   models.FileID(secondaryFileID),
			Path: "secondary.mp4",
			Fingerprints: models.Fingerprints{
				{Type: models.FingerprintTypeOshash, Fingerprint: newHash},
			},
		},
	}

	primaryFileIDCopy := primaryFile.ID
	scene := &models.Scene{
		ID:            testSceneID,
		PrimaryFileID: &primaryFileIDCopy,
		Files:         models.NewRelatedVideoFiles([]*models.VideoFile{primaryFile, oldFile}),
	}

	tmpDir := t.TempDir()
	p := paths.NewPaths(tmpDir, filepath.Join(tmpDir, "blobs"))

	staleSprite := p.Scene.GetSpriteImageFilePath(oldHash)
	seedStaleFile(t, staleSprite)

	db := mocks.NewDatabase()
	db.File.On("GetCaptions", mock.Anything, mock.Anything).Return(nil, nil)
	db.Scene.On("FindByFileID", mock.Anything, models.FileID(secondaryFileID)).Return([]*models.Scene{scene}, nil)
	db.Scene.On("GetFiles", mock.Anything, testSceneID).Return([]*models.VideoFile{primaryFile, oldFile}, nil)
	db.Scene.On("UpdatePartial", mock.Anything, testSceneID, mock.Anything).
		Return(&models.Scene{ID: testSceneID}, nil)
	db.Scene.On("UpdateCover", mock.Anything, testSceneID, mock.Anything).Return(nil)
	db.Gallery.On("FindByPath", mock.Anything, mock.Anything).Return(nil, nil)

	h := newTestScanHandler(db, &p)

	withCommittedTxn(t, db, func(ctx context.Context) error {
		return h.Handle(ctx, newFile, oldFile)
	})

	// nothing scene-level (sprite/preview/transcode/cover) was ever
	// generated from the secondary file's hash - only the primary file's
	// hash drives that - so a secondary file's edit must leave everything
	// alone, in case its old hash happens to still be in use by the
	// (unchanged) primary file
	assert.FileExists(t, staleSprite)
	db.Scene.AssertNotCalled(t, "UpdateCover", mock.Anything, mock.Anything, mock.Anything)
}

func TestHandle_ContentChangedAtSamePath_SharedHashSibling(t *testing.T) {
	const testSceneID = 1
	const primaryFileID = 100
	const secondaryFileID = 200

	// primary and secondary start out as byte-identical duplicates,
	// sharing a hash - a real, reachable state via FindByFingerprints
	// matching a new file to an existing scene by content
	const sharedHash = "sharedhash00000000000000000000000"

	primaryFile := &models.VideoFile{
		BaseFile: &models.BaseFile{
			ID:   models.FileID(primaryFileID),
			Path: "primary.mp4",
			Fingerprints: models.Fingerprints{
				{Type: models.FingerprintTypeOshash, Fingerprint: sharedHash},
			},
		},
	}
	oldFile := &models.VideoFile{
		BaseFile: &models.BaseFile{
			ID:   models.FileID(secondaryFileID),
			Path: "secondary.mp4",
			Fingerprints: models.Fingerprints{
				{Type: models.FingerprintTypeOshash, Fingerprint: sharedHash},
			},
		},
	}
	newFile := &models.VideoFile{
		BaseFile: &models.BaseFile{
			ID:   models.FileID(secondaryFileID),
			Path: "secondary.mp4",
			Fingerprints: models.Fingerprints{
				{Type: models.FingerprintTypeOshash, Fingerprint: newHash},
			},
		},
	}

	primaryFileIDCopy := primaryFile.ID
	scene := &models.Scene{
		ID:            testSceneID,
		PrimaryFileID: &primaryFileIDCopy,
		Files:         models.NewRelatedVideoFiles([]*models.VideoFile{primaryFile, oldFile}),
	}

	tmpDir := t.TempDir()
	p := paths.NewPaths(tmpDir, filepath.Join(tmpDir, "blobs"))

	// generated content sitting at the shared hash - this is what the
	// still-unedited primary file depends on
	sharedSprite := p.Scene.GetSpriteImageFilePath(sharedHash)
	seedStaleFile(t, sharedSprite)

	db := mocks.NewDatabase()
	db.File.On("GetCaptions", mock.Anything, mock.Anything).Return(nil, nil)
	db.Scene.On("FindByFileID", mock.Anything, models.FileID(secondaryFileID)).Return([]*models.Scene{scene}, nil)
	db.Scene.On("GetFiles", mock.Anything, testSceneID).Return([]*models.VideoFile{primaryFile, oldFile}, nil)
	db.Scene.On("UpdatePartial", mock.Anything, testSceneID, mock.Anything).
		Return(&models.Scene{ID: testSceneID}, nil)
	db.Scene.On("UpdateCover", mock.Anything, testSceneID, mock.Anything).Return(nil)
	db.Gallery.On("FindByPath", mock.Anything, mock.Anything).Return(nil, nil)

	h := newTestScanHandler(db, &p)

	withCommittedTxn(t, db, func(ctx context.Context) error {
		return h.Handle(ctx, newFile, oldFile)
	})

	// editing the secondary file must not destroy generated content the
	// still-unedited primary file depends on, even though they used to
	// share a hash
	assert.FileExists(t, sharedSprite)
	db.Scene.AssertNotCalled(t, "UpdateCover", mock.Anything, mock.Anything, mock.Anything)
}

// TestHandle_ContentChangedAtSamePath_RollbackDoesNotDeleteFiles verifies
// InvalidateGeneratedFiles is genuinely deferred to a post-commit hook: if
// something later in the same Handle() call fails and the transaction
// rolls back, the generated files it would have deleted must still exist
// on disk afterward, since the DB itself also reverted to believing the
// old hash (and its files) are still current.
func TestHandle_ContentChangedAtSamePath_RollbackDoesNotDeleteFiles(t *testing.T) {
	const testSceneID = 1
	const testFileID = 100

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
				{Type: models.FingerprintTypeOshash, Fingerprint: newHash},
			},
		},
	}

	primaryFileID := oldFile.ID
	scene := &models.Scene{
		ID:            testSceneID,
		PrimaryFileID: &primaryFileID,
		Files:         models.NewRelatedVideoFiles([]*models.VideoFile{oldFile}),
	}

	tmpDir := t.TempDir()
	p := paths.NewPaths(tmpDir, filepath.Join(tmpDir, "blobs"))

	staleSprite := p.Scene.GetSpriteImageFilePath(oldHash)
	seedStaleFile(t, staleSprite)

	db := mocks.NewDatabase()
	db.File.On("GetCaptions", mock.Anything, mock.Anything).Return(nil, nil)
	db.Scene.On("FindByFileID", mock.Anything, models.FileID(testFileID)).Return([]*models.Scene{scene}, nil)
	db.Scene.On("GetFiles", mock.Anything, testSceneID).Return([]*models.VideoFile{oldFile}, nil)
	db.Scene.On("UpdatePartial", mock.Anything, testSceneID, mock.Anything).
		Return(&models.Scene{ID: testSceneID}, nil)
	db.Scene.On("UpdateCover", mock.Anything, testSceneID, mock.Anything).Return(nil)
	// a later step (gallery association) fails, forcing Handle() to
	// return an error and the transaction to roll back
	db.Gallery.On("FindByPath", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	h := newTestScanHandler(db, &p)

	err := txn.WithTxn(context.Background(), db, func(ctx context.Context) error {
		return h.Handle(ctx, newFile, oldFile)
	})
	assert.Error(t, err)

	// the transaction rolled back, so the post-commit deletion hook must
	// never have run - the DB reverted to the old hash, and its
	// generated files must still be there to match
	assert.FileExists(t, staleSprite)
}

type noopScanGenerator struct{}

func (noopScanGenerator) Generate(ctx context.Context, s *models.Scene, f *models.VideoFile) error {
	return nil
}
