package users

import (
	"chat-go/internal/app"
	"chat-go/internal/middleware"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func RouteUsers(api fiber.Router, state *app.State, logger *slog.Logger) {
	repo := NewRepository(state.DB, logger)
	service := NewService(repo, logger)
	handler := NewHandler(service, state, logger)

	users := api.Group("/users")

	users.Use(func(c fiber.Ctx) error {
		return middleware.CheckAuth(
			c,
			state.Config.JWT.Secret,
			*logger,
		)
	})

	users.Get("/profile", handler.Profile)

	users.Patch("/profile", handler.EditProfile)

	users.Post("/profile/avatar", handler.UploadProfile)
}
