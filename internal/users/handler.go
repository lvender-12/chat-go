package users

import (
	"chat-go/internal/app"
	"chat-go/internal/utils"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
	state   *app.State
	logger  *slog.Logger
}

func NewHandler(
	service *Service,
	state *app.State,
	logger *slog.Logger,
) *Handler {
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
// @Success 200 {object} app.Response
// @Failure 401 {object} app.Response
// @Failure 404 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/users/profile [get]
func (h *Handler) Profile(c fiber.Ctx) error {
	h.logger.Debug("profile request received")

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
		"user authenticated",
		"user_id", userID,
	)

	user, err := h.service.Profile(
		userID,
		c,
	)
	if err != nil {
		return err
	}

	return app.JSON(
		c,
		fiber.StatusOK,
		"profile retrieved successfully",
		user,
	)
}

// EditProfile
// @Summary Edit current user profile
// @Description Update the profile of the currently authenticated user. Only non-empty fields will be updated.
// @Tags users
// @Accept json
// @Produce json
// @Param request body UserEdit true "User profile data"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 401 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/users/profile [patch]
func (h *Handler) EditProfile(c fiber.Ctx) error {
	var input UserEdit

	if err := c.Bind().Body(&input); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid request body",
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

	user, err := h.service.EditUser(
		userID,
		input,
		c,
	)
	if err != nil {
		return err
	}

	return app.JSON(
		c,
		fiber.StatusOK,
		"profile updated successfully",
		user,
	)
}

// UploadProfile
// @Summary Upload profile avatar
// @Description Upload or replace the current user's profile avatar. Maximum file size is 5 MB.
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Profile avatar"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 401 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/users/profile/avatar [post]
func (h *Handler) UploadProfile(c fiber.Ctx) error {
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

	file, err := c.FormFile("avatar")
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"avatar is required",
		)
	}

	const maxFileSize = 5 * 1024 * 1024

	if file.Size > maxFileSize {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"avatar must be smaller than 5 MB",
		)
	}

	allowedExtensions := map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
		".webp": {},
		".gif":  {},
		".bmp":  {},
		".tiff": {},
		".tif":  {},
		".avif": {},
		".heic": {},
		".heif": {},
	}

	ext := strings.ToLower(
		filepath.Ext(file.Filename),
	)

	if _, ok := allowedExtensions[ext]; !ok {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"avatar file is not allowed",
		)
	}

	profileDir := filepath.Join(
		h.state.Config.Storage.Path,
		"profile",
	)

	if err := os.MkdirAll(
		profileDir,
		0755,
	); err != nil {
		h.logger.Error(
			"failed to create avatar directory",
			"path", profileDir,
			"error", err,
		)

		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to create avatar directory",
		)
	}

	filename := fmt.Sprintf(
		"%d%s",
		userID,
		ext,
	)

	path := filepath.Join(
		profileDir,
		filename,
	)

	// Remove previous avatar files using any supported extension.
	for oldExt := range allowedExtensions {
		oldPath := filepath.Join(
			profileDir,
			fmt.Sprintf("%d%s", userID, oldExt),
		)

		if oldPath == path {
			continue
		}

		if err := os.Remove(oldPath); err != nil &&
			!os.IsNotExist(err) {
			h.logger.Warn(
				"failed to remove old avatar",
				"path", oldPath,
				"error", err,
			)
		}
	}

	if err := c.SaveFile(
		file,
		path,
	); err != nil {
		h.logger.Error(
			"failed to save avatar",
			"path", path,
			"error", err,
		)

		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to save avatar",
		)
	}

	avatarPath, err := filepath.Rel(
		h.state.Config.Storage.Path,
		path,
	)
	if err != nil {
		os.Remove(path)

		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to build avatar path",
		)
	}

	avatarPath = "/" + filepath.ToSlash(
		avatarPath,
	)

	if err := h.service.UpdateAvatar(
		userID,
		avatarPath,
		c,
	); err != nil {
		os.Remove(path)

		return err
	}

	user, err := h.service.Profile(
		userID,
		c,
	)
	if err != nil {
		return err
	}

	return app.JSON(
		c,
		fiber.StatusOK,
		"avatar uploaded successfully",
		user,
	)
}
