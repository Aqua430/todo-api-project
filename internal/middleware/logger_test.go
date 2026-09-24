package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"todo-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestStructuredLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		handlerStatus int
		expectedLevel string
		expectedMsg   string
		expectedPath  string
		expectedQuery string
	}{
		{
			name:          "Status 200 - Info Level",
			handlerStatus: http.StatusOK,
			expectedLevel: "INFO",
			expectedMsg:   "HTTP Request Success",
			expectedPath:  "/test",
			expectedQuery: "param=value",
		},
		{
			name:          "Status 400 - Warn Level",
			handlerStatus: http.StatusBadRequest,
			expectedLevel: "WARN",
			expectedMsg:   "HTTP Request Client Error",
			expectedPath:  "/test",
			expectedQuery: "",
		},
		{
			name:          "Status 500 - Error Level",
			handlerStatus: http.StatusInternalServerError,
			expectedLevel: "ERROR",
			expectedMsg:   "HTTP Request Server Error",
			expectedPath:  "/test",
			expectedQuery: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))
			slog.SetDefault(logger)

			router := gin.New()
			router.Use(middleware.StructuredLoggerMiddleware())

			router.GET("/test", func(c *gin.Context) {
				c.Status(tt.handlerStatus)
			})

			url := "/test"
			if tt.expectedQuery != "" {
				url += "?" + tt.expectedQuery
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("User-Agent", "TestAgent")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.handlerStatus {
				t.Fatalf("got status = %d, want %d", w.Code, tt.handlerStatus)
			}

			logOutput := buf.String()

			if !strings.Contains(logOutput, `"level":"`+tt.expectedLevel+`"`) {
				t.Errorf("log output missing level %s. Got: %s", tt.expectedLevel, logOutput)
			}

			if !strings.Contains(logOutput, `"msg":"`+tt.expectedMsg+`"`) {
				t.Errorf("log output missing msg %s. Got: %s", tt.expectedMsg, logOutput)
			}

			if !strings.Contains(logOutput, `"path":"`+tt.expectedPath+`"`) {
				t.Errorf("log output missing path %s. Got: %s", tt.expectedPath, logOutput)
			}

			if tt.expectedQuery != "" && !strings.Contains(logOutput, `"query":"`+tt.expectedQuery+`"`) {
				t.Errorf("log output missing query %s. Got: %s", tt.expectedQuery, logOutput)
			}

			if !strings.Contains(logOutput, `"user_agent":"TestAgent"`) {
				t.Errorf("log output missing user_agent. Got: %s", logOutput)
			}
		})
	}
}
