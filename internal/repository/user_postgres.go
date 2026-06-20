package repository

import (
	"context"
	"errors"
	"todo-api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (name) VALUES ($1) RETURNING id"

	err := r.db.QueryRow(ctx, query, user.Name).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := "SELECT id, name FROM users"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User

		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := "SELECT id, name FROM users WHERE id = $1"

	var u models.User

	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) DeleteUserByID(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUserByID(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET name = $1 WHERE id = $2 RETURNING id, name"

	err := r.db.QueryRow(ctx, query, user.Name, user.ID).Scan(&user.ID, &user.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrUserNotFound
		}
		return err
	}

	return nil
}
