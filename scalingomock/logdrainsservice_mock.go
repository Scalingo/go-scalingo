// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-scalingo/v4 (interfaces: LogDrainsService)

// Package scalingomock is a generated GoMock package.
package scalingomock

import (
	reflect "reflect"

	scalingo "github.com/Scalingo/go-scalingo/v4"
	gomock "github.com/golang/mock/gomock"
)

// MockLogDrainsService is a mock of LogDrainsService interface
type MockLogDrainsService struct {
	ctrl     *gomock.Controller
	recorder *MockLogDrainsServiceMockRecorder
}

// MockLogDrainsServiceMockRecorder is the mock recorder for MockLogDrainsService
type MockLogDrainsServiceMockRecorder struct {
	mock *MockLogDrainsService
}

// NewMockLogDrainsService creates a new mock instance
func NewMockLogDrainsService(ctrl *gomock.Controller) *MockLogDrainsService {
	mock := &MockLogDrainsService{ctrl: ctrl}
	mock.recorder = &MockLogDrainsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogDrainsService) EXPECT() *MockLogDrainsServiceMockRecorder {
	return m.recorder
}

// LogDrainAdd mocks base method
func (m *MockLogDrainsService) LogDrainAdd(arg0 string, arg1 scalingo.LogDrainAddParams) (*scalingo.LogDrainRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogDrainAdd", arg0, arg1)
	ret0, _ := ret[0].(*scalingo.LogDrainRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogDrainAdd indicates an expected call of LogDrainAdd
func (mr *MockLogDrainsServiceMockRecorder) LogDrainAdd(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogDrainAdd", reflect.TypeOf((*MockLogDrainsService)(nil).LogDrainAdd), arg0, arg1)
}

// LogDrainAddonAdd mocks base method
func (m *MockLogDrainsService) LogDrainAddonAdd(arg0, arg1 string, arg2 scalingo.LogDrainAddParams) (*scalingo.LogDrainRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogDrainAddonAdd", arg0, arg1, arg2)
	ret0, _ := ret[0].(*scalingo.LogDrainRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogDrainAddonAdd indicates an expected call of LogDrainAddonAdd
func (mr *MockLogDrainsServiceMockRecorder) LogDrainAddonAdd(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogDrainAddonAdd", reflect.TypeOf((*MockLogDrainsService)(nil).LogDrainAddonAdd), arg0, arg1, arg2)
}

// LogDrainAddonRemove mocks base method
func (m *MockLogDrainsService) LogDrainAddonRemove(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogDrainAddonRemove", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogDrainAddonRemove indicates an expected call of LogDrainAddonRemove
func (mr *MockLogDrainsServiceMockRecorder) LogDrainAddonRemove(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogDrainAddonRemove", reflect.TypeOf((*MockLogDrainsService)(nil).LogDrainAddonRemove), arg0, arg1, arg2)
}

// LogDrainRemove mocks base method
func (m *MockLogDrainsService) LogDrainRemove(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogDrainRemove", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogDrainRemove indicates an expected call of LogDrainRemove
func (mr *MockLogDrainsServiceMockRecorder) LogDrainRemove(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogDrainRemove", reflect.TypeOf((*MockLogDrainsService)(nil).LogDrainRemove), arg0, arg1)
}

// LogDrainsAddonList mocks base method
func (m *MockLogDrainsService) LogDrainsAddonList(arg0, arg1 string) (scalingo.LogDrainsRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogDrainsAddonList", arg0, arg1)
	ret0, _ := ret[0].(scalingo.LogDrainsRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogDrainsAddonList indicates an expected call of LogDrainsAddonList
func (mr *MockLogDrainsServiceMockRecorder) LogDrainsAddonList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogDrainsAddonList", reflect.TypeOf((*MockLogDrainsService)(nil).LogDrainsAddonList), arg0, arg1)
}

// LogDrainsList mocks base method
func (m *MockLogDrainsService) LogDrainsList(arg0 string) ([]scalingo.LogDrain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogDrainsList", arg0)
	ret0, _ := ret[0].([]scalingo.LogDrain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogDrainsList indicates an expected call of LogDrainsList
func (mr *MockLogDrainsServiceMockRecorder) LogDrainsList(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogDrainsList", reflect.TypeOf((*MockLogDrainsService)(nil).LogDrainsList), arg0)
}