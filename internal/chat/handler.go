package chat

import (
	"chat-go/internal/app"
	"chat-go/internal/utils"
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

func NewHandler(service *Service, state *app.State, hub *Hub, logger *slog.Logger) *Handler {
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

	idUint64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ws.WriteJSON(WsResponse{
			Status:  fiber.StatusBadRequest,
			Message: "invalid conversation id",
		})
	}

	h.hub.Register(idUint64, ws)

	defer h.hub.Unregister(idUint64, ws)

	chat, err := h.service.GetChat(idUint64)
	if err != nil {
		return ws.WriteJSON(WsResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if err := ws.WriteJSON(WsResponse{
		Status:  fiber.StatusOK,
		Message: "success",
		Data:    chat,
	}); err != nil {
		return err
	}

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			return err
		}
	}
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

	conversationID, err := strconv.ParseUint(c.Params("id"), 10, 64)
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

	h.hub.Broadcast(conversationID, chat)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "message sent",
	})
}
