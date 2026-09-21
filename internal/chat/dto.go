package chat

import "time"

type MessageResponse struct {
	ID      uint64     `json:"id"`
	Message []Messages `json:"message"`
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

type WsResponse struct {
	Status  int
	Message string
	Data    interface{}
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
