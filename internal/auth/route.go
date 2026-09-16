package auth

import (
	"chat-go/internal/app"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func RouteAuth(app fiber.Router, state *app.State, logger *slog.Logger) {
	repo := NewRepository(state.DB, logger)
	service := NewService(repo, logger)
	handler := NewHandler(service, logger)

	auth := app.Group("/auth")
	auth.Get("/", handler.HandleAuth)
}
