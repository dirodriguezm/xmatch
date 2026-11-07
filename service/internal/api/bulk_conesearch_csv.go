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
	"encoding/csv"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	bulkConesearchCSVMaxColumns      = 10
	bulkConesearchCSVInitialCapacity = 1024
)

var bulkConesearchCSVContentTypes = []string{"text/csv", "application/csv", "application/vnd.ms-excel"}

func isBulkConesearchCSVContentType(contentType string) bool {
	return slices.Contains(bulkConesearchCSVContentTypes, strings.ToLower(contentType))
}

func isBulkConesearchMultipartContentType(contentType string) bool {
	return strings.EqualFold(contentType, "multipart/form-data")
}

func fileError(format string, args ...any) error {
	return NewParseError("", "file", fmt.Sprintf(format, args...))
}

func (api *API) parseBulkConesearchCSVFile(c *gin.Context) ([]float64, []float64, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, nil, fileError("Could not parse multipart form: %s", err)
	}

	files := form.File["file"]
	if len(files) == 0 {
		return nil, nil, fileError("Multipart request must contain a CSV file in the file field")
	}

	file, err := files[0].Open()
	if err != nil {
		return nil, nil, fileError("Could not open CSV file: %s", err)
	}
	defer func() { _ = file.Close() }()

	return parseBulkConesearchCSV(file, api.config.BulkConesearchMaxRows)
}

func parseBulkConesearchCSV(reader io.Reader, maxRows int) ([]float64, []float64, error) {
	csvReader := csv.NewReader(reader)

	header, err := csvReader.Read()
	if err != nil {
		return nil, nil, fileError("CSV must start with a header row: %s", err)
	}
	if len(header) > bulkConesearchCSVMaxColumns {
		return nil, nil, fileError("CSV can't have more than %d columns", bulkConesearchCSVMaxColumns)
	}

	columns := make([]string, len(header))
	for i, column := range header {
		columns[i] = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(column, "\ufeff")))
	}
	raIndex, decIndex := slices.Index(columns, "ra"), slices.Index(columns, "dec")
	if raIndex == -1 || decIndex == -1 {
		return nil, nil, NewParseError("", "header", "CSV header must contain ra and dec columns")
	}

	capacity := min(maxRows, bulkConesearchCSVInitialCapacity)
	ra := make([]float64, 0, capacity)
	dec := make([]float64, 0, capacity)

	for len(ra) < maxRows {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fileError("Could not read CSV row: %s", err)
		}

		parsedRa, err := strconv.ParseFloat(strings.TrimSpace(record[raIndex]), 64)
		if err != nil {
			return nil, nil, NewParseError(record[raIndex], "RA", "Could not parse float.")
		}
		parsedDec, err := strconv.ParseFloat(strings.TrimSpace(record[decIndex]), 64)
		if err != nil {
			return nil, nil, NewParseError(record[decIndex], "Dec", "Could not parse float.")
		}

		ra = append(ra, parsedRa)
		dec = append(dec, parsedDec)
	}

	return ra, dec, nil
}
