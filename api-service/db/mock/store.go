// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m1thrandir225/galore-services/db/sqlc (interfaces: Store)
//
// Generated by this command:
//
//	mockgen -package mockdb -destination db/mock/store.go github.com/m1thrandir225/galore-services/db/sqlc Store
//

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateCategory mocks base method.
func (m *MockStore) CreateCategory(arg0 context.Context, arg1 db.CreateCategoryParams) (db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCategory", arg0, arg1)
	ret0, _ := ret[0].(db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCategory indicates an expected call of CreateCategory.
func (mr *MockStoreMockRecorder) CreateCategory(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCategory", reflect.TypeOf((*MockStore)(nil).CreateCategory), arg0, arg1)
}

// CreateCocktail mocks base method.
func (m *MockStore) CreateCocktail(arg0 context.Context, arg1 db.CreateCocktailParams) (db.Cocktail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCocktail", arg0, arg1)
	ret0, _ := ret[0].(db.Cocktail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCocktail indicates an expected call of CreateCocktail.
func (mr *MockStoreMockRecorder) CreateCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCocktail", reflect.TypeOf((*MockStore)(nil).CreateCocktail), arg0, arg1)
}

// CreateCocktailCategory mocks base method.
func (m *MockStore) CreateCocktailCategory(arg0 context.Context, arg1 db.CreateCocktailCategoryParams) (db.CocktailCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCocktailCategory", arg0, arg1)
	ret0, _ := ret[0].(db.CocktailCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCocktailCategory indicates an expected call of CreateCocktailCategory.
func (mr *MockStoreMockRecorder) CreateCocktailCategory(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCocktailCategory", reflect.TypeOf((*MockStore)(nil).CreateCocktailCategory), arg0, arg1)
}

// CreateFCMToken mocks base method.
func (m *MockStore) CreateFCMToken(arg0 context.Context, arg1 db.CreateFCMTokenParams) (db.FcmToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFCMToken", arg0, arg1)
	ret0, _ := ret[0].(db.FcmToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFCMToken indicates an expected call of CreateFCMToken.
func (mr *MockStoreMockRecorder) CreateFCMToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFCMToken", reflect.TypeOf((*MockStore)(nil).CreateFCMToken), arg0, arg1)
}

// CreateFlavour mocks base method.
func (m *MockStore) CreateFlavour(arg0 context.Context, arg1 string) (db.Flavour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFlavour", arg0, arg1)
	ret0, _ := ret[0].(db.Flavour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFlavour indicates an expected call of CreateFlavour.
func (mr *MockStoreMockRecorder) CreateFlavour(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFlavour", reflect.TypeOf((*MockStore)(nil).CreateFlavour), arg0, arg1)
}

// CreateNotification mocks base method.
func (m *MockStore) CreateNotification(arg0 context.Context, arg1 db.CreateNotificationParams) (db.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNotification", arg0, arg1)
	ret0, _ := ret[0].(db.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNotification indicates an expected call of CreateNotification.
func (mr *MockStoreMockRecorder) CreateNotification(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNotification", reflect.TypeOf((*MockStore)(nil).CreateNotification), arg0, arg1)
}

// CreateNotificationType mocks base method.
func (m *MockStore) CreateNotificationType(arg0 context.Context, arg1 db.CreateNotificationTypeParams) (db.NotificationType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNotificationType", arg0, arg1)
	ret0, _ := ret[0].(db.NotificationType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNotificationType indicates an expected call of CreateNotificationType.
func (mr *MockStoreMockRecorder) CreateNotificationType(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNotificationType", reflect.TypeOf((*MockStore)(nil).CreateNotificationType), arg0, arg1)
}

// CreateSession mocks base method.
func (m *MockStore) CreateSession(arg0 context.Context, arg1 db.CreateSessionParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockStoreMockRecorder) CreateSession(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockStore)(nil).CreateSession), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.CreateUserRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.CreateUserRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// CreateUserCocktail mocks base method.
func (m *MockStore) CreateUserCocktail(arg0 context.Context, arg1 db.CreateUserCocktailParams) (db.CreatedCocktail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserCocktail", arg0, arg1)
	ret0, _ := ret[0].(db.CreatedCocktail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserCocktail indicates an expected call of CreateUserCocktail.
func (mr *MockStoreMockRecorder) CreateUserCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserCocktail", reflect.TypeOf((*MockStore)(nil).CreateUserCocktail), arg0, arg1)
}

// DeleteCategory mocks base method.
func (m *MockStore) DeleteCategory(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategory", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategory indicates an expected call of DeleteCategory.
func (mr *MockStoreMockRecorder) DeleteCategory(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategory", reflect.TypeOf((*MockStore)(nil).DeleteCategory), arg0, arg1)
}

// DeleteCocktail mocks base method.
func (m *MockStore) DeleteCocktail(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCocktail", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCocktail indicates an expected call of DeleteCocktail.
func (mr *MockStoreMockRecorder) DeleteCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCocktail", reflect.TypeOf((*MockStore)(nil).DeleteCocktail), arg0, arg1)
}

// DeleteCocktailCategory mocks base method.
func (m *MockStore) DeleteCocktailCategory(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCocktailCategory", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCocktailCategory indicates an expected call of DeleteCocktailCategory.
func (mr *MockStoreMockRecorder) DeleteCocktailCategory(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCocktailCategory", reflect.TypeOf((*MockStore)(nil).DeleteCocktailCategory), arg0, arg1)
}

// DeleteFCMToken mocks base method.
func (m *MockStore) DeleteFCMToken(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFCMToken", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFCMToken indicates an expected call of DeleteFCMToken.
func (mr *MockStoreMockRecorder) DeleteFCMToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFCMToken", reflect.TypeOf((*MockStore)(nil).DeleteFCMToken), arg0, arg1)
}

// DeleteFlavour mocks base method.
func (m *MockStore) DeleteFlavour(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFlavour", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFlavour indicates an expected call of DeleteFlavour.
func (mr *MockStoreMockRecorder) DeleteFlavour(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFlavour", reflect.TypeOf((*MockStore)(nil).DeleteFlavour), arg0, arg1)
}

// DeleteNotificationType mocks base method.
func (m *MockStore) DeleteNotificationType(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNotificationType", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNotificationType indicates an expected call of DeleteNotificationType.
func (mr *MockStoreMockRecorder) DeleteNotificationType(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNotificationType", reflect.TypeOf((*MockStore)(nil).DeleteNotificationType), arg0, arg1)
}

// DeleteSession mocks base method.
func (m *MockStore) DeleteSession(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockStoreMockRecorder) DeleteSession(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockStore)(nil).DeleteSession), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0, arg1)
}

// DeleteUserCocktail mocks base method.
func (m *MockStore) DeleteUserCocktail(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserCocktail", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserCocktail indicates an expected call of DeleteUserCocktail.
func (mr *MockStoreMockRecorder) DeleteUserCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserCocktail", reflect.TypeOf((*MockStore)(nil).DeleteUserCocktail), arg0, arg1)
}

// GetAllCategories mocks base method.
func (m *MockStore) GetAllCategories(arg0 context.Context) ([]db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCategories", arg0)
	ret0, _ := ret[0].([]db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCategories indicates an expected call of GetAllCategories.
func (mr *MockStoreMockRecorder) GetAllCategories(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCategories", reflect.TypeOf((*MockStore)(nil).GetAllCategories), arg0)
}

// GetAllFlavours mocks base method.
func (m *MockStore) GetAllFlavours(arg0 context.Context) ([]db.Flavour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllFlavours", arg0)
	ret0, _ := ret[0].([]db.Flavour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllFlavours indicates an expected call of GetAllFlavours.
func (mr *MockStoreMockRecorder) GetAllFlavours(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllFlavours", reflect.TypeOf((*MockStore)(nil).GetAllFlavours), arg0)
}

// GetAllTypes mocks base method.
func (m *MockStore) GetAllTypes(arg0 context.Context) ([]db.NotificationType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTypes", arg0)
	ret0, _ := ret[0].([]db.NotificationType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTypes indicates an expected call of GetAllTypes.
func (mr *MockStoreMockRecorder) GetAllTypes(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTypes", reflect.TypeOf((*MockStore)(nil).GetAllTypes), arg0)
}

// GetAllUserSessions mocks base method.
func (m *MockStore) GetAllUserSessions(arg0 context.Context, arg1 string) ([]db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUserSessions", arg0, arg1)
	ret0, _ := ret[0].([]db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUserSessions indicates an expected call of GetAllUserSessions.
func (mr *MockStoreMockRecorder) GetAllUserSessions(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUserSessions", reflect.TypeOf((*MockStore)(nil).GetAllUserSessions), arg0, arg1)
}

// GetCategoriesForCocktail mocks base method.
func (m *MockStore) GetCategoriesForCocktail(arg0 context.Context, arg1 uuid.UUID) ([]db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoriesForCocktail", arg0, arg1)
	ret0, _ := ret[0].([]db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoriesForCocktail indicates an expected call of GetCategoriesForCocktail.
func (mr *MockStoreMockRecorder) GetCategoriesForCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoriesForCocktail", reflect.TypeOf((*MockStore)(nil).GetCategoriesForCocktail), arg0, arg1)
}

// GetCategoryById mocks base method.
func (m *MockStore) GetCategoryById(arg0 context.Context, arg1 uuid.UUID) (db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoryById", arg0, arg1)
	ret0, _ := ret[0].(db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoryById indicates an expected call of GetCategoryById.
func (mr *MockStoreMockRecorder) GetCategoryById(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoryById", reflect.TypeOf((*MockStore)(nil).GetCategoryById), arg0, arg1)
}

// GetCategoryByTag mocks base method.
func (m *MockStore) GetCategoryByTag(arg0 context.Context, arg1 string) (db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoryByTag", arg0, arg1)
	ret0, _ := ret[0].(db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoryByTag indicates an expected call of GetCategoryByTag.
func (mr *MockStoreMockRecorder) GetCategoryByTag(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoryByTag", reflect.TypeOf((*MockStore)(nil).GetCategoryByTag), arg0, arg1)
}

// GetCocktail mocks base method.
func (m *MockStore) GetCocktail(arg0 context.Context, arg1 uuid.UUID) (db.Cocktail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCocktail", arg0, arg1)
	ret0, _ := ret[0].(db.Cocktail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCocktail indicates an expected call of GetCocktail.
func (mr *MockStoreMockRecorder) GetCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCocktail", reflect.TypeOf((*MockStore)(nil).GetCocktail), arg0, arg1)
}

// GetCocktailAndSimilar mocks base method.
func (m *MockStore) GetCocktailAndSimilar(arg0 context.Context, arg1 uuid.UUID) ([]db.GetCocktailAndSimilarRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCocktailAndSimilar", arg0, arg1)
	ret0, _ := ret[0].([]db.GetCocktailAndSimilarRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCocktailAndSimilar indicates an expected call of GetCocktailAndSimilar.
func (mr *MockStoreMockRecorder) GetCocktailAndSimilar(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCocktailAndSimilar", reflect.TypeOf((*MockStore)(nil).GetCocktailAndSimilar), arg0, arg1)
}

// GetCocktailCategory mocks base method.
func (m *MockStore) GetCocktailCategory(arg0 context.Context, arg1 uuid.UUID) (db.CocktailCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCocktailCategory", arg0, arg1)
	ret0, _ := ret[0].(db.CocktailCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCocktailCategory indicates an expected call of GetCocktailCategory.
func (mr *MockStoreMockRecorder) GetCocktailCategory(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCocktailCategory", reflect.TypeOf((*MockStore)(nil).GetCocktailCategory), arg0, arg1)
}

// GetCocktailsForCategory mocks base method.
func (m *MockStore) GetCocktailsForCategory(arg0 context.Context, arg1 uuid.UUID) ([]db.GetCocktailsForCategoryRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCocktailsForCategory", arg0, arg1)
	ret0, _ := ret[0].([]db.GetCocktailsForCategoryRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCocktailsForCategory indicates an expected call of GetCocktailsForCategory.
func (mr *MockStoreMockRecorder) GetCocktailsForCategory(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCocktailsForCategory", reflect.TypeOf((*MockStore)(nil).GetCocktailsForCategory), arg0, arg1)
}

// GetFCMTokenById mocks base method.
func (m *MockStore) GetFCMTokenById(arg0 context.Context, arg1 uuid.UUID) (db.FcmToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFCMTokenById", arg0, arg1)
	ret0, _ := ret[0].(db.FcmToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFCMTokenById indicates an expected call of GetFCMTokenById.
func (mr *MockStoreMockRecorder) GetFCMTokenById(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFCMTokenById", reflect.TypeOf((*MockStore)(nil).GetFCMTokenById), arg0, arg1)
}

// GetFlavourId mocks base method.
func (m *MockStore) GetFlavourId(arg0 context.Context, arg1 uuid.UUID) (db.Flavour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlavourId", arg0, arg1)
	ret0, _ := ret[0].(db.Flavour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlavourId indicates an expected call of GetFlavourId.
func (mr *MockStoreMockRecorder) GetFlavourId(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlavourId", reflect.TypeOf((*MockStore)(nil).GetFlavourId), arg0, arg1)
}

// GetFlavourName mocks base method.
func (m *MockStore) GetFlavourName(arg0 context.Context, arg1 string) (db.Flavour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlavourName", arg0, arg1)
	ret0, _ := ret[0].(db.Flavour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlavourName indicates an expected call of GetFlavourName.
func (mr *MockStoreMockRecorder) GetFlavourName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlavourName", reflect.TypeOf((*MockStore)(nil).GetFlavourName), arg0, arg1)
}

// GetLikedCocktail mocks base method.
func (m *MockStore) GetLikedCocktail(arg0 context.Context, arg1 db.GetLikedCocktailParams) (db.GetLikedCocktailRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedCocktail", arg0, arg1)
	ret0, _ := ret[0].(db.GetLikedCocktailRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedCocktail indicates an expected call of GetLikedCocktail.
func (mr *MockStoreMockRecorder) GetLikedCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedCocktail", reflect.TypeOf((*MockStore)(nil).GetLikedCocktail), arg0, arg1)
}

// GetLikedCocktails mocks base method.
func (m *MockStore) GetLikedCocktails(arg0 context.Context, arg1 uuid.UUID) ([]db.GetLikedCocktailsRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedCocktails", arg0, arg1)
	ret0, _ := ret[0].([]db.GetLikedCocktailsRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedCocktails indicates an expected call of GetLikedCocktails.
func (mr *MockStoreMockRecorder) GetLikedCocktails(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedCocktails", reflect.TypeOf((*MockStore)(nil).GetLikedCocktails), arg0, arg1)
}

// GetLikedFlavour mocks base method.
func (m *MockStore) GetLikedFlavour(arg0 context.Context, arg1 db.GetLikedFlavourParams) (db.Flavour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedFlavour", arg0, arg1)
	ret0, _ := ret[0].(db.Flavour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedFlavour indicates an expected call of GetLikedFlavour.
func (mr *MockStoreMockRecorder) GetLikedFlavour(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedFlavour", reflect.TypeOf((*MockStore)(nil).GetLikedFlavour), arg0, arg1)
}

// GetNotificationType mocks base method.
func (m *MockStore) GetNotificationType(arg0 context.Context, arg1 uuid.UUID) (db.NotificationType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotificationType", arg0, arg1)
	ret0, _ := ret[0].(db.NotificationType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotificationType indicates an expected call of GetNotificationType.
func (mr *MockStoreMockRecorder) GetNotificationType(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotificationType", reflect.TypeOf((*MockStore)(nil).GetNotificationType), arg0, arg1)
}

// GetSession mocks base method.
func (m *MockStore) GetSession(arg0 context.Context, arg1 uuid.UUID) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockStoreMockRecorder) GetSession(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockStore)(nil).GetSession), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 uuid.UUID) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}

// GetUserByEmail mocks base method.
func (m *MockStore) GetUserByEmail(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockStoreMockRecorder) GetUserByEmail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockStore)(nil).GetUserByEmail), arg0, arg1)
}

// GetUserCocktail mocks base method.
func (m *MockStore) GetUserCocktail(arg0 context.Context, arg1 uuid.UUID) (db.CreatedCocktail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCocktail", arg0, arg1)
	ret0, _ := ret[0].(db.CreatedCocktail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCocktail indicates an expected call of GetUserCocktail.
func (mr *MockStoreMockRecorder) GetUserCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCocktail", reflect.TypeOf((*MockStore)(nil).GetUserCocktail), arg0, arg1)
}

// GetUserLikedFlavours mocks base method.
func (m *MockStore) GetUserLikedFlavours(arg0 context.Context, arg1 uuid.UUID) ([]db.Flavour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserLikedFlavours", arg0, arg1)
	ret0, _ := ret[0].([]db.Flavour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserLikedFlavours indicates an expected call of GetUserLikedFlavours.
func (mr *MockStoreMockRecorder) GetUserLikedFlavours(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserLikedFlavours", reflect.TypeOf((*MockStore)(nil).GetUserLikedFlavours), arg0, arg1)
}

// GetUserNotifications mocks base method.
func (m *MockStore) GetUserNotifications(arg0 context.Context, arg1 uuid.UUID) ([]db.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserNotifications", arg0, arg1)
	ret0, _ := ret[0].([]db.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserNotifications indicates an expected call of GetUserNotifications.
func (mr *MockStoreMockRecorder) GetUserNotifications(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserNotifications", reflect.TypeOf((*MockStore)(nil).GetUserNotifications), arg0, arg1)
}

// InvalidateSession mocks base method.
func (m *MockStore) InvalidateSession(arg0 context.Context, arg1 uuid.UUID) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InvalidateSession indicates an expected call of InvalidateSession.
func (mr *MockStoreMockRecorder) InvalidateSession(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateSession", reflect.TypeOf((*MockStore)(nil).InvalidateSession), arg0, arg1)
}

// LikeCocktail mocks base method.
func (m *MockStore) LikeCocktail(arg0 context.Context, arg1 db.LikeCocktailParams) (db.LikedCocktail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikeCocktail", arg0, arg1)
	ret0, _ := ret[0].(db.LikedCocktail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LikeCocktail indicates an expected call of LikeCocktail.
func (mr *MockStoreMockRecorder) LikeCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikeCocktail", reflect.TypeOf((*MockStore)(nil).LikeCocktail), arg0, arg1)
}

// LikeFlavour mocks base method.
func (m *MockStore) LikeFlavour(arg0 context.Context, arg1 db.LikeFlavourParams) (db.LikedFlavour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikeFlavour", arg0, arg1)
	ret0, _ := ret[0].(db.LikedFlavour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LikeFlavour indicates an expected call of LikeFlavour.
func (mr *MockStoreMockRecorder) LikeFlavour(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikeFlavour", reflect.TypeOf((*MockStore)(nil).LikeFlavour), arg0, arg1)
}

// SearchCocktails mocks base method.
func (m *MockStore) SearchCocktails(arg0 context.Context, arg1 string) ([]db.Cocktail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchCocktails", arg0, arg1)
	ret0, _ := ret[0].([]db.Cocktail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchCocktails indicates an expected call of SearchCocktails.
func (mr *MockStoreMockRecorder) SearchCocktails(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchCocktails", reflect.TypeOf((*MockStore)(nil).SearchCocktails), arg0, arg1)
}

// UnlikeCocktail mocks base method.
func (m *MockStore) UnlikeCocktail(arg0 context.Context, arg1 db.UnlikeCocktailParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlikeCocktail", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlikeCocktail indicates an expected call of UnlikeCocktail.
func (mr *MockStoreMockRecorder) UnlikeCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlikeCocktail", reflect.TypeOf((*MockStore)(nil).UnlikeCocktail), arg0, arg1)
}

// UnlikeFlavour mocks base method.
func (m *MockStore) UnlikeFlavour(arg0 context.Context, arg1 db.UnlikeFlavourParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlikeFlavour", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlikeFlavour indicates an expected call of UnlikeFlavour.
func (mr *MockStoreMockRecorder) UnlikeFlavour(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlikeFlavour", reflect.TypeOf((*MockStore)(nil).UnlikeFlavour), arg0, arg1)
}

// UpdateCategory mocks base method.
func (m *MockStore) UpdateCategory(arg0 context.Context, arg1 db.UpdateCategoryParams) (db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCategory", arg0, arg1)
	ret0, _ := ret[0].(db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCategory indicates an expected call of UpdateCategory.
func (mr *MockStoreMockRecorder) UpdateCategory(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCategory", reflect.TypeOf((*MockStore)(nil).UpdateCategory), arg0, arg1)
}

// UpdateCocktail mocks base method.
func (m *MockStore) UpdateCocktail(arg0 context.Context, arg1 db.UpdateCocktailParams) (db.Cocktail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCocktail", arg0, arg1)
	ret0, _ := ret[0].(db.Cocktail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCocktail indicates an expected call of UpdateCocktail.
func (mr *MockStoreMockRecorder) UpdateCocktail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCocktail", reflect.TypeOf((*MockStore)(nil).UpdateCocktail), arg0, arg1)
}

// UpdateFlavour mocks base method.
func (m *MockStore) UpdateFlavour(arg0 context.Context, arg1 db.UpdateFlavourParams) (db.Flavour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFlavour", arg0, arg1)
	ret0, _ := ret[0].(db.Flavour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFlavour indicates an expected call of UpdateFlavour.
func (mr *MockStoreMockRecorder) UpdateFlavour(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFlavour", reflect.TypeOf((*MockStore)(nil).UpdateFlavour), arg0, arg1)
}

// UpdateNotificationType mocks base method.
func (m *MockStore) UpdateNotificationType(arg0 context.Context, arg1 db.UpdateNotificationTypeParams) (db.NotificationType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNotificationType", arg0, arg1)
	ret0, _ := ret[0].(db.NotificationType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNotificationType indicates an expected call of UpdateNotificationType.
func (mr *MockStoreMockRecorder) UpdateNotificationType(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNotificationType", reflect.TypeOf((*MockStore)(nil).UpdateNotificationType), arg0, arg1)
}

// UpdateUserEmailNotifications mocks base method.
func (m *MockStore) UpdateUserEmailNotifications(arg0 context.Context, arg1 db.UpdateUserEmailNotificationsParams) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserEmailNotifications", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserEmailNotifications indicates an expected call of UpdateUserEmailNotifications.
func (mr *MockStoreMockRecorder) UpdateUserEmailNotifications(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserEmailNotifications", reflect.TypeOf((*MockStore)(nil).UpdateUserEmailNotifications), arg0, arg1)
}

// UpdateUserInformation mocks base method.
func (m *MockStore) UpdateUserInformation(arg0 context.Context, arg1 db.UpdateUserInformationParams) (db.UpdateUserInformationRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserInformation", arg0, arg1)
	ret0, _ := ret[0].(db.UpdateUserInformationRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserInformation indicates an expected call of UpdateUserInformation.
func (mr *MockStoreMockRecorder) UpdateUserInformation(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserInformation", reflect.TypeOf((*MockStore)(nil).UpdateUserInformation), arg0, arg1)
}

// UpdateUserNotification mocks base method.
func (m *MockStore) UpdateUserNotification(arg0 context.Context, arg1 db.UpdateUserNotificationParams) (db.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserNotification", arg0, arg1)
	ret0, _ := ret[0].(db.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserNotification indicates an expected call of UpdateUserNotification.
func (mr *MockStoreMockRecorder) UpdateUserNotification(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserNotification", reflect.TypeOf((*MockStore)(nil).UpdateUserNotification), arg0, arg1)
}

// UpdateUserPassword mocks base method.
func (m *MockStore) UpdateUserPassword(arg0 context.Context, arg1 db.UpdateUserPasswordParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockStoreMockRecorder) UpdateUserPassword(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockStore)(nil).UpdateUserPassword), arg0, arg1)
}

// UpdateUserPushNotifications mocks base method.
func (m *MockStore) UpdateUserPushNotifications(arg0 context.Context, arg1 db.UpdateUserPushNotificationsParams) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPushNotifications", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserPushNotifications indicates an expected call of UpdateUserPushNotifications.
func (mr *MockStoreMockRecorder) UpdateUserPushNotifications(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPushNotifications", reflect.TypeOf((*MockStore)(nil).UpdateUserPushNotifications), arg0, arg1)
}
