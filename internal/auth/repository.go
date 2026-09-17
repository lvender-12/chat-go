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

func (r *Repository) CreateUser(input *UserInput, ctx fiber.Ctx) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (username, email, password_hash, display_name) VALUES (?, ?, ?, ?)",
		input.Username, input.Email, input.Password, input.DisplayName)
	r.logger.Debug("Register Repository", "input", input)
	return err
}
