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

package csv_reader

type CsvReaderOption func(r *CsvReader)

func WithHeader(header []string) CsvReaderOption {
	return func(r *CsvReader) {
		r.Header = header
	}
}

func WithFirstLineHeader(firstLineHeader bool) CsvReaderOption {
	return func(r *CsvReader) {
		r.FirstLineHeader = firstLineHeader
	}
}

func WithCsvBatchSize(size int) CsvReaderOption {
	return func(r *CsvReader) {
		if size <= 0 {
			size = 1
		}
		r.BatchSize = size
	}
}

func WithComment(comment string) CsvReaderOption {
	return func(r *CsvReader) {
		for _, char := range comment {
			r.currentReader.Comment = rune(char)
			break
		}
	}
}
