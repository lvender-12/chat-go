package auth

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

func (r *Repository) CreateUser(input *UserRegister, ctx fiber.Ctx) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (
			username,
			email,
			password_hash,
			password_salt,
			display_name
		) VALUES (?, ?, ?, ?, ?)`,
		input.Username,
		input.Email,
		input.PasswordHash,
		input.PasswordSalt,
		input.DisplayName,
	)

	if err != nil {
		r.logger.Error(
			"failed to insert user",
			"username", input.Username,
			"email", input.Email,
			"error", err,
		)

		return err
	}

	r.logger.Debug(
		"user inserted into database",
		"username", input.Username,
		"email", input.Email,
	)

	return nil
}

func (r *Repository) GetUserByUsernameOrEmail(usernameOrEmail string, ctx fiber.Ctx) (*User, error) {
	r.logger.Debug(
		"querying user",
		"identifier", usernameOrEmail,
	)

	var user User

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, username, email, password_hash, display_name, password_salt FROM users WHERE username = ? OR email = ?",
		usernameOrEmail,
		usernameOrEmail,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.PasswordSalt,
	)

	if err != nil {
		r.logger.Warn(
			"user query failed",
			"identifier", usernameOrEmail,
			"error", err,
		)

		return nil, err
	}

	r.logger.Debug(
		"user found",
		"user_id", user.ID,
		"username", user.Username,
	)

	return &user, nil
}
