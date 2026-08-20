package utils

import (
	"errors"
	"net/http"
	"strconv"
	"todo-api/internal/pkg/apperrors"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func HandleError(ctx *gin.Context, err error) {
	var appErr *apperrors.AppError

	if errors.As(err, &appErr) {
		ctx.JSON(appErr.StatusCode, ErrorResponse{
			Error:  appErr.Message,
			Fields: appErr.Fields,
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "internal server error",
	})
}

func MustGetID(ctx *gin.Context, param string) (int, bool) {
	idStr := ctx.Param(param)

	id, err := strconv.Atoi(idStr)
	if err != nil && id <= 0 {
		HandleError(ctx, apperrors.NewBadRequestError("invalid id parameter"))
		return 0, false
	}

	return id, true
}

func MustBind(ctx *gin.Context, obj any) bool {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		fields := FormatValidationErrors(err)
		if len(fields) > 0 {
			HandleError(ctx, apperrors.NewValidationError(fields))
			return false
		}
		HandleError(ctx, apperrors.NewBadRequestError("invalid request body"))
		return false
	}
	return true
}

func GetUserID(c *gin.Context) (int, error) {
	id, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user id is not found in context")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
