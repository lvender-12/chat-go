package test

import (
	"testing"

	"chat-go/internal/utils"
)

func TestHashingAndVerifyPassword(t *testing.T) {
	t.Parallel()

	password := "secret-password"

	hash, salt, err := utils.HashingPassword(password)
	if err != nil {
		t.Fatalf("HashingPassword() error = %v", err)
	}
	if hash == "" {
		t.Fatal("hash is empty")
	}
	if salt == "" {
		t.Fatal("salt is empty")
	}

	if !utils.VerifyPassword(password, hash, salt) {
		t.Fatal("VerifyPassword() = false for correct password")
	}

	if utils.VerifyPassword("wrong-password", hash, salt) {
		t.Fatal("VerifyPassword() = true for wrong password")
	}
}

func TestHashingPassword_UniqueSalt(t *testing.T) {
	t.Parallel()

	_, salt1, err := utils.HashingPassword("same")
	if err != nil {
		t.Fatalf("HashingPassword() error = %v", err)
	}

	_, salt2, err := utils.HashingPassword("same")
	if err != nil {
		t.Fatalf("HashingPassword() error = %v", err)
	}

	if salt1 == salt2 {
		t.Fatal("expected different salts for each hash")
	}
}

func TestVerifyPassword_InvalidHex(t *testing.T) {
	t.Parallel()

	if utils.VerifyPassword("password", "not-hex", "00") {
		t.Fatal("VerifyPassword() = true for invalid hash hex")
	}

	if utils.VerifyPassword("password", "00", "not-hex") {
		t.Fatal("VerifyPassword() = true for invalid salt hex")
	}
}
