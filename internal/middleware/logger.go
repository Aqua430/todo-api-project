package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func StrucutredLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		ctx.Next()

		slog.Info("HTTP Request",
			slog.Int("status", ctx.Writer.Status()),
			slog.String("method", ctx.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", ctx.ClientIP()),
			slog.String("user_agent", ctx.Request.UserAgent()),
		)
	}
}
