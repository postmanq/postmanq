// Code generated by mockery v1.0.0. DO NOT EDIT.

package module

import mock "github.com/stretchr/testify/mock"

// InitComponent is an autogenerated mock type for the InitComponent type
type InitComponent struct {
	mock.Mock
}

// OnInit provides a mock function with given fields:
func (_m *InitComponent) OnInit() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
