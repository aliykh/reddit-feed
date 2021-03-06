// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/aliykh/reddit-feed/internal/posts/models"
	pagination "github.com/aliykh/reddit-feed/pkg/pagination"
	gomock "github.com/golang/mock/gomock"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUseCase) Create(arg0 context.Context, arg1 *models.Post) (*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUseCaseMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUseCase)(nil).Create), arg0, arg1)
}

// GenerateFeeds mocks base method.
func (m *MockUseCase) GenerateFeeds(arg0 context.Context, arg1 *pagination.Query) (*models.Feed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateFeeds", arg0, arg1)
	ret0, _ := ret[0].(*models.Feed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateFeeds indicates an expected call of GenerateFeeds.
func (mr *MockUseCaseMockRecorder) GenerateFeeds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateFeeds", reflect.TypeOf((*MockUseCase)(nil).GenerateFeeds), arg0, arg1)
}
