// Code generated by MockGen. DO NOT EDIT.
// Source: internal/billing_config/repositories/billing_config_repository.go

// Package billing_config_mock is a generated GoMock package.
package billing_config_mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/okiww/billing-loan-system/internal/billing_config/models"
)

// MockBillingConfigRepositoryInterface is a mock of BillingConfigRepositoryInterface interface.
type MockBillingConfigRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockBillingConfigRepositoryInterfaceMockRecorder
}

// MockBillingConfigRepositoryInterfaceMockRecorder is the mock recorder for MockBillingConfigRepositoryInterface.
type MockBillingConfigRepositoryInterfaceMockRecorder struct {
	mock *MockBillingConfigRepositoryInterface
}

// NewMockBillingConfigRepositoryInterface creates a new mock instance.
func NewMockBillingConfigRepositoryInterface(ctrl *gomock.Controller) *MockBillingConfigRepositoryInterface {
	mock := &MockBillingConfigRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockBillingConfigRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBillingConfigRepositoryInterface) EXPECT() *MockBillingConfigRepositoryInterfaceMockRecorder {
	return m.recorder
}

// GetBillingConfigByName mocks base method.
func (m *MockBillingConfigRepositoryInterface) GetBillingConfigByName(ctx context.Context, name string) (*models.BillingConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBillingConfigByName", ctx, name)
	ret0, _ := ret[0].(*models.BillingConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBillingConfigByName indicates an expected call of GetBillingConfigByName.
func (mr *MockBillingConfigRepositoryInterfaceMockRecorder) GetBillingConfigByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBillingConfigByName", reflect.TypeOf((*MockBillingConfigRepositoryInterface)(nil).GetBillingConfigByName), ctx, name)
}
