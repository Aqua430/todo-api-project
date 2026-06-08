package main

import (
	"todo-api/internal/handlers"
	"todo-api/internal/logger"
	"todo-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitLogger()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.StrucutredLoggerMiddleware())

	r.GET("/todos", handlers.GetTodos)
	r.GET("/todos/:id", handlers.GetTodoByID)
	r.GET("/users/:id", handlers.GetUserByID)
	r.GET("/users", handlers.GetUsers)
	r.GET("/users/:id/todos", handlers.GetTodosByUserID)
	r.POST("/users/:id/todos", handlers.PostTodo)
	r.POST("/users", handlers.PostUser)
	r.DELETE("/todos/:id", handlers.DeleteTodo)
	r.PATCH("todos/:id", handlers.PatchTodo)

	r.Run()
}
