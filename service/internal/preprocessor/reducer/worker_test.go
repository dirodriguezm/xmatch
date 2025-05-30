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

package reducer

import (
	"sync"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	partition_reader "github.com/dirodriguezm/xmatch/service/internal/preprocessor/reader"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

func testInput() partition_reader.Records {
	ra1, dec1 := 2.0, 2.0
	ra2, dec2 := 3.0, 4.0
	eRa1, eDec1 := 0.1, 0.1
	eRa2, eDec2 := 0.2, 0.2
	flux1, eFlux1 := 1.0, 0.1
	flux2, eFlux2 := 2.0, 0.2
	oid := "oid"
	input := partition_reader.Records{
		&repository.VlassInputSchema{
			Component_name: &oid,
			RA:             &ra1,
			DEC:            &dec1,
			ERA:            &eRa1,
			EDEC:           &eDec1,
			TotalFlux:      &flux1,
			ETotalFlux:     &eFlux1,
		},
		&repository.VlassInputSchema{
			Component_name: &oid,
			RA:             &ra2,
			DEC:            &dec2,
			ERA:            &eRa2,
			EDEC:           &eDec2,
			TotalFlux:      &flux2,
			ETotalFlux:     &eFlux2,
		},
	}
	return input
}

func TestWorker_GetVlassObject(t *testing.T) {
	input := testInput()

	worker := &Worker{
		schema: config.VlassSchema,
	}

	result := worker.getObject(input)

	require.Equal(t, "oid", result.GetId())
	require.Equal(t, 2.5, *result.(*repository.VlassObjectSchema).Ra)
	require.Equal(t, 3.0, *result.(*repository.VlassObjectSchema).Dec)
}

func TestWorker_TestStart(t *testing.T) {
	input := testInput()
	inCh := make(chan partition_reader.Records, 1)
	outCh := make(chan writer.WriterInput[repository.InputSchema], 1)

	worker := &Worker{
		schema: config.VlassSchema,
		inCh:   inCh,
		outCh:  outCh,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go worker.Start(&wg)
	inCh <- input
	close(inCh)
	wg.Wait()
	close(outCh)

	batch := <-outCh
	result := batch.Rows[0]
	require.Equal(t, "oid", result.GetId())
	require.Equal(t, 2.5, *result.(*repository.VlassObjectSchema).Ra)
	require.Equal(t, 3.0, *result.(*repository.VlassObjectSchema).Dec)
}
