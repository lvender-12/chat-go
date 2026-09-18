package friends

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

func (s *Service) AddFriend(sender uint64, name string, ctx fiber.Ctx) (string, error) {
	s.logger.Debug(
		"processing add friend",
		"sender_id", sender,
		"identifier", name,
	)

	user, err := s.repo.GetUserByUsernameOrEmail(name, ctx)
	if err != nil {
		s.logger.Warn(
			"user lookup failed",
			"name", name,
			"error", err,
		)

		return "", fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid username/email or password",
		)
	}

	receiver := user.ID

	s.logger.Debug(
		"friend target found",
		"sender_id", sender,
		"receiver_id", receiver,
	)

	if sender == receiver {
		s.logger.Warn(
			"user attempted to add themselves",
			"user_id", sender,
		)

		return "", fiber.NewError(
			fiber.StatusBadRequest,
			"cannot add yourself",
		)
	}

	requestID, err := s.repo.CheckFriendRequest(sender, receiver, ctx)

	if err == nil {
		return "friend request accepted", s.repo.AcceptFriendRequest(requestID, ctx)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	return "friend request sent", s.repo.AddFriend(sender, receiver, ctx)
}

func (s *Service) GetFriends(userID uint64, ctx fiber.Ctx) (FriendsResponse, error) {
	s.logger.Debug("Hit Get Friends Handler")
	friends, err := s.repo.GetFriends(userID, ctx)
	if err != nil {
		return FriendsResponse{}, fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to get friends",
		)
	}

	return friends, nil
}
