// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"

import mock "github.com/stretchr/testify/mock"

// Processor is an autogenerated mock type for the Processor type
type Processor struct {
	mock.Mock
}

// Process provides a mock function with given fields: ctx, payload
func (_m *Processor) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	ret := _m.Called(ctx, payload)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) interface{}); ok {
		r0 = rf(ctx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
