package metadata

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMetadata_ValidateCatalog(t *testing.T) {
	m := &MetadataService{}

	err := m.validateCatalog("allwise")
	require.Nil(t, err)
	err = m.validateCatalog("AllWise")
	require.Nil(t, err)

	err = m.validateCatalog("vlass")
	require.Nil(t, err)
	err = m.validateCatalog("Vlass")
	require.Nil(t, err)

	err = m.validateCatalog("ztf")
	require.Nil(t, err)
	err = m.validateCatalog("ZTF")
	require.Nil(t, err)

	err = m.validateCatalog("invalid")
	require.NotNil(t, err)
	require.Equal(t, "Could not parse field catalog with value invalid: Allowed catalogs are [allwise vlass ztf]", err.Error())
}

func TestMetadata_FindByID(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetAllwise", mock.Anything, "allwise1").Return(repository.Allwise{ID: "allwise1"}, nil)

	m := &MetadataService{
		repository: repo,
	}

	result, err := m.FindByID(context.TODO(), "allwise1", "allwise")
	require.Nil(t, err)
	repo.AssertExpectations(t)
	require.Equal(t, "allwise1", *result.(repository.AllwiseMetadata).Source_id)
}

func TestMetadata_EmptyResult(t *testing.T) {
	repo := &mocks.Repository{}

	repo.
		On("GetAllwise", mock.Anything, "allwise1").
		Return(repository.Allwise{}, sql.ErrNoRows) // sql.ErrNoRows is returned when no rows are found

	m := &MetadataService{
		repository: repo,
	}

	_, err := m.FindByID(context.TODO(), "allwise1", "allwise")
	require.NotNil(t, err)
	repo.AssertExpectations(t)
	require.EqualError(t, err, "sql: no rows in result set")
}

func TestMetadata_SomeDBError(t *testing.T) {
	repo := &mocks.Repository{}

	repo.
		On("GetAllwise", mock.Anything, "allwise1").
		Return(repository.Allwise{}, fmt.Errorf("db error"))

	m := &MetadataService{
		repository: repo,
	}

	_, err := m.FindByID(context.TODO(), "allwise1", "allwise")
	require.NotNil(t, err)
	repo.AssertExpectations(t)
	require.EqualError(t, err, "db error")
}
