package service

import (
	"context"
	"todo-api/internal/models"
	"todo-api/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	if id <= 0 {
		return nil, models.ErrUserNotFound
	}

	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return models.ErrUserNotFound
	}

	return s.userRepo.DeleteUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id int, newName string) (*models.User, error) {
	if id <= 0 {
		return nil, models.ErrUserNotFound
	}

	user := &models.User{
		ID:   id,
		Name: newName,
	}

	err := s.userRepo.UpdateUserByID(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
