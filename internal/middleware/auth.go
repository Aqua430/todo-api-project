package middleware

import (
	"net/http"
	"strings"
	"todo-api/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const (
	userCtxKey = "userID"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abortWithStatusUnauthorized(c, "empty auth header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			abortWithStatusUnauthorized(c, "invalid auth header format")
			return
		}

		tokenString := parts[1]

		userID, err := jwt.ParseToken(tokenString)
		if err != nil {
			abortWithStatusUnauthorized(c, err.Error())
			return
		}

		c.Set(userCtxKey, userID)

		c.Next()
	}
}

func abortWithStatusUnauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msg})
}
