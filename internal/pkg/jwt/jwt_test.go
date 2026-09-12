package jwt_test

import (
	"testing"

	jwtpkg "todo-api/internal/pkg/jwt"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := jwtpkg.NewJWTManager("test-secret-key", 24)

	token, err := manager.GenerateToken(100)
	if err != nil {
		t.Fatalf("unexpected error during token generation: %v", err)
	}

	if token == "" {
		t.Fatal("expected non-empty token string")
	}
}
