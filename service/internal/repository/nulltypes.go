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
)

// NullFloat64 embeds sql.NullFloat64 and provides custom JSON marshaling.
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON implements json.Marshaler.
// Returns null if Valid is false, otherwise the float64 value.
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON implements json.Unmarshaler.
func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var v *float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		nf.Valid = false
		return nil
	}
	nf.Valid = true
	nf.Float64 = *v
	return nil
}

// NullInt64 embeds sql.NullInt64 and provides custom JSON marshaling.
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON implements json.Marshaler.
// Returns null if Valid is false, otherwise the int64 value.
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON implements json.Unmarshaler.
func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var v *int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		ni.Valid = false
		return nil
	}
	ni.Valid = true
	ni.Int64 = *v
	return nil
}

// NullString embeds sql.NullString and provides custom JSON marshaling.
type NullString struct {
	sql.NullString
}

// MarshalJSON implements json.Marshaler.
// Returns null if Valid is false, otherwise the string value.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON implements json.Unmarshaler.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var v *string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	ns.String = *v
	return nil
}

// Ensure our types implement driver.Valuer and sql.Scanner (already via embedding).
