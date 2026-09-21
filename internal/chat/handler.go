package chat

import (
	"chat-go/internal/app"
	"log/slog"
	"strconv"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
	state   *app.State
	logger  *slog.Logger
}

func NewHandler(service *Service, state *app.State, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		state:   state,
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

	if id == "" {
		return ws.WriteJSON(WsResponse{
			Status:  fiber.StatusBadRequest,
			Message: "id is required",
		})
	}

	h.logger.Debug("GetChat", "id", id)
	idUint64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ws.WriteJSON(WsResponse{
			Status:  fiber.StatusBadRequest,
			Message: "id is required",
		})
	}
	chat, err := h.service.GetChat(idUint64)
	if err != nil {
		return ws.WriteJSON(WsResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ws.WriteJSON(WsResponse{
		Status:  fiber.StatusOK,
		Message: "success",
		Data:    chat,
	})
}
