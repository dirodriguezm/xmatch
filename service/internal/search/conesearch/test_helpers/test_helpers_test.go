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

package test_helpers_test

import (
	"database/sql"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch/test_helpers"
	"testing"
)

func TestInsertAllwiseMastercat(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		nobjects int
		db       *sql.DB
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := test_helpers.InsertAllwiseMastercat(tt.nobjects, tt.db)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InsertAllwiseMastercat() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InsertAllwiseMastercat() succeeded unexpectedly")
			}
		})
	}
}
