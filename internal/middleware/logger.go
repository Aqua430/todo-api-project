package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func StructuredLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		ctx.Next()

		status := ctx.Writer.Status()
		duration := time.Since(start)

		args := []any{
			slog.Int("status", status),
			slog.String("method", ctx.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Duration("duration", duration),
			slog.String("ip", ctx.ClientIP()),
			slog.String("user_agent", ctx.Request.UserAgent()),
		}

		switch {
		case status >= 500:
			slog.Error("HTTP Request Server Error", args...)
		case status >= 400:
			slog.Warn("HTTP Request Client Error", args...)
		default:
			slog.Info("HTTP Request Success", args...)
		}
	}
}
