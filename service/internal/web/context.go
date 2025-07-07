// Copyright 2024-2025 Matías Medina Silva
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
package web

import (
	"context"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Key int

const (
	Localizer Key = iota
	TraceID
	Route
)

func localizerFrom(ctx context.Context) (*i18n.Localizer, error) {
	loc, ok := ctx.Value(Localizer).(*i18n.Localizer)
	if !ok {
		return nil, fmt.Errorf("could not get localizer from context")
	}
	return loc, nil
}

func routeFrom(ctx context.Context) (string, error) {
	selected, ok := ctx.Value(Route).(string)
	if !ok {
		return "", fmt.Errorf("could not get route from context")
	}
	return selected, nil
}
