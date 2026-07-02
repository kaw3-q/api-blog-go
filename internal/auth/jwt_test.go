package auth

import (
	"github.com/kaw3-q/api-blog-go/internal/models"
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "my_secret_password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Errorf("Password verification failed for correct password")
	}

	if CheckPasswordHash("wrong_password", hash) {
		t.Errorf("Password verification succeeded for incorrect password")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	user := models.User{
		ID:       123,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     models.RoleUser,
	}

	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %d, got %d", user.ID, claims.UserID)
	}

	if claims.Role != user.Role {
		t.Errorf("Expected Role %s, got %s", user.Role, claims.Role)
	}
}
