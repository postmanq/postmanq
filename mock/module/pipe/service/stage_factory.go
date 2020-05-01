// Code generated by mockery v1.0.0. DO NOT EDIT.

package service

import (
	entity "github.com/postmanq/postmanq/module/pipe/entity"

	mock "github.com/stretchr/testify/mock"

	stage "github.com/postmanq/postmanq/module/pipe/service/stage"
)

// StageFactory is an autogenerated mock type for the StageFactory type
type StageFactory struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *StageFactory) Create(_a0 *entity.Stage) (stage.Stage, error) {
	ret := _m.Called(_a0)

	var r0 stage.Stage
	if rf, ok := ret.Get(0).(func(*entity.Stage) stage.Stage); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(stage.Stage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*entity.Stage) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
