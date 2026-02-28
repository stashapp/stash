package scene

import (
	"context"
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// sentinel error to force rollback, preventing post-commit hooks from
// trying to execute plugins on the zero-value Cache in tests.
var errRollback = errors.New("rollback")

type mockSceneScanCU struct {
	mock.Mock
}

func (m *mockSceneScanCU) FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Scene, error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).([]*models.Scene), args.Error(1)
}

func (m *mockSceneScanCU) FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Scene, error) {
	args := m.Called(ctx, fp)
	return args.Get(0).([]*models.Scene), args.Error(1)
}

func (m *mockSceneScanCU) GetFiles(ctx context.Context, relatedID int) ([]*models.VideoFile, error) {
	args := m.Called(ctx, relatedID)
	return args.Get(0).([]*models.VideoFile), args.Error(1)
}

func (m *mockSceneScanCU) Create(ctx context.Context, newScene *models.Scene, fileIDs []models.FileID) error {
	args := m.Called(ctx, newScene, fileIDs)
	return args.Error(0)
}

func (m *mockSceneScanCU) UpdatePartial(ctx context.Context, id int, updatedScene models.ScenePartial) (*models.Scene, error) {
	args := m.Called(ctx, id, updatedScene)
	return args.Get(0).(*models.Scene), args.Error(1)
}

func (m *mockSceneScanCU) AddFileID(ctx context.Context, id int, fileID models.FileID) error {
	args := m.Called(ctx, id, fileID)
	return args.Error(0)
}

type noopTxnManager struct{}

func (m *noopTxnManager) Begin(ctx context.Context, writable bool) (context.Context, error) {
	return ctx, nil
}
func (m *noopTxnManager) Commit(ctx context.Context) error   { return nil }
func (m *noopTxnManager) Rollback(ctx context.Context) error { return nil }
func (m *noopTxnManager) IsLocked(err error) bool            { return false }

// withTxnCtx runs fn inside a transaction context so that txn.AddPostCommitHook
// (called by plugin.Cache.RegisterPostHooks) doesn't panic on a nil hookManager.
func withTxnCtx(fn func(ctx context.Context)) {
	_ = txn.WithTxn(context.Background(), &noopTxnManager{}, func(ctx context.Context) error {
		fn(ctx)
		return errRollback
	})
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
			cu := &mockSceneScanCU{}
			cu.On("GetFiles", mock.Anything, testSceneID).Return([]*models.VideoFile{existingFile}, nil)

			if tt.expectUpdate {
				cu.On("UpdatePartial", mock.Anything, testSceneID, mock.Anything).
					Return(&models.Scene{ID: testSceneID}, nil)
			}

			h := &ScanHandler{
				CreatorUpdater: cu,
				PluginCache:    &plugin.Cache{},
			}

			withTxnCtx(func(ctx context.Context) {
				err := h.associateExisting(ctx, []*models.Scene{makeScene()}, existingFile, tt.updateExisting)
				assert.NoError(t, err)
			})

			if tt.expectUpdate {
				cu.AssertCalled(t, "UpdatePartial", mock.Anything, testSceneID, mock.Anything)
			} else {
				cu.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
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

	cu := &mockSceneScanCU{}
	cu.On("GetFiles", mock.Anything, testSceneID).Return([]*models.VideoFile{existingFile}, nil)
	cu.On("AddFileID", mock.Anything, testSceneID, models.FileID(newFileID)).Return(nil)
	cu.On("UpdatePartial", mock.Anything, testSceneID, mock.Anything).
		Return(&models.Scene{ID: testSceneID}, nil)

	h := &ScanHandler{
		CreatorUpdater: cu,
		PluginCache:    &plugin.Cache{},
	}

	withTxnCtx(func(ctx context.Context) {
		err := h.associateExisting(ctx, []*models.Scene{scene}, newFile, false)
		assert.NoError(t, err)
	})

	cu.AssertCalled(t, "AddFileID", mock.Anything, testSceneID, models.FileID(newFileID))
	cu.AssertCalled(t, "UpdatePartial", mock.Anything, testSceneID, mock.Anything)
}
