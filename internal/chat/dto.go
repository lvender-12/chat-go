package chat

import "time"

type MessageResponse struct {
	ConversationID uint64     `json:"conversation_id"`
	Messages       []Messages `json:"messages"`
}

type Messages struct {
	ID          uint64     `json:"id"`
	SenderID    uint64     `json:"sender_id"`
	Content     *string    `json:"content,omitempty"`
	MessageType string     `json:"message_type"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MessageRequest struct {
	Content string `json:"content"`
}

type MessageWithID struct {
	ID      uint64 `json:"ID"`
	Content string `json:"content"`
}

type DeleteMessageRequest struct {
	ID uint64 `json:"id"`
}
