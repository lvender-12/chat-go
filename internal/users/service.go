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
	user, err := s.repo.GetUserByID(
		userID,
		ctx,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(
				fiber.StatusNotFound,
				"user not found",
			)
		}

		s.logger.Error(
			"failed to get user profile",
			"user_id", userID,
			"error", err,
		)

		return nil, fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to get user profile",
		)
	}

	return user, nil
}

func (s *Service) EditUser(userID uint64, userEdit UserEdit, ctx fiber.Ctx) (*UserProfile, error) {
	user, err := s.repo.EditUser(
		userID,
		userEdit,
		ctx,
	)
	if err != nil {
		s.logger.Error(
			"failed to edit user",
			"user_id", userID,
			"error", err,
		)

		return nil, fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to update user profile",
		)
	}

	return user, nil
}

func (s *Service) UpdateAvatar(userID uint64, avatarPath string, ctx fiber.Ctx) error {
	if err := s.repo.UpdateAvatar(
		userID,
		avatarPath,
		ctx,
	); err != nil {
		s.logger.Error(
			"failed to update avatar",
			"user_id", userID,
			"error", err,
		)

		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to update avatar",
		)
	}

	return nil
}
