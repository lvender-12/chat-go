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
	handler := NewHandler(service, state, logger)

	chat := app.Group("/chat")

	chat.Use(func(c fiber.Ctx) error {
		return middleware.CheckAuth(c, state.Config.JWT.Secret, *logger)
	})

	chat.Use(func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			logger.Debug("Request to open websocket channel")
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	chat.Get("/:id", websocket.New(func(ws *websocket.Conn) {
		handler.GetChat(ws)
	}))
}
