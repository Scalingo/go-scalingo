// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-scalingo/v4 (interfaces: LogsArchivesService)

// Package scalingomock is a generated GoMock package.
package scalingomock

import (
	reflect "reflect"

	scalingo "github.com/Scalingo/go-scalingo/v4"
	gomock "github.com/golang/mock/gomock"
)

// MockLogsArchivesService is a mock of LogsArchivesService interface
type MockLogsArchivesService struct {
	ctrl     *gomock.Controller
	recorder *MockLogsArchivesServiceMockRecorder
}

// MockLogsArchivesServiceMockRecorder is the mock recorder for MockLogsArchivesService
type MockLogsArchivesServiceMockRecorder struct {
	mock *MockLogsArchivesService
}

// NewMockLogsArchivesService creates a new mock instance
func NewMockLogsArchivesService(ctrl *gomock.Controller) *MockLogsArchivesService {
	mock := &MockLogsArchivesService{ctrl: ctrl}
	mock.recorder = &MockLogsArchivesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogsArchivesService) EXPECT() *MockLogsArchivesServiceMockRecorder {
	return m.recorder
}

// LogsArchives mocks base method
func (m *MockLogsArchivesService) LogsArchives(arg0 string, arg1 int) (*scalingo.LogsArchivesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogsArchives", arg0, arg1)
	ret0, _ := ret[0].(*scalingo.LogsArchivesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogsArchives indicates an expected call of LogsArchives
func (mr *MockLogsArchivesServiceMockRecorder) LogsArchives(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogsArchives", reflect.TypeOf((*MockLogsArchivesService)(nil).LogsArchives), arg0, arg1)
}

// LogsArchivesByCursor mocks base method
func (m *MockLogsArchivesService) LogsArchivesByCursor(arg0, arg1 string) (*scalingo.LogsArchivesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogsArchivesByCursor", arg0, arg1)
	ret0, _ := ret[0].(*scalingo.LogsArchivesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogsArchivesByCursor indicates an expected call of LogsArchivesByCursor
func (mr *MockLogsArchivesServiceMockRecorder) LogsArchivesByCursor(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogsArchivesByCursor", reflect.TypeOf((*MockLogsArchivesService)(nil).LogsArchivesByCursor), arg0, arg1)
}
