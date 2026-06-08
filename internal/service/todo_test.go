package service

import (
	"testing"
	"todo-api/internal/models"
)

func TestCreateTodo(t *testing.T) {
	tests := []struct {
		name   string
		input  models.CreateTodoRequest
		userID int
	}{
		{"normal todo", models.CreateTodoRequest{Title: "Купить хлеб"}, 1},
		{"another todo", models.CreateTodoRequest{Title: "Сделать дз"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models.Todos = []models.Todo{}
			models.Users = []models.User{{ID: tt.userID}}
			models.NextID = 1

			result := CreateTodo(tt.input, tt.userID)

			if result.ID != 1 {
				t.Errorf("expected ID 1, got %d", result.ID)
			}

			if result.Title != tt.input.Title {
				t.Errorf("expected title %s, got %s", tt.input.Title, result.Title)
			}

			if result.Done != false {
				t.Errorf("expected done=false, got %t", result.Done)
			}

			if result.UserID != tt.userID {
				t.Errorf("expected userID %d, got %d", tt.userID, result.UserID)
			}
		})
	}
}

func TestGetTodoByID(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   []models.Todo
		wantErr bool
	}{
		{"todo exists", 1, []models.Todo{{ID: 1, Title: "test"}}, false},
		{"not found", 1, []models.Todo{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models.Todos = tt.setup

			todo, err := GetTodoByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got err=%v", tt.wantErr, err)
			}

			if !tt.wantErr {
				if todo.ID != tt.setup[0].ID {
					t.Errorf("expected ID %d, got %d", tt.setup[0].ID, todo.ID)
				}

				if todo.Title != tt.setup[0].Title {
					t.Errorf("expected Title %s, got %s", tt.setup[0].Title, todo.Title)
				}
			}
		})

	}
}
