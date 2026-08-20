package router

import (
	"todo-api/internal/handlers"
	"todo-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, authHandler *handlers.AuthHandler, todoHandler *handlers.TodoHandler) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/sign-up", authHandler.SignUp)
		authGroup.POST("/sign-in", authHandler.SignIn)
	}

	v1Group := r.Group("/api/v1", middleware.AuthMiddleware())
	{
		v1Group.POST("/todos", todoHandler.CreateTodo)
		v1Group.GET("/todos", todoHandler.GetAllTodos)
		v1Group.DELETE("/todos/:id", todoHandler.DeleteTodo)
		v1Group.PATCH("/todos/:id", todoHandler.ToggleCompleted)
		v1Group.PUT("/todos/:id", todoHandler.UpdateTodoTitle)
	}
}
