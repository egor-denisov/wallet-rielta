// Code generated by MockGen. DO NOT EDIT.
// Source: m.go

// Package mock_gateway is a generated GoMock package.
package mock_gateway

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockWallet is a mock of Wallet interface.
type MockWallet struct {
	ctrl     *gomock.Controller
	recorder *MockWalletMockRecorder
}

// MockWalletMockRecorder is the mock recorder for MockWallet.
type MockWalletMockRecorder struct {
	mock *MockWallet
}

// NewMockWallet creates a new mock instance.
func NewMockWallet(ctrl *gomock.Controller) *MockWallet {
	mock := &MockWallet{ctrl: ctrl}
	mock.recorder = &MockWalletMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWallet) EXPECT() *MockWalletMockRecorder {
	return m.recorder
}

// RemoteCall mocks base method.
func (m *MockWallet) RemoteCall(ctx context.Context, handler string, request, response interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteCall", ctx, handler, request, response)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoteCall indicates an expected call of RemoteCall.
func (mr *MockWalletMockRecorder) RemoteCall(ctx, handler, request, response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteCall", reflect.TypeOf((*MockWallet)(nil).RemoteCall), ctx, handler, request, response)
}