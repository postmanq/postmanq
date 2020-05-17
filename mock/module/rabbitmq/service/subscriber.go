// Code generated by mockery v1.0.0. DO NOT EDIT.

package service

import (
	context "context"

	amqp "github.com/streadway/amqp"

	mock "github.com/stretchr/testify/mock"
)

// Subscriber is an autogenerated mock type for the Subscriber type
type Subscriber struct {
	mock.Mock
}

// Bind provides a mock function with given fields: _a0, _a1
func (_m *Subscriber) Bind(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *Subscriber) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Remove provides a mock function with given fields:
func (_m *Subscriber) Remove() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: _a0
func (_m *Subscriber) Subscribe(_a0 context.Context) (<-chan amqp.Delivery, error) {
	ret := _m.Called(_a0)

	var r0 <-chan amqp.Delivery
	if rf, ok := ret.Get(0).(func(context.Context) <-chan amqp.Delivery); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan amqp.Delivery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Unbind provides a mock function with given fields: _a0, _a1
func (_m *Subscriber) Unbind(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}