// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	file "github.com/stashapp/stash/pkg/file"
	mock "github.com/stretchr/testify/mock"

	models "github.com/stashapp/stash/pkg/models"
)

// SceneReaderWriter is an autogenerated mock type for the SceneReaderWriter type
type SceneReaderWriter struct {
	mock.Mock
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

// Create provides a mock function with given fields: ctx, newScene, fileIDs
func (_m *SceneReaderWriter) Create(ctx context.Context, newScene *models.Scene, fileIDs []file.ID) error {
	ret := _m.Called(ctx, newScene, fileIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Scene, []file.ID) error); ok {
		r0 = rf(ctx, newScene, fileIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DecrementOCounter provides a mock function with given fields: ctx, id
func (_m *SceneReaderWriter) DecrementOCounter(ctx context.Context, id int) (int, error) {
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

// DestroyCover provides a mock function with given fields: ctx, sceneID
func (_m *SceneReaderWriter) DestroyCover(ctx context.Context, sceneID int) error {
	ret := _m.Called(ctx, sceneID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, sceneID)
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

// FindDuplicates provides a mock function with given fields: ctx, distance
func (_m *SceneReaderWriter) FindDuplicates(ctx context.Context, distance int) ([][]*models.Scene, error) {
	ret := _m.Called(ctx, distance)

	var r0 [][]*models.Scene
	if rf, ok := ret.Get(0).(func(context.Context, int) [][]*models.Scene); ok {
		r0 = rf(ctx, distance)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]*models.Scene)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, distance)
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

// GetCaptions provides a mock function with given fields: ctx, sceneID
func (_m *SceneReaderWriter) GetCaptions(ctx context.Context, sceneID int) ([]*models.SceneCaption, error) {
	ret := _m.Called(ctx, sceneID)

	var r0 []*models.SceneCaption
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.SceneCaption); ok {
		r0 = rf(ctx, sceneID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SceneCaption)
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

// IncrementOCounter provides a mock function with given fields: ctx, id
func (_m *SceneReaderWriter) IncrementOCounter(ctx context.Context, id int) (int, error) {
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

// ResetOCounter provides a mock function with given fields: ctx, id
func (_m *SceneReaderWriter) ResetOCounter(ctx context.Context, id int) (int, error) {
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

// UpdateCaptions provides a mock function with given fields: ctx, id, captions
func (_m *SceneReaderWriter) UpdateCaptions(ctx context.Context, id int, captions []*models.SceneCaption) error {
	ret := _m.Called(ctx, id, captions)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, []*models.SceneCaption) error); ok {
		r0 = rf(ctx, id, captions)
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
