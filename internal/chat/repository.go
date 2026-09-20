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
	messages := MessageResponse{
		ConversationID: id,
		Messages:       make([]Messages, 0),
	}

	rows, err := r.db.Query(`
		SELECT
			id,
			sender_id,
			content,
			message_type,
			created_at,
			updated_at,
			deleted_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY created_at ASC
	`, id)

	if err != nil {
		return MessageResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var message Messages

		if err := rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.Content,
			&message.MessageType,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.DeletedAt,
		); err != nil {
			return MessageResponse{}, err
		}

		messages.Messages = append(messages.Messages, message)
	}

	if err := rows.Err(); err != nil {
		return MessageResponse{}, err
	}

	return messages, nil
}

func (r *Repository) SendMessage(idUser uint64, idConversation uint64, content string) error {
	_, err := r.db.Exec(`
		INSERT INTO messages (
			sender_id,
			conversation_id,
			content
		)
		VALUES (?, ?, ?)
	`, idUser, idConversation, content)

	return err
}
