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

package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToInsertParams(t *testing.T) {
	id := "id"
	cat := "cat"
	ra := 1.0
	dec := 1.0
	ipix := int64(1)
	m := ParquetMastercat{
		ID:   &id,
		Ipix: &ipix,
		Ra:   &ra,
		Dec:  &dec,
		Cat:  &cat,
	}
	result := m.ToInsertParams()
	require.Equal(t, *m.ID, result.(InsertObjectParams).ID)
	require.Equal(t, *m.Ipix, result.(InsertObjectParams).Ipix)
	require.Equal(t, *m.Ra, result.(InsertObjectParams).Ra)
	require.Equal(t, *m.Dec, result.(InsertObjectParams).Dec)
	require.Equal(t, *m.Cat, result.(InsertObjectParams).Cat)
}
