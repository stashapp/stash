// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	models "github.com/stashapp/stash/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// SavedFilterReaderWriter is an autogenerated mock type for the SavedFilterReaderWriter type
type SavedFilterReaderWriter struct {
	mock.Mock
}

// Create provides a mock function with given fields: obj
func (_m *SavedFilterReaderWriter) Create(obj models.SavedFilter) (*models.SavedFilter, error) {
	ret := _m.Called(obj)

	var r0 *models.SavedFilter
	if rf, ok := ret.Get(0).(func(models.SavedFilter) *models.SavedFilter); ok {
		r0 = rf(obj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.SavedFilter) error); ok {
		r1 = rf(obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Destroy provides a mock function with given fields: id
func (_m *SavedFilterReaderWriter) Destroy(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *SavedFilterReaderWriter) Find(id int) (*models.SavedFilter, error) {
	ret := _m.Called(id)

	var r0 *models.SavedFilter
	if rf, ok := ret.Get(0).(func(int) *models.SavedFilter); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByMode provides a mock function with given fields: mode
func (_m *SavedFilterReaderWriter) FindByMode(mode models.FilterMode) ([]*models.SavedFilter, error) {
	ret := _m.Called(mode)

	var r0 []*models.SavedFilter
	if rf, ok := ret.Get(0).(func(models.FilterMode) []*models.SavedFilter); ok {
		r0 = rf(mode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.FilterMode) error); ok {
		r1 = rf(mode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindDefault provides a mock function with given fields: mode
func (_m *SavedFilterReaderWriter) FindDefault(mode models.FilterMode) (*models.SavedFilter, error) {
	ret := _m.Called(mode)

	var r0 *models.SavedFilter
	if rf, ok := ret.Get(0).(func(models.FilterMode) *models.SavedFilter); ok {
		r0 = rf(mode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.FilterMode) error); ok {
		r1 = rf(mode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetDefault provides a mock function with given fields: obj
func (_m *SavedFilterReaderWriter) SetDefault(obj models.SavedFilter) (*models.SavedFilter, error) {
	ret := _m.Called(obj)

	var r0 *models.SavedFilter
	if rf, ok := ret.Get(0).(func(models.SavedFilter) *models.SavedFilter); ok {
		r0 = rf(obj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.SavedFilter) error); ok {
		r1 = rf(obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: obj
func (_m *SavedFilterReaderWriter) Update(obj models.SavedFilter) (*models.SavedFilter, error) {
	ret := _m.Called(obj)

	var r0 *models.SavedFilter
	if rf, ok := ret.Get(0).(func(models.SavedFilter) *models.SavedFilter); ok {
		r0 = rf(obj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.SavedFilter) error); ok {
		r1 = rf(obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindRecommended provides a mock function with given fields: mode
func (_m *SavedFilterReaderWriter) FindRecommended() ([]*models.SavedFilter, error) {
	ret := _m.Called()

	var r0 []*models.SavedFilter
	if rf, ok := ret.Get(0).(func() []*models.SavedFilter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SavedFilter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
