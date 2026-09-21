package test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"chat-go/internal/app"
	"chat-go/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func TestCheckAuth_MissingToken(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: app.ErrorHandler,
	})

	fiberApp.Use(func(c fiber.Ctx) error {
		return middleware.CheckAuth(c, "secret", *logger)
	})
	fiberApp.Get("/secure", func(c fiber.Ctx) error {
		return app.JSON(c, fiber.StatusOK, "ok", nil)
	})

	resp, err := fiberApp.Test(httptest.NewRequest(http.MethodGet, "/secure", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusUnauthorized)
	}

	var body app.Response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Message != "no auth token" {
		t.Fatalf("message = %q, want %q", body.Message, "no auth token")
	}
}

func TestCheckAuth_ValidToken(t *testing.T) {
	t.Parallel()

	secret := "auth-secret"
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: app.ErrorHandler,
	})

	fiberApp.Use(func(c fiber.Ctx) error {
		return middleware.CheckAuth(c, secret, *logger)
	})
	fiberApp.Get("/secure", func(c fiber.Ctx) error {
		return app.JSON(c, fiber.StatusOK, "ok", nil)
	})

	token, err := mustToken(t, 7, []byte(secret))
	if err != nil {
		t.Fatalf("token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.AddCookie(&http.Cookie{Name: "AuthToken", Value: token})

	resp, err := fiberApp.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}
