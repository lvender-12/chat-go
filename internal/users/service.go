package users

import (
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
		return nil, err
	}
	return user, nil
}
