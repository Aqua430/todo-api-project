package handlers

import (
	"errors"
	"net/http"
	"todo-api/internal/models"
	"todo-api/internal/service"
	"todo-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var user models.User

	if !utils.MustBind(ctx, &user) {
		return
	}

	if err := h.userService.CreateUser(ctx.Request.Context(), &user); err != nil {
		utils.WriteError(ctx, http.StatusInternalServerError, "Не удалось сохранить пользователя")
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	users, err := h.userService.GetUsers(ctx)
	if err != nil {
		utils.WriteError(ctx, http.StatusInternalServerError, "Не удалось получить список пользователей")
		return
	}

	ctx.JSON(200, gin.H{
		"users": users,
	})
}

func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	id, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	user, err := h.userService.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			utils.WriteError(ctx, http.StatusNotFound, models.ErrUserNotFound.Error())
			return
		}

		utils.WriteError(ctx, http.StatusInternalServerError, "Не удалось получить пользователя")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	id, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	err := h.userService.DeleteUser(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			utils.WriteError(ctx, http.StatusNotFound, models.ErrUserNotFound.Error())
			return
		}

		utils.WriteError(ctx, http.StatusInternalServerError, "Не удалось удалить пользователя")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Пользователь успешно удален",
	})
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	id, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	var updateReq struct {
		Name string `json:"name" binding:"required,min=1"`
	}

	ok = utils.MustBind(ctx, &updateReq)
	if !ok {
		return
	}

	updatedUser, err := h.userService.UpdateUser(ctx.Request.Context(), id, updateReq.Name)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			utils.WriteError(ctx, http.StatusNotFound, models.ErrUserNotFound.Error())
			return
		}
		utils.WriteError(ctx, http.StatusInternalServerError, "Не удалось поменять пользователя")
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}
