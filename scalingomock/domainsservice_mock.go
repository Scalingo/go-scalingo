// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-scalingo/v7 (interfaces: DomainsService)

// Package scalingomock is a generated GoMock package.
package scalingomock

import (
	context "context"
	reflect "reflect"

	scalingo "github.com/Scalingo/go-scalingo/v7"
	gomock "github.com/golang/mock/gomock"
)

// MockDomainsService is a mock of DomainsService interface.
type MockDomainsService struct {
	ctrl     *gomock.Controller
	recorder *MockDomainsServiceMockRecorder
}

// MockDomainsServiceMockRecorder is the mock recorder for MockDomainsService.
type MockDomainsServiceMockRecorder struct {
	mock *MockDomainsService
}

// NewMockDomainsService creates a new mock instance.
func NewMockDomainsService(ctrl *gomock.Controller) *MockDomainsService {
	mock := &MockDomainsService{ctrl: ctrl}
	mock.recorder = &MockDomainsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDomainsService) EXPECT() *MockDomainsServiceMockRecorder {
	return m.recorder
}

// DomainSetCanonical mocks base method.
func (m *MockDomainsService) DomainSetCanonical(arg0 context.Context, arg1, arg2 string) (scalingo.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DomainSetCanonical", arg0, arg1, arg2)
	ret0, _ := ret[0].(scalingo.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DomainSetCanonical indicates an expected call of DomainSetCanonical.
func (mr *MockDomainsServiceMockRecorder) DomainSetCanonical(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainSetCanonical", reflect.TypeOf((*MockDomainsService)(nil).DomainSetCanonical), arg0, arg1, arg2)
}

// DomainSetCertificate mocks base method.
func (m *MockDomainsService) DomainSetCertificate(arg0 context.Context, arg1, arg2, arg3, arg4 string) (scalingo.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DomainSetCertificate", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(scalingo.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DomainSetCertificate indicates an expected call of DomainSetCertificate.
func (mr *MockDomainsServiceMockRecorder) DomainSetCertificate(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainSetCertificate", reflect.TypeOf((*MockDomainsService)(nil).DomainSetCertificate), arg0, arg1, arg2, arg3, arg4)
}

// DomainUnsetCanonical mocks base method.
func (m *MockDomainsService) DomainUnsetCanonical(arg0 context.Context, arg1 string) (scalingo.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DomainUnsetCanonical", arg0, arg1)
	ret0, _ := ret[0].(scalingo.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DomainUnsetCanonical indicates an expected call of DomainUnsetCanonical.
func (mr *MockDomainsServiceMockRecorder) DomainUnsetCanonical(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainUnsetCanonical", reflect.TypeOf((*MockDomainsService)(nil).DomainUnsetCanonical), arg0, arg1)
}

// DomainUnsetCertificate mocks base method.
func (m *MockDomainsService) DomainUnsetCertificate(arg0 context.Context, arg1, arg2 string) (scalingo.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DomainUnsetCertificate", arg0, arg1, arg2)
	ret0, _ := ret[0].(scalingo.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DomainUnsetCertificate indicates an expected call of DomainUnsetCertificate.
func (mr *MockDomainsServiceMockRecorder) DomainUnsetCertificate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainUnsetCertificate", reflect.TypeOf((*MockDomainsService)(nil).DomainUnsetCertificate), arg0, arg1, arg2)
}

// DomainsAdd mocks base method.
func (m *MockDomainsService) DomainsAdd(arg0 context.Context, arg1 string, arg2 scalingo.Domain) (scalingo.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DomainsAdd", arg0, arg1, arg2)
	ret0, _ := ret[0].(scalingo.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DomainsAdd indicates an expected call of DomainsAdd.
func (mr *MockDomainsServiceMockRecorder) DomainsAdd(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainsAdd", reflect.TypeOf((*MockDomainsService)(nil).DomainsAdd), arg0, arg1, arg2)
}

// DomainsList mocks base method.
func (m *MockDomainsService) DomainsList(arg0 context.Context, arg1 string) ([]scalingo.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DomainsList", arg0, arg1)
	ret0, _ := ret[0].([]scalingo.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DomainsList indicates an expected call of DomainsList.
func (mr *MockDomainsServiceMockRecorder) DomainsList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainsList", reflect.TypeOf((*MockDomainsService)(nil).DomainsList), arg0, arg1)
}

// DomainsRemove mocks base method.
func (m *MockDomainsService) DomainsRemove(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DomainsRemove", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DomainsRemove indicates an expected call of DomainsRemove.
func (mr *MockDomainsServiceMockRecorder) DomainsRemove(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainsRemove", reflect.TypeOf((*MockDomainsService)(nil).DomainsRemove), arg0, arg1, arg2)
}
