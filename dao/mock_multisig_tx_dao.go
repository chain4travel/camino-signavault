// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/chain4travel/camino-signavault/dao (interfaces: MultisigTxDao)

// Package dao is a generated GoMock package.
package dao

import (
	reflect "reflect"
	time "time"

	model "github.com/chain4travel/camino-signavault/model"
	gomock "github.com/golang/mock/gomock"
)

// MockMultisigTxDao is a mock of MultisigTxDao interface.
type MockMultisigTxDao struct {
	ctrl     *gomock.Controller
	recorder *MockMultisigTxDaoMockRecorder
}

// MockMultisigTxDaoMockRecorder is the mock recorder for MockMultisigTxDao.
type MockMultisigTxDaoMockRecorder struct {
	mock *MockMultisigTxDao
}

// NewMockMultisigTxDao creates a new mock instance.
func NewMockMultisigTxDao(ctrl *gomock.Controller) *MockMultisigTxDao {
	mock := &MockMultisigTxDao{ctrl: ctrl}
	mock.recorder = &MockMultisigTxDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMultisigTxDao) EXPECT() *MockMultisigTxDaoMockRecorder {
	return m.recorder
}

// AddSigner mocks base method.
func (m *MockMultisigTxDao) AddSigner(arg0, arg1, arg2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSigner", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSigner indicates an expected call of AddSigner.
func (mr *MockMultisigTxDaoMockRecorder) AddSigner(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSigner", reflect.TypeOf((*MockMultisigTxDao)(nil).AddSigner), arg0, arg1, arg2)
}

// CreateMultisigTx mocks base method.
func (m *MockMultisigTxDao) CreateMultisigTx(arg0, arg1 string, arg2 int, arg3, arg4, arg5, arg6, arg7 string, arg8 []string, arg9 *time.Time) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMultisigTx", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMultisigTx indicates an expected call of CreateMultisigTx.
func (mr *MockMultisigTxDaoMockRecorder) CreateMultisigTx(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMultisigTx", reflect.TypeOf((*MockMultisigTxDao)(nil).CreateMultisigTx), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
}

// GetMultisigTx mocks base method.
func (m *MockMultisigTxDao) GetMultisigTx(arg0, arg1, arg2 string) (*[]model.MultisigTx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMultisigTx", arg0, arg1, arg2)
	ret0, _ := ret[0].(*[]model.MultisigTx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMultisigTx indicates an expected call of GetMultisigTx.
func (mr *MockMultisigTxDaoMockRecorder) GetMultisigTx(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMultisigTx", reflect.TypeOf((*MockMultisigTxDao)(nil).GetMultisigTx), arg0, arg1, arg2)
}

// PendingAliasExists mocks base method.
func (m *MockMultisigTxDao) PendingAliasExists(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PendingAliasExists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PendingAliasExists indicates an expected call of PendingAliasExists.
func (mr *MockMultisigTxDaoMockRecorder) PendingAliasExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingAliasExists", reflect.TypeOf((*MockMultisigTxDao)(nil).PendingAliasExists), arg0)
}

// UpdateTransactionId mocks base method.
func (m *MockMultisigTxDao) UpdateTransactionId(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransactionId", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTransactionId indicates an expected call of UpdateTransactionId.
func (mr *MockMultisigTxDaoMockRecorder) UpdateTransactionId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransactionId", reflect.TypeOf((*MockMultisigTxDao)(nil).UpdateTransactionId), arg0, arg1)
}
