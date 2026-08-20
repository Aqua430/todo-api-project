package handlers

import (
	"context"
	"net/http"
	"todo-api/internal/models"
	"todo-api/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TodoService interface {
	CreateTodo(ctx context.Context, userID int, title string) (int, error)
	GetAllTodos(ctx context.Context, userID int) ([]models.TodoItem, error)
	DeleteTodo(ctx context.Context, todoID, userID int) error
	ToggleCompleted(ctx context.Context, todoID, userID int) error
	UpdateTodoTitle(ctx context.Context, todoID, userID int, todoTitle string) error
}

type TodoHandler struct {
	todoService TodoService
}

func NewTodoHandler(todoService TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var input models.CreateTodoInput
	if !utils.MustBind(c, &input) {
		return
	}

	id, err := h.todoService.CreateTodo(c.Request.Context(), userID, input.Title)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	todos, err := h.todoService.GetAllTodos(c.Request.Context(), userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	todoID, ok := utils.MustGetID(c, "id")
	if !ok {
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err := h.todoService.DeleteTodo(c.Request.Context(), todoID, userID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}

func (h *TodoHandler) ToggleCompleted(c *gin.Context) {
	todoID, ok := utils.MustGetID(c, "id")
	if !ok {
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err := h.todoService.ToggleCompleted(c.Request.Context(), todoID, userID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
	c.Writer.WriteHeaderNow()
}

func (h *TodoHandler) UpdateTodoTitle(c *gin.Context) {
	todoID, ok := utils.MustGetID(c, "id")
	if !ok {
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	var input models.UpdateTodoInput
	if !utils.MustBind(c, &input) {
		return
	}

	if err := h.todoService.UpdateTodoTitle(c.Request.Context(), todoID, userID, input.Title); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
	c.Writer.WriteHeaderNow()
}
