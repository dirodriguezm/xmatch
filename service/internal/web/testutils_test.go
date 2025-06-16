
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
	"log/slog"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func Equal[T comparable](t *testing.T, this, that T) {
	t.Helper()
	if this != that {
		t.Fatalf("expected %v to be equal to %v", this, that)
	}
}

func EqualErr(t *testing.T, this, that error) {
	t.Helper()
	if this == nil && that == nil {
		return
	}

	if (this == nil && that != nil) || (this != nil && that == nil) {
		t.Fatalf("expected %v to be equal to %v", this, that)
	}

	if that.Error() != this.Error() {
		t.Fatalf("expected %v to be equal to %v", this, that)
	}

	return
}

func NotNil(t *testing.T, object any, msgAndArgs ...any) {
	t.Helper()
	if object == nil {
		t.Fatalf("expected non nil value, go %v", object)
	}
}

func SetupRouter(t *testing.T, stdout *strings.Builder) *gin.Engine {
	t.Helper()

	loadTranslations()

	logger := slog.New(slog.NewTextHandler(
		stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	slog.SetDefault(logger)

	gin.SetMode(gin.TestMode)

	r := gin.New(func(e *gin.Engine) {
		e.Use(gin.LoggerWithWriter(stdout))
	})

	return r
}
