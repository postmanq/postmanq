// Code generated by MockGen. DO NOT EDIT.
// Source: config/services.go
//
// Generated by this command:
//
//	mockgen -source config/services.go -destination config_mock/services.go
//
// Package mock_config is a generated GoMock package.
package mock_config

import (
	reflect "reflect"

	config "github.com/postmanq/postmanq/pkg/configfx/config"
	gomock "go.uber.org/mock/gomock"
)

// MockProviderFactory is a mock of ProviderFactory interface.
type MockProviderFactory struct {
	ctrl     *gomock.Controller
	recorder *MockProviderFactoryMockRecorder
}

// MockProviderFactoryMockRecorder is the mock recorder for MockProviderFactory.
type MockProviderFactoryMockRecorder struct {
	mock *MockProviderFactory
}

// NewMockProviderFactory creates a new mock instance.
func NewMockProviderFactory(ctrl *gomock.Controller) *MockProviderFactory {
	mock := &MockProviderFactory{ctrl: ctrl}
	mock.recorder = &MockProviderFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProviderFactory) EXPECT() *MockProviderFactoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProviderFactory) Create(options ...config.Option) (config.Provider, error) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(config.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockProviderFactoryMockRecorder) Create(options ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProviderFactory)(nil).Create), options...)
}

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// Populate mocks base method.
func (m *MockProvider) Populate(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Populate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Populate indicates an expected call of Populate.
func (mr *MockProviderMockRecorder) Populate(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Populate", reflect.TypeOf((*MockProvider)(nil).Populate), arg0)
}

// PopulateByKey mocks base method.
func (m *MockProvider) PopulateByKey(arg0 string, arg1 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PopulateByKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PopulateByKey indicates an expected call of PopulateByKey.
func (mr *MockProviderMockRecorder) PopulateByKey(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopulateByKey", reflect.TypeOf((*MockProvider)(nil).PopulateByKey), arg0, arg1)
}
