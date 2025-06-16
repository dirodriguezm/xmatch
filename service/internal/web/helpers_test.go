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

// func (web *Web) render(c *gin.Context, status int, page string, data templateData) {
// ts, ok := web.templateCache[page]
// if !ok {
// err := fmt.Errorf("the template %s does not exist", page)
// web.serverError(c, err)
// return
// }

// buf := new(bytes.Buffer)
// if err := ts.ExecuteTemplate(buf, "base", data); err != nil {
// web.serverError(c, err)
// return
// }

// c.Data(status, "text/html; charset=utf-8", buf.Bytes())
// }

// func TestRender(t *testing.T) {
// tests := []struct {
// name string
// }{}
// }
func TestServerError(t *testing.T) {
	// Setup
	stdout := &strings.Builder{}
	r := SetupRouter(t, stdout)

	tc, err := newTemplateCache()
	if err != nil {
		t.Fatal("Could not create template cache")
	}

	// Create minimal Web struct with just what we need for testing
	web := &Web{
		templateCache: tc,
	}

	// Add a simple route to test serverError
	r.GET("/test", func(c *gin.Context) {
		web.serverError(c, fmt.Errorf("test error"))
	})

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify HTTP response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
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
			stdout := new(strings.Builder)
			r := SetupRouter(t, stdout)

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
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if body := w.Body.String(); body != tt.expectedBody {
				t.Errorf("expected body '%s', got '%s'", tt.expectedBody, body)
			}
		})
	}
}
