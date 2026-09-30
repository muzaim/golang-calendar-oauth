package auth

import (
	"testing"
)

func TestJWTGenerationAndValidation(t *testing.T) {
	secret := "testsecretkey123"
	userID := int64(42)
	email := "user@example.com"

	token, err := GenerateAccessToken(userID, email, secret, 1)
	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}

	_, err = ValidateToken(token, "wrongsecret")
	if err == nil {
		t.Errorf("Expected validation error for invalid secret")
	}
}
