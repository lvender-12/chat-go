package auth

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

// Register
// @Summary Register user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body UserRegister true "User Register"
// @Success 200 {string} string
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c fiber.Ctx) error {
	var input UserRegister
	if err := c.Bind().Body(&input); err != nil {
		return err
	}

	h.logger.Debug(
		"register request received",
		"username", input.Username,
		"email", input.Email,
	)

	if !utils.IsValidEmail(input.Email) {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"email is invalid",
		)
	}

	if err := h.service.Register(&input, c); err != nil {
		h.logger.Warn(
			"user registration failed",
			"username", input.Username,
			"email", input.Email,
			"error", err,
		)

		return err
	}

	h.logger.Info(
		"user registered successfully",
		"username", input.Username,
		"email", input.Email,
	)

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	c.Status(fiber.StatusCreated)
	return c.JSON(fiber.Map{
		"message": "Registration handled successfully",
	})
}

// Login
// @Summary Login user
// @Description Authenticate a user and create an authentication cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body UserLogin true "Login credentials"
// @Success 200 {string} string "Login successful"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Invalid username or password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c fiber.Ctx) error {
	var input UserLogin

	h.logger.Debug("login request received")

	if err := c.Bind().Body(&input); err != nil {
		h.logger.Warn(
			"failed to parse login request",
			"error", err,
		)

		return err
	}

	h.logger.Debug(
		"login credentials received",
		"name", input.Name,
	)

	user, err := h.service.Login(&input, c)
	if err != nil {
		h.logger.Warn(
			"user login failed",
			"name", input.Name,
			"error", err,
		)

		return err
	}

	h.logger.Debug(
		"user credentials verified",
		"user_id", user.ID,
		"username", user.Username,
	)

	token, err := utils.GenerateToken(
		user.ID,
		[]byte(h.state.Config.JWT.Secret),
		h.state.Config.JWT.Exp,
	)
	if err != nil {
		h.logger.Error(
			"failed to generate authentication token",
			"user_id", user.ID,
			"error", err,
		)

		return err
	}

	h.logger.Debug(
		"authentication token generated",
		"user_id", user.ID,
	)

	c.Cookie(&fiber.Cookie{
		Name:     "AuthToken",
		Value:    token,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	h.logger.Info(
		"user logged in successfully",
		"user_id", user.ID,
		"username", user.Username,
	)

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "Auth handled successfully",
	})
}
