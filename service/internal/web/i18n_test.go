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
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestLoadTransactions(t *testing.T) {
	w := &Web{}
	err := w.loadTranslations()
	if err != nil {
		t.Fatalf("loading translation failed: %v", err)
	}

	if w.defaultLocalizer == nil {
		t.Fatalf("default lozalizer did not get set")
	}

	if w.translations == nil {
		t.Fatalf("translations did not get set")
	}

	// NOTE: testing just with english for the moment
	localizer := i18n.NewLocalizer(w.translations, "en")

	msgID := "time_layout"
	str, err := localizer.LocalizeMessage(&i18n.Message{ID: msgID})
	if err != nil {
		t.Fatalf("localizing message failed %v", err)
	}

	if len(str) == 0 {
		t.Fatalf("transalted text is empty, check if %s exits in .toml", msgID)
	}
}

func createTestLocalizer(t *testing.T) *i18n.Localizer {
	t.Helper()

	bundle := i18n.NewBundle(language.English)
	messages := []*i18n.Message{
		{
			ID:    "welcome_message",
			Other: "Welcome to the test app",
		},
		{
			ID:    "error_not_found",
			Other: "Item not found",
		},
		{
			ID:    "greeting_with_message",
			Other: "Hello {{.v0}}, hello also to {{.v1}}",
		},
		{
			ID:    "message_count",
			Other: "{{.ct}} messages available",
			One:   "1 message available",
		},
		{
			ID:    "greeting_with_count",
			Other: "Hello {{.v0}}, you have {{.ct}} messages",
			One:   "Hello {{.v0}}, you have 1 message",
		},
	}
	err := bundle.AddMessages(language.English, messages...)

	if err != nil {
		t.Fatalf("failed to add messages to the bundle: %v", err)
	}

	return i18n.NewLocalizer(bundle, language.English.String())
}

func TestTranslate(t *testing.T) {
	tests := []struct {
		name    string
		args    []any
		wantErr bool
		msgID   string
	}{
		{
			name:  "No args",
			msgID: "welcome_message",
		},
		{
			name:    "Error, message does not exists",
			msgID:   "this_does_not_exist",
			wantErr: true,
		},
		{
			name:  "With args",
			msgID: "greeting_with_message",
			args:  []any{"John", "Maria"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localizer := createTestLocalizer(t)
			msg := translate(localizer, tt.msgID, tt.args...)

			if len(msg) == 0 {
				t.Fatal("translation failed, got empty message")
			}

			errMsg := fmt.Sprintf(
				"[TL err: message \"%s\" not found in language \"en\"]",
				tt.msgID,
			)
			if tt.wantErr {
				if strings.Compare(errMsg, msg) != 0 {
					t.Fatalf("expected msg err to be: %s, got: %s", errMsg, msg)
				}
				return
			}

			if strings.Contains(msg, "err") {
				t.Fatalf("translation got an error: %s", msg)
			}

			if len(tt.args) >= 0 {
				if strings.Contains(msg, "<no value>") {
					t.Fatalf("translation got an error: %s", msg)
				}
			}
		})
	}
}

func TestTranslateCount(t *testing.T) {
	tests := []struct {
		name    string
		msgID   string
		count   int
		args    []any
		wantErr bool
	}{
		{
			name:  "Singular with no args",
			msgID: "message_count",
			count: 1,
		},
		{
			name:  "Plural with no args",
			msgID: "message_count",
			count: 5,
		},
		{
			name:  "Singular with args",
			msgID: "greeting_with_count",
			count: 1,
			args:  []any{"John"},
		},
		{
			name:  "Plural with args",
			msgID: "greeting_with_count",
			count: 3,
			args:  []any{"Team"},
		},
		{
			name:    "Error, message does not exist",
			msgID:   "nonexistent_message",
			count:   1,
			wantErr: true,
		},
		{
			name:  "Zero count",
			msgID: "message_count",
			count: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localizer := createTestLocalizer(t)
			msg := translateCount(localizer, tt.msgID, tt.count, tt.args...)

			if len(msg) == 0 {
				t.Fatal("translation failed, got empty message")
			}

			errMsg := fmt.Sprintf(
				"[TL err: message \"%s\" not found in language \"en\"]",
				tt.msgID,
			)

			if tt.wantErr {
				if strings.Compare(msg, errMsg) != 0 {
					t.Fatalf("expected msg err to be: %s, got: %s", errMsg, msg)
				}
				return
			}

			if strings.Contains(msg, "err") {
				t.Fatalf("unexpected translation error: %s", msg)
			}

			if strings.Contains(msg, "<no value>") {
				t.Fatalf("missing template values in translation: %s", msg)
			}

			if !strings.Contains(msg, strconv.Itoa(tt.count)) && tt.count != 0 {
				t.Errorf("translation output should reflect count %d: %s", tt.count, msg)
			}
		})
	}
}
