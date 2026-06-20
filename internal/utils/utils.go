package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIDFromParametre(ctx *gin.Context, paramName string) (int, error) {
	idStr := ctx.Param(paramName)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("invalid id")
	}

	return id, nil
}

func WriteError(ctx *gin.Context, statusCode int, err any) {
	ctx.JSON(statusCode, gin.H{
		"error": err,
	})
}

func MustGetID(ctx *gin.Context, param string) (int, bool) {
	id, err := GetIDFromParametre(ctx, param)
	if err != nil {
		WriteError(ctx, http.StatusBadRequest, err.Error())
		return 0, false
	}

	return id, true
}

func MustBind(ctx *gin.Context, obj any) bool {
	err := BindAndValidate(ctx, obj)
	if err != nil {
		WriteError(ctx, http.StatusBadRequest, err)
		return false
	}

	return true
}
