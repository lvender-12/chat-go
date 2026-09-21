package chat

import (
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
)

type Client struct {
	UserID         uint64
	Conn           *websocket.Conn
	ConversationID uint64
	Send           chan []byte
}

type Broadcast struct {
	ConversationID uint64
	Msg            []byte
}

type Hub struct {
	Clients map[uint64]map[*Client]struct{}

	mu sync.RWMutex

	register   chan *Client
	unregister chan *Client
	broadcast  chan Broadcast
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint64]map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Broadcast),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case broadcast := <-h.broadcast:
			h.sendBroadcast(broadcast)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.Clients[client.ConversationID] == nil {
		h.Clients[client.ConversationID] = make(map[*Client]struct{})
	}

	h.Clients[client.ConversationID][client] = struct{}{}
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, ok := h.Clients[client.ConversationID]
	if !ok {
		return
	}

	delete(clients, client)

	if len(clients) == 0 {
		delete(h.Clients, client.ConversationID)
	}
}

func (h *Hub) sendBroadcast(broadcast Broadcast) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := h.Clients[broadcast.ConversationID]

	for client := range clients {
		select {
		case client.Send <- broadcast.Msg:
		default:
		}
	}
}

func (h *Hub) Broadcast(broadcast Broadcast) {
	h.broadcast <- broadcast
}
