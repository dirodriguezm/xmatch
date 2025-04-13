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

type ParquetMastercat struct {
	ID   *string  `parquet:"name=id, type=BYTE_ARRAY"`
	Ipix *int64   `parquet:"name=ipix, type=INT64"`
	Ra   *float64 `parquet:"name=ra, type=DOUBLE"`
	Dec  *float64 `parquet:"name=dec, type=DOUBLE"`
	Cat  *string  `parquet:"name=cat, type=BYTE_ARRAY"`
}

func (m ParquetMastercat) ToInsertParams() any {
	return InsertObjectParams{
		ID:   *m.ID,
		Ipix: *m.Ipix,
		Ra:   *m.Ra,
		Dec:  *m.Dec,
		Cat:  *m.Cat,
	}
}
