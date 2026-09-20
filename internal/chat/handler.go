package chat

import (
	"chat-go/internal/app"
	"chat-go/internal/utils"
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
	state   *app.State
	hub     *Hub
	logger  *slog.Logger
}

func NewHandler(
	service *Service,
	state *app.State,
	hub *Hub,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		service: service,
		state:   state,
		hub:     hub,
		logger:  logger,
	}
}

// GetChat godoc
// @Summary Get chat history
// @Description Open a WebSocket connection and get chat history for a conversation
// @Tags Chat
// @Param id path uint64 true "Conversation ID"
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 426 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/chat/{id} [get]
func (h *Handler) GetChat(ws *websocket.Conn) error {
	id := ws.Params("id")

	conversationID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ws.WriteJSON(WsResponse{
			Status:  fiber.StatusBadRequest,
			Message: "invalid conversation id",
		})
	}

	userIDValue := ws.Locals("user_id")

	userID, ok := userIDValue.(uint64)
	if !ok {
		return ws.WriteJSON(WsResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "unauthorized",
		})
	}

	chat, err := h.service.GetChat(conversationID)
	if err != nil {
		return ws.WriteJSON(WsResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	data, err := json.Marshal(WsResponse{
		Status:  fiber.StatusOK,
		Message: "success",
		Data:    chat,
	})
	if err != nil {
		return err
	}

	client := &Client{
		UserID:         userID,
		Conn:           ws,
		ConversationID: conversationID,
		Send:           make(chan []byte, 256),
	}

	h.hub.register <- client

	defer func() {
		h.hub.unregister <- client
	}()

	go client.WritePump(h.logger)

	client.Send <- data

	client.ReadPump(h.hub, h.logger)

	return nil
}

// SendMessage godoc
// @Summary Send a message
// @Description Send a message to a conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Param id path uint64 true "Conversation ID"
// @Param request body MessageRequest true "Message content"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/chat/{id} [post]
func (h *Handler) SendMessage(c fiber.Ctx) error {
	var input MessageRequest

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

	conversationID, err := strconv.ParseUint(
		c.Params("id"),
		10,
		64,
	)
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid conversation id",
		)
	}

	if err := h.service.SendMessage(
		userID,
		conversationID,
		input.Content,
	); err != nil {
		return err
	}

	chat, err := h.service.GetChat(conversationID)
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}

	data, err := json.Marshal(WsResponse{
		Status:  fiber.StatusOK,
		Message: "message",
		Data:    chat,
	})
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"failed to encode message",
		)
	}

	h.hub.broadcast <- Broadcast{
		ConversationID: conversationID,
		Msg:            data,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "message sent",
	})
}
