package auth

import (
	"chat-go/internal/utils"
	"encoding/hex"
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

func (s *Service) Register(input *UserInput, ctx fiber.Ctx) error {
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing required fields")
	}

	hashedPassword, err := utils.HashingPassword(input.Password)
	if err != nil {
		return err
	}
	input.Password = hex.EncodeToString(hashedPassword)

	s.logger.Debug("Register Service", "input", input)
	if err := s.repo.CreateUser(input, ctx); err != nil {
		return err
	}

	return nil
}
