package handlers

import (
	"todo-api/internal/models"
	"todo-api/internal/service"
	"todo-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetTodos(ctx *gin.Context) {
	todos := service.GetAllTodos()
	ctx.JSON(200, gin.H{
		"all todos": todos,
	})
}

func PostTodo(ctx *gin.Context) {
	userID, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	err := service.CheckUserID(userID)
	if err != nil {
		utils.WriteError(ctx, 404, err.Error())
		return
	}

	var todo models.CreateTodoRequest

	if !utils.MustBind(ctx, &todo) {
		return
	}

	createdTodo := service.CreateTodo(todo, userID)

	ctx.JSON(201, gin.H{
		"status":       "created",
		"created_todo": createdTodo,
	})
}

func GetTodoByID(ctx *gin.Context) {
	id, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	todoByID, err := service.GetTodoByID(id)
	if err != nil {
		utils.WriteError(ctx, 404, err.Error())
		return
	}

	ctx.JSON(200, gin.H{
		"todo": todoByID,
	})
}

func DeleteTodo(ctx *gin.Context) {
	id, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	err := service.DeleteTodoByID(id)
	if err != nil {
		utils.WriteError(ctx, 404, err.Error())
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "deleted",
		"todo_id": id,
	})
}

func PatchTodo(ctx *gin.Context) {
	id, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	todoByID, err := service.GetTodoByID(id)
	if err != nil {
		utils.WriteError(ctx, 404, err.Error())
		return
	}

	var patchReq models.PatchTodoRequest

	if !utils.MustBind(ctx, &patchReq) {
		return
	}

	if patchReq.Done == nil && patchReq.Title == nil {
		utils.WriteError(ctx, 400, "no fields to update")
		return
	}

	service.PatchTodo(todoByID, patchReq)

	ctx.JSON(200, gin.H{
		"status": "patched",
	})
}

func PostUser(ctx *gin.Context) {
	var req models.User

	if !utils.MustBind(ctx, &req) {
		return
	}

	createdUser := service.CreateUser(req)

	ctx.JSON(201, gin.H{
		"status":       "created",
		"created_user": createdUser,
	})
}

func GetUserByID(ctx *gin.Context) {
	id, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	userByID, err := service.GetUserByID(id)
	if err != nil {
		utils.WriteError(ctx, 404, err.Error())
		return
	}

	ctx.JSON(200, gin.H{
		"user": userByID,
	})
}

func GetUsers(ctx *gin.Context) {
	users := service.GetAllUsers()
	ctx.JSON(200, gin.H{
		"users": users,
	})
}

func GetTodosByUserID(ctx *gin.Context) {
	id, ok := utils.MustGetID(ctx, "id")
	if !ok {
		return
	}

	todosByUserID, err := service.GetTodosByUserID(id)
	if err != nil {
		utils.WriteError(ctx, 404, err.Error())
		return
	}

	ctx.JSON(200, gin.H{
		"todos": todosByUserID,
	})
}
