// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m1thrandir225/galore-services/scheduler (interfaces: SchedulerService)
//
// Generated by this command:
//
//	mockgen -package mockscheduler -destination scheduler/mock/scheduler.go github.com/m1thrandir225/galore-services/scheduler SchedulerService
//

// Package mockscheduler is a generated GoMock package.
package mockscheduler

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSchedulerService is a mock of SchedulerService interface.
type MockSchedulerService struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerServiceMockRecorder
}

// MockSchedulerServiceMockRecorder is the mock recorder for MockSchedulerService.
type MockSchedulerServiceMockRecorder struct {
	mock *MockSchedulerService
}

// NewMockSchedulerService creates a new mock instance.
func NewMockSchedulerService(ctrl *gomock.Controller) *MockSchedulerService {
	mock := &MockSchedulerService{ctrl: ctrl}
	mock.recorder = &MockSchedulerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchedulerService) EXPECT() *MockSchedulerServiceMockRecorder {
	return m.recorder
}

// EnqueueJob mocks base method.
func (m *MockSchedulerService) EnqueueJob(arg0 string, arg1 map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnqueueJob", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnqueueJob indicates an expected call of EnqueueJob.
func (mr *MockSchedulerServiceMockRecorder) EnqueueJob(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnqueueJob", reflect.TypeOf((*MockSchedulerService)(nil).EnqueueJob), arg0, arg1)
}

// RegisterCronJob mocks base method.
func (m *MockSchedulerService) RegisterCronJob(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterCronJob", arg0, arg1)
}

// RegisterCronJob indicates an expected call of RegisterCronJob.
func (mr *MockSchedulerServiceMockRecorder) RegisterCronJob(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterCronJob", reflect.TypeOf((*MockSchedulerService)(nil).RegisterCronJob), arg0, arg1)
}

// RegisterJob mocks base method.
func (m *MockSchedulerService) RegisterJob(arg0 string, arg1 bool, arg2 func(map[string]any) error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterJob", arg0, arg1, arg2)
}

// RegisterJob indicates an expected call of RegisterJob.
func (mr *MockSchedulerServiceMockRecorder) RegisterJob(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterJob", reflect.TypeOf((*MockSchedulerService)(nil).RegisterJob), arg0, arg1, arg2)
}

// Start mocks base method.
func (m *MockSchedulerService) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockSchedulerServiceMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockSchedulerService)(nil).Start))
}

// Stop mocks base method.
func (m *MockSchedulerService) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockSchedulerServiceMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockSchedulerService)(nil).Stop))
}
