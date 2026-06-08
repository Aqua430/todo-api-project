package utils

import "github.com/gin-gonic/gin"

func BindAndValidate(ctx *gin.Context, req any) *ErrorResponse {
	if err := ctx.ShouldBindJSON(req); err != nil {
		fields := FormatValidationErrors(err)

		if len(fields) > 0 {
			res := NewValidationError(fields)
			return &res
		}

		res := NewBadRequest("invalid request body")
		return &res
	}

	return nil
}
