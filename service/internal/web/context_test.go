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
)

func TestLocalizerFrom(t *testing.T) {
	tests := []struct {
		name    string
		ctx     func() context.Context
		wantErr error
	}{
		{
			name: "No Localizer in CTX",
			ctx: func() context.Context {
				return context.Background()
			},
			wantErr: fmt.Errorf("could not get localizer from context"),
		},
		{
			name: "Good Localizer in CTX",
			ctx: func() context.Context {
				ctx := context.Background()
				localizer := i18n.NewLocalizer(translations, "en")
				ctx = context.WithValue(ctx, Localizer, localizer)

				return ctx
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := localizerFrom(tt.ctx())

			EqualErr(t, err, tt.wantErr)
			if tt.wantErr != nil {
				NotNil(t, l)
			}
		})
	}
}
