// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-scalingo/v6 (interfaces: DatabasesService)

// Package scalingomock is a generated GoMock package.
package scalingomock

import (
	context "context"
	reflect "reflect"

	scalingo "github.com/Scalingo/go-scalingo/v6"
	gomock "github.com/golang/mock/gomock"
)

// MockDatabasesService is a mock of DatabasesService interface.
type MockDatabasesService struct {
	ctrl     *gomock.Controller
	recorder *MockDatabasesServiceMockRecorder
}

// MockDatabasesServiceMockRecorder is the mock recorder for MockDatabasesService.
type MockDatabasesServiceMockRecorder struct {
	mock *MockDatabasesService
}

// NewMockDatabasesService creates a new mock instance.
func NewMockDatabasesService(ctrl *gomock.Controller) *MockDatabasesService {
	mock := &MockDatabasesService{ctrl: ctrl}
	mock.recorder = &MockDatabasesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabasesService) EXPECT() *MockDatabasesServiceMockRecorder {
	return m.recorder
}

// DatabaseDisableFeature mocks base method.
func (m *MockDatabasesService) DatabaseDisableFeature(arg0 context.Context, arg1, arg2, arg3 string) (scalingo.DatabaseDisableFeatureResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatabaseDisableFeature", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(scalingo.DatabaseDisableFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DatabaseDisableFeature indicates an expected call of DatabaseDisableFeature.
func (mr *MockDatabasesServiceMockRecorder) DatabaseDisableFeature(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatabaseDisableFeature", reflect.TypeOf((*MockDatabasesService)(nil).DatabaseDisableFeature), arg0, arg1, arg2, arg3)
}

// DatabaseEnableFeature mocks base method.
func (m *MockDatabasesService) DatabaseEnableFeature(arg0 context.Context, arg1, arg2, arg3 string) (scalingo.DatabaseEnableFeatureResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatabaseEnableFeature", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(scalingo.DatabaseEnableFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DatabaseEnableFeature indicates an expected call of DatabaseEnableFeature.
func (mr *MockDatabasesServiceMockRecorder) DatabaseEnableFeature(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatabaseEnableFeature", reflect.TypeOf((*MockDatabasesService)(nil).DatabaseEnableFeature), arg0, arg1, arg2, arg3)
}

// DatabaseShow mocks base method.
func (m *MockDatabasesService) DatabaseShow(arg0 context.Context, arg1, arg2 string) (scalingo.Database, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatabaseShow", arg0, arg1, arg2)
	ret0, _ := ret[0].(scalingo.Database)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DatabaseShow indicates an expected call of DatabaseShow.
func (mr *MockDatabasesServiceMockRecorder) DatabaseShow(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatabaseShow", reflect.TypeOf((*MockDatabasesService)(nil).DatabaseShow), arg0, arg1, arg2)
}

// DatabaseUpdatePeriodicBackupsConfig mocks base method.
func (m *MockDatabasesService) DatabaseUpdatePeriodicBackupsConfig(arg0 context.Context, arg1, arg2 string, arg3 scalingo.DatabaseUpdatePeriodicBackupsConfigParams) (scalingo.Database, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatabaseUpdatePeriodicBackupsConfig", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(scalingo.Database)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DatabaseUpdatePeriodicBackupsConfig indicates an expected call of DatabaseUpdatePeriodicBackupsConfig.
func (mr *MockDatabasesServiceMockRecorder) DatabaseUpdatePeriodicBackupsConfig(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatabaseUpdatePeriodicBackupsConfig", reflect.TypeOf((*MockDatabasesService)(nil).DatabaseUpdatePeriodicBackupsConfig), arg0, arg1, arg2, arg3)
}
