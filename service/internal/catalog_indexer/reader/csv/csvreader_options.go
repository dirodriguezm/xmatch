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

package csvreader

type CsvReaderOption[T any] func(r *CsvReader[T])

func WithHeader[T any](header []string) CsvReaderOption[T] {
	return func(r *CsvReader[T]) {
		r.Header = header
	}
}

func WithFirstLineHeader[T any](firstLineHeader bool) CsvReaderOption[T] {
	return func(r *CsvReader[T]) {
		r.FirstLineHeader = firstLineHeader
	}
}

func WithCsvBatchSize[T any](size int) CsvReaderOption[T] {
	return func(r *CsvReader[T]) {
		if size <= 0 {
			size = 1
		}
		r.batchSize = size
	}
}

func WithComment[T any](comment string) CsvReaderOption[T] {
	return func(r *CsvReader[T]) {
		for _, char := range comment {
			r.currentReader.Comment = rune(char)
			break
		}
	}
}
