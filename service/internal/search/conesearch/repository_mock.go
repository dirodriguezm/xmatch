package conesearch

import (
	"context"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
	Objects []repository.Mastercat
	Error   error
}

func (m *MockRepository) FindObjects(ctx context.Context, pixelList []int64) ([]repository.Mastercat, error) {
	m.Called(pixelList)
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Objects, nil
}

func (m *MockRepository) InsertObject(ctx context.Context, object repository.InsertObjectParams) (repository.Mastercat, error) {
	m.Called(object)
	if m.Error != nil {
		return repository.Mastercat{}, m.Error
	}
	return repository.Mastercat{}, nil
}

func (m *MockRepository) GetAllObjects(ctx context.Context) ([]repository.Mastercat, error) {
	m.Called()
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Objects, nil
}
