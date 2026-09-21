package users

import (
	"database/sql"
	"log/slog"

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

func (r *Repository) GetUserByID(
	userID uint64,
	ctx fiber.Ctx,
) (*UserProfile, error) {
	r.logger.Debug(
		"querying user",
		"identifier", userID,
	)

	var user UserProfile

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, username, email, display_name FROM users WHERE id = ?",
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
	)

	if err != nil {
		r.logger.Warn(
			"user query failed",
			"identifier", userID,
			"error", err,
		)

		return nil, err
	}

	return &user, nil
}
