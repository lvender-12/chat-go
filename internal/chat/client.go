package chat

import (
	"log/slog"

	"github.com/gofiber/contrib/v3/websocket"
)

func (c *Client) ReadPump(hub *Hub, logger *slog.Logger) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			logger.Debug(
				"websocket connection closed",
				"user_id", c.UserID,
				"conversation_id", c.ConversationID,
				"error", err,
			)

			return
		}
	}
}

func (c *Client) WritePump(logger *slog.Logger) {
	defer c.Conn.Close()

	for message := range c.Send {
		if err := c.Conn.WriteMessage(
			websocket.TextMessage,
			message,
		); err != nil {
			logger.Debug(
				"websocket write failed",
				"user_id", c.UserID,
				"conversation_id", c.ConversationID,
				"error", err,
			)

			return
		}
	}
}
