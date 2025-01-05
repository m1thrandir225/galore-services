// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m1thrandir225/galore-services/notifications (interfaces: NotificationService)
//
// Generated by this command:
//
//	mockgen -package mocknotifications -destination notifications/mock/notifications.go github.com/m1thrandir225/galore-services/notifications NotificationService
//

// Package mocknotifications is a generated GoMock package.
package mocknotifications

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockNotificationService is a mock of NotificationService interface.
type MockNotificationService struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceMockRecorder
}

// MockNotificationServiceMockRecorder is the mock recorder for MockNotificationService.
type MockNotificationServiceMockRecorder struct {
	mock *MockNotificationService
}

// NewMockNotificationService creates a new mock instance.
func NewMockNotificationService(ctrl *gomock.Controller) *MockNotificationService {
	mock := &MockNotificationService{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationService) EXPECT() *MockNotificationServiceMockRecorder {
	return m.recorder
}

// SendNotification mocks base method.
func (m *MockNotificationService) SendNotification(arg0, arg1 string, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendNotification", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendNotification indicates an expected call of SendNotification.
func (mr *MockNotificationServiceMockRecorder) SendNotification(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendNotification", reflect.TypeOf((*MockNotificationService)(nil).SendNotification), arg0, arg1, arg2)
}