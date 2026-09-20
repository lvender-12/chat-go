package chat

import (
	"database/sql"
	"log/slog"
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

func (r *Repository) GetChat(id uint64) (MessageResponse, error) {
	var messages MessageResponse
	rows, err := r.db.Query("SELECT id, sender_id, content, created_at, updated_at, deleted_at FROM messages WHERE conversation_id = ?", id)
	if err != nil {
		return MessageResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Messages

		err := rows.Scan(&message.ID, &message.SenderID, &message.Content, &message.CreatedAt, &message.UpdatedAt, &message.DeletedAt)
		if err != nil {
			return MessageResponse{}, err
		}

		messages.ID = message.ID
		messages.Message = append(messages.Message, message)

	}

	if err := rows.Err(); err != nil {
		return MessageResponse{}, err
	}

	return messages, nil

}
