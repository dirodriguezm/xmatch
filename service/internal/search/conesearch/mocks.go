package conesearch

import (
	"context"
	"database/sql"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	mock "github.com/stretchr/testify/mock"
)

type MockMastercatStore struct {
	mock.Mock
}

func NewMockMastercatStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMastercatStore {
	m := &MockMastercatStore{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *MockMastercatStore) FindObjects(ctx context.Context, pixels []int64) ([]repository.Mastercat, error) {
	args := m.Called(ctx, pixels)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.Mastercat), args.Error(1)
}

func (m *MockMastercatStore) InsertMastercat(ctx context.Context, arg repository.Mastercat) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockMastercatStore) GetAllObjects(ctx context.Context) ([]repository.Mastercat, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.Mastercat), args.Error(1)
}

func (m *MockMastercatStore) RemoveAllObjects(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMastercatStore) BulkInsertObject(ctx context.Context, db *sql.DB, rows []any) error {
	args := m.Called(ctx, db, rows)
	return args.Error(0)
}
