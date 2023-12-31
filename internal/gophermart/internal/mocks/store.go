// Code generated by MockGen. DO NOT EDIT.
// Source: internal/gophermart/internal/storage/storage.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/kripsy/gophermart/internal/gophermart/internal/models"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// GetBalance mocks base method.
func (m *MockStore) GetBalance(ctx context.Context, userName interface{}) (models.ResponseBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, userName)
	ret0, _ := ret[0].(models.ResponseBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockStoreMockRecorder) GetBalance(ctx, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockStore)(nil).GetBalance), ctx, userName)
}

// GetNewOrders mocks base method.
func (m *MockStore) GetNewOrders(ctx context.Context) ([]models.ResponseOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewOrders", ctx)
	ret0, _ := ret[0].([]models.ResponseOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewOrders indicates an expected call of GetNewOrders.
func (mr *MockStoreMockRecorder) GetNewOrders(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewOrders", reflect.TypeOf((*MockStore)(nil).GetNewOrders), ctx)
}

// GetOrders mocks base method.
func (m *MockStore) GetOrders(ctx context.Context, username interface{}) ([]models.ResponseOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", ctx, username)
	ret0, _ := ret[0].([]models.ResponseOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockStoreMockRecorder) GetOrders(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockStore)(nil).GetOrders), ctx, username)
}

// GetProcessingOrders mocks base method.
func (m *MockStore) GetProcessingOrders(ctx context.Context) ([]models.ResponseOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProcessingOrders", ctx)
	ret0, _ := ret[0].([]models.ResponseOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProcessingOrders indicates an expected call of GetProcessingOrders.
func (mr *MockStoreMockRecorder) GetProcessingOrders(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProcessingOrders", reflect.TypeOf((*MockStore)(nil).GetProcessingOrders), ctx)
}

// GetWithdraws mocks base method.
func (m *MockStore) GetWithdraws(ctx context.Context, userName interface{}) ([]models.ResponseWithdrawals, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdraws", ctx, userName)
	ret0, _ := ret[0].([]models.ResponseWithdrawals)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdraws indicates an expected call of GetWithdraws.
func (mr *MockStoreMockRecorder) GetWithdraws(ctx, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdraws", reflect.TypeOf((*MockStore)(nil).GetWithdraws), ctx, userName)
}

// PutOrder mocks base method.
func (m *MockStore) PutOrder(ctx context.Context, userName interface{}, number int64) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutOrder", ctx, userName, number)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutOrder indicates an expected call of PutOrder.
func (mr *MockStoreMockRecorder) PutOrder(ctx, userName, number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutOrder", reflect.TypeOf((*MockStore)(nil).PutOrder), ctx, userName, number)
}

// PutWithdraw mocks base method.
func (m *MockStore) PutWithdraw(ctx context.Context, userName interface{}, number int64, accrual int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutWithdraw", ctx, userName, number, accrual)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutWithdraw indicates an expected call of PutWithdraw.
func (mr *MockStoreMockRecorder) PutWithdraw(ctx, userName, number, accrual interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutWithdraw", reflect.TypeOf((*MockStore)(nil).PutWithdraw), ctx, userName, number, accrual)
}

// UpdateStatusOrder mocks base method.
func (m *MockStore) UpdateStatusOrder(ctx context.Context, number, status string, accrual int) (models.ResponseOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatusOrder", ctx, number, status, accrual)
	ret0, _ := ret[0].(models.ResponseOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStatusOrder indicates an expected call of UpdateStatusOrder.
func (mr *MockStoreMockRecorder) UpdateStatusOrder(ctx, number, status, accrual interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatusOrder", reflect.TypeOf((*MockStore)(nil).UpdateStatusOrder), ctx, number, status, accrual)
}
