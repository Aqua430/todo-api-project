package service

import (
	"context"
	"todo-api/internal/models"
	"todo-api/internal/pkg/apperrors"

	"errors"
)

type TodoRepositoryInterface interface {
	Create(ctx context.Context, userID int, title string) (int, error)
	GetAll(ctx context.Context, userID int) ([]models.TodoItem, error)
	DeleteTodo(ctx context.Context, todoID, userID int) error
	ToggleCompleted(ctx context.Context, todoID, userID int) error
	UpdateTodoTitle(ctx context.Context, todoID, userID int, todoTitle string) error
}

type TodoService struct {
	repo TodoRepositoryInterface
}

func NewTodoService(repo TodoRepositoryInterface) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(ctx context.Context, userID int, title string) (int, error) {
	return s.repo.Create(ctx, userID, title)
}

func (s *TodoService) GetAllTodos(ctx context.Context, userID int) ([]models.TodoItem, error) {
	return s.repo.GetAll(ctx, userID)
}

func (s *TodoService) DeleteTodo(ctx context.Context, todoID, userID int) error {
	err := s.repo.DeleteTodo(ctx, todoID, userID)
	if errors.Is(err, models.ErrTodoNotFound) {
		return apperrors.NewNotFoundError(models.ErrTodoNotFound.Error())
	}
	return err
}

func (s *TodoService) ToggleCompleted(ctx context.Context, todoID, userID int) error {
	err := s.repo.ToggleCompleted(ctx, todoID, userID)
	if errors.Is(err, models.ErrTodoNotFound) {
		return apperrors.NewNotFoundError(models.ErrTodoNotFound.Error())
	}
	return err
}

func (s *TodoService) UpdateTodoTitle(ctx context.Context, todoID, userID int, todoTitle string) error {
	err := s.repo.UpdateTodoTitle(ctx, todoID, userID, todoTitle)
	if errors.Is(err, models.ErrTodoNotFound) {
		return apperrors.NewNotFoundError(models.ErrTodoNotFound.Error())
	}
	return err
}
