// Code generated by MockGen. DO NOT EDIT.
// Source: internal/payment/repositories/payment_repository.go

// Package payment_mock is a generated GoMock package.
package payment_mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/okiww/billing-loan-system/internal/payment/models"
)

// MockPaymentRepositoryInterface is a mock of PaymentRepositoryInterface interface.
type MockPaymentRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentRepositoryInterfaceMockRecorder
}

// MockPaymentRepositoryInterfaceMockRecorder is the mock recorder for MockPaymentRepositoryInterface.
type MockPaymentRepositoryInterfaceMockRecorder struct {
	mock *MockPaymentRepositoryInterface
}

// NewMockPaymentRepositoryInterface creates a new mock instance.
func NewMockPaymentRepositoryInterface(ctrl *gomock.Controller) *MockPaymentRepositoryInterface {
	mock := &MockPaymentRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockPaymentRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPaymentRepositoryInterface) EXPECT() *MockPaymentRepositoryInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPaymentRepositoryInterface) Create(ctx context.Context, payment *models.Payment) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, payment)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPaymentRepositoryInterfaceMockRecorder) Create(ctx, payment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPaymentRepositoryInterface)(nil).Create), ctx, payment)
}

// GetPaymentByID mocks base method.
func (m *MockPaymentRepositoryInterface) GetPaymentByID(ctx context.Context, id int32) (*models.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPaymentByID", ctx, id)
	ret0, _ := ret[0].(*models.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPaymentByID indicates an expected call of GetPaymentByID.
func (mr *MockPaymentRepositoryInterfaceMockRecorder) GetPaymentByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPaymentByID", reflect.TypeOf((*MockPaymentRepositoryInterface)(nil).GetPaymentByID), ctx, id)
}

// UpdatePaymentStatus mocks base method.
func (m *MockPaymentRepositoryInterface) UpdatePaymentStatus(ctx context.Context, id int32, status, note string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePaymentStatus", ctx, id, status, note)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePaymentStatus indicates an expected call of UpdatePaymentStatus.
func (mr *MockPaymentRepositoryInterfaceMockRecorder) UpdatePaymentStatus(ctx, id, status, note interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePaymentStatus", reflect.TypeOf((*MockPaymentRepositoryInterface)(nil).UpdatePaymentStatus), ctx, id, status, note)
}
