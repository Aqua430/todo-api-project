package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-api/internal/middleware"
	"todo-api/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret-key"
	manager := jwt.NewJWTManager(secret, 24)

	validToken, err := manager.GenerateToken(42)
	if err != nil {
		t.Fatalf("failed to generate valid token for setup: %v", err)
	}

	wrongKeyManager := jwt.NewJWTManager("wrong-secret-key", 24)
	wrongKeyToken, err := wrongKeyManager.GenerateToken(42)
	if err != nil {
		t.Fatalf("failed to generate wrong key token for setup: %v", err)
	}

	tests := []struct {
		name           string
		requestHeader  string
		expectedErrMsg string
		expectedStatus int
		expectedID     int
	}{
		{
			name:           "Missed Authorization Header",
			expectedErrMsg: `{"error":"empty auth header"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedID:     0,
		},
		{
			name:           "Incorrect Authorization Header Format: Missing Bearer",
			requestHeader:  "wrongBearer tokenExample",
			expectedErrMsg: `{"error":"invalid auth header format"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedID:     0,
		},
		{
			name:           "Incorrect Authorization Header Format: Exceeding Normal Length Of Authorization Header",
			requestHeader:  "Bearer tokenExample UnnecessaryText",
			expectedErrMsg: `{"error":"invalid auth header format"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedID:     0,
		},
		{
			name:           "Invalid Token In Header",
			requestHeader:  "Bearer " + wrongKeyToken,
			expectedErrMsg: `{"error":"invalid or expired token"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedID:     0,
		},
		{
			name:           "Valid Token In Header",
			requestHeader:  "Bearer " + validToken,
			expectedErrMsg: "",
			expectedStatus: http.StatusOK,
			expectedID:     42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotUserID int

			router := gin.New()
			router.Use(middleware.AuthMiddleware(manager))

			router.GET("/test", func(c *gin.Context) {
				if id, ok := c.Get(middleware.UserCtxKey); ok {
					gotUserID = id.(int)
				}
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.requestHeader != "" {
				req.Header.Set("Authorization", tt.requestHeader)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("got status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if w.Body.String() != tt.expectedErrMsg {
				t.Fatalf("got body = %s, want %s", w.Body.String(), tt.expectedErrMsg)
			}

			if gotUserID != tt.expectedID {
				t.Fatalf("user_id = %d, want %d", gotUserID, tt.expectedID)
			}
		})
	}
}
