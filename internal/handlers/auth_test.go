package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo-api/internal/handlers"
	"todo-api/internal/models"
	"todo-api/internal/pkg/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAuthService struct {
	signUpFunc func(ctx context.Context, email, password string) (int, error)
	signInFunc func(ctx context.Context, email, password string) (string, error)
}

func (m *mockAuthService) SignUp(ctx context.Context, email, password string) (int, error) {
	return m.signUpFunc(ctx, email, password)
}

func (m *mockAuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	return m.signInFunc(ctx, email, password)
}

func TestAuthHandler_SignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		bodyJSON       string
		mockBehavior   func(t *testing.T, m *mockAuthService)
		expectedStatus int
		expectedBody   string
	}

	tests := []testCase{
		{
			name:     "Success SignUp",
			bodyJSON: `{"email": "test@example.com", "password": "password123"}`,
			mockBehavior: func(t *testing.T, m *mockAuthService) {
				m.signUpFunc = func(ctx context.Context, email, password string) (int, error) {
					assert.Equal(t, "test@example.com", email)
					assert.Equal(t, "password123", password)
					return 1, nil
				}
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id": 1}`,
		},
		{
			name:           "Validation Error",
			bodyJSON:       `{"email": "invalid-email", "password": ""}`,
			mockBehavior:   func(t *testing.T, m *mockAuthService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error": "validation error","fields":{"email":"is invalid", "password":"is required"}}`,
		},
		{
			name:     "Email already exists",
			bodyJSON: `{"email": "existing@example.com", "password": "password123"}`,
			mockBehavior: func(t *testing.T, m *mockAuthService) {
				m.signUpFunc = func(ctx context.Context, email, password string) (int, error) {
					return 0, apperrors.NewBadRequestError("email is already taken")
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "email is already taken"}`,
		},
		{
			name:     "Internal Server Error",
			bodyJSON: `{"email": "test@example.com", "password": "password123"}`,
			mockBehavior: func(t *testing.T, m *mockAuthService) {
				m.signUpFunc = func(ctx context.Context, email, password string) (int, error) {
					return 0, models.ErrInternal
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "internal server error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockAuthService{}
			tc.mockBehavior(t, mockSvc)

			handler := handlers.NewAuthHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBufferString(tc.bodyJSON))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.SignUp(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestAuthHandler_SignIn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		bodyJSON       string
		mockBehavior   func(t *testing.T, m *mockAuthService)
		expectedStatus int
		expectedBody   string
	}

	tests := []testCase{
		{
			name:     "Success SignIn",
			bodyJSON: `{"email":"test@example.com", "password":"password123"}`,
			mockBehavior: func(t *testing.T, m *mockAuthService) {
				m.signInFunc = func(ctx context.Context, email, password string) (string, error) {
					assert.Equal(t, "test@example.com", email)
					assert.Equal(t, "password123", password)
					return "valid-jwt-token", nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token": "valid-jwt-token"}`,
		},
		{
			name:           "Validation Error",
			bodyJSON:       `{"email": "invalid-email", "password": ""}`,
			mockBehavior:   func(t *testing.T, m *mockAuthService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error": "validation error","fields":{"email":"is invalid", "password":"is required"}}`,
		},
		{
			name:     "Unauthorized Error: Invalid Credentials",
			bodyJSON: `{"email": "test@example.com", "password": "wrong-password"}`,
			mockBehavior: func(t *testing.T, m *mockAuthService) {
				m.signInFunc = func(ctx context.Context, email, password string) (string, error) {
					return "", apperrors.NewUnauthorizedError("invalid email or password")
				}
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error": "invalid email or password"}`,
		},
		{
			name:     "Internal Server Error",
			bodyJSON: `{"email": "test@example.com", "password":"password123"}`,
			mockBehavior: func(t *testing.T, m *mockAuthService) {
				m.signInFunc = func(ctx context.Context, email, password string) (string, error) {
					return "", models.ErrInternal
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "internal server error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockAuthService{}
			tc.mockBehavior(t, mockSvc)

			handler := handlers.NewAuthHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest("POST", "/auth/sign-in", bytes.NewBufferString(tc.bodyJSON))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.SignIn(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}
