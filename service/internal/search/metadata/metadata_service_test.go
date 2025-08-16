// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	result, err := m.FindByID(context.Background(), "allwise1", "allwise")
	require.Nil(t, err)
	repo.AssertExpectations(t)
	require.Equal(t, "allwise1", result.(repository.Allwise).ID)
}

func TestMetadata_BulkFindByID(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("BulkGetAllwise", mock.Anything, []string{"allwise1", "allwise2"}).Return([]repository.Allwise{{ID: "allwise1"}, {ID: "allwise2"}}, nil)

	m := &MetadataService{
		repository: repo,
	}

	result, err := m.BulkFindByID(context.Background(), []string{"allwise1", "allwise2"}, "allwise")
	require.Nil(t, err)
	repo.AssertExpectations(t)
	expectedIds := []string{"allwise1", "allwise2"}
	for i := 0; i < len(result.([]repository.Allwise)); i++ {
		require.Equal(t, expectedIds[i], result.([]repository.Allwise)[i].ID)
	}
}

func TestMetadata_Bulk_EmptyResult(t *testing.T) {
	repo := &mocks.Repository{}

	repo.
		On("GetAllwise", mock.Anything, "allwise1").
		Return(repository.Allwise{}, sql.ErrNoRows) // sql.ErrNoRows is returned when no rows are found

	m := &MetadataService{
		repository: repo,
	}

	_, err := m.FindByID(context.Background(), "allwise1", "allwise")
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

	_, err := m.FindByID(context.Background(), "allwise1", "allwise")
	require.NotNil(t, err)
	repo.AssertExpectations(t)
	require.EqualError(t, err, "db error")
}
