package hash_test

import (
	"testing"
	"todo-api/internal/pkg/hash"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		expectErr bool
	}{
		{
			name:      "Valid password",
			password:  "secret123",
			expectErr: false,
		},
		{
			name:      "Empty password",
			password:  "",
			expectErr: false,
		},
		{
			name:      "Password exceeds bcrypt 72 bytes limit",
			password:  string(make([]byte, 73)),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := hash.HashPassword(tt.password)

			if (err != nil) != tt.expectErr {
				t.Fatalf("got err = %v, want expectErr = %v", err, tt.expectErr)
			}

			if !tt.expectErr {
				if hashed == "" {
					t.Error("expected non-empty hash string")
				}
				if hashed == tt.password {
					t.Error("hash must not match plain text password")
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	validPassword := "my_secure_password"
	validHash, err := hash.HashPassword(validPassword)
	if err != nil {
		t.Fatalf("failed to generate hash for setup: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "Correct password and hash",
			password: validPassword,
			hash:     validHash,
			expected: true,
		},
		{
			name:     "Incorrect password",
			password: "wrong_password",
			hash:     validHash,
			expected: false,
		},
		{
			name:     "Malformed hash string",
			password: validPassword,
			hash:     "invalid_bcrypt_hash_format",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hash.CheckPasswordHash(tt.password, tt.hash)
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}
