package friends

import (
	"chat-go/internal/app"
	"chat-go/internal/middleware"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func RouteFriends(app fiber.Router, state *app.State, logger *slog.Logger) {
	repo := NewRepository(state.DB, logger)
	service := NewService(repo, logger)
	handler := NewHandler(service, state, logger)

	friends := app.Group("/friend")
	friends.Use(func(c fiber.Ctx) error {
		return middleware.CheckAuth(c, state.Config.JWT.Secret, *logger)
	})

	friends.Post("/add-friend", handler.AddFriend)
	friends.Get("/friends", handler.GetFriends)
}
