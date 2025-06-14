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
	"testing"
	"time"

	"github.com/dirodriguezm/xmatch/service/internal/assertions"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2023, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2023 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2023, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2023 at 09:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assertions.Equal(hd, tt.want)
		})
	}
}

func TestNewTemplateCache(t *testing.T) {
	cache, err := newTemplateCache()

	if err != nil {
		t.Fatalf("new template failed with err: %v", err)
	}

	if len(cache) == 0 {
		t.Fatalf("new template did not load any templates")
	}

	tests := []struct {
		name string
		file string
	}{
		{
			name: "error",
			file: "error.tmpl.html",
		},
		{
			name: "home",
			file: "home.tmpl.html",
		},
		{
			name: "htmxtest",
			file: "htmxtest.tmpl.html",
		},
		{
			name: "notfound",
			file: "notfound.tmpl.html",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			_, ok := cache[test.file]
			if !ok {
				t.Fatalf("the template %s does not exist", test.file)
			}
		})
	}
}
