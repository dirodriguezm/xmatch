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
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
)

func TestValueFromDelays(t *testing.T) {
	tests := []struct {
		name string
		ctx  func() context.Context
		err  error
		want delay
	}{
		{
			name: "No Delays in CTX",
			ctx: func() context.Context {
				return context.Background()
			},
			err: fmt.Errorf("could not get value from context for key %d", Delays),
		},
		{
			name: "Good Delays in CTX",
			ctx: func() context.Context {
				ctx := context.Background()
				delays := delay{
					Slow:   -2400,
					Medium: -1200,
					Fast:   -600,
				}
				ctx = context.WithValue(ctx, Delays, delays)

				return ctx
			},
			want: delay{
				Slow:   -2400,
				Medium: -1200,
				Fast:   -600,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ValueContext[delay](tt.ctx(), Delays)

			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Errorf("expected error %v, got %v", tt.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Equal(t, val, tt.want)
		})
	}
}

func TestValueFromLocalizer(t *testing.T) {
	w := Web{}
	tests := []struct {
		name string
		ctx  func() context.Context
		err  error
		want *i18n.Localizer
	}{
		{
			name: "No Localizer in CTX",
			ctx: func() context.Context {
				return context.Background()
			},
			err: fmt.Errorf("could not get value from context for key %d", Localizer),
		},
		{
			name: "Good Localizer in CTX",
			ctx: func() context.Context {
				ctx := context.Background()
				localizer := i18n.NewLocalizer(w.translations, "en")
				ctx = context.WithValue(ctx, Localizer, localizer)

				return ctx
			},
			want: i18n.NewLocalizer(w.translations, "en"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ValueContext[*i18n.Localizer](tt.ctx(), Localizer)

			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Errorf("expected error %v, got %v", tt.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Equal(t, val, tt.want)
		})
	}
}

func TestValueFromRoute(t *testing.T) {
	tests := []struct {
		name string
		ctx  func() context.Context
		err  error
		want string
	}{
		{
			name: "No Route in CTX",
			ctx: func() context.Context {
				return context.Background()
			},
			err: fmt.Errorf("could not get value from context for key %d", Route),
		},
		{
			name: "Good Route in CTX",
			want: "home",
			ctx: func() context.Context {
				ctx := context.Background()
				route := "home"
				ctx = context.WithValue(ctx, Route, route)

				return ctx
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ValueContext[string](tt.ctx(), Route)

			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Errorf("expected error %v, got %v", tt.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Equal(t, val, tt.want)
		})
	}
}
