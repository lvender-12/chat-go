package router

import (
	"chat-go/internal/app"
	"chat-go/internal/auth"
	"chat-go/internal/chat"
	"chat-go/internal/friends"
	"chat-go/internal/users"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func SetupRouter(
	fiberApp *fiber.App,
	state *app.State,
	logger *slog.Logger,
) {
	logger.Debug(
		"static storage path",
		"path", state.Config.Storage.Path,
	)

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     state.Config.CORS.AllowOrigins,
		AllowMethods:     state.Config.CORS.AllowMethods,
		AllowHeaders:     state.Config.CORS.AllowHeaders,
		AllowCredentials: state.Config.CORS.AllowCredentials,
	}))

	fiberApp.Use("/", static.New(state.Config.Storage.Path))

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
	chat.RouteWs(api, state, logger)
}
