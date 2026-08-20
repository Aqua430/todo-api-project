package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todo-api/internal/config"
	"todo-api/internal/database"
	"todo-api/internal/handlers"
	"todo-api/internal/logger"
	"todo-api/internal/middleware"
	"todo-api/internal/repository"
	"todo-api/internal/router"
	"todo-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	logger.InitLogger()
	slog.Info("configuration successfully loaded", slog.String("env", cfg.Env))

	dbPool, err := database.NewPostgresPool(cfg.Database.URL)
	if err != nil {
		slog.Error("failed to initialize database pool", "error", err)
		os.Exit(1)
	}

	slog.Info("successfully connected to PostgreSQL")

	userRepo := repository.NewUserRepository(dbPool)
	authService := service.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	todoRepo := repository.NewTodoRepository(dbPool)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.StrucutredLoggerMiddleware())

	router.SetupRouter(r, authHandler, todoHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPServer.Port,
		Handler: r,
	}

	go func() {
		slog.Info("starting server", slog.String("port", cfg.HTTPServer.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server crashed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	dbPool.Close()
	slog.Info("server stopped gracefully")
}
