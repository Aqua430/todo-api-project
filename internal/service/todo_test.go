package service_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"todo-api/internal/models"
	"todo-api/internal/pkg/apperrors"
	"todo-api/internal/service"

	"github.com/stretchr/testify/assert"
)

type mockTodoRepository struct {
	deleteFunc          func(ctx context.Context, todoID, userID int) error
	getAllFunc          func(ctx context.Context, userID int) ([]models.TodoItem, error)
	createFunc          func(ctx context.Context, userID int, title string) (int, error)
	toggleCompletedFunc func(ctx context.Context, todoID, userID int) error
	updateTodoTitleFunc func(ctx context.Context, todoID, userID int, todoTitle string) error
}

func (m *mockTodoRepository) DeleteTodo(ctx context.Context, todoID, userID int) error {
	return m.deleteFunc(ctx, todoID, userID)
}

func (m *mockTodoRepository) GetAll(ctx context.Context, userID int) ([]models.TodoItem, error) {
	return m.getAllFunc(ctx, userID)
}

func (m *mockTodoRepository) Create(ctx context.Context, userID int, title string) (int, error) {
	return m.createFunc(ctx, userID, title)
}

func (m *mockTodoRepository) ToggleCompleted(ctx context.Context, todoID, userID int) error {
	return m.toggleCompletedFunc(ctx, todoID, userID)
}

func (m *mockTodoRepository) UpdateTodoTitle(ctx context.Context, todoID, userID int, todoTitle string) error {
	return m.updateTodoTitleFunc(ctx, todoID, userID, todoTitle)
}

func TestTodoService_CreateTodo(t *testing.T) {
	type testCase struct {
		name         string
		userID       int
		title        string
		mockBehavior func(t *testing.T, m *mockTodoRepository)
		expectedID   int
		expectedErr  error
	}

	tests := []testCase{
		{
			name:   "Success: creation transparent pass-through",
			userID: 1,
			title:  "Купить хлеб",
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.createFunc = func(ctx context.Context, userID int, title string) (int, error) {
					assert.Equal(t, 1, userID)
					assert.Equal(t, "Купить хлеб", title)
					return 99, nil
				}
			},
			expectedID:  99,
			expectedErr: nil,
		},
		{
			name:   "Failure: repository error",
			userID: 1,
			title:  "Купить хлеб",
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.createFunc = func(ctx context.Context, userID int, title string) (int, error) {
					return 0, models.ErrInternal
				}
			},
			expectedID:  0,
			expectedErr: models.ErrInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockTodoRepository{}
			tc.mockBehavior(t, mockRepo)

			service := service.NewTodoService(mockRepo)
			id, err := service.CreateTodo(context.Background(), tc.userID, tc.title)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedID, id)
		})
	}
}

func TestTodoService_GetAllTodos(t *testing.T) {
	type testCase struct {
		name           string
		userID         int
		mockBehavior   func(t *testing.T, m *mockTodoRepository)
		expectedResult []models.TodoItem
		expectedErr    error
	}

	tests := []testCase{
		{
			name:   "Success: multiple todos found",
			userID: 1,
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.getAllFunc = func(ctx context.Context, userID int) ([]models.TodoItem, error) {
					return []models.TodoItem{
						{ID: 1, Title: "Купить молоко", UserID: 1},
						{ID: 2, Title: "Выучить Go", UserID: 1},
					}, nil
				}
			},
			expectedResult: []models.TodoItem{
				{ID: 1, Title: "Купить молоко", UserID: 1},
				{ID: 2, Title: "Выучить Go", UserID: 1},
			},
			expectedErr: nil,
		},
		{
			name:   "Success: no todos found (empty list)",
			userID: 1,
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.getAllFunc = func(ctx context.Context, userID int) ([]models.TodoItem, error) {
					return []models.TodoItem{}, nil
				}
			},
			expectedResult: []models.TodoItem{},
			expectedErr:    nil,
		},
		{
			name:   "Repository Fail: Internal DB Error",
			userID: 1,
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.getAllFunc = func(ctx context.Context, userID int) ([]models.TodoItem, error) {
					return nil, models.ErrInternal
				}
			},
			expectedResult: nil,
			expectedErr:    models.ErrInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockTodoRepository{}
			tc.mockBehavior(t, mockRepo)

			service := service.NewTodoService(mockRepo)
			result, err := service.GetAllTodos(context.Background(), tc.userID)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestTodoService_DeleteTodo(t *testing.T) {
	type testCase struct {
		name         string
		todoID       int
		userID       int
		mockBehavior func(t *testing.T, m *mockTodoRepository)
		checkErrFunc func(t *testing.T, err error)
	}

	tests := []testCase{
		{
			name:   "Success deletion",
			todoID: 42,
			userID: 1,
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.deleteFunc = func(ctx context.Context, todoID, userID int) error {
					return nil
				}
			},
			checkErrFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "Repository Fail: Todo Not Found transformed to AppError",
			todoID: 999,
			userID: 1,
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.deleteFunc = func(ctx context.Context, todoID, userID int) error {
					return models.ErrTodoNotFound
				}
			},
			checkErrFunc: func(t *testing.T, err error) {
				var appErr *apperrors.AppError
				assert.True(t, errors.As(err, &appErr))
				assert.Equal(t, http.StatusNotFound, appErr.StatusCode)
				assert.Equal(t, models.ErrTodoNotFound.Error(), appErr.Message)
			},
		},
		{
			name:   "Repository Fail: Generic Error pass-through",
			todoID: 42,
			userID: 1,
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.deleteFunc = func(ctx context.Context, todoID, userID int) error {
					return models.ErrInternal
				}
			},
			checkErrFunc: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, models.ErrInternal)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockTodoRepository{}
			tc.mockBehavior(t, mockRepo)

			service := service.NewTodoService(mockRepo)
			err := service.DeleteTodo(context.Background(), tc.todoID, tc.userID)

			tc.checkErrFunc(t, err)
		})
	}
}

func TestTodoService_ToggleCompleted(t *testing.T) {
	type testCase struct {
		name         string
		todoID       int
		userID       int
		mockBehavior func(t *testing.T, m *mockTodoRepository)
		checkErrFunc func(t *testing.T, err error)
	}

	tests := []testCase{
		{
			name:   "Success toggle",
			todoID: 5,
			userID: 1,
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.toggleCompletedFunc = func(ctx context.Context, todoID, userID int) error {
					assert.Equal(t, 5, todoID)
					assert.Equal(t, 1, userID)
					return nil
				}
			},
			checkErrFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "Repository Fail: Todo Not Found transformed to AppError",
			todoID: 999,
			userID: 1,
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.toggleCompletedFunc = func(ctx context.Context, todoID, userID int) error {
					return models.ErrTodoNotFound
				}
			},
			checkErrFunc: func(t *testing.T, err error) {
				var appErr *apperrors.AppError
				assert.True(t, errors.As(err, &appErr))
				assert.Equal(t, http.StatusNotFound, appErr.StatusCode)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockTodoRepository{}
			tc.mockBehavior(t, mockRepo)

			service := service.NewTodoService(mockRepo)
			err := service.ToggleCompleted(context.Background(), tc.todoID, tc.userID)

			tc.checkErrFunc(t, err)
		})
	}
}

func TestTodoService_UpdateTodoTitle(t *testing.T) {
	type testCase struct {
		name         string
		todoID       int
		userID       int
		todoTitle    string
		mockBehavior func(t *testing.T, m *mockTodoRepository)
		checkErrFunc func(t *testing.T, err error)
	}

	tests := []testCase{
		{
			name:      "Success title update",
			todoID:    10,
			userID:    2,
			todoTitle: "Купить новые кроссовки",
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.updateTodoTitleFunc = func(ctx context.Context, todoID, userID int, todoTitle string) error {
					assert.Equal(t, 10, todoID)
					assert.Equal(t, 2, userID)
					assert.Equal(t, "Купить новые кроссовки", todoTitle)
					return nil
				}
			},
			checkErrFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:      "Repository Fail: Todo Not Found transformed to AppError",
			todoID:    999,
			userID:    2,
			todoTitle: "Обновление",
			mockBehavior: func(t *testing.T, m *mockTodoRepository) {
				m.updateTodoTitleFunc = func(ctx context.Context, todoID, userID int, todoTitle string) error {
					return models.ErrTodoNotFound
				}
			},
			checkErrFunc: func(t *testing.T, err error) {
				var appErr *apperrors.AppError
				assert.True(t, errors.As(err, &appErr))
				assert.Equal(t, http.StatusNotFound, appErr.StatusCode)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockTodoRepository{}
			tc.mockBehavior(t, mockRepo)

			service := service.NewTodoService(mockRepo)
			err := service.UpdateTodoTitle(context.Background(), tc.todoID, tc.userID, tc.todoTitle)

			tc.checkErrFunc(t, err)
		})
	}
}
