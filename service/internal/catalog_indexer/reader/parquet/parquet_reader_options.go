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

package parquet_reader

type ParquetReaderOption[T any] func(r *ParquetReader[T])

func WithParquetBatchSize[T any](size int) ParquetReaderOption[T] {
	return func(r *ParquetReader[T]) {
		if size <= 0 {
			size = 1
		}
		r.BatchSize = size
	}
}
