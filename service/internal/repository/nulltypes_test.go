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
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullFloat64_MarshalJSON_Valid(t *testing.T) {
	nf := NullFloat64{sql.NullFloat64{Float64: 9.963, Valid: true}}
	data, err := json.Marshal(nf)
	require.NoError(t, err)
	assert.Equal(t, "9.963", string(data))
}

func TestNullFloat64_MarshalJSON_Null(t *testing.T) {
	nf := NullFloat64{sql.NullFloat64{Float64: 0, Valid: false}}
	data, err := json.Marshal(nf)
	require.NoError(t, err)
	assert.Equal(t, "null", string(data))
}

func TestNullFloat64_UnmarshalJSON_Valid(t *testing.T) {
	var nf NullFloat64
	err := json.Unmarshal([]byte("9.963"), &nf)
	require.NoError(t, err)
	assert.True(t, nf.Valid)
	assert.Equal(t, 9.963, nf.Float64)
}

func TestNullFloat64_UnmarshalJSON_Null(t *testing.T) {
	var nf NullFloat64
	err := json.Unmarshal([]byte("null"), &nf)
	require.NoError(t, err)
	assert.False(t, nf.Valid)
}

func TestNullInt64_MarshalJSON_Valid(t *testing.T) {
	ni := NullInt64{sql.NullInt64{Int64: 42, Valid: true}}
	data, err := json.Marshal(ni)
	require.NoError(t, err)
	assert.Equal(t, "42", string(data))
}

func TestNullInt64_MarshalJSON_Null(t *testing.T) {
	ni := NullInt64{sql.NullInt64{Int64: 0, Valid: false}}
	data, err := json.Marshal(ni)
	require.NoError(t, err)
	assert.Equal(t, "null", string(data))
}

func TestNullString_MarshalJSON_Valid(t *testing.T) {
	ns := NullString{sql.NullString{String: "test", Valid: true}}
	data, err := json.Marshal(ns)
	require.NoError(t, err)
	assert.Equal(t, "\"test\"", string(data))
}

func TestNullString_MarshalJSON_Null(t *testing.T) {
	ns := NullString{sql.NullString{String: "", Valid: false}}
	data, err := json.Marshal(ns)
	require.NoError(t, err)
	assert.Equal(t, "null", string(data))
}
