package auth

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Register
// @Summary Register user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body UserInput true "Register input"
// @Success 200 {string} string
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c fiber.Ctx) error {
	var input UserInput
	if err := c.Bind().Body(&input); err != nil {
		return err
	}

	if err := h.service.Register(&input, c); err != nil {
		return err
	}

	c.Set("Content-Type", "application/json")
	h.logger.Debug("Register Handler", "input", input)
	c.Status(fiber.StatusCreated)
	return c.SendString("Auth handled successfully")
}
