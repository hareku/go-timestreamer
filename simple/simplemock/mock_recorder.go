// Code generated by MockGen. DO NOT EDIT.
// Source: recorder.go

// Package simplemock is a generated GoMock package.
package simplemock

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockRecorder is a mock of Recorder interface.
type MockRecorder struct {
	ctrl     *gomock.Controller
	recorder *MockRecorderMockRecorder
}

// MockRecorderMockRecorder is the mock recorder for MockRecorder.
type MockRecorderMockRecorder struct {
	mock *MockRecorder
}

// NewMockRecorder creates a new mock instance.
func NewMockRecorder(ctrl *gomock.Controller) *MockRecorder {
	mock := &MockRecorder{ctrl: ctrl}
	mock.recorder = &MockRecorderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRecorder) EXPECT() *MockRecorderMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockRecorder) Do(ctx context.Context, t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// Do indicates an expected call of Do.
func (mr *MockRecorderMockRecorder) Do(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockRecorder)(nil).Do), ctx, t)
}
