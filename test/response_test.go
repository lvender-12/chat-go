package test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"chat-go/internal/app"

	"github.com/gofiber/fiber/v3"
)

func TestJSON_WithData(t *testing.T) {
	t.Parallel()

	fiberApp := fiber.New()
	fiberApp.Get("/ok", func(c fiber.Ctx) error {
		return app.JSON(c, fiber.StatusOK, "success", map[string]string{"id": "1"})
	})

	resp, err := fiberApp.Test(httptest.NewRequest(http.MethodGet, "/ok", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var body app.Response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Status != fiber.StatusOK {
		t.Errorf("body.Status = %d, want %d", body.Status, fiber.StatusOK)
	}
	if body.Message != "success" {
		t.Errorf("body.Message = %q, want %q", body.Message, "success")
	}
	if body.Data == nil {
		t.Fatal("body.Data is nil, want data")
	}
}

func TestJSON_WithoutData(t *testing.T) {
	t.Parallel()

	fiberApp := fiber.New()
	fiberApp.Get("/empty", func(c fiber.Ctx) error {
		return app.JSON(c, fiber.StatusCreated, "created", nil)
	})

	resp, err := fiberApp.Test(httptest.NewRequest(http.MethodGet, "/empty", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusCreated)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var body app.Response
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Status != fiber.StatusCreated {
		t.Errorf("body.Status = %d, want %d", body.Status, fiber.StatusCreated)
	}
	if body.Message != "created" {
		t.Errorf("body.Message = %q, want %q", body.Message, "created")
	}
	if body.Data != nil {
		t.Errorf("body.Data = %#v, want nil", body.Data)
	}
}

func TestErrorHandler_FiberError(t *testing.T) {
	t.Parallel()

	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: app.ErrorHandler,
	})
	fiberApp.Get("/bad", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	})

	resp, err := fiberApp.Test(httptest.NewRequest(http.MethodGet, "/bad", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}

	var body app.Response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Status != fiber.StatusBadRequest {
		t.Errorf("body.Status = %d, want %d", body.Status, fiber.StatusBadRequest)
	}
	if body.Message != "invalid input" {
		t.Errorf("body.Message = %q, want %q", body.Message, "invalid input")
	}
	if body.Data != nil {
		t.Errorf("body.Data = %#v, want nil", body.Data)
	}
}

func TestErrorHandler_GenericError(t *testing.T) {
	t.Parallel()

	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: app.ErrorHandler,
	})
	fiberApp.Get("/fail", func(c fiber.Ctx) error {
		return errors.New("boom")
	})

	resp, err := fiberApp.Test(httptest.NewRequest(http.MethodGet, "/fail", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusInternalServerError)
	}

	var body app.Response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Status != fiber.StatusInternalServerError {
		t.Errorf("body.Status = %d, want %d", body.Status, fiber.StatusInternalServerError)
	}
	if body.Message != "boom" {
		t.Errorf("body.Message = %q, want %q", body.Message, "boom")
	}
	if body.Data != nil {
		t.Errorf("body.Data = %#v, want nil", body.Data)
	}
}
