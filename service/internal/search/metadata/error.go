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

package metadata

import "fmt"

type ValidationError struct {
	Field  string
	Reason string
	Value  string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf(
		"Could not parse field %s with value %s: %s",
		e.Field,
		e.Value,
		e.Reason,
	)
}

type ArgumentError struct {
	Reason string
	Value  string
	Name   string
}

func (e ArgumentError) Error() string {
	return fmt.Sprintf(
		"Could not parse argument %s with value %s: %s",
		e.Name,
		e.Value,
		e.Reason,
	)
}
