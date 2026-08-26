package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Hub manages all active WebSocket connections
type Hub struct {
	mu      sync.RWMutex
	clients map[int64]*Client // userID → client (one active session per user)
}

// NewHub creates a hub
func NewHub() *Hub {
	return &Hub{clients: make(map[int64]*Client)}
}

// Register adds a client
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if old, ok := h.clients[c.UserID]; ok {
		// kick old connection (new login wins)
		old.Close()
	}
	h.clients[c.UserID] = c
}

// Unregister removes a client
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if cur, ok := h.clients[c.UserID]; ok && cur == c {
		delete(h.clients, c.UserID)
	}
}

// SendToUser delivers an event to a specific user (if online)
func (h *Hub) SendToUser(userID int64, event Event) {
	h.mu.RLock()
	c, ok := h.clients[userID]
	h.mu.RUnlock()
	if ok {
		c.Send(event)
	}
}

// BroadcastToRoom delivers an event to all online members of a room.
// memberIDs must be the visible (non-hidden) members.
func (h *Hub) BroadcastToRoom(roomID int64, memberIDs []int64, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	// quick check: only users with active clients
	for _, uid := range memberIDs {
		if c, ok := h.clients[uid]; ok {
			ev := event
			ev.RoomID = roomID
			c.Send(ev)
		}
	}
}

// BroadcastToAll sends an event to every connected client
func (h *Hub) BroadcastToAll(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.clients {
		c.Send(event)
	}
}

// OnlineUserIDs returns ids of currently connected users
func (h *Hub) OnlineUserIDs() []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ids := make([]int64, 0, len(h.clients))
	for id := range h.clients {
		ids = append(ids, id)
	}
	return ids
}

// IsOnline checks if a user has an active connection
func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// ---------- Client ----------

// Client wraps a single websocket connection
type Client struct {
	UserID   int64
	Username string
	Display  string
	conn     *websocket.Conn
	send     chan Event
	hub      *Hub
	once     sync.Once
}

// NewClient creates a client
func NewClient(hub *Hub, userID int64, username, display string, conn *websocket.Conn) *Client {
	return &Client{
		UserID: userID, Username: username, Display: display,
		conn: conn, send: make(chan Event, 64), hub: hub,
	}
}

// Send queues an event (non-blocking with drop on full)
func (c *Client) Send(e Event) {
	select {
	case c.send <- e:
	default:
		// drop if slow consumer — client can re-fetch history
	}
}

// Close terminates the connection
func (c *Client) Close() {
	c.once.Do(func() {
		c.hub.Unregister(c)
		_ = c.conn.Close()
		close(c.send)
	})
}

// WritePump writes events from the send channel to the socket
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case ev, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			data, err := json.Marshal(ev)
			if err != nil {
				continue
			}
			if err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump reads client events and dispatches to the handler
func (c *Client) ReadPump(handler func(ClientEvent)) {
	defer c.Close()
	c.conn.SetReadLimit(1 << 20) // 1MB max message
	_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var ev ClientEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			log.Printf("ws: bad event from %d: %v", c.UserID, err)
			continue
		}
		ev.UserID = c.UserID
		handler(ev)
	}
}
