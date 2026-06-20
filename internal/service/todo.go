package service

import (
	"context"
	"todo-api/internal/models"
	"todo-api/internal/repository"
)

type TodoService struct {
	todoRepo *repository.TodoRepository
}

func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

func (s *TodoService) CreateTodo(ctx context.Context, todo *models.Todo) error {
	if todo.UserID <= 0 {
		return models.ErrUserNotFound
	}

	return s.todoRepo.Create(ctx, todo)
}

func (s *TodoService) GetTodosByUserID(ctx context.Context, userID int) ([]models.Todo, error) {
	if userID <= 0 {
		return nil, models.ErrUserNotFound
	}

	return s.todoRepo.GetTodosByUserID(ctx, userID)
}

func (s *TodoService) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	return s.todoRepo.GetAllTodos(ctx)
}
