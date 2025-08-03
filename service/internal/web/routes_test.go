package web_test

// import (
// "strings"
// "testing"

// "github.com/dirodriguezm/xmatch/service/internal/web"
// )

// func TestHome(t *testing.T) {
// // t.Parallel()

// // r, stdout := web.SetupTestRouter(t)

// // // Create minimal Web struct with just what we need for testing
// // w := &web.Web{
// // templateCache: tc,
// // }
// // w.loadTranslations()

// // var ctx context.Context
// // r.Use(extractRoute())
// // r.GET(tt.route, func(c *gin.Context) {
// // ctx = c.Request.Context()
// // c.String(200, "ok")
// // })
// // w := httptest.NewRecorder()
// // req, _ := http.NewRequest("GET", tt.route, nil)

// // r.ServeHTTP(w, req)

// // route, ok := ctx.Value(Route).(string)
// // if !ok {
// // t.Fatal("could not get route from context")
// // }
// // Equal(t, route, tt.expected)
// // r, stdout := SetupTestRouter(t)

// // tc, err := newTemplateCache()
// // if err != nil {
// // t.Fatal("Could not create template cache")
// // }

// // // Add a simple route to test serverError
// // r.GET("/test", func(c *gin.Context) {
// // w.serverError(c, fmt.Errorf("test error"))
// // })

// // // Create test request
// // req := httptest.NewRequest("GET", "/test", nil)
// // rec := httptest.NewRecorder()
// // r.ServeHTTP(rec, req)

// // // Verify HTTP response
// // if rec.Code != http.StatusInternalServerError {
// // t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
// // }

// // // Verify error was logged
// // logOutput := stdout.String()
// // if !strings.Contains(logOutput, "test error") {
// // t.Errorf("expected error log to contain 'test error', got: %s", logOutput)
// // }
// // if !strings.Contains(logOutput, "method=GET") {
// // t.Error("expected log to contain request method")
// // }

// }
