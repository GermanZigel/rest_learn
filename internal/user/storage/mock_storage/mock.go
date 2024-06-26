// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"
	userProxy "rest/internal/userProxy"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository) Create(ctx context.Context, user userProxy.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, user)
}

// DeleteOnce mocks base method.
func (m *MockRepository) DeleteOnce(ctx context.Context, id int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOnce", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteOnce indicates an expected call of DeleteOnce.
func (mr *MockRepositoryMockRecorder) DeleteOnce(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOnce", reflect.TypeOf((*MockRepository)(nil).DeleteOnce), ctx, id)
}

// GetList mocks base method.
func (m *MockRepository) GetList(ctx context.Context) ([]userProxy.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx)
	ret0, _ := ret[0].([]userProxy.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockRepositoryMockRecorder) GetList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockRepository)(nil).GetList), ctx)
}

// GetOnce mocks base method.
func (m *MockRepository) GetOnce(ctx context.Context, id int) (userProxy.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOnce", ctx, id)
	ret0, _ := ret[0].(userProxy.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOnce indicates an expected call of GetOnce.
func (mr *MockRepositoryMockRecorder) GetOnce(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOnce", reflect.TypeOf((*MockRepository)(nil).GetOnce), ctx, id)
}

// Update mocks base method.
func (m *MockRepository) Update(arg0 context.Context, arg1 userProxy.User) (userProxy.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(userProxy.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), arg0, arg1)
}
