package chat

import (
	"log/slog"
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
