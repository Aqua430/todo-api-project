package repository

import (
	"context"
	"todo-api/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
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

func (r *TodoRepository) GetAll(ctx context.Context, userID, limit, offset int) ([]models.TodoItem, int, error) {
	query := `WITH total AS (
    	SELECT COUNT(*) AS amount
    	FROM todos
    	WHERE user_id = $1
	),
	page AS (
    	SELECT id, user_id, title, completed, created_at
    	FROM todos
    	WHERE user_id = $1
    	ORDER BY created_at DESC, id DESC
    	LIMIT $2 OFFSET $3
	) SELECT
    	page.id,
    	page.user_id,
    	page.title,
    	page.completed,
    	page.created_at,
    	total.amount
	FROM total
	LEFT JOIN page ON TRUE;`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	todos := make([]models.TodoItem, 0)

	var total int64
	for rows.Next() {
		var (
			id        pgtype.Int4
			userID    pgtype.Int4
			title     pgtype.Text
			completed pgtype.Bool
			createdAt pgtype.Timestamptz
		)
		err := rows.Scan(&id, &userID, &title, &completed, &createdAt, &total)
		if err != nil {
			return nil, 0, err
		}

		if id.Valid {
			item := models.TodoItem{
				ID:        int(id.Int32),
				UserID:    int(userID.Int32),
				Title:     title.String,
				Completed: completed.Bool,
				CreatedAt: createdAt.Time,
			}
			todos = append(todos, item)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return todos, int(total), nil
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
