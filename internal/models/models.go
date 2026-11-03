package models

import (
	"errors"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type TodoItem struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTodoInput struct {
	Title string `json:"title" binding:"required,min=1"`
}

type UpdateTodoInput struct {
	Title string `json:"title" binding:"required,min=1"`
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PaginationInput struct {
	Page  int `form:"page" binding:"gte=1"`
	Limit int `form:"limit" binding:"gte=1,lte=100"`
}

type PaginationParams struct {
	Page  int
	Limit int
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

var (
	ErrUserNotFound       = errors.New("user is not found")
	ErrTodoNotFound       = errors.New("task is not found or does not belongs to you")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInternal           = errors.New("internal server error")
)
