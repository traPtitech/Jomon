// Code generated by MockGen. DO NOT EDIT.
// Source: transaction.go

// Package mock_model is a generated GoMock package.
package mock_model

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	model "github.com/traPtitech/Jomon/model"
	reflect "reflect"
)

// MockTransactionRepository is a mock of TransactionRepository interface.
type MockTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRepositoryMockRecorder
}

// MockTransactionRepositoryMockRecorder is the mock recorder for MockTransactionRepository.
type MockTransactionRepositoryMockRecorder struct {
	mock *MockTransactionRepository
}

// NewMockTransactionRepository creates a new mock instance.
func NewMockTransactionRepository(ctrl *gomock.Controller) *MockTransactionRepository {
	mock := &MockTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionRepository) EXPECT() *MockTransactionRepositoryMockRecorder {
	return m.recorder
}

// GetTransactions mocks base method
func (m *MockTransactionRepository) GetTransactions(ctx context.Context, query model.TransactionQuery) ([]*model.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactions", ctx, query)
	ret0, _ := ret[0].([]*model.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactions indicates an expected call of GetTransactions
func (mr *MockTransactionRepositoryMockRecorder) GetTransactions(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactions", reflect.TypeOf((*MockTransactionRepository)(nil).GetTransactions), ctx, query)
}

// CreateTransaction mocks base method
func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, Amount int, Target string, tags []*uuid.UUID, group *uuid.UUID) (*model.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", ctx, Amount, Target, tags, group)
	ret0, _ := ret[0].(*model.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction
func (mr *MockTransactionRepositoryMockRecorder) CreateTransaction(ctx, Amount, Target, tags, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).CreateTransaction), ctx, Amount, Target, tags, group)
}

// GetTransaction mocks base method
func (m *MockTransactionRepository) GetTransaction(ctx context.Context, transactionID uuid.UUID) (*model.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", ctx, transactionID)
	ret0, _ := ret[0].(*model.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction
func (mr *MockTransactionRepositoryMockRecorder) GetTransaction(ctx, transactionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).GetTransaction), ctx, transactionID)
}

// UpdateTransaction mocks base method
func (m *MockTransactionRepository) UpdateTransaction(ctx context.Context, transactionID uuid.UUID, Amount int, Target string, tags []*uuid.UUID, group *uuid.UUID) (*model.TransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", ctx, transactionID, Amount, Target, tags, group)
	ret0, _ := ret[0].(*model.TransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTransaction indicates an expected call of UpdateTransaction
func (mr *MockTransactionRepositoryMockRecorder) UpdateTransaction(ctx, transactionID, Amount, Target, tags, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).UpdateTransaction), ctx, transactionID, Amount, Target, tags, group)
}
