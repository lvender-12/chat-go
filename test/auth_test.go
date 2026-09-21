package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"chat-go/internal/utils"

	"github.com/gofiber/fiber/v3"
)

func TestGetUserIDFromToken_Success(t *testing.T) {
	t.Parallel()

	secret := []byte("auth-secret")
	wantID := uint64(99)

	token, err := utils.GenerateToken(wantID, secret, 1)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	fiberApp := fiber.New()
	var gotID uint64
	var gotErr error

	fiberApp.Get("/me", func(c fiber.Ctx) error {
		gotID, gotErr = utils.GetUserIDFromToken(c, secret)
		return gotErr
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.AddCookie(&http.Cookie{Name: "AuthToken", Value: token})

	resp, err := fiberApp.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
	if gotErr != nil {
		t.Fatalf("GetUserIDFromToken() error = %v", gotErr)
	}
	if gotID != wantID {
		t.Fatalf("userID = %d, want %d", gotID, wantID)
	}
}

func TestGetUserIDFromToken_MissingCookie(t *testing.T) {
	t.Parallel()

	fiberApp := fiber.New()
	var gotErr error

	fiberApp.Get("/me", func(c fiber.Ctx) error {
		_, gotErr = utils.GetUserIDFromToken(c, []byte("secret"))
		return nil
	})

	resp, err := fiberApp.Test(httptest.NewRequest(http.MethodGet, "/me", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if gotErr == nil {
		t.Fatal("GetUserIDFromToken() error = nil, want error")
	}
}

func TestGetUserIDFromToken_InvalidToken(t *testing.T) {
	t.Parallel()

	fiberApp := fiber.New()
	var gotErr error

	fiberApp.Get("/me", func(c fiber.Ctx) error {
		_, gotErr = utils.GetUserIDFromToken(c, []byte("secret"))
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.AddCookie(&http.Cookie{Name: "AuthToken", Value: "invalid.token.value"})

	resp, err := fiberApp.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if gotErr == nil {
		t.Fatal("GetUserIDFromToken() error = nil, want error")
	}
}
