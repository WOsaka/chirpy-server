package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hashedPassword == password {
		t.Error("Hashed password should not be the same as the original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testPassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = CheckPasswordHash(hashedPassword, password)
	if err != nil {
		t.Errorf("CheckPasswordHash failed: %v", err)
	}

	// Test with an incorrect password
	err = CheckPasswordHash(hashedPassword, "wrongPassword")
	if err == nil {
		t.Error("CheckPasswordHash should have failed with an incorrect password")
	}
}

