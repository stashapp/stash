package image

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

type mockImageScanCU struct {
	mock.Mock
}

func (m *mockImageScanCU) FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Image, error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).([]*models.Image), args.Error(1)
}

func (m *mockImageScanCU) FindByFolderID(ctx context.Context, folderID models.FolderID) ([]*models.Image, error) {
	args := m.Called(ctx, folderID)
	return args.Get(0).([]*models.Image), args.Error(1)
}

func (m *mockImageScanCU) FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Image, error) {
	args := m.Called(ctx, fp)
	return args.Get(0).([]*models.Image), args.Error(1)
}

func (m *mockImageScanCU) GetFiles(ctx context.Context, relatedID int) ([]models.File, error) {
	args := m.Called(ctx, relatedID)
	return args.Get(0).([]models.File), args.Error(1)
}

func (m *mockImageScanCU) GetGalleryIDs(ctx context.Context, relatedID int) ([]int, error) {
	args := m.Called(ctx, relatedID)
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockImageScanCU) Create(ctx context.Context, newImage *models.CreateImageInput) error {
	args := m.Called(ctx, newImage)
	return args.Error(0)
}

func (m *mockImageScanCU) UpdatePartial(ctx context.Context, id int, updatedImage models.ImagePartial) (*models.Image, error) {
	args := m.Called(ctx, id, updatedImage)
	return args.Get(0).(*models.Image), args.Error(1)
}

func (m *mockImageScanCU) AddFileID(ctx context.Context, id int, fileID models.FileID) error {
	args := m.Called(ctx, id, fileID)
	return args.Error(0)
}

type mockScanConfig struct{}

func (m *mockScanConfig) GetCreateGalleriesFromFolders() bool { return false }

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

type mockGalleryFinderCreator struct {
	mock.Mock
}

func (m *mockGalleryFinderCreator) FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Gallery, error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).([]*models.Gallery), args.Error(1)
}

func (m *mockGalleryFinderCreator) FindByFolderID(ctx context.Context, folderID models.FolderID) ([]*models.Gallery, error) {
	args := m.Called(ctx, folderID)
	return args.Get(0).([]*models.Gallery), args.Error(1)
}

func (m *mockGalleryFinderCreator) Create(ctx context.Context, input *models.CreateGalleryInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *mockGalleryFinderCreator) UpdatePartial(ctx context.Context, id int, updatedGallery models.GalleryPartial) (*models.Gallery, error) {
	args := m.Called(ctx, id, updatedGallery)
	return args.Get(0).(*models.Gallery), args.Error(1)
}

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
			cu := &mockImageScanCU{}
			cu.On("GetFiles", mock.Anything, testImageID).Return([]models.File{existingFile}, nil)
			cu.On("GetGalleryIDs", mock.Anything, testImageID).Return([]int{}, nil)

			if tt.expectUpdate {
				cu.On("UpdatePartial", mock.Anything, testImageID, mock.Anything).
					Return(&models.Image{ID: testImageID}, nil)
			}

			h := &ScanHandler{
				CreatorUpdater: cu,
				GalleryFinder:  &mockGalleryFinderCreator{},
				ScanConfig:     &mockScanConfig{},
				PluginCache:    &plugin.Cache{},
			}

			withTxnCtx(func(ctx context.Context) {
				err := h.associateExisting(ctx, []*models.Image{makeImage()}, existingFile, tt.updateExisting)
				assert.NoError(t, err)
			})

			if tt.expectUpdate {
				cu.AssertCalled(t, "UpdatePartial", mock.Anything, testImageID, mock.Anything)
			} else {
				cu.AssertNotCalled(t, "UpdatePartial", mock.Anything, mock.Anything, mock.Anything)
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

	cu := &mockImageScanCU{}
	cu.On("GetFiles", mock.Anything, testImageID).Return([]models.File{existingFile}, nil)
	cu.On("GetGalleryIDs", mock.Anything, testImageID).Return([]int{}, nil)
	cu.On("AddFileID", mock.Anything, testImageID, models.FileID(newFileID)).Return(nil)
	cu.On("UpdatePartial", mock.Anything, testImageID, mock.Anything).
		Return(&models.Image{ID: testImageID}, nil)

	h := &ScanHandler{
		CreatorUpdater: cu,
		GalleryFinder:  &mockGalleryFinderCreator{},
		ScanConfig:     &mockScanConfig{},
		PluginCache:    &plugin.Cache{},
	}

	withTxnCtx(func(ctx context.Context) {
		err := h.associateExisting(ctx, []*models.Image{image}, newFile, false)
		assert.NoError(t, err)
	})

	cu.AssertCalled(t, "AddFileID", mock.Anything, testImageID, models.FileID(newFileID))
	cu.AssertCalled(t, "UpdatePartial", mock.Anything, testImageID, mock.Anything)
}
