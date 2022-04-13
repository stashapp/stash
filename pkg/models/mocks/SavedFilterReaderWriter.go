// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/stashapp/stash/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// SavedFilterReaderWriter is an autogenerated mock type for the SavedFilterReaderWriter type
type SavedFilterReaderWriter struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, obj
func (_m *SavedFilterReaderWriter) Create(ctx context.Context, obj models.SavedFilter) (*models.SavedFilter, error) {
	ret := _m.Called(ctx, obj)

	var r0 *models.SavedFilter
	if rf, ok := ret.Get(0).(func(context.Context, models.SavedFilter) *models.SavedFilter); ok {
		r0 = rf(ctx, obj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.SavedFilter) error); ok {
		r1 = rf(ctx, obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Destroy provides a mock function with given fields: ctx, id
func (_m *SavedFilterReaderWriter) Destroy(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: ctx, id
func (_m *SavedFilterReaderWriter) Find(ctx context.Context, id int) (*models.SavedFilter, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.SavedFilter
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.SavedFilter); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedFilter)
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

// FindByMode provides a mock function with given fields: ctx, mode
func (_m *SavedFilterReaderWriter) FindByMode(ctx context.Context, mode models.FilterMode) ([]*models.SavedFilter, error) {
	ret := _m.Called(ctx, mode)

	var r0 []*models.SavedFilter
	if rf, ok := ret.Get(0).(func(context.Context, models.FilterMode) []*models.SavedFilter); ok {
		r0 = rf(ctx, mode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.FilterMode) error); ok {
		r1 = rf(ctx, mode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindDefault provides a mock function with given fields: ctx, mode
func (_m *SavedFilterReaderWriter) FindDefault(ctx context.Context, mode models.FilterMode) (*models.SavedFilter, error) {
	ret := _m.Called(ctx, mode)

	var r0 *models.SavedFilter
	if rf, ok := ret.Get(0).(func(context.Context, models.FilterMode) *models.SavedFilter); ok {
		r0 = rf(ctx, mode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.FilterMode) error); ok {
		r1 = rf(ctx, mode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetDefault provides a mock function with given fields: ctx, obj
func (_m *SavedFilterReaderWriter) SetDefault(ctx context.Context, obj models.SavedFilter) (*models.SavedFilter, error) {
	ret := _m.Called(ctx, obj)

	var r0 *models.SavedFilter
	if rf, ok := ret.Get(0).(func(context.Context, models.SavedFilter) *models.SavedFilter); ok {
		r0 = rf(ctx, obj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.SavedFilter) error); ok {
		r1 = rf(ctx, obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, obj
func (_m *SavedFilterReaderWriter) Update(ctx context.Context, obj models.SavedFilter) (*models.SavedFilter, error) {
	ret := _m.Called(ctx, obj)

	var r0 *models.SavedFilter
	if rf, ok := ret.Get(0).(func(context.Context, models.SavedFilter) *models.SavedFilter); ok {
		r0 = rf(ctx, obj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.SavedFilter) error); ok {
		r1 = rf(ctx, obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
