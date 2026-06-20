package router

import (
	"todo-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userHandler *handlers.UserHandler, todoHandler *handlers.TodoHandler) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/users", userHandler.CreateUser)
		v1.GET("/users", userHandler.GetUsers)
		v1.GET("/users/:id", userHandler.GetUserByID)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
		v1.PATCH("/users/:id", userHandler.UpdateUser)

		v1.POST("/users/:id/todos", todoHandler.CreateTodo)
		v1.GET("/users/:id/todos", todoHandler.GetTodosByUserID)
		v1.GET("/todos", todoHandler.GetAllTodos)
	}
}
