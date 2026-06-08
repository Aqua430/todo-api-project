package service

import (
	"errors"
	"todo-api/internal/models"
)

func GetAllTodos() []models.Todo {
	return models.Todos
}

func CreateTodo(todo models.CreateTodoRequest, userID int) models.Todo {
	newTodo := models.Todo{
		ID:     models.NextID,
		Title:  todo.Title,
		Done:   false,
		UserID: userID,
	}

	models.Todos = append(models.Todos, newTodo)

	models.NextID++

	return newTodo
}

func GetTodoByID(id int) (*models.Todo, error) {
	for i := range models.Todos {
		if models.Todos[i].ID == id {
			return &models.Todos[i], nil
		}
	}

	return nil, errors.New("Todo не найдено")
}

func DeleteTodoByID(id int) error {
	for i := range models.Todos {
		if models.Todos[i].ID == id {
			models.Todos = append(models.Todos[:i], models.Todos[i+1:]...)
			return nil
		}
	}

	return errors.New("Todo не найдено")
}

func PatchTodo(todoByID *models.Todo, newReq models.PatchTodoRequest) {
	if newReq.Done != nil {
		todoByID.Done = *newReq.Done
	}

	if newReq.Title != nil {
		todoByID.Title = *newReq.Title
	}
}

func GetTodosByUserID(id int) ([]models.Todo, error) {
	err := CheckUserID(id)
	if err != nil {
		return nil, err
	}

	var todosByUserID []models.Todo

	for i := range models.Todos {
		if models.Todos[i].UserID == id {
			todosByUserID = append(todosByUserID, models.Todos[i])
		}
	}

	return todosByUserID, nil
}

func CreateUser(user models.User) models.User {
	newUser := models.User{
		ID:   models.NextUserID,
		Name: user.Name,
	}

	models.Users = append(models.Users, newUser)

	models.NextUserID++

	return newUser
}

func GetUserByID(id int) (*models.User, error) {
	for i := range models.Users {
		if models.Users[i].ID == id {
			return &models.Users[i], nil
		}
	}

	return nil, errors.New("User не найден")
}

func GetAllUsers() []models.User {
	return models.Users
}

func CheckUserID(id int) error {
	for i := range models.Users {
		if models.Users[i].ID == id {
			return nil
		}
	}

	return errors.New("User не найден")
}
