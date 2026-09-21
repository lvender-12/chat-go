package users

import (
	"database/sql"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type Repository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewRepository(db *sql.DB, logger *slog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) GetUserByID(userID uint64, ctx fiber.Ctx) (*UserProfile, error) {
	r.logger.Debug(
		"querying user",
		"user_id", userID,
	)

	var user UserProfile

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			id,
			username,
			email,
			display_name,
			COALESCE(avatar_path, '') AS avatar_path
		FROM users
		WHERE id = ?
		`,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&user.AvatarPath,
	)

	if err != nil {
		r.logger.Warn(
			"user query failed",
			"user_id", userID,
			"error", err,
		)

		return nil, err
	}

	return &user, nil
}

func (r *Repository) EditUser(userID uint64, user UserEdit, ctx fiber.Ctx) (*UserProfile, error) {
	r.logger.Debug(
		"updating user",
		"user_id", userID,
	)

	updates := make([]string, 0, 3)
	args := make([]any, 0, 4)

	if user.Username != nil {
		updates = append(updates, "username = ?")
		args = append(args, *user.Username)
	}

	if user.Email != nil {
		updates = append(updates, "email = ?")
		args = append(args, *user.Email)
	}

	if user.DisplayName != nil {
		updates = append(updates, "display_name = ?")
		args = append(args, *user.DisplayName)
	}

	if len(updates) == 0 {
		return r.GetUserByID(userID, ctx)
	}

	query := `
		UPDATE users
		SET ` + strings.Join(updates, ", ") + `
		WHERE id = ?
	`

	args = append(args, userID)

	_, err := r.db.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		r.logger.Warn(
			"user update failed",
			"user_id", userID,
			"error", err,
		)

		return nil, err
	}

	return r.GetUserByID(userID, ctx)
}

func (r *Repository) UpdateAvatar(userID uint64, avatarPath string, ctx fiber.Ctx) error {
	r.logger.Debug(
		"updating user avatar",
		"user_id", userID,
		"avatar_path", avatarPath,
	)

	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE users
		SET avatar_path = ?
		WHERE id = ?
		`,
		avatarPath,
		userID,
	)

	if err != nil {
		r.logger.Warn(
			"user avatar update failed",
			"user_id", userID,
			"error", err,
		)

		return err
	}

	return nil
}
