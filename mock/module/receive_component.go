// Code generated by mockery v1.0.0. DO NOT EDIT.

package module

import (
	module "github.com/postmanq/postmanq/module"
	mock "github.com/stretchr/testify/mock"
)

// ReceiveComponent is an autogenerated mock type for the ReceiveComponent type
type ReceiveComponent struct {
	mock.Mock
}

// OnReceive provides a mock function with given fields: _a0, _a1
func (_m *ReceiveComponent) OnReceive(_a0 chan module.Delivery, _a1 chan module.Delivery) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(chan module.Delivery, chan module.Delivery) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}