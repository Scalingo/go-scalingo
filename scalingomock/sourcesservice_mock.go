// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-scalingo/v4 (interfaces: SourcesService)

// Package scalingomock is a generated GoMock package.
package scalingomock

import (
	reflect "reflect"

	scalingo "github.com/Scalingo/go-scalingo/v4"
	gomock "github.com/golang/mock/gomock"
)

// MockSourcesService is a mock of SourcesService interface
type MockSourcesService struct {
	ctrl     *gomock.Controller
	recorder *MockSourcesServiceMockRecorder
}

// MockSourcesServiceMockRecorder is the mock recorder for MockSourcesService
type MockSourcesServiceMockRecorder struct {
	mock *MockSourcesService
}

// NewMockSourcesService creates a new mock instance
func NewMockSourcesService(ctrl *gomock.Controller) *MockSourcesService {
	mock := &MockSourcesService{ctrl: ctrl}
	mock.recorder = &MockSourcesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSourcesService) EXPECT() *MockSourcesServiceMockRecorder {
	return m.recorder
}

// SourcesCreate mocks base method
func (m *MockSourcesService) SourcesCreate() (*scalingo.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SourcesCreate")
	ret0, _ := ret[0].(*scalingo.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SourcesCreate indicates an expected call of SourcesCreate
func (mr *MockSourcesServiceMockRecorder) SourcesCreate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SourcesCreate", reflect.TypeOf((*MockSourcesService)(nil).SourcesCreate))
}
