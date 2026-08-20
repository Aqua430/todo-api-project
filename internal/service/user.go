package service

import (
	"context"
	"errors"
	"todo-api/internal/models"
	"todo-api/internal/pkg/apperrors"
	"todo-api/internal/pkg/hash"
	"todo-api/internal/pkg/jwt"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, email, passwordHash string) (int, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
}

type AuthService struct {
	repo UserRepositoryInterface
}

func NewAuthService(repo UserRepositoryInterface) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) SignUp(ctx context.Context, email, password string) (int, error) {
	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.Create(ctx, email, hashedPassword)
	if err != nil {
		if errors.Is(err, models.ErrEmailAlreadyExists) {
			return 0, apperrors.NewBadRequestError("email is already taken")
		}
		return 0, err
	}

	return id, nil
}

func (s *AuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return "", apperrors.NewUnauthorizedError("invalid email or password")
		}
		return "", err
	}

	ok := hash.CheckPasswordHash(password, user.PasswordHash)
	if !ok {
		return "", apperrors.NewUnauthorizedError("invalid email or password")
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
