package utils_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-api/internal/middleware"
	"todo-api/internal/pkg/apperrors"
	"todo-api/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func TestMustGetID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		paramValue     string
		expectedID     int
		expectedOk     bool
		expectedStatus int
	}{
		{
			name:           "Valid positive ID",
			paramValue:     "42",
			expectedID:     42,
			expectedOk:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-numeric ID",
			paramValue:     "abc",
			expectedID:     0,
			expectedOk:     false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Zero ID",
			paramValue:     "0",
			expectedID:     0,
			expectedOk:     false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Negative ID",
			paramValue:     "-10",
			expectedID:     0,
			expectedOk:     false,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = []gin.Param{{Key: "id", Value: tt.paramValue}}

			id, ok := utils.MustGetID(c, "id")

			if ok != tt.expectedOk {
				t.Errorf("got ok = %v, want %v", ok, tt.expectedOk)
			}
			if id != tt.expectedID {
				t.Errorf("got id = %d, want %d", id, tt.expectedID)
			}
			if w.Code != tt.expectedStatus {
				t.Errorf("got status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		setupContext func(c *gin.Context)
		expectedID   int
		expectedErr  bool
	}{
		{
			name: "Valid user ID in context",
			setupContext: func(c *gin.Context) {
				c.Set(middleware.UserCtxKey, 1)
			},
			expectedID:  1,
			expectedErr: false,
		},
		{
			name:        "Missing user ID in context",
			expectedID:  0,
			expectedErr: true,
		},
		{
			name: "Invalid type in context (string instead of int)",
			setupContext: func(c *gin.Context) {
				c.Set(middleware.UserCtxKey, "1")
			},
			expectedID:  0,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setupContext != nil {
				tt.setupContext(c)
			}

			id, err := utils.GetUserID(c)

			if (err != nil) != tt.expectedErr {
				t.Errorf("got error presence = %v, want %v", err != nil, tt.expectedErr)
			}

			if id != tt.expectedID {
				t.Errorf("got id = %d, want %d", id, tt.expectedID)
			}
		})
	}
}

type sampleStruct struct {
	Email    string `validate:"required,email"`
	Password string `validate:"min=6"`
	Age      int    `validate:"max=10"`
}

func TestFormatValidationErrors(t *testing.T) {
	validate := validator.New()

	validationErr := validate.Struct(sampleStruct{
		Email:    "invalid-email",
		Password: "123",
		Age:      20,
	})

	tests := []struct {
		name     string
		err      error
		expected map[string]string
	}{
		{
			name:     "Nil error",
			err:      nil,
			expected: map[string]string{},
		},
		{
			name:     "Non-validation error",
			err:      errors.New("generic error"),
			expected: map[string]string{},
		},
		{
			name: "Validation errors formatting",
			err:  validationErr,
			expected: map[string]string{
				"email":    "is invalid",
				"password": "must be at least 6 characters long",
				"age":      "invalid value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.FormatValidationErrors(tt.err)

			if len(got) != len(tt.expected) {
				t.Fatalf("got map length = %d, want %d", len(got), len(tt.expected))
			}

			for key, expectedMsg := range tt.expected {
				gotMsg, exists := got[key]
				if !exists {
					t.Errorf("missing key %q in result map", key)
					continue
				}
				if gotMsg != expectedMsg {
					t.Errorf("for key %q got %q, want %q", key, gotMsg, expectedMsg)
				}
			}
		})
	}
}

func TestHandleError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   utils.ErrorResponse
	}{
		{
			name:           "AppError without fields",
			err:            apperrors.NewBadRequestError("bad request message"),
			expectedStatus: http.StatusBadRequest,
			expectedBody: utils.ErrorResponse{
				Error: "bad request message",
			},
		},
		{
			name: "AppError with validation fields",
			err: apperrors.NewValidationError(map[string]string{
				"email": "is invalid",
			}),
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: utils.ErrorResponse{
				Error: "validation error",
				Fields: map[string]string{
					"email": "is invalid",
				},
			},
		},
		{
			name:           "Generic error returns 500",
			err:            errors.New("unexpected database driver error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: utils.ErrorResponse{
				Error: "internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			utils.HandleError(c, tt.err)

			if w.Code != tt.expectedStatus {
				t.Errorf("got status = %d, want %d", w.Code, tt.expectedStatus)
			}

			var actualBody utils.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &actualBody); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			if actualBody.Error != tt.expectedBody.Error {
				t.Errorf("got error message = %q, want %q", actualBody.Error, tt.expectedBody.Error)
			}

			if len(actualBody.Fields) != len(tt.expectedBody.Fields) {
				t.Fatalf("got fields map len = %d, want %d", len(actualBody.Fields), len(tt.expectedBody.Fields))
			}

			for k, expectedVal := range tt.expectedBody.Fields {
				gotVal, exists := actualBody.Fields[k]
				if !exists {
					t.Errorf("missing key %q in response fields", k)
					continue
				}
				if gotVal != expectedVal {
					t.Errorf("for field key %q got %q, want %q", k, gotVal, expectedVal)
				}
			}
		})
	}
}
