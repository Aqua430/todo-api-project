package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-api/internal/models"

	"github.com/gin-gonic/gin"
)

func TestGetTodosHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	models.Todos = []models.Todo{{
		ID:    1,
		Title: "Тестовое задание",
		Done:  false,
	}}

	r := gin.New()
	r.GET("/todos", GetTodos)

	req, _ := http.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, but received %d", w.Code)
	}

	var response map[string][]models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed parsing JSON of the object: %v", err)
	}

	todos := response["all todos"]
	if len(todos) != 1 || todos[0].Title != "Тестовое задание" {
		t.Errorf("Handler responded incorrect data: %v", response)
	}
}

func TestPostTodoHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	models.Todos = []models.Todo{}
	models.NextID = 1

	r := gin.New()
	r.POST("/todos", PostTodo)

	jsonBody := []byte(`{"title": "Тестовое задание"}`)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, received %d", w.Code)
	}

	var response struct {
		Status      string      `json:"status"`
		CreatedTodo models.Todo `json:"created_todo"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed parsing JSON of the object: %v", err)
	}

	if response.Status != "created" {
		t.Errorf("Expected Status 'created', but received '%s'", response.Status)
	}

	if response.CreatedTodo.ID != 1 {
		t.Errorf("Expected ID 1, but received %d", response.CreatedTodo.ID)
	}

	if response.CreatedTodo.Title != "Тестовое задание" {
		t.Errorf("Expected Title 'Тестовое задание', but received '%s'", response.CreatedTodo.Title)
	}

	if response.CreatedTodo.Done != false {
		t.Errorf("Expected Done false, but received %t", response.CreatedTodo.Done)
	}
}

func TestGetTodoByIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	models.Todos = []models.Todo{{
		ID:    1,
		Title: "Тестовое задание 1",
		Done:  false,
	}, {
		ID:    2,
		Title: "Тестовое задание 2",
		Done:  false,
	}}

	r := gin.New()

	r.GET("/todos/:id", GetTodoByID)

	req, _ := http.NewRequest("GET", "/todos/2", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	t.Log(w.Body.String())

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, but received %d", w.Code)
	}

	var response map[string]models.Todo

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed parsing JSON of the object: %v", err)
	}

	todo := response["todo"]

	if todo.ID != 2 {
		t.Errorf("Expected ID 2, but got %d", todo.ID)
	}
	if todo.Title != "Тестовое задание 2" {
		t.Errorf("Expected Title 'Тестовое задание 2', but got '%s'", todo.Title)
	}
	if todo.Done != false {
		t.Errorf("Expected Done to be false, but got %t", todo.Done)
	}
}

func TestGetTodoById_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	models.Todos = []models.Todo{}

	r := gin.New()
	r.GET("todos/:id", GetTodoByID)

	req, _ := http.NewRequest("GET", "/todos/999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	t.Log(w.Body.String())

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, but got %d", w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON of response: %v", err)
	}

	expectedError := "Todo не найдено"
	if response["error"] != expectedError {
		t.Errorf("Expected text error '%s', but got '%s'", expectedError, response["error"])
	}
}
