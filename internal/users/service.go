package users

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

type Service struct {
	repo   *Repository
	logger *slog.Logger
}

func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Profile(userID uint64, ctx fiber.Ctx) (*UserProfile, error) {
	user, err := s.repo.GetUserByID(userID, ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(
				fiber.StatusNotFound,
				"user not found",
			)
		}

		return nil, fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to get user profile",
		)
	}

	return user, nil
}
