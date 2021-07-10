// Code generated by MockGen. DO NOT EDIT.
// Source: model/request_file.go

// Package mock_model is a generated GoMock package.
package mock_model

import (
	gomock "github.com/golang/mock/gomock"
)

// MockRequestFileRepository is a mock of RequestFileRepository interface
type MockRequestFileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRequestFileRepositoryMockRecorder
}

// MockRequestFileRepositoryMockRecorder is the mock recorder for MockRequestFileRepository
type MockRequestFileRepositoryMockRecorder struct {
	mock *MockRequestFileRepository
}

// NewMockRequestFileRepository creates a new mock instance
func NewMockRequestFileRepository(ctrl *gomock.Controller) *MockRequestFileRepository {
	mock := &MockRequestFileRepository{ctrl: ctrl}
	mock.recorder = &MockRequestFileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequestFileRepository) EXPECT() *MockRequestFileRepositoryMockRecorder {
	return m.recorder
}
