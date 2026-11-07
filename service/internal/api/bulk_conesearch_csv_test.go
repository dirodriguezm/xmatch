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

package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBulkConesearchCSV(t *testing.T) {
	t.Run("reads ra and dec columns case-insensitively with annotations", func(t *testing.T) {
		ra, dec, err := parseBulkConesearchCSV(
			strings.NewReader("RA, name, Dec\n1.5, first, -2.5\n3, second, 4\n"),
			10,
		)

		require.NoError(t, err)
		require.Equal(t, []float64{1.5, 3}, ra)
		require.Equal(t, []float64{-2.5, 4}, dec)
	})

	t.Run("stops reading at the row limit", func(t *testing.T) {
		ra, dec, err := parseBulkConesearchCSV(
			strings.NewReader("ra,dec\n1,2\n3,4\n5,6\n"),
			2,
		)

		require.NoError(t, err)
		require.Equal(t, []float64{1, 3}, ra)
		require.Equal(t, []float64{2, 4}, dec)
	})

	t.Run("rejects a CSV without a header", func(t *testing.T) {
		_, _, err := parseBulkConesearchCSV(strings.NewReader(""), 10)

		require.Error(t, err)
	})

	t.Run("rejects a header without ra and dec", func(t *testing.T) {
		_, _, err := parseBulkConesearchCSV(strings.NewReader("ra,name\n1,first\n"), 10)

		require.Error(t, err)
	})

	t.Run("rejects more than ten columns", func(t *testing.T) {
		header := "ra,dec,c1,c2,c3,c4,c5,c6,c7,c8,c9"
		_, _, err := parseBulkConesearchCSV(strings.NewReader(header+"\n"), 10)

		require.Error(t, err)
	})

	t.Run("rejects rows with more than ten columns", func(t *testing.T) {
		csv := "ra,dec,c1,c2,c3,c4,c5,c6,c7,c8\n" +
			"1,2,3,4,5,6,7,8,9,10,11\n"
		_, _, err := parseBulkConesearchCSV(strings.NewReader(csv), 10)

		require.Error(t, err)
	})

	t.Run("rejects rows that do not match the header", func(t *testing.T) {
		_, _, err := parseBulkConesearchCSV(strings.NewReader("ra,dec,name\n1,2\n"), 10)

		require.Error(t, err)
	})

	t.Run("rejects invalid coordinates", func(t *testing.T) {
		_, _, err := parseBulkConesearchCSV(strings.NewReader("ra,dec\nnot-a-number,2\n"), 10)

		require.Error(t, err)
	})
}
