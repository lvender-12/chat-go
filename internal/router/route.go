package router

import (
	"chat-go/internal/app"
	"chat-go/internal/auth"
	"chat-go/internal/friends"
	"chat-go/internal/users"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func SetupRouter(fiberApp *fiber.App, state *app.State, logger *slog.Logger) {
	api := fiberApp.Group("/api/v1")

	api.Use(func(c fiber.Ctx) error {
		logger.Info(
			"request",
			"method", c.Method(),
			"path", c.Path(),
		)

		return c.Next()
	})

	auth.RouteAuth(api, state, logger)
	users.RouteUsers(api, state, logger)
	friends.RouteFriends(api, state, logger)
}
