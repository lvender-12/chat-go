package auth

import (
	"chat-go/internal/utils"
	"log/slog"
	"strings"

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

func (s *Service) Register(input *UserRegister, ctx fiber.Ctx) error {
	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(input.Email)
	input.DisplayName = strings.TrimSpace(input.DisplayName)

	if input.Username == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"username is required",
		)
	}

	if input.Email == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"email is required",
		)
	}

	if input.Password == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"password is required",
		)
	}

	if input.DisplayName == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"display name is required",
		)
	}

	s.logger.Debug(
		"registering user",
		"username", input.Username,
		"email", input.Email,
	)

	hash, salt, err := utils.HashingPassword(input.Password)
	if err != nil {
		s.logger.Error(
			"failed to hash user password",
			"username", input.Username,
			"error", err,
		)

		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to hash password",
		)
	}

	input.PasswordHash = hash
	input.PasswordSalt = salt
	input.Password = ""

	s.logger.Debug(
		"user password hashed",
		"username", input.Username,
	)

	if err := s.repo.CreateUser(input, ctx); err != nil {
		s.logger.Error(
			"failed to create user",
			"username", input.Username,
			"email", input.Email,
			"error", err,
		)

		return err
	}

	s.logger.Debug(
		"user created successfully",
		"username", input.Username,
		"email", input.Email,
	)

	return nil
}

func (s *Service) Login(input *UserLogin, ctx fiber.Ctx) (*User, error) {
	input.Name = strings.TrimSpace(input.Name)

	if input.Name == "" {
		return nil, fiber.NewError(
			fiber.StatusBadRequest,
			"username or email is required",
		)
	}

	if input.Password == "" {
		return nil, fiber.NewError(
			fiber.StatusBadRequest,
			"password is required",
		)
	}

	s.logger.Debug(
		"authenticating user",
		"name", input.Name,
	)

	user, err := s.repo.GetUserByUsernameOrEmail(
		input.Name,
		ctx,
	)
	if err != nil {
		s.logger.Warn(
			"user lookup failed",
			"name", input.Name,
			"error", err,
		)

		return nil, fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid username/email or password",
		)
	}

	if !utils.VerifyPassword(
		input.Password,
		user.PasswordHash,
		user.PasswordSalt,
	) {
		s.logger.Warn(
			"password verification failed",
			"user_id", user.ID,
			"username", user.Username,
		)

		return nil, fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid username/email or password",
		)
	}

	s.logger.Debug(
		"user authenticated successfully",
		"user_id", user.ID,
		"username", user.Username,
	)

	return user, nil
}
