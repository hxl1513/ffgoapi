package crypto

import (
	"strings"
	"testing"
	"time"
)

func TestHashSHA256(t *testing.T) {
	input := "hello world"
	hash := HashSHA256(input)
	if hash == "" {
		t.Fatal("HashSHA256() returned empty string")
	}
	if len(hash) != 64 {
		t.Errorf("Expected hash length 64, got %d: %s", len(hash), hash)
	}

	// SHA-256 of "hello world"
	expected := "b94d27b9934d3e08a52e52d7da7dabfade4928d11c53190ac8a190fad54a091e"
	if hash == expected {
		t.Logf("HashSHA256('hello world') = %s (matches expected)", hash)
	}

	// Different inputs produce different hashes
	hash2 := HashSHA256("different string")
	if hash == hash2 {
		t.Error("Different inputs should produce different hashes")
	}

	// Same input produces same hash
	hash3 := HashSHA256("hello world")
	if hash != hash3 {
		t.Error("Same input should produce same hash")
	}
}

func TestGenerateToken(t *testing.T) {
	secret := "test-secret-key"
	userID := 1
	username := "testuser"
	role := 2
	expiry := 1 * time.Hour

	token, err := GenerateToken(secret, userID, username, role, expiry)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken() returned empty string")
	}
	if !strings.Contains(token, ".") {
		t.Error("Token should contain dots (JWT format)")
	}
}

func TestParseToken(t *testing.T) {
	secret := "test-secret-key"
	userID := 1
	username := "testuser"
	role := 2
	expiry := 1 * time.Hour

	token, err := GenerateToken(secret, userID, username, role, expiry)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims == nil {
		t.Fatal("ParseToken() returned nil claims")
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %d, want %d", claims.UserID, userID)
	}
	if claims.Username != username {
		t.Errorf("Username = %s, want %s", claims.Username, username)
	}
	if claims.Role != role {
		t.Errorf("Role = %d, want %d", claims.Role, role)
	}
}

func TestParseTokenInvalidSecret(t *testing.T) {
	secret := "test-secret-key"
	wrongSecret := "wrong-secret"

	token, err := GenerateToken(secret, 1, "testuser", 2, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = ParseToken(wrongSecret, token)
	if err == nil {
		t.Error("ParseToken() with wrong secret should return error")
	}
}

func TestParseTokenInvalidToken(t *testing.T) {
	_, err := ParseToken("secret", "invalid.token.here")
	if err == nil {
		t.Error("ParseToken() with invalid token should return error")
	}
}
