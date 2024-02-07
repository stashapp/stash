// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/stashapp/stash/pkg/models"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// SceneReaderWriter is an autogenerated mock type for the SceneReaderWriter type
type SceneReaderWriter struct {
	mock.Mock
}

// AddFileID provides a mock function with given fields: ctx, id, fileID
func (_m *SceneReaderWriter) AddFileID(ctx context.Context, id int, fileID models.FileID) error {
	ret := _m.Called(ctx, id, fileID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, models.FileID) error); ok {
		r0 = rf(ctx, id, fileID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddGalleryIDs provides a mock function with given fields: ctx, sceneID, galleryIDs
func (_m *SceneReaderWriter) AddGalleryIDs(ctx context.Context, sceneID int, galleryIDs []int) error {
	ret := _m.Called(ctx, sceneID, galleryIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, []int) error); ok {
		r0 = rf(ctx, sceneID, galleryIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddO provides a mock function with given fields: ctx, id, date
func (_m *SceneReaderWriter) AddO(ctx context.Context, id int, date *time.Time) (int, error) {
	ret := _m.Called(ctx, id, date)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int, *time.Time) int); ok {
		r0 = rf(ctx, id, date)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, *time.Time) error); ok {
		r1 = rf(ctx, id, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddView provides a mock function with given fields: ctx, sceneID, date
func (_m *SceneReaderWriter) AddView(ctx context.Context, sceneID int, date *time.Time) (int, error) {
	ret := _m.Called(ctx, sceneID, date)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int, *time.Time) int); ok {
		r0 = rf(ctx, sceneID, date)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, *time.Time) error); ok {
		r1 = rf(ctx, sceneID, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// All provides a mock function with given fields: ctx
func (_m *SceneReaderWriter) All(ctx context.Context) ([]*models.Scene, error) {
	ret := _m.Called(ctx)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context) []*models.Scene); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AssignFiles provides a mock function with given fields: ctx, sceneID, fileID
func (_m *SceneReaderWriter) AssignFiles(ctx context.Context, sceneID int, fileID []models.FileID) error {
	ret := _m.Called(ctx, sceneID, fileID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, []models.FileID) error); ok {
		r0 = rf(ctx, sceneID, fileID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Count provides a mock function with given fields: ctx
func (_m *SceneReaderWriter) Count(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountByFileID provides a mock function with given fields: ctx, fileID
func (_m *SceneReaderWriter) CountByFileID(ctx context.Context, fileID models.FileID) (int, error) {
	ret := _m.Called(ctx, fileID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, models.FileID) int); ok {
		r0 = rf(ctx, fileID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.FileID) error); ok {
		r1 = rf(ctx, fileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountByMovieID provides a mock function with given fields: ctx, movieID
func (_m *SceneReaderWriter) CountByMovieID(ctx context.Context, movieID int) (int, error) {
	ret := _m.Called(ctx, movieID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, movieID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, movieID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountByPerformerID provides a mock function with given fields: ctx, performerID
func (_m *SceneReaderWriter) CountByPerformerID(ctx context.Context, performerID int) (int, error) {
	ret := _m.Called(ctx, performerID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, performerID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, performerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountByStudioID provides a mock function with given fields: ctx, studioID
func (_m *SceneReaderWriter) CountByStudioID(ctx context.Context, studioID int) (int, error) {
	ret := _m.Called(ctx, studioID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, studioID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, studioID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountByTagID provides a mock function with given fields: ctx, tagID
func (_m *SceneReaderWriter) CountByTagID(ctx context.Context, tagID int) (int, error) {
	ret := _m.Called(ctx, tagID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, tagID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, tagID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountMissingChecksum provides a mock function with given fields: ctx
func (_m *SceneReaderWriter) CountMissingChecksum(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountMissingOSHash provides a mock function with given fields: ctx
func (_m *SceneReaderWriter) CountMissingOSHash(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountViews provides a mock function with given fields: ctx, id
func (_m *SceneReaderWriter) CountViews(ctx context.Context, id int) (int, error) {
	ret := _m.Called(ctx, id)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, newScene, fileIDs
func (_m *SceneReaderWriter) Create(ctx context.Context, newScene *models.Scene, fileIDs []models.FileID) error {
	ret := _m.Called(ctx, newScene, fileIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Scene, []models.FileID) error); ok {
		r0 = rf(ctx, newScene, fileIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAllViews provides a mock function with given fields: ctx, id
func (_m *SceneReaderWriter) DeleteAllViews(ctx context.Context, id int) (int, error) {
	ret := _m.Called(ctx, id)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteO provides a mock function with given fields: ctx, id, date
func (_m *SceneReaderWriter) DeleteO(ctx context.Context, id int, date *time.Time) (int, error) {
	ret := _m.Called(ctx, id, date)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int, *time.Time) int); ok {
		r0 = rf(ctx, id, date)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, *time.Time) error); ok {
		r1 = rf(ctx, id, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteView provides a mock function with given fields: ctx, id, date
func (_m *SceneReaderWriter) DeleteView(ctx context.Context, id int, date *time.Time) (int, error) {
	ret := _m.Called(ctx, id, date)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int, *time.Time) int); ok {
		r0 = rf(ctx, id, date)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, *time.Time) error); ok {
		r1 = rf(ctx, id, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Destroy provides a mock function with given fields: ctx, id
func (_m *SceneReaderWriter) Destroy(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Duration provides a mock function with given fields: ctx
func (_m *SceneReaderWriter) Duration(ctx context.Context) (float64, error) {
	ret := _m.Called(ctx)

	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context) float64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: ctx, id
func (_m *SceneReaderWriter) Find(ctx context.Context, id int) (*models.Scene, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.Scene); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByChecksum provides a mock function with given fields: ctx, checksum
func (_m *SceneReaderWriter) FindByChecksum(ctx context.Context, checksum string) ([]*models.Scene, error) {
	ret := _m.Called(ctx, checksum)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, string) []*models.Scene); ok {
		r0 = rf(ctx, checksum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, checksum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByFileID provides a mock function with given fields: ctx, fileID
func (_m *SceneReaderWriter) FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Scene, error) {
	ret := _m.Called(ctx, fileID)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, models.FileID) []*models.Scene); ok {
		r0 = rf(ctx, fileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.FileID) error); ok {
		r1 = rf(ctx, fileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByFingerprints provides a mock function with given fields: ctx, fp
func (_m *SceneReaderWriter) FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Scene, error) {
	ret := _m.Called(ctx, fp)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, []models.Fingerprint) []*models.Scene); ok {
		r0 = rf(ctx, fp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []models.Fingerprint) error); ok {
		r1 = rf(ctx, fp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByGalleryID provides a mock function with given fields: ctx, performerID
func (_m *SceneReaderWriter) FindByGalleryID(ctx context.Context, performerID int) ([]*models.Scene, error) {
	ret := _m.Called(ctx, performerID)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.Scene); ok {
		r0 = rf(ctx, performerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, performerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByMovieID provides a mock function with given fields: ctx, movieID
func (_m *SceneReaderWriter) FindByMovieID(ctx context.Context, movieID int) ([]*models.Scene, error) {
	ret := _m.Called(ctx, movieID)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.Scene); ok {
		r0 = rf(ctx, movieID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, movieID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByOSHash provides a mock function with given fields: ctx, oshash
func (_m *SceneReaderWriter) FindByOSHash(ctx context.Context, oshash string) ([]*models.Scene, error) {
	ret := _m.Called(ctx, oshash)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, string) []*models.Scene); ok {
		r0 = rf(ctx, oshash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, oshash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByPath provides a mock function with given fields: ctx, path
func (_m *SceneReaderWriter) FindByPath(ctx context.Context, path string) ([]*models.Scene, error) {
	ret := _m.Called(ctx, path)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, string) []*models.Scene); ok {
		r0 = rf(ctx, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByPerformerID provides a mock function with given fields: ctx, performerID
func (_m *SceneReaderWriter) FindByPerformerID(ctx context.Context, performerID int) ([]*models.Scene, error) {
	ret := _m.Called(ctx, performerID)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.Scene); ok {
		r0 = rf(ctx, performerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, performerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByPrimaryFileID provides a mock function with given fields: ctx, fileID
func (_m *SceneReaderWriter) FindByPrimaryFileID(ctx context.Context, fileID models.FileID) ([]*models.Scene, error) {
	ret := _m.Called(ctx, fileID)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, models.FileID) []*models.Scene); ok {
		r0 = rf(ctx, fileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.FileID) error); ok {
		r1 = rf(ctx, fileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindDuplicates provides a mock function with given fields: ctx, distance, durationDiff
func (_m *SceneReaderWriter) FindDuplicates(ctx context.Context, distance int, durationDiff float64) ([][]*models.Scene, error) {
	ret := _m.Called(ctx, distance, durationDiff)

	var r0 [][]*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, int, float64) [][]*models.Scene); ok {
		r0 = rf(ctx, distance, durationDiff)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, float64) error); ok {
		r1 = rf(ctx, distance, durationDiff)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindMany provides a mock function with given fields: ctx, ids
func (_m *SceneReaderWriter) FindMany(ctx context.Context, ids []int) ([]*models.Scene, error) {
	ret := _m.Called(ctx, ids)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, []int) []*models.Scene); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCover provides a mock function with given fields: ctx, sceneID
func (_m *SceneReaderWriter) GetCover(ctx context.Context, sceneID int) ([]byte, error) {
	ret := _m.Called(ctx, sceneID)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, int) []byte); ok {
		r0 = rf(ctx, sceneID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, sceneID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFiles provides a mock function with given fields: ctx, relatedID
func (_m *SceneReaderWriter) GetFiles(ctx context.Context, relatedID int) ([]*models.VideoFile, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []*models.VideoFile
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.VideoFile); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.VideoFile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGalleryIDs provides a mock function with given fields: ctx, relatedID
func (_m *SceneReaderWriter) GetGalleryIDs(ctx context.Context, relatedID int) ([]int, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []int
	if rf, ok := ret.Get(0).(func(context.Context, int) []int); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManyFileIDs provides a mock function with given fields: ctx, ids
func (_m *SceneReaderWriter) GetManyFileIDs(ctx context.Context, ids []int) ([][]models.FileID, error) {
	ret := _m.Called(ctx, ids)

	var r0 [][]models.FileID
	if rf, ok := ret.Get(0).(func(context.Context, []int) [][]models.FileID); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]models.FileID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManyLastViewed provides a mock function with given fields: ctx, ids
func (_m *SceneReaderWriter) GetManyLastViewed(ctx context.Context, ids []int) ([]*time.Time, error) {
	ret := _m.Called(ctx, ids)

	var r0 []*time.Time
	if rf, ok := ret.Get(0).(func(context.Context, []int) []*time.Time); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*time.Time)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManyOCount provides a mock function with given fields: ctx, ids
func (_m *SceneReaderWriter) GetManyOCount(ctx context.Context, ids []int) ([]int, error) {
	ret := _m.Called(ctx, ids)

	var r0 []int
	if rf, ok := ret.Get(0).(func(context.Context, []int) []int); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManyODates provides a mock function with given fields: ctx, ids
func (_m *SceneReaderWriter) GetManyODates(ctx context.Context, ids []int) ([][]time.Time, error) {
	ret := _m.Called(ctx, ids)

	var r0 [][]time.Time
	if rf, ok := ret.Get(0).(func(context.Context, []int) [][]time.Time); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]time.Time)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManyViewCount provides a mock function with given fields: ctx, ids
func (_m *SceneReaderWriter) GetManyViewCount(ctx context.Context, ids []int) ([]int, error) {
	ret := _m.Called(ctx, ids)

	var r0 []int
	if rf, ok := ret.Get(0).(func(context.Context, []int) []int); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManyViewDates provides a mock function with given fields: ctx, ids
func (_m *SceneReaderWriter) GetManyViewDates(ctx context.Context, ids []int) ([][]time.Time, error) {
	ret := _m.Called(ctx, ids)

	var r0 [][]time.Time
	if rf, ok := ret.Get(0).(func(context.Context, []int) [][]time.Time); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]time.Time)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMovies provides a mock function with given fields: ctx, id
func (_m *SceneReaderWriter) GetMovies(ctx context.Context, id int) ([]models.MoviesScenes, error) {
	ret := _m.Called(ctx, id)

	var r0 []models.MoviesScenes
	if rf, ok := ret.Get(0).(func(context.Context, int) []models.MoviesScenes); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.MoviesScenes)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOCount provides a mock function with given fields: ctx, id
func (_m *SceneReaderWriter) GetOCount(ctx context.Context, id int) (int, error) {
	ret := _m.Called(ctx, id)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetODates provides a mock function with given fields: ctx, relatedID
func (_m *SceneReaderWriter) GetODates(ctx context.Context, relatedID int) ([]time.Time, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []time.Time
	if rf, ok := ret.Get(0).(func(context.Context, int) []time.Time); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]time.Time)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPerformerIDs provides a mock function with given fields: ctx, relatedID
func (_m *SceneReaderWriter) GetPerformerIDs(ctx context.Context, relatedID int) ([]int, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []int
	if rf, ok := ret.Get(0).(func(context.Context, int) []int); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStashIDs provides a mock function with given fields: ctx, relatedID
func (_m *SceneReaderWriter) GetStashIDs(ctx context.Context, relatedID int) ([]models.StashID, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []models.StashID
	if rf, ok := ret.Get(0).(func(context.Context, int) []models.StashID); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.StashID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTagIDs provides a mock function with given fields: ctx, relatedID
func (_m *SceneReaderWriter) GetTagIDs(ctx context.Context, relatedID int) ([]int, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []int
	if rf, ok := ret.Get(0).(func(context.Context, int) []int); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetURLs provides a mock function with given fields: ctx, relatedID
func (_m *SceneReaderWriter) GetURLs(ctx context.Context, relatedID int) ([]string, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, int) []string); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetViewDates provides a mock function with given fields: ctx, relatedID
func (_m *SceneReaderWriter) GetViewDates(ctx context.Context, relatedID int) ([]time.Time, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []time.Time
	if rf, ok := ret.Get(0).(func(context.Context, int) []time.Time); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]time.Time)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HasCover provides a mock function with given fields: ctx, sceneID
func (_m *SceneReaderWriter) HasCover(ctx context.Context, sceneID int) (bool, error) {
	ret := _m.Called(ctx, sceneID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int) bool); ok {
		r0 = rf(ctx, sceneID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, sceneID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OCount provides a mock function with given fields: ctx
func (_m *SceneReaderWriter) OCount(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OCountByPerformerID provides a mock function with given fields: ctx, performerID
func (_m *SceneReaderWriter) OCountByPerformerID(ctx context.Context, performerID int) (int, error) {
	ret := _m.Called(ctx, performerID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, performerID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, performerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlayCount provides a mock function with given fields: ctx
func (_m *SceneReaderWriter) PlayCount(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlayDuration provides a mock function with given fields: ctx
func (_m *SceneReaderWriter) PlayDuration(ctx context.Context) (float64, error) {
	ret := _m.Called(ctx)

	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context) float64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: ctx, options
func (_m *SceneReaderWriter) Query(ctx context.Context, options models.SceneQueryOptions) (*models.SceneQueryResult, error) {
	ret := _m.Called(ctx, options)

	var r0 *models.SceneQueryResult
	if rf, ok := ret.Get(0).(func(context.Context, models.SceneQueryOptions) *models.SceneQueryResult); ok {
		r0 = rf(ctx, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SceneQueryResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.SceneQueryOptions) error); ok {
		r1 = rf(ctx, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryCount provides a mock function with given fields: ctx, sceneFilter, findFilter
func (_m *SceneReaderWriter) QueryCount(ctx context.Context, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType) (int, error) {
	ret := _m.Called(ctx, sceneFilter, findFilter)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, *models.SceneFilterType, *models.FindFilterType) int); ok {
		r0 = rf(ctx, sceneFilter, findFilter)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.SceneFilterType, *models.FindFilterType) error); ok {
		r1 = rf(ctx, sceneFilter, findFilter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResetO provides a mock function with given fields: ctx, id
func (_m *SceneReaderWriter) ResetO(ctx context.Context, id int) (int, error) {
	ret := _m.Called(ctx, id)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveActivity provides a mock function with given fields: ctx, sceneID, resumeTime, playDuration
func (_m *SceneReaderWriter) SaveActivity(ctx context.Context, sceneID int, resumeTime *float64, playDuration *float64) (bool, error) {
	ret := _m.Called(ctx, sceneID, resumeTime, playDuration)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int, *float64, *float64) bool); ok {
		r0 = rf(ctx, sceneID, resumeTime, playDuration)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, *float64, *float64) error); ok {
		r1 = rf(ctx, sceneID, resumeTime, playDuration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Size provides a mock function with given fields: ctx
func (_m *SceneReaderWriter) Size(ctx context.Context) (float64, error) {
	ret := _m.Called(ctx)

	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context) float64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UniqueScenePlayCount provides a mock function with given fields: ctx
func (_m *SceneReaderWriter) UniqueScenePlayCount(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, updatedScene
func (_m *SceneReaderWriter) Update(ctx context.Context, updatedScene *models.Scene) error {
	ret := _m.Called(ctx, updatedScene)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Scene) error); ok {
		r0 = rf(ctx, updatedScene)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateCover provides a mock function with given fields: ctx, sceneID, cover
func (_m *SceneReaderWriter) UpdateCover(ctx context.Context, sceneID int, cover []byte) error {
	ret := _m.Called(ctx, sceneID, cover)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, []byte) error); ok {
		r0 = rf(ctx, sceneID, cover)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePartial provides a mock function with given fields: ctx, id, updatedScene
func (_m *SceneReaderWriter) UpdatePartial(ctx context.Context, id int, updatedScene models.ScenePartial) (*models.Scene, error) {
	ret := _m.Called(ctx, id, updatedScene)

	var r0 *models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, int, models.ScenePartial) *models.Scene); ok {
		r0 = rf(ctx, id, updatedScene)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, models.ScenePartial) error); ok {
		r1 = rf(ctx, id, updatedScene)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Wall provides a mock function with given fields: ctx, q
func (_m *SceneReaderWriter) Wall(ctx context.Context, q *string) ([]*models.Scene, error) {
	ret := _m.Called(ctx, q)

	var r0 []*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, *string) []*models.Scene); ok {
		r0 = rf(ctx, q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *string) error); ok {
		r1 = rf(ctx, q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
