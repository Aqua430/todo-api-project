package models

import "errors"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required,min=1"`
}

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Done   bool   `json:"done"`
	UserID int    `json:"user_id"`
}

var (
	ErrUserNotFound = errors.New("Пользователь не найден")
)

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,min=1"`
}
