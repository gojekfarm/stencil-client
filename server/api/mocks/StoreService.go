// Code generated by mockery v2.6.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/odpf/stencil/server/models"
	mock "github.com/stretchr/testify/mock"
)

// StoreService is an autogenerated mock type for the StoreService type
type StoreService struct {
	mock.Mock
}

// Download provides a mock function with given fields: _a0, _a1
func (_m *StoreService) Download(_a0 context.Context, _a1 *models.FileMetadata) (*models.FileData, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *models.FileData
	if rf, ok := ret.Get(0).(func(context.Context, *models.FileMetadata) *models.FileData); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FileData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.FileMetadata) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListNames provides a mock function with given fields: _a0
func (_m *StoreService) ListNames(_a0 ...string) []string {
	_va := make([]interface{}, len(_a0))
	for _i := range _a0 {
		_va[_i] = _a0[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []string
	if rf, ok := ret.Get(0).(func(...string) []string); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// ListVersions provides a mock function with given fields: _a0
func (_m *StoreService) ListVersions(_a0 ...string) []string {
	_va := make([]interface{}, len(_a0))
	for _i := range _a0 {
		_va[_i] = _a0[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []string
	if rf, ok := ret.Get(0).(func(...string) []string); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Upload provides a mock function with given fields: _a0, _a1
func (_m *StoreService) Upload(_a0 context.Context, _a1 *models.DescriptorPayload) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.DescriptorPayload) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}