package chat

import (
	"chat-go/internal/app"
	"chat-go/internal/middleware"
	"log/slog"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func RouteWs(app fiber.Router, state *app.State, logger *slog.Logger) {
	repo := NewRepository(state.DB, logger)
	service := NewService(repo, logger)
	hub := NewHub()
	handler := NewHandler(service, state, hub, logger)

	chat := app.Group("/chat")

	chat.Use(func(c fiber.Ctx) error {
		return middleware.CheckAuth(
			c,
			state.Config.JWT.Secret,
			*logger,
		)
	})

	chat.Get(
		"/:id",
		func(c fiber.Ctx) error {
			if !websocket.IsWebSocketUpgrade(c) {
				return fiber.ErrUpgradeRequired
			}

			logger.Debug(
				"Request to open websocket channel",
			)

			return c.Next()
		},
		websocket.New(func(ws *websocket.Conn) {
			if err := handler.GetChat(ws); err != nil {
				logger.Error(
					"websocket handler error",
					"error",
					err,
				)
			}
		}),
	)

	chat.Post("/:id", handler.SendMessage)
}
