// Code generated by MockGen. DO NOT EDIT.
// Source: db/sqlc/querier.go
//
// Generated by this command:
//
//	mockgen -source=db/sqlc/querier.go -destination=db/mock/entity.go Entity
//

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	db "mole/db/sqlc"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
	isgomock struct{}
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// CreateEntity mocks base method.
func (m *MockQuerier) CreateEntity(ctx context.Context, arg db.CreateEntityParams) (db.CreateEntityRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEntity", ctx, arg)
	ret0, _ := ret[0].(db.CreateEntityRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEntity indicates an expected call of CreateEntity.
func (mr *MockQuerierMockRecorder) CreateEntity(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEntity", reflect.TypeOf((*MockQuerier)(nil).CreateEntity), ctx, arg)
}

// DeleteEntity mocks base method.
func (m *MockQuerier) DeleteEntity(ctx context.Context, entityID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEntity", ctx, entityID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEntity indicates an expected call of DeleteEntity.
func (mr *MockQuerierMockRecorder) DeleteEntity(ctx, entityID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEntity", reflect.TypeOf((*MockQuerier)(nil).DeleteEntity), ctx, entityID)
}

// GetEntitiesByNames mocks base method.
func (m *MockQuerier) GetEntitiesByNames(ctx context.Context, dollar_1 []string) ([]db.GetEntitiesByNamesRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntitiesByNames", ctx, dollar_1)
	ret0, _ := ret[0].([]db.GetEntitiesByNamesRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntitiesByNames indicates an expected call of GetEntitiesByNames.
func (mr *MockQuerierMockRecorder) GetEntitiesByNames(ctx, dollar_1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntitiesByNames", reflect.TypeOf((*MockQuerier)(nil).GetEntitiesByNames), ctx, dollar_1)
}

// GetEntity mocks base method.
func (m *MockQuerier) GetEntity(ctx context.Context, entityID uuid.UUID) (db.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntity", ctx, entityID)
	ret0, _ := ret[0].(db.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntity indicates an expected call of GetEntity.
func (mr *MockQuerierMockRecorder) GetEntity(ctx, entityID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntity", reflect.TypeOf((*MockQuerier)(nil).GetEntity), ctx, entityID)
}

// GetEntityByNameAndIntegrationSource mocks base method.
func (m *MockQuerier) GetEntityByNameAndIntegrationSource(ctx context.Context, arg db.GetEntityByNameAndIntegrationSourceParams) (db.GetEntityByNameAndIntegrationSourceRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntityByNameAndIntegrationSource", ctx, arg)
	ret0, _ := ret[0].(db.GetEntityByNameAndIntegrationSourceRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntityByNameAndIntegrationSource indicates an expected call of GetEntityByNameAndIntegrationSource.
func (mr *MockQuerierMockRecorder) GetEntityByNameAndIntegrationSource(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntityByNameAndIntegrationSource", reflect.TypeOf((*MockQuerier)(nil).GetEntityByNameAndIntegrationSource), ctx, arg)
}

// GetEntityByNames mocks base method.
func (m *MockQuerier) GetEntityByNames(ctx context.Context, name string) ([]db.GetEntityByNamesRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntityByNames", ctx, name)
	ret0, _ := ret[0].([]db.GetEntityByNamesRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntityByNames indicates an expected call of GetEntityByNames.
func (mr *MockQuerierMockRecorder) GetEntityByNames(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntityByNames", reflect.TypeOf((*MockQuerier)(nil).GetEntityByNames), ctx, name)
}

// GetInstances mocks base method.
func (m *MockQuerier) GetInstances(ctx context.Context) ([]db.GetInstancesRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstances", ctx)
	ret0, _ := ret[0].([]db.GetInstancesRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstances indicates an expected call of GetInstances.
func (mr *MockQuerierMockRecorder) GetInstances(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstances", reflect.TypeOf((*MockQuerier)(nil).GetInstances), ctx)
}

// InsertInstance mocks base method.
func (m *MockQuerier) InsertInstance(ctx context.Context, arg db.InsertInstanceParams) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertInstance", ctx, arg)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertInstance indicates an expected call of InsertInstance.
func (mr *MockQuerierMockRecorder) InsertInstance(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertInstance", reflect.TypeOf((*MockQuerier)(nil).InsertInstance), ctx, arg)
}

// InsertPosition mocks base method.
func (m *MockQuerier) InsertPosition(ctx context.Context, arg db.InsertPositionParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertPosition", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertPosition indicates an expected call of InsertPosition.
func (mr *MockQuerierMockRecorder) InsertPosition(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertPosition", reflect.TypeOf((*MockQuerier)(nil).InsertPosition), ctx, arg)
}

// ListEntities mocks base method.
func (m *MockQuerier) ListEntities(ctx context.Context, arg db.ListEntitiesParams) ([]db.ListEntitiesRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEntities", ctx, arg)
	ret0, _ := ret[0].([]db.ListEntitiesRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEntities indicates an expected call of ListEntities.
func (mr *MockQuerierMockRecorder) ListEntities(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEntities", reflect.TypeOf((*MockQuerier)(nil).ListEntities), ctx, arg)
}

// UpdateEntityByName mocks base method.
func (m *MockQuerier) UpdateEntityByName(ctx context.Context, arg db.UpdateEntityByNameParams) (db.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEntityByName", ctx, arg)
	ret0, _ := ret[0].(db.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEntityByName indicates an expected call of UpdateEntityByName.
func (mr *MockQuerierMockRecorder) UpdateEntityByName(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEntityByName", reflect.TypeOf((*MockQuerier)(nil).UpdateEntityByName), ctx, arg)
}

// UpdateEntityIntegrationSourceByNameAndSource mocks base method.
func (m *MockQuerier) UpdateEntityIntegrationSourceByNameAndSource(ctx context.Context, arg db.UpdateEntityIntegrationSourceByNameAndSourceParams) (db.UpdateEntityIntegrationSourceByNameAndSourceRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEntityIntegrationSourceByNameAndSource", ctx, arg)
	ret0, _ := ret[0].(db.UpdateEntityIntegrationSourceByNameAndSourceRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEntityIntegrationSourceByNameAndSource indicates an expected call of UpdateEntityIntegrationSourceByNameAndSource.
func (mr *MockQuerierMockRecorder) UpdateEntityIntegrationSourceByNameAndSource(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEntityIntegrationSourceByNameAndSource", reflect.TypeOf((*MockQuerier)(nil).UpdateEntityIntegrationSourceByNameAndSource), ctx, arg)
}
