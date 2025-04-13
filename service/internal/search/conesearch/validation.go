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

package conesearch

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type ValidationError struct {
	Reason   string
	ErrValue string
	Field    string
}

func (e ValidationError) Error() string {
	err := "Invalid field %s with value %s. %w"
	err = fmt.Sprintf(err, e.Field, e.ErrValue, e.Reason)
	return err
}

func NewValidationError(e string, errValue string, field string) error {
	return ValidationError{
		Reason:   e,
		ErrValue: errValue,
		Field:    field,
	}
}

func ValidateRadius(rad float64) error {
	if rad <= 0 {
		err := "Radius can't be lower than 0"
		return NewValidationError(err, strconv.FormatFloat(rad, 'f', 3, 64), "radius")
	}
	return nil
}

func ValidateRa(ra float64) error {
	err := ValidationError{
		ErrValue: strconv.FormatFloat(ra, 'f', 3, 64),
		Field:    "RA",
	}
	if ra < 0 {
		err.Reason = "RA can't be lower than 0"
		return err
	}
	if ra > 360 {
		err.Reason = "RA can't be greater than 360"
		return err
	}
	return nil
}

func ValidateDec(dec float64) error {
	err := ValidationError{
		ErrValue: strconv.FormatFloat(dec, 'f', 3, 64),
		Field:    "Dec",
	}
	if dec < -90 {
		err.Reason = "Dec can't be lower than -90"
		return err
	}
	if dec > 90 {
		err.Reason = "Dec can't be greater than 90"
		return err
	}
	return nil
}

func ValidateNneighbor(nneighbor int) error {
	if nneighbor <= 0 {
		return NewValidationError("Nneighbor must be a positive integer", strconv.FormatInt(int64(nneighbor), 10), "nneighbor")
	}
	return nil
}

func ValidateCatalog(catalog string) error {
	available := []string{"vlass", "ztf", "allwise", "all"}
	cat := strings.ToLower(catalog)
	if !slices.Contains(available, cat) {
		err := NewValidationError("Catalog not available", catalog, "catalog")
		return err
	}
	return nil
}

func ValidateArguments(ra, dec, radius float64, nneighbor int, catalog string) error {
	if err := ValidateRa(ra); err != nil {
		return err
	}
	if err := ValidateDec(dec); err != nil {
		return err
	}
	if err := ValidateRadius(radius); err != nil {
		return err
	}
	if err := ValidateNneighbor(nneighbor); err != nil {
		return err
	}
	if err := ValidateCatalog(catalog); err != nil {
		return err
	}
	return nil
}

func ValidateBulkArguments(ra, dec []float64, radius float64, nneighbor int, catalog string) error {
	if len(ra) != len(dec) {
		return NewValidationError("Ra and Dec must have the same length", fmt.Sprintf("%d", len(ra)), "ra")
	}
	for i := 0; i < len(ra); i++ {
		if err := ValidateRa(ra[i]); err != nil {
			return err
		}
		if err := ValidateDec(dec[i]); err != nil {
			return err
		}
	}
	if err := ValidateRadius(radius); err != nil {
		return err
	}
	if err := ValidateNneighbor(nneighbor); err != nil {
		return err
	}
	if err := ValidateCatalog(catalog); err != nil {
		return err
	}
	return nil
}
