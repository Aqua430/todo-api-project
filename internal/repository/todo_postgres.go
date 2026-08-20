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

func (r *TodoRepository) Create(ctx context.Context, userID int, title string) (int, error) {
	var id int
	query := "INSERT INTO todos (user_id, title) VALUES ($1, $2) RETURNING id"

	err := r.db.QueryRow(ctx, query, userID, title).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TodoRepository) GetAll(ctx context.Context, userID int) ([]models.TodoItem, error) {
	query := "SELECT id, user_id, title, completed, created_at FROM todos WHERE user_id = $1 ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]models.TodoItem, 0)

	for rows.Next() {
		var item models.TodoItem
		err := rows.Scan(&item.ID, &item.UserID, &item.Title, &item.Completed, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) DeleteTodo(ctx context.Context, todoID, userID int) error {
	query := "DELETE FROM todos WHERE id = $1 AND user_id = $2"

	result, err := r.db.Exec(ctx, query, todoID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrTodoNotFound
	}

	return nil
}

func (r *TodoRepository) ToggleCompleted(ctx context.Context, todoID, userID int) error {
	query := "UPDATE todos SET completed = NOT completed WHERE id = $1 AND user_id = $2"

	result, err := r.db.Exec(ctx, query, todoID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrTodoNotFound
	}

	return nil
}

func (r *TodoRepository) UpdateTodoTitle(ctx context.Context, todoID, userID int, todoTitle string) error {
	query := "UPDATE todos SET title = $1 WHERE id = $2 AND user_id = $3"

	result, err := r.db.Exec(ctx, query, todoTitle, todoID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrTodoNotFound
	}

	return nil
}
