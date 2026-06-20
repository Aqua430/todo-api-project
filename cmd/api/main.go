package main

import (
	"fmt"
	"log"
	"os"
	"todo-api/internal/database"
	"todo-api/internal/handlers"
	"todo-api/internal/logger"
	"todo-api/internal/middleware"
	"todo-api/internal/repository"
	"todo-api/internal/router"
	"todo-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Переменная DATABASE_URL не найдена в окружении")
	}

	fmt.Println("[INIT] Конфигурация успешно загружена")

	dbPool, err := database.NewPostgresPool(connStr)
	if err != nil {
		log.Fatalf("[FATAL] Ошибка инициализации базы данных: %v\n", err)
	}
	defer dbPool.Close()

	fmt.Println("[INIT] Успешное подключение к PostgreSQL")

	userRepo := repository.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	todoRepo := repository.NewTodoRepository(dbPool)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	logger.InitLogger()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.StrucutredLoggerMiddleware())

	router.SetupRouter(r, userHandler, todoHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
