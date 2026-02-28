package gallery

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

var errRollback = errors.New("rollback")

type mockGalleryScanCU struct {
	mock.Mock
}

func (m *mockGalleryScanCU) FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Gallery, error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).([]*models.Gallery), args.Error(1)
}

func (m *mockGalleryScanCU) FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Gallery, error) {
	args := m.Called(ctx, fp)
	return args.Get(0).([]*models.Gallery), args.Error(1)
}

func (m *mockGalleryScanCU) GetFiles(ctx context.Context, relatedID int) ([]models.File, error) {
	args := m.Called(ctx, relatedID)
	return args.Get(0).([]models.File), args.Error(1)
}

func (m *mockGalleryScanCU) Create(ctx context.Context, input *models.CreateGalleryInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *mockGalleryScanCU) UpdatePartial(ctx context.Context, id int, updatedGallery models.GalleryPartial) (*models.Gallery, error) {
	args := m.Called(ctx, id, updatedGallery)
	return args.Get(0).(*models.Gallery), args.Error(1)
}

func (m *mockGalleryScanCU) AddFileID(ctx context.Context, id int, fileID models.FileID) error {
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

func withTxnCtx(fn func(ctx context.Context)) {
	_ = txn.WithTxn(context.Background(), &noopTxnManager{}, func(ctx context.Context) error {
		fn(ctx)
		return errRollback
	})
}

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
			cu := &mockGalleryScanCU{}
			cu.On("GetFiles", mock.Anything, testGalleryID).Return([]models.File{existingFile}, nil)

			if tt.expectUpdate {
				cu.On("UpdatePartial", mock.Anything, testGalleryID, mock.Anything).
					Return(&models.Gallery{ID: testGalleryID}, nil)
			}

			h := &ScanHandler{
				CreatorUpdater: cu,
				PluginCache:    &plugin.Cache{},
			}

			withTxnCtx(func(ctx context.Context) {
				err := h.associateExisting(ctx, []*models.Gallery{makeGallery()}, existingFile, tt.updateExisting)
				assert.NoError(t, err)
			})

			if tt.expectUpdate {
				cu.AssertCalled(t, "UpdatePartial", mock.Anything, testGalleryID, mock.Anything)
			} else {
				cu.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
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

	cu := &mockGalleryScanCU{}
	cu.On("GetFiles", mock.Anything, testGalleryID).Return([]models.File{existingFile}, nil)
	cu.On("AddFileID", mock.Anything, testGalleryID, models.FileID(newFileID)).Return(nil)
	cu.On("UpdatePartial", mock.Anything, testGalleryID, mock.Anything).
		Return(&models.Gallery{ID: testGalleryID}, nil)

	h := &ScanHandler{
		CreatorUpdater: cu,
		PluginCache:    &plugin.Cache{},
	}

	withTxnCtx(func(ctx context.Context) {
		err := h.associateExisting(ctx, []*models.Gallery{gallery}, newFile, false)
		assert.NoError(t, err)
	})

	cu.AssertCalled(t, "AddFileID", mock.Anything, testGalleryID, models.FileID(newFileID))
	cu.AssertCalled(t, "UpdatePartial", mock.Anything, testGalleryID, mock.Anything)
}
