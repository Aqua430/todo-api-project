package handlers

import (
	"errors"
	"net/http"
	"todo-api/internal/models"
	"todo-api/internal/service"
	"todo-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoService *service.TodoService
}

func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

func (h *TodoHandler) CreateTodo(ctx *gin.Context) {
	userID, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	var todoReq models.CreateTodoRequest
	ok = utils.MustBind(ctx, &todoReq)
	if !ok {
		return
	}

	todo := models.Todo{
		Title:  todoReq.Title,
		UserID: userID,
		Done:   false,
	}

	if err := h.todoService.CreateTodo(ctx.Request.Context(), &todo); err != nil {
		utils.WriteError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodosByUserID(ctx *gin.Context) {
	userID, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	todos, err := h.todoService.GetTodosByUserID(ctx.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			utils.WriteError(ctx, http.StatusNotFound, models.ErrUserNotFound.Error())
			return
		}
		utils.WriteError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})
}

func (h *TodoHandler) GetAllTodos(ctx *gin.Context) {
	todos, err := h.todoService.GetAllTodos(ctx.Request.Context())
	if err != nil {
		utils.WriteError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})
}
