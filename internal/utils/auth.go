package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromToken(c fiber.Ctx, secret []byte) (uint64, error) {
	cookie := c.Cookies("AuthToken")

	if cookie == "" {
		return 0, fmt.Errorf("auth cookie is required")
	}

	token, err := ParseToken(cookie, secret)
	if err != nil {
		return 0, fmt.Errorf("invalid auth token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user id")
	}

	return uint64(userID), nil
}
