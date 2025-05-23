// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-scalingo/v8 (interfaces: AutoscalersService)

// Package scalingomock is a generated GoMock package.
package scalingomock

import (
	context "context"
	reflect "reflect"

	scalingo "github.com/Scalingo/go-scalingo/v8"
	gomock "github.com/golang/mock/gomock"
)

// MockAutoscalersService is a mock of AutoscalersService interface.
type MockAutoscalersService struct {
	ctrl     *gomock.Controller
	recorder *MockAutoscalersServiceMockRecorder
}

// MockAutoscalersServiceMockRecorder is the mock recorder for MockAutoscalersService.
type MockAutoscalersServiceMockRecorder struct {
	mock *MockAutoscalersService
}

// NewMockAutoscalersService creates a new mock instance.
func NewMockAutoscalersService(ctrl *gomock.Controller) *MockAutoscalersService {
	mock := &MockAutoscalersService{ctrl: ctrl}
	mock.recorder = &MockAutoscalersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAutoscalersService) EXPECT() *MockAutoscalersServiceMockRecorder {
	return m.recorder
}

// AutoscalerAdd mocks base method.
func (m *MockAutoscalersService) AutoscalerAdd(arg0 context.Context, arg1 string, arg2 scalingo.AutoscalerAddParams) (*scalingo.Autoscaler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AutoscalerAdd", arg0, arg1, arg2)
	ret0, _ := ret[0].(*scalingo.Autoscaler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AutoscalerAdd indicates an expected call of AutoscalerAdd.
func (mr *MockAutoscalersServiceMockRecorder) AutoscalerAdd(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AutoscalerAdd", reflect.TypeOf((*MockAutoscalersService)(nil).AutoscalerAdd), arg0, arg1, arg2)
}

// AutoscalerRemove mocks base method.
func (m *MockAutoscalersService) AutoscalerRemove(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AutoscalerRemove", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AutoscalerRemove indicates an expected call of AutoscalerRemove.
func (mr *MockAutoscalersServiceMockRecorder) AutoscalerRemove(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AutoscalerRemove", reflect.TypeOf((*MockAutoscalersService)(nil).AutoscalerRemove), arg0, arg1, arg2)
}

// AutoscalersList mocks base method.
func (m *MockAutoscalersService) AutoscalersList(arg0 context.Context, arg1 string) ([]scalingo.Autoscaler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AutoscalersList", arg0, arg1)
	ret0, _ := ret[0].([]scalingo.Autoscaler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AutoscalersList indicates an expected call of AutoscalersList.
func (mr *MockAutoscalersServiceMockRecorder) AutoscalersList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AutoscalersList", reflect.TypeOf((*MockAutoscalersService)(nil).AutoscalersList), arg0, arg1)
}
