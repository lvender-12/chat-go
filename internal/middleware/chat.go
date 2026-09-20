package middleware

import (
	"chat-go/internal/utils"
	"database/sql"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func CheckUserIdOnConversations(c fiber.Ctx, secret string, logger slog.Logger, db *sql.DB) error {
	conversationID := c.Params("id")
	var exists bool

	logger.Debug("checking auth")

	if value := c.Cookies("AuthToken"); value == "" {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"no auth token",
		)
	}

	userID, err := utils.GetUserIDFromToken(
		c,
		[]byte(secret),
	)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			err.Error(),
		)
	}

	logger.Info("authenticated", "userID", userID)
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM conversations
			WHERE id = ?
			AND (user_one_id = ? OR user_two_id = ?)
		)
	`, conversationID, userID, userID).Scan(&exists)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			err.Error(),
		)
	}

	if !exists {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"user not found in conversation",
		)
	}

	return c.Next()
}
