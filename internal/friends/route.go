package friends

import (
	"chat-go/internal/app"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func RouteFriends(app fiber.Router, state *app.State, logger *slog.Logger) {
	repo := NewRepository(state.DB, logger)
	service := NewService(repo, logger)
	handler := NewHandler(service, state, logger)

	friends := app.Group("/friend")
	friends.Post("/add-friend", handler.AddFriend)
	friends.Get("/friends", handler.GetFriends)
}
