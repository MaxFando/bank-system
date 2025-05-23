// Code generated by MockGen. DO NOT EDIT.
// Source: credit.go

// Package bank is a generated GoMock package.
package bank

import (
	context "context"
	reflect "reflect"

	entity "github.com/MaxFando/bank-system/internal/core/bank/entity"
	transaction "github.com/MaxFando/bank-system/pkg/sqlext/transaction"
	gomock "github.com/golang/mock/gomock"
)

// MockCreditRepository is a mock of CreditRepository interface.
type MockCreditRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCreditRepositoryMockRecorder
}

// MockCreditRepositoryMockRecorder is the mock recorder for MockCreditRepository.
type MockCreditRepositoryMockRecorder struct {
	mock *MockCreditRepository
}

// NewMockCreditRepository creates a new mock instance.
func NewMockCreditRepository(ctrl *gomock.Controller) *MockCreditRepository {
	mock := &MockCreditRepository{ctrl: ctrl}
	mock.recorder = &MockCreditRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreditRepository) EXPECT() *MockCreditRepositoryMockRecorder {
	return m.recorder
}

// CreatePaymentSchedule mocks base method.
func (m *MockCreditRepository) CreatePaymentSchedule(ctx context.Context, paymentSchedule *entity.PaymentSchedule) (*entity.PaymentSchedule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePaymentSchedule", ctx, paymentSchedule)
	ret0, _ := ret[0].(*entity.PaymentSchedule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePaymentSchedule indicates an expected call of CreatePaymentSchedule.
func (mr *MockCreditRepositoryMockRecorder) CreatePaymentSchedule(ctx, paymentSchedule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePaymentSchedule", reflect.TypeOf((*MockCreditRepository)(nil).CreatePaymentSchedule), ctx, paymentSchedule)
}

// GetCreditByID mocks base method.
func (m *MockCreditRepository) GetCreditByID(ctx context.Context, creditID int32) (*entity.Credit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCreditByID", ctx, creditID)
	ret0, _ := ret[0].(*entity.Credit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCreditByID indicates an expected call of GetCreditByID.
func (mr *MockCreditRepositoryMockRecorder) GetCreditByID(ctx, creditID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCreditByID", reflect.TypeOf((*MockCreditRepository)(nil).GetCreditByID), ctx, creditID)
}

// GetPaymentSchedule mocks base method.
func (m *MockCreditRepository) GetPaymentSchedule(ctx context.Context, creditID int32) ([]entity.PaymentSchedule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPaymentSchedule", ctx, creditID)
	ret0, _ := ret[0].([]entity.PaymentSchedule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPaymentSchedule indicates an expected call of GetPaymentSchedule.
func (mr *MockCreditRepositoryMockRecorder) GetPaymentSchedule(ctx, creditID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPaymentSchedule", reflect.TypeOf((*MockCreditRepository)(nil).GetPaymentSchedule), ctx, creditID)
}

// Save mocks base method.
func (m *MockCreditRepository) Save(ctx context.Context, credit *entity.Credit) (*entity.Credit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, credit)
	ret0, _ := ret[0].(*entity.Credit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockCreditRepositoryMockRecorder) Save(ctx, credit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockCreditRepository)(nil).Save), ctx, credit)
}

// WithTx mocks base method.
func (m *MockCreditRepository) WithTx(ctx context.Context, fn transaction.AtomicFn, opts ...transaction.TxOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, fn}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WithTx", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockCreditRepositoryMockRecorder) WithTx(ctx, fn interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, fn}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockCreditRepository)(nil).WithTx), varargs...)
}
