package jwt_test

import (
	"testing"
	"time"

	jwtpkg "todo-api/internal/pkg/jwt"

	golangjwt "github.com/golang-jwt/jwt/v5"
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

func TestJWTManager_ParseToken(t *testing.T) {
	secret := "test-secret-key"
	manager := jwtpkg.NewJWTManager(secret, 24)

	validToken, err := manager.GenerateToken(42)
	if err != nil {
		t.Fatalf("failed to generate valid token for setup: %v", err)
	}

	wrongKeyManager := jwtpkg.NewJWTManager("wrong-secret-key", 24)
	wrongKeyToken, err := wrongKeyManager.GenerateToken(42)
	if err != nil {
		t.Fatalf("failed to generate wrong key token for setup: %v", err)
	}

	expiredManager := jwtpkg.NewJWTManager(secret, -1)
	expiredToken, err := expiredManager.GenerateToken(42)
	if err != nil {
		t.Fatalf("failed to generate expired token for setup: %v", err)
	}

	noneTokenObj := golangjwt.NewWithClaims(golangjwt.SigningMethodNone, golangjwt.MapClaims{
		"user_id": 42,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	noneToken, _ := noneTokenObj.SignedString(golangjwt.UnsafeAllowNoneSignatureType)

	tests := []struct {
		name        string
		token       string
		expectedID  int
		expectError bool
	}{
		{
			name:        "Valid token",
			token:       validToken,
			expectedID:  42,
			expectError: false,
		},
		{
			name:        "Token signed with wrong key",
			token:       wrongKeyToken,
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "Expired token",
			token:       expiredToken,
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "Malformed token string",
			token:       "invalid.jwt.string",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "Empty token string",
			token:       "",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "Invalid signing method (none)",
			token:       noneToken,
			expectedID:  0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := manager.ParseToken(tt.token)

			if (err != nil) != tt.expectError {
				t.Fatalf("got error presence = %v, want expectError = %v (err: %v)", err != nil, tt.expectError, err)
			}

			if userID != tt.expectedID {
				t.Errorf("got userID = %d, want %d", userID, tt.expectedID)
			}
		})
	}
}
