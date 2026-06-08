package models

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

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,min=1"`
}

var (
	Todos []Todo
	Users []User
)

var (
	NextID     = 1
	NextUserID = 1
)

type PatchTodoRequest struct {
	Done  *bool   `json:"done,omitempty"`
	Title *string `json:"title" binding:"omitempty,min=1"`
}
