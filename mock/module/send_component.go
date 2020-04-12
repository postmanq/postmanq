// Code generated by mockery v1.0.0. DO NOT EDIT.

package module

import (
	module "github.com/postmanq/postmanq/module"
	mock "github.com/stretchr/testify/mock"
)

// SendComponent is an autogenerated mock type for the SendComponent type
type SendComponent struct {
	mock.Mock
}

// GetName provides a mock function with given fields:
func (_m *SendComponent) GetName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// OnSend provides a mock function with given fields: _a0
func (_m *SendComponent) OnSend(_a0 module.Delivery) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(module.Delivery) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
