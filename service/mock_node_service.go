// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/chain4travel/camino-signavault/service (interfaces: NodeService)

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	ids "github.com/ava-labs/avalanchego/ids"
	platformvm "github.com/ava-labs/avalanchego/vms/platformvm"
	model "github.com/chain4travel/camino-signavault/model"
	gomock "github.com/golang/mock/gomock"
)

// MockNodeService is a mock of NodeService interface.
type MockNodeService struct {
	ctrl     *gomock.Controller
	recorder *MockNodeServiceMockRecorder
}

// MockNodeServiceMockRecorder is the mock recorder for MockNodeService.
type MockNodeServiceMockRecorder struct {
	mock *MockNodeService
}

// NewMockNodeService creates a new mock instance.
func NewMockNodeService(ctrl *gomock.Controller) *MockNodeService {
	mock := &MockNodeService{ctrl: ctrl}
	mock.recorder = &MockNodeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeService) EXPECT() *MockNodeServiceMockRecorder {
	return m.recorder
}

// GetAllDepositOffers mocks base method.
func (m *MockNodeService) GetAllDepositOffers(arg0 *platformvm.GetAllDepositOffersArgs) (*platformvm.GetAllDepositOffersReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDepositOffers", arg0)
	ret0, _ := ret[0].(*platformvm.GetAllDepositOffersReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDepositOffers indicates an expected call of GetAllDepositOffers.
func (mr *MockNodeServiceMockRecorder) GetAllDepositOffers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDepositOffers", reflect.TypeOf((*MockNodeService)(nil).GetAllDepositOffers), arg0)
}

// GetMultisigAlias mocks base method.
func (m *MockNodeService) GetMultisigAlias(arg0 string) (*model.AliasInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMultisigAlias", arg0)
	ret0, _ := ret[0].(*model.AliasInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMultisigAlias indicates an expected call of GetMultisigAlias.
func (mr *MockNodeServiceMockRecorder) GetMultisigAlias(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMultisigAlias", reflect.TypeOf((*MockNodeService)(nil).GetMultisigAlias), arg0)
}

// IssueTx mocks base method.
func (m *MockNodeService) IssueTx(arg0 []byte) (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IssueTx", arg0)
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IssueTx indicates an expected call of IssueTx.
func (mr *MockNodeServiceMockRecorder) IssueTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IssueTx", reflect.TypeOf((*MockNodeService)(nil).IssueTx), arg0)
}
