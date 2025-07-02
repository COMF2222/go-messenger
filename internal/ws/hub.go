package ws

import "sync"

// Client - подключение конктретного пользователя
type Client struct {
	UserID int
	Conn   *Connection
	Send   chan []byte
}

// Hub - главный менеджер WS соединений
type Hub struct {
	Clients    map[int]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan WSMessage
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan WSMessage),
	}
}

type WSMessage struct {
	FromUserID int
	ToUserID   int
	Message    []byte
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Clients[client.UserID] == nil {
				h.Clients[client.UserID] = make(map[*Client]bool)
			}
			h.Clients[client.UserID][client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID][client]; ok {
				delete(h.Clients[client.UserID], client)
				close(client.Send)
				if len(h.Clients[client.UserID]) == 0 {
					delete(h.Clients, client.UserID)
				}
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.Lock()
			clients := h.Clients[message.ToUserID]
			for client := range clients {
				select {
				case client.Send <- message.Message:
				default:
					close(client.Send)
					delete(clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}
