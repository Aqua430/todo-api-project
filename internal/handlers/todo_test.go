package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todo-api/internal/handlers"
	"todo-api/internal/models"
	"todo-api/internal/pkg/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockTodoService struct {
	createTodoFunc      func(ctx context.Context, userID int, title string) (int, error)
	getAllTodosFunc     func(ctx context.Context, userID int) ([]models.TodoItem, error)
	deleteTodoFunc      func(ctx context.Context, todoID, userID int) error
	toggleCompletedFunc func(ctx context.Context, todoID, userID int) error
	updateTodoTitleFunc func(ctx context.Context, todoID, userID int, todoTitle string) error
}

func (m *mockTodoService) CreateTodo(ctx context.Context, userID int, title string) (int, error) {
	return m.createTodoFunc(ctx, userID, title)
}

func (m *mockTodoService) GetAllTodos(ctx context.Context, userID int) ([]models.TodoItem, error) {
	return m.getAllTodosFunc(ctx, userID)
}

func (m *mockTodoService) DeleteTodo(ctx context.Context, todoID, userID int) error {
	return m.deleteTodoFunc(ctx, todoID, userID)
}

func (m *mockTodoService) ToggleCompleted(ctx context.Context, todoID, userID int) error {
	return m.toggleCompletedFunc(ctx, todoID, userID)
}

func (m *mockTodoService) UpdateTodoTitle(ctx context.Context, todoID, userID int, todoTitle string) error {
	return m.updateTodoTitleFunc(ctx, todoID, userID, todoTitle)
}

func TestTodoHandler_CreateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		setupContext   func(c *gin.Context)
		bodyJSON       string
		mockBehavior   func(t *testing.T, m *mockTodoService)
		expectedStatus int
		expectedBody   string
	}

	tests := []testCase{
		{
			name: "Success creation",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
			},
			bodyJSON: `{"title": "Купить хлеб"}`,
			mockBehavior: func(t *testing.T, m *mockTodoService) {
				m.createTodoFunc = func(ctx context.Context, userID int, title string) (int, error) {
					assert.Equal(t, 1, userID)
					assert.Equal(t, "Купить хлеб", title)
					return 42, nil
				}
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":42}`,
		},
		{
			name:           "Unauthorized - Missing userID in context",
			setupContext:   func(c *gin.Context) {},
			bodyJSON:       `{"title": "Купить хлеб"}`,
			mockBehavior:   func(t *testing.T, m *mockTodoService) {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
		{
			name: "Validation error",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
			},
			bodyJSON:       `{"title": ""}`,
			mockBehavior:   func(t *testing.T, m *mockTodoService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error": "validation error","fields":{"title":"is required"}}`,
		},
		{
			name: "Internal service error",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
			},
			bodyJSON: `{"title": "Купить хлеб"}`,
			mockBehavior: func(t *testing.T, m *mockTodoService) {
				m.createTodoFunc = func(ctx context.Context, userID int, title string) (int, error) {
					return 0, models.ErrInternal
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "internal server error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockTodoService{}
			tc.mockBehavior(t, mockSvc)

			handler := handlers.NewTodoHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tc.setupContext != nil {
				tc.setupContext(c)
			}

			req := httptest.NewRequest("POST", "/api/v1/todos", bytes.NewBufferString(tc.bodyJSON))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.CreateTodo(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestTodoHandler_GetAllTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixedTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	type testCase struct {
		name           string
		setupContext   func(c *gin.Context)
		mockBehavior   func(t *testing.T, m *mockTodoService)
		expectedStatus int
		expectedBody   string
	}

	tests := []testCase{
		{
			name: "Success fetching todos",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
			},
			mockBehavior: func(t *testing.T, m *mockTodoService) {
				m.getAllTodosFunc = func(ctx context.Context, userID int) ([]models.TodoItem, error) {
					assert.Equal(t, 1, userID)
					return []models.TodoItem{
						{
							ID:        1,
							UserID:    1,
							Title:     "Купить хлеб",
							Completed: false,
							CreatedAt: fixedTime,
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"user_id":1,"title":"Купить хлеб","completed":false,"created_at":"2026-01-01T00:00:00Z"}]`,
		},
		{
			name:           "Context error",
			setupContext:   func(c *gin.Context) {},
			mockBehavior:   func(t *testing.T, m *mockTodoService) {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "internal server error"}`,
		},
		{
			name: "Internal error",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
			},
			mockBehavior: func(t *testing.T, m *mockTodoService) {
				m.getAllTodosFunc = func(ctx context.Context, userID int) ([]models.TodoItem, error) {
					return nil, models.ErrInternal
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "internal server error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockTodoService{}
			tc.mockBehavior(t, mockSvc)

			handler := handlers.NewTodoHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tc.setupContext != nil {
				tc.setupContext(c)
			}

			req := httptest.NewRequest("GET", "/api/v1/todos", nil)
			c.Request = req

			handler.GetAllTodos(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestTodoHandler_DeleteTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		setupContext   func(c *gin.Context)
		mockBehavior   func(t *testing.T, m *mockTodoService)
		expectedStatus int
		expectedBody   string
	}

	tests := []testCase{
		{
			name: "Success deletion",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
				c.Params = []gin.Param{{Key: "id", Value: "42"}}
			},
			mockBehavior: func(t *testing.T, m *mockTodoService) {
				m.deleteTodoFunc = func(ctx context.Context, todoID, userID int) error {
					assert.Equal(t, 1, userID)
					assert.Equal(t, 42, todoID)
					return nil
				}
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name: "Invalid ID parameter",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
				c.Params = []gin.Param{{Key: "id", Value: "invalid"}}
			},
			mockBehavior:   func(t *testing.T, m *mockTodoService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "invalid id parameter"}`,
		},
		{
			name: "Context error",
			setupContext: func(c *gin.Context) {
				c.Params = []gin.Param{{Key: "id", Value: "42"}}
			},
			mockBehavior:   func(t *testing.T, m *mockTodoService) {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "internal server error"}`,
		},
		{
			name: "Todo not found",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
				c.Params = []gin.Param{{Key: "id", Value: "42"}}
			},
			mockBehavior: func(t *testing.T, m *mockTodoService) {
				m.deleteTodoFunc = func(ctx context.Context, todoID, userID int) error {
					return apperrors.NewNotFoundError("task is not found or does not belongs to you")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error": "task is not found or does not belongs to you"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockTodoService{}
			tc.mockBehavior(t, mockSvc)

			handler := handlers.NewTodoHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tc.setupContext != nil {
				tc.setupContext(c)
			}

			req := httptest.NewRequest("DELETE", "/api/v1/todos/42", nil)
			c.Request = req

			handler.DeleteTodo(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedBody == "" {
				assert.Empty(t, w.Body.String())
			} else {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestTodoHandler_ToggleCompleted(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		setupContext   func(c *gin.Context)
		mockBehavior   func(t *testing.T, m *mockTodoService)
		expectedStatus int
		expectedBody   string
	}

	tests := []testCase{
		{
			name: "Success toggle",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
				c.Params = []gin.Param{{Key: "id", Value: "42"}}
			},
			mockBehavior: func(t *testing.T, m *mockTodoService) {
				m.toggleCompletedFunc = func(ctx context.Context, todoID, userID int) error {
					assert.Equal(t, 1, userID)
					assert.Equal(t, 42, todoID)
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name: "Todo not found",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
				c.Params = []gin.Param{{Key: "id", Value: "42"}}
			},
			mockBehavior: func(t *testing.T, m *mockTodoService) {
				m.toggleCompletedFunc = func(ctx context.Context, todoID, userID int) error {
					return apperrors.NewNotFoundError("task is not found or does not belongs to you")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error": "task is not found or does not belongs to you"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockTodoService{}
			tc.mockBehavior(t, mockSvc)

			handler := handlers.NewTodoHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tc.setupContext != nil {
				tc.setupContext(c)
			}

			req := httptest.NewRequest("PATCH", "/api/v1/todos/42", nil)
			c.Request = req

			handler.ToggleCompleted(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedBody == "" {
				assert.Empty(t, w.Body.String())
			} else {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestTodoHandler_UpdateTodoTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		setupContext   func(c *gin.Context)
		bodyJSON       string
		mockBehavior   func(t *testing.T, m *mockTodoService)
		expectedStatus int
		expectedBody   string
	}

	tests := []testCase{
		{
			name: "Success title update",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
				c.Params = []gin.Param{{Key: "id", Value: "42"}}
			},
			bodyJSON: `{"title": "Купить молоко"}`,
			mockBehavior: func(t *testing.T, m *mockTodoService) {
				m.updateTodoTitleFunc = func(ctx context.Context, todoID, userID int, title string) error {
					assert.Equal(t, 1, userID)
					assert.Equal(t, 42, todoID)
					assert.Equal(t, "Купить молоко", title)
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name: "Validation error",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
				c.Params = []gin.Param{{Key: "id", Value: "42"}}
			},
			bodyJSON:       `{"title": ""}`,
			mockBehavior:   func(t *testing.T, m *mockTodoService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error": "validation error","fields":{"title":"is required"}}`,
		},
		{
			name: "Todo not found",
			setupContext: func(c *gin.Context) {
				c.Set("userID", 1)
				c.Params = []gin.Param{{Key: "id", Value: "42"}}
			},
			bodyJSON: `{"title": "Купить хлеб"}`,
			mockBehavior: func(t *testing.T, m *mockTodoService) {
				m.updateTodoTitleFunc = func(ctx context.Context, todoID, userID int, title string) error {
					return apperrors.NewNotFoundError("task is not found or does not belongs to you")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error": "task is not found or does not belongs to you"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockTodoService{}
			tc.mockBehavior(t, mockSvc)

			handler := handlers.NewTodoHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tc.setupContext != nil {
				tc.setupContext(c)
			}

			req := httptest.NewRequest("PUT", "/api/v1/todos/42", bytes.NewBufferString(tc.bodyJSON))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.UpdateTodoTitle(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedBody == "" {
				assert.Empty(t, w.Body.String())
			} else {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}
