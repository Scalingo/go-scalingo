// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-scalingo/v6 (interfaces: ContainerSizesService)

// Package scalingomock is a generated GoMock package.
package scalingomock

import (
	context "context"
	reflect "reflect"

	scalingo "github.com/Scalingo/go-scalingo/v6"
	gomock "github.com/golang/mock/gomock"
)

// MockContainerSizesService is a mock of ContainerSizesService interface.
type MockContainerSizesService struct {
	ctrl     *gomock.Controller
	recorder *MockContainerSizesServiceMockRecorder
}

// MockContainerSizesServiceMockRecorder is the mock recorder for MockContainerSizesService.
type MockContainerSizesServiceMockRecorder struct {
	mock *MockContainerSizesService
}

// NewMockContainerSizesService creates a new mock instance.
func NewMockContainerSizesService(ctrl *gomock.Controller) *MockContainerSizesService {
	mock := &MockContainerSizesService{ctrl: ctrl}
	mock.recorder = &MockContainerSizesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContainerSizesService) EXPECT() *MockContainerSizesServiceMockRecorder {
	return m.recorder
}

// ContainerSizesList mocks base method.
func (m *MockContainerSizesService) ContainerSizesList(arg0 context.Context) ([]scalingo.ContainerSize, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerSizesList", arg0)
	ret0, _ := ret[0].([]scalingo.ContainerSize)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerSizesList indicates an expected call of ContainerSizesList.
func (mr *MockContainerSizesServiceMockRecorder) ContainerSizesList(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerSizesList", reflect.TypeOf((*MockContainerSizesService)(nil).ContainerSizesList), arg0)
}
