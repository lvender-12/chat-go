package middleware

import (
	"chat-go/internal/utils"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func CheckAuth(c fiber.Ctx, secret string, logger slog.Logger) error {

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

	logger.Info("authenticated", "userID", userID)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			err.Error(),
		)
	}

	return c.Next()
}
