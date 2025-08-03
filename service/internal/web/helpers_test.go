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
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestServerError(t *testing.T) {
	r, stdout := SetupTestRouter(t)

	tc, err := newTemplateCache()
	if err != nil {
		t.Fatal("Could not create template cache")
	}

	// Create minimal Web struct with just what we need for testing
	w := &Web{
		templateCache: tc,
	}
	w.loadTranslations()

	// Add a simple route to test serverError
	r.GET("/test", func(c *gin.Context) {
		w.serverError(c, fmt.Errorf("test error"))
	})

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Verify HTTP response
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	// Verify error was logged
	logOutput := stdout.String()
	if !strings.Contains(logOutput, "test error") {
		t.Errorf("expected error log to contain 'test error', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "method=GET") {
		t.Error("expected log to contain request method")
	}
}

func TestRender(t *testing.T) {
	tests := []struct {
		name           string
		templateName   string
		templateCache  map[string]*template.Template
		expectedStatus int
		expectedBody   string
		wantErr        string
	}{
		{
			name:         "successful render",
			templateName: "test.tmpl.html",
			templateCache: map[string]*template.Template{
				"test.tmpl.html": template.Must(
					template.New("test.tmpl.html").Parse(`{{define "base"}}Test{{end}}`),
				),
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Test",
			wantErr:        "",
		},
		{
			name:           "template not found",
			templateName:   "missing.tmpl.html",
			templateCache:  map[string]*template.Template{},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
			wantErr:        "the template missing.tmpl.html does not exist",
		},
		{
			name:         "template execution error",
			templateName: "error.tmpl.html",
			templateCache: map[string]*template.Template{
				"error.tmpl.html": template.Must(
					template.New("error.tmpl.html").Parse(
						`{{define "base"}}{{.NonExistentField}}{{end}}`,
					),
				),
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
			wantErr:        "can't evaluate field NonExistentField",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := SetupTestRouter(t)

			web := &Web{
				templateCache: tt.templateCache,
			}

			r.GET("/test", func(c *gin.Context) {
				err := web.render(c, tt.expectedStatus, tt.templateName, templateData{})
				if len(tt.wantErr) != 0 {
					if err == nil {
						t.Fatal("expected render to fail")
					}

					if !strings.Contains(err.Error(), tt.wantErr) {
						t.Fatalf(
							"expected log to contain '%s', got: %s",
							tt.wantErr,
							err.Error(),
						)
					}
					return
				}

				if err != nil {
					t.Fatal("expected render not to fail")
				}
			})

			if len(tt.wantErr) != 0 {
				return
			}

			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if body := rec.Body.String(); body != tt.expectedBody {
				t.Errorf("expected body '%s', got '%s'", tt.expectedBody, body)
			}
		})
	}
}
