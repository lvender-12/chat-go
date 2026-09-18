package users

import (
	"chat-go/internal/app"
	"chat-go/internal/utils"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
	state   *app.State
	logger  *slog.Logger
}

func NewHandler(service *Service, state *app.State, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		state:   state,
		logger:  logger,
	}
}

// Profile
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} UserProfile
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users/profile [get]
func (h *Handler) Profile(c fiber.Ctx) error {
	h.logger.Debug("profile request received")

	userID, err := utils.GetUserIDFromToken(
		c,
		[]byte(h.state.Config.JWT.Secret),
	)
	h.logger.Debug("Id", "userID", userID)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			err.Error(),
		)
	}

	user, err := h.service.Profile(userID, c)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return c.JSON(user)
}

// AddFriend
// @Summary Add a friend
// @Description Send a friend request to a user by username or email. If the target user has already sent a friend request, the request will be accepted instead.
// @Tags friend
// @Accept json
// @Produce json
// @Param request body AddFriendRequest true "Friend request"
// @Success 200 {object} map[string]string "Friend request sent or accepted"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users/add-friend [post]
func (h *Handler) AddFriend(c fiber.Ctx) error {
	var input AddFriendRequest

	if err := c.Bind().Body(&input); err != nil {
		h.logger.Warn(
			"failed to bind add friend request",
			"error", err,
		)

		return fiber.NewError(
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	h.logger.Debug(
		"add friend request received",
		"identifier", input.Name,
	)

	sender, err := utils.GetUserIDFromToken(
		c,
		[]byte(h.state.Config.JWT.Secret),
	)
	if err != nil {
		h.logger.Warn(
			"failed to get user id from token",
			"error", err,
		)

		return fiber.NewError(
			fiber.StatusUnauthorized,
			err.Error(),
		)
	}

	h.logger.Debug(
		"user identified",
		"sender_id", sender,
	)

	message, err := h.service.AddFriend(sender, input.Name, c)
	if err != nil {
		h.logger.Warn(
			"failed to add friend",
			"sender_id", sender,
			"identifier", input.Name,
			"error", err,
		)

		return err
	}

	h.logger.Debug(
		"friend operation completed",
		"sender_id", sender,
		"identifier", input.Name,
	)

	return c.Status(fiber.StatusOK).JSON(message)
}
