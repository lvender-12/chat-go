package chat

import (
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
)

type Hub struct {
	mu    sync.RWMutex
	rooms map[uint64]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[uint64]map[*websocket.Conn]struct{}),
	}
}

func (h *Hub) Register(conversationID uint64, ws *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[conversationID] == nil {
		h.rooms[conversationID] = make(map[*websocket.Conn]struct{})
	}

	h.rooms[conversationID][ws] = struct{}{}
}

func (h *Hub) Unregister(conversationID uint64, ws *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.rooms[conversationID], ws)

	if len(h.rooms[conversationID]) == 0 {
		delete(h.rooms, conversationID)
	}
}

func (h *Hub) Broadcast(conversationID uint64, data interface{}) {
	h.mu.RLock()

	connections := make([]*websocket.Conn, 0, len(h.rooms[conversationID]))

	for ws := range h.rooms[conversationID] {
		connections = append(connections, ws)
	}

	h.mu.RUnlock()

	for _, ws := range connections {
		_ = ws.WriteJSON(WsResponse{
			Status:  200,
			Message: "message",
			Data:    data,
		})
	}
}
