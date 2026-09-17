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
// @Tags Auth
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

	return c.JSON(user)
}
