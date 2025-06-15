package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
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


func TestMakeJWTAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	expiresIn := 1 * time.Hour

	// Test MakeJWT
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	if token == "" {
		t.Fatal("MakeJWT returned an empty token")
	}

	// Test ValidateJWT with correct secret
	parsedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if parsedUserID != userID {
		t.Errorf("ValidateJWT returned wrong userID: got %v, want %v", parsedUserID, userID)
	}

	// Test ValidateJWT with wrong secret
	_, err = ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Error("ValidateJWT should fail with wrong secret")
	}
}

func TestValidateJWTExpired(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	expiresIn := -1 * time.Second // already expired

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("ValidateJWT should fail for expired token")
	}
}

func TestJWT_ValidToken(t *testing.T) {
	userID := uuid.New()
	secret := "supersecret"
	expiresIn := 10 * time.Minute

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	parsedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if parsedUserID != userID {
		t.Errorf("ValidateJWT returned wrong userID: got %v, want %v", parsedUserID, userID)
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "supersecret"
	expiresIn := -1 * time.Second // already expired

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("ValidateJWT should fail for expired token")
	}
}

func TestJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "supersecret"
	wrongSecret := "nottherightsecret"
	expiresIn := 10 * time.Minute

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("ValidateJWT should fail with wrong secret")
	}
}

func TestJWT_TamperedToken(t *testing.T) {
	userID := uuid.New()
	secret := "supersecret"
	expiresIn := 10 * time.Minute

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Tamper with the token by changing a character
	tampered := token[:len(token)-1] + "x"

	_, err = ValidateJWT(tampered, secret)
	if err == nil {
		t.Error("ValidateJWT should fail for tampered token")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer testtoken123")
	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken failed: %v", err)
	}
	if token != "testtoken123" {
		t.Errorf("GetBearerToken returned wrong token: got %s, want %s", token, "testtoken123")
	}
	// Test with missing Authorization header
	headers = http.Header{}
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Error("GetBearerToken should fail when Authorization header is missing")
	}
}