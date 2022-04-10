// Code generated by MockGen. DO NOT EDIT.
// Source: group.go

// Package mock_model is a generated GoMock package.
package mock_model

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	model "github.com/traPtitech/Jomon/model"
	reflect "reflect"
)

// MockGroupRepository is a mock of GroupRepository interface
type MockGroupRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGroupRepositoryMockRecorder
}

// MockGroupRepositoryMockRecorder is the mock recorder for MockGroupRepository
type MockGroupRepositoryMockRecorder struct {
	mock *MockGroupRepository
}

// NewMockGroupRepository creates a new mock instance
func NewMockGroupRepository(ctrl *gomock.Controller) *MockGroupRepository {
	mock := &MockGroupRepository{ctrl: ctrl}
	mock.recorder = &MockGroupRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGroupRepository) EXPECT() *MockGroupRepositoryMockRecorder {
	return m.recorder
}

// GetGroups mocks base method
func (m *MockGroupRepository) GetGroups(ctx context.Context) ([]*model.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroups", ctx)
	ret0, _ := ret[0].([]*model.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroups indicates an expected call of GetGroups
func (mr *MockGroupRepositoryMockRecorder) GetGroups(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroups", reflect.TypeOf((*MockGroupRepository)(nil).GetGroups), ctx)
}

// GetGroup mocks base method
func (m *MockGroupRepository) GetGroup(ctx context.Context, groupID uuid.UUID) (*model.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", ctx, groupID)
	ret0, _ := ret[0].(*model.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup
func (mr *MockGroupRepositoryMockRecorder) GetGroup(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockGroupRepository)(nil).GetGroup), ctx, groupID)
}

// CreateGroup mocks base method
func (m *MockGroupRepository) CreateGroup(ctx context.Context, name, description string, budget *int) (*model.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", ctx, name, description, budget)
	ret0, _ := ret[0].(*model.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup
func (mr *MockGroupRepositoryMockRecorder) CreateGroup(ctx, name, description, budget interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockGroupRepository)(nil).CreateGroup), ctx, name, description, budget)
}

// UpdateGroup mocks base method
func (m *MockGroupRepository) UpdateGroup(ctx context.Context, groupID uuid.UUID, name, description string, budget *int) (*model.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroup", ctx, groupID, name, description, budget)
	ret0, _ := ret[0].(*model.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateGroup indicates an expected call of UpdateGroup
func (mr *MockGroupRepositoryMockRecorder) UpdateGroup(ctx, groupID, name, description, budget interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroup", reflect.TypeOf((*MockGroupRepository)(nil).UpdateGroup), ctx, groupID, name, description, budget)
}

// DeleteGroup mocks base method
func (m *MockGroupRepository) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroup", ctx, groupID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup
func (mr *MockGroupRepositoryMockRecorder) DeleteGroup(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockGroupRepository)(nil).DeleteGroup), ctx, groupID)
}

// GetOwners mocks base method
func (m *MockGroupRepository) GetOwners(ctx context.Context, groupID uuid.UUID) ([]*model.Owner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOwners", ctx, groupID)
	ret0, _ := ret[0].([]*model.Owner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOwners indicates an expected call of GetOwners
func (mr *MockGroupRepositoryMockRecorder) GetOwners(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwners", reflect.TypeOf((*MockGroupRepository)(nil).GetOwners), ctx, groupID)
}

// AddOwner mocks base method
func (m *MockGroupRepository) AddOwner(ctx context.Context, groupID, ownerID uuid.UUID) (*model.Owner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOwner", ctx, groupID, ownerID)
	ret0, _ := ret[0].(*model.Owner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddOwner indicates an expected call of AddOwner
func (mr *MockGroupRepositoryMockRecorder) AddOwner(ctx, groupID, ownerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOwner", reflect.TypeOf((*MockGroupRepository)(nil).AddOwner), ctx, groupID, ownerID)
}

// DeleteOwner mocks base method
func (m *MockGroupRepository) DeleteOwner(ctx context.Context, groupID, ownerID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOwner", ctx, groupID, ownerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOwner indicates an expected call of DeleteOwner
func (mr *MockGroupRepositoryMockRecorder) DeleteOwner(ctx, groupID, ownerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOwner", reflect.TypeOf((*MockGroupRepository)(nil).DeleteOwner), ctx, groupID, ownerID)
}

// GetMembers mocks base method
func (m *MockGroupRepository) GetMembers(ctx context.Context, groupID uuid.UUID) ([]*model.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMembers", ctx, groupID)
	ret0, _ := ret[0].([]*model.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMembers indicates an expected call of GetMembers
func (mr *MockGroupRepositoryMockRecorder) GetMembers(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMembers", reflect.TypeOf((*MockGroupRepository)(nil).GetMembers), ctx, groupID)
}

// AddMember mocks base method
func (m *MockGroupRepository) AddMember(ctx context.Context, groupID, userID uuid.UUID) (*model.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMember", ctx, groupID, userID)
	ret0, _ := ret[0].(*model.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMember indicates an expected call of AddMember
func (mr *MockGroupRepositoryMockRecorder) AddMember(ctx, groupID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMember", reflect.TypeOf((*MockGroupRepository)(nil).AddMember), ctx, groupID, userID)
}

// DeleteMember mocks base method
func (m *MockGroupRepository) DeleteMember(ctx context.Context, groupID, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMember", ctx, groupID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMember indicates an expected call of DeleteMember
func (mr *MockGroupRepositoryMockRecorder) DeleteMember(ctx, groupID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMember", reflect.TypeOf((*MockGroupRepository)(nil).DeleteMember), ctx, groupID, userID)
}
