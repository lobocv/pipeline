// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"

import mock "github.com/stretchr/testify/mock"

// ErrorHandler is an autogenerated mock type for the ErrorHandler type
type ErrorHandler struct {
	mock.Mock
}

// HandleError provides a mock function with given fields: _a0, _a1
func (_m *ErrorHandler) HandleError(_a0 context.Context, _a1 error) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, error) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}