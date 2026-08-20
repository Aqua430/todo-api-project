package handlers

import (
	"context"
	"net/http"
	"todo-api/internal/models"
	"todo-api/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	SignUp(ctx context.Context, email, password string) (int, error)
	SignIn(ctx context.Context, email, password string) (string, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req models.RegisterInput
	if !utils.MustBind(c, &req) {
		return
	}

	id, err := h.authService.SignUp(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var req models.LoginInput
	if !utils.MustBind(c, &req) {
		return
	}

	token, err := h.authService.SignIn(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
