package service_test

import (
	"context"
	"testing"
	"todo-api/internal/models"
	"todo-api/internal/pkg/apperrors"
	"todo-api/internal/pkg/hash"
	"todo-api/internal/service"

	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	createFunc     func(ctx context.Context, email, passwordHash string) (int, error)
	getByEmailFunc func(ctx context.Context, email string) (models.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, email, passwordHash string) (int, error) {
	return m.createFunc(ctx, email, passwordHash)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	return m.getByEmailFunc(ctx, email)
}

func TestAuthService_SignUp(t *testing.T) {
	type testCase struct {
		name          string
		email         string
		password      string
		mockBehavior  func(m *mockUserRepository)
		expectedID    int
		expectedError error
	}

	tests := []testCase{
		{
			name:     "Success SignUp",
			email:    "test@example.com",
			password: "password123",
			mockBehavior: func(m *mockUserRepository) {
				m.createFunc = func(ctx context.Context, email, passwordHash string) (int, error) {
					assert.Equal(t, "test@example.com", email)
					assert.True(t, hash.CheckPasswordHash("password123", passwordHash))
					return 1, nil
				}
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:     "Email already exists",
			email:    "existing@example.com",
			password: "password123",
			mockBehavior: func(m *mockUserRepository) {
				m.createFunc = func(ctx context.Context, email, passwordHash string) (int, error) {
					return 0, models.ErrEmailAlreadyExists
				}
			},
			expectedID:    0,
			expectedError: apperrors.NewBadRequestError("email is already taken"),
		},
		{
			name:     "Internal repository error",
			email:    "test@example.com",
			password: "password123",
			mockBehavior: func(m *mockUserRepository) {
				m.createFunc = func(ctx context.Context, email, passwordHash string) (int, error) {
					return 0, models.ErrInternal
				}
			},
			expectedID:    0,
			expectedError: models.ErrInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tc.mockBehavior(mockRepo)

			authService := service.NewAuthService(mockRepo)

			id, err := authService.SignUp(context.Background(), tc.email, tc.password)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedID, id)
		})
	}
}

func TestAuthService_SignIn(t *testing.T) {
	type testCase struct {
		name         string
		email        string
		password     string
		mockBehavior func(m *mockUserRepository)
		expectedErr  error
		expectToken  bool
	}

	rawPassword := "password123"
	hashedPassword, err := hash.HashPassword(rawPassword)
	assert.NoError(t, err)

	tests := []testCase{
		{
			name:     "Success SignIn",
			email:    "test@example.com",
			password: rawPassword,
			mockBehavior: func(m *mockUserRepository) {
				m.getByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					assert.Equal(t, "test@example.com", email)
					return models.User{
						ID:           1,
						Email:        email,
						PasswordHash: hashedPassword,
					}, nil
				}
			},
			expectedErr: nil,
			expectToken: true,
		},
		{
			name:     "User not found",
			email:    "notfound@example.com",
			password: rawPassword,
			mockBehavior: func(m *mockUserRepository) {
				m.getByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					return models.User{}, models.ErrUserNotFound
				}
			},
			expectedErr: apperrors.NewUnauthorizedError("invalid email or password"),
			expectToken: false,
		},
		{
			name:     "Invalid password",
			email:    "test@example.com",
			password: "wrong_password",
			mockBehavior: func(m *mockUserRepository) {
				m.getByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					return models.User{
						ID:           1,
						Email:        "test@example.com",
						PasswordHash: hashedPassword,
					}, nil
				}
			},
			expectedErr: apperrors.NewUnauthorizedError("invalid email or password"),
			expectToken: false,
		},
		{
			name:     "Internal repository error",
			email:    "test@example.com",
			password: "password123",
			mockBehavior: func(m *mockUserRepository) {
				m.getByEmailFunc = func(ctx context.Context, email string) (models.User, error) {
					return models.User{}, models.ErrInternal
				}
			},
			expectedErr: models.ErrInternal,
			expectToken: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tc.mockBehavior(mockRepo)

			authService := service.NewAuthService(mockRepo)

			token, err := authService.SignIn(context.Background(), tc.email, tc.password)

			assert.Equal(t, tc.expectedErr, err)
			if tc.expectToken {
				assert.NotEmpty(t, token)
			} else {
				assert.Empty(t, token)
			}
		})
	}
}
