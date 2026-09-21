package chat

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

func (s *Service) GetChat(id uint64) (MessageResponse, error) {
	messages, err := s.repo.GetChat(id)
	if err != nil {
		return MessageResponse{}, err
	}

	return messages, nil
}

func (s *Service) SendMessage(idUser uint64, idConversation uint64, content string) error {
	if content == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"content is required",
		)
	}

	if err := s.repo.SendMessage(
		idUser,
		idConversation,
		content,
	); err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to send message",
		)
	}

	return nil
}

func (s *Service) EditMessage(idUser uint64, idMessage uint64, content string, ctx fiber.Ctx) error {
	if content == "" {
		return fiber.NewError(fiber.StatusBadRequest,
			"Content is requered")
	}

	if err := s.repo.EditMessage(idUser, idMessage, content, ctx); err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to send message",
		)
	}

	return nil
}

func (s *Service) DeleteMessage(idUser uint64, idMessage uint64, ctx fiber.Ctx) error {
	if err := s.repo.DeleteMessage(idUser, idMessage, ctx); err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to send message",
		)
	}

	return nil
}
