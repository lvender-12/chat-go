package test

import (
	"testing"

	"chat-go/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseToken(t *testing.T) {
	t.Parallel()

	secret := []byte("test-secret")
	userID := uint64(42)

	tokenString, err := utils.GenerateToken(userID, secret, 1)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if tokenString == "" {
		t.Fatal("token is empty")
	}

	token, err := utils.ParseToken(tokenString, secret)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if !token.Valid {
		t.Fatal("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims type assertion failed")
	}

	gotID, ok := claims["user_id"].(float64)
	if !ok {
		t.Fatal("user_id claim missing or wrong type")
	}
	if uint64(gotID) != userID {
		t.Fatalf("user_id = %v, want %d", gotID, userID)
	}
}

func TestParseToken_InvalidSecret(t *testing.T) {
	t.Parallel()

	tokenString, err := utils.GenerateToken(1, []byte("correct-secret"), 1)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = utils.ParseToken(tokenString, []byte("wrong-secret"))
	if err == nil {
		t.Fatal("ParseToken() error = nil, want error")
	}
}

func TestParseToken_Malformed(t *testing.T) {
	t.Parallel()

	_, err := utils.ParseToken("not.a.token", []byte("secret"))
	if err == nil {
		t.Fatal("ParseToken() error = nil, want error")
	}
}
