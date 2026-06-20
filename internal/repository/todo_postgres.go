package repository

import (
	"context"
	"todo-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, todo *models.Todo) error {
	query := "INSERT INTO todos (user_id, title, done) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRow(ctx, query, todo.UserID, todo.Title, todo.Done).Scan(&todo.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TodoRepository) GetTodosByUserID(ctx context.Context, userID int) ([]models.Todo, error) {
	query := "SELECT id, title, done, user_id FROM todos WHERE user_id = $1"

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo

	for rows.Next() {
		var todo models.Todo

		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Done, &todo.UserID); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	query := "SELECT id, title, done, user_id FROM todos"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo

		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Done, &todo.UserID); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
