// Code generated by MockGen. DO NOT EDIT.
// Source: request_target.go

// Package mock_model is a generated GoMock package.
package mock_model

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	model "github.com/traPtitech/Jomon/model"
	reflect "reflect"
)

// MockRequestTargetRepository is a mock of RequestTargetRepository interface
type MockRequestTargetRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRequestTargetRepositoryMockRecorder
}

// MockRequestTargetRepositoryMockRecorder is the mock recorder for MockRequestTargetRepository
type MockRequestTargetRepositoryMockRecorder struct {
	mock *MockRequestTargetRepository
}

// NewMockRequestTargetRepository creates a new mock instance
func NewMockRequestTargetRepository(ctrl *gomock.Controller) *MockRequestTargetRepository {
	mock := &MockRequestTargetRepository{ctrl: ctrl}
	mock.recorder = &MockRequestTargetRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequestTargetRepository) EXPECT() *MockRequestTargetRepositoryMockRecorder {
	return m.recorder
}

// GetRequestTargets mocks base method
func (m *MockRequestTargetRepository) GetRequestTargets(ctx context.Context, requestID uuid.UUID) ([]*model.RequestTargetDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRequestTargets", ctx, requestID)
	ret0, _ := ret[0].([]*model.RequestTargetDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRequestTargets indicates an expected call of GetRequestTargets
func (mr *MockRequestTargetRepositoryMockRecorder) GetRequestTargets(ctx, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRequestTargets", reflect.TypeOf((*MockRequestTargetRepository)(nil).GetRequestTargets), ctx, requestID)
}
