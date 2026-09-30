package auth

import (
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	password := "SecretPass123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" || hash == password {
		t.Fatalf("Hash should not be empty or equal to plain password")
	}

	if !CheckPasswordHash(password, hash) {
		t.Fatalf("Password verification failed for correct password")
	}

	if CheckPasswordHash("WrongPass123!", hash) {
		t.Fatalf("Password verification succeeded for wrong password")
	}
}
