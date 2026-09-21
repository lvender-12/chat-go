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

	if c.Cookies("AuthToken") == "" {
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

	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM conversations
			WHERE id = ?
			  AND (user_one_id = ? OR user_two_id = ?)
		)
	`, conversationID, userID, userID).Scan(&exists)

	if err != nil {
		logger.Error(
			"failed to check conversation membership",
			"error", err,
			"conversationID", conversationID,
			"userID", userID,
		)

		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to check conversation access",
		)
	}

	logger.Info(
		"conversation membership result",
		"conversationID", conversationID,
		"userID", userID,
		"exists", exists,
	)

	if !exists {
		logger.Warn(
			"user is not a member of conversation",
			"conversationID", conversationID,
			"userID", userID,
		)

		return fiber.NewError(
			fiber.StatusForbidden,
			"user is not a member of this conversation",
		)
	}

	return c.Next()
}
