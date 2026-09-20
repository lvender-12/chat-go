package friends

import (
	"chat-go/internal/app"
	"chat-go/internal/utils"
	"log/slog"
	"strconv"

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
// @Router /api/v1/friend/add-friend [post]
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

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

// GetFriends
// @Summary Get friends
// @Description Get all friends of the currently authenticated user
// @Tags friend
// @Produce json
// @Success 200 {object} FriendsResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/friend/friends [get]
func (h *Handler) GetFriends(c fiber.Ctx) error {
	h.logger.Debug("Hit Get Friends Handler")
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

	friends, err := h.service.GetFriends(userID, c)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(friends)
}

// GetFriends
// @Summary Get friends request
// @Description Get all friends of the currently authenticated user
// @Tags friend
// @Produce json
// @Success 200 {object} []FriendRequestDto
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/friend/friends-request [get]
func (h *Handler) GetRequests(c fiber.Ctx) error {
	h.logger.Debug("Hit Get Friends Handler")
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

	requests, err := h.service.GetRequests(userID, c)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(requests)
}

// RejectRequest
// @Summary Reject a friend request
// @Description Reject a friend request received by the currently authenticated user
// @Tags friend
// @Produce json
// @Param id path uint64 true "Friend request ID"
// @Success 200 {object} map[string]string "Friend request rejected successfully"
// @Failure 400 {object} map[string]string "Invalid friend request ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Friend request not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/friend/reject-friend/{id} [post]
func (h *Handler) RejectRequest(c fiber.Ctx) error {
	h.logger.Debug("hit reject request handler")

	requestID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid friend request id",
		)
	}

	userID, err := utils.GetUserIDFromToken(
		c,
		[]byte(h.state.Config.JWT.Secret),
	)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			err.Error(),
		)
	}

	h.logger.Debug(
		"rejecting friend request",
		"request_id", requestID,
		"user_id", userID,
	)

	if err := h.service.RejectRequest(requestID, userID, c); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "request rejected successfully",
	})
}
