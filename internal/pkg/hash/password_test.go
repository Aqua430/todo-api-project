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
