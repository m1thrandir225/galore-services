// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m1thrandir225/galore-services/cache (interfaces: KeyValueStore)
//
// Generated by this command:
//
//	mockgen -package mockcache -destination cache/mock/cache.go github.com/m1thrandir225/galore-services/cache KeyValueStore
//

// Package mockcache is a generated GoMock package.
package mockcache

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockKeyValueStore is a mock of KeyValueStore interface.
type MockKeyValueStore struct {
	ctrl     *gomock.Controller
	recorder *MockKeyValueStoreMockRecorder
}

// MockKeyValueStoreMockRecorder is the mock recorder for MockKeyValueStore.
type MockKeyValueStoreMockRecorder struct {
	mock *MockKeyValueStore
}

// NewMockKeyValueStore creates a new mock instance.
func NewMockKeyValueStore(ctrl *gomock.Controller) *MockKeyValueStore {
	mock := &MockKeyValueStore{ctrl: ctrl}
	mock.recorder = &MockKeyValueStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyValueStore) EXPECT() *MockKeyValueStoreMockRecorder {
	return m.recorder
}

// DeleteItem mocks base method.
func (m *MockKeyValueStore) DeleteItem(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItem", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItem indicates an expected call of DeleteItem.
func (mr *MockKeyValueStoreMockRecorder) DeleteItem(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItem", reflect.TypeOf((*MockKeyValueStore)(nil).DeleteItem), arg0, arg1)
}

// GetItem mocks base method.
func (m *MockKeyValueStore) GetItem(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockKeyValueStoreMockRecorder) GetItem(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockKeyValueStore)(nil).GetItem), arg0, arg1)
}

// SaveItem mocks base method.
func (m *MockKeyValueStore) SaveItem(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveItem", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveItem indicates an expected call of SaveItem.
func (mr *MockKeyValueStoreMockRecorder) SaveItem(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveItem", reflect.TypeOf((*MockKeyValueStore)(nil).SaveItem), arg0, arg1, arg2)
}
