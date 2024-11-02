// sockets/socket.go
package sockets

import (
	"encoding/json"
	"log"
	"my-app/models"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID      string
	Conn    *websocket.Conn
	Send    chan []byte
	Room    string
	IsAdmin bool
	UserID  int // Có thể lưu userID nếu người dùng đã đăng nhập
}

type Message struct {
	UserID  int    `json:"userID,omitempty"`
	Text    string `json:"text"`
	Room    string `json:"room"`
	IsAdmin bool   `json:"isAdmin"`
}

type Hub struct {
	Clients    map[*Client]bool
	Rooms      map[string]map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Exported HubInstance for access in main.go
var HubInstance = NewHub()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWs handles WebSocket requests from clients
func ServeWs(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Generate a unique ID for the client
	clientID := uuid.NewString()

	// Get room from query parameter
	room := r.URL.Query().Get("room")
	if room == "" {
		room = "default"
	}

	// Check if the client is admin
	isAdmin := r.URL.Query().Get("admin") == "true"

	// Get userID if available (for logged-in users)
	userIDStr := r.URL.Query().Get("userID")
	var userID int
	if userIDStr != "" {
		userID, _ = strconv.Atoi(userIDStr)
	}

	client := &Client{
		ID:      clientID,
		Conn:    conn,
		Send:    make(chan []byte),
		Room:    room,
		IsAdmin: isAdmin,
		UserID:  userID,
	}

	HubInstance.Register <- client

	// Start read and write pumps
	go client.readPump()
	go client.writePump()
}

// readPump pumps messages from the WebSocket connection to the Hub
func (c *Client) readPump() {
	defer func() {
		HubInstance.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}

		// Gắn thêm thông tin vào tin nhắn
		msg.Room = c.Room
		msg.IsAdmin = c.IsAdmin
		msg.UserID = c.UserID

		// Lưu tin nhắn vào database nếu cần
		if !c.IsAdmin {
			err = models.AddMessage(msg.UserID, msg.Text)
			if err != nil {
				log.Println("Failed to save message:", err)
				continue
			}
		}

		// Broadcast message to the room
		HubInstance.Broadcast <- msg
	}
}

// writePump pumps messages from the Hub to the WebSocket connection
func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Write error:", err)
				return
			}
		}
	}
}

// Run starts the Hub to listen for register, unregister, and broadcast events
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			if h.Rooms[client.Room] == nil {
				h.Rooms[client.Room] = make(map[*Client]bool)
			}
			h.Rooms[client.Room][client] = true
			log.Printf("Client %s joined room %s\n", client.ID, client.Room)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				if h.Rooms[client.Room] != nil {
					delete(h.Rooms[client.Room], client)
					if len(h.Rooms[client.Room]) == 0 {
						delete(h.Rooms, client.Room)
					}
				}
				close(client.Send)
				log.Printf("Client %s disconnected\n", client.ID)
			}
		case message := <-h.Broadcast:
			if clients, ok := h.Rooms[message.Room]; ok {
				msgBytes, err := json.Marshal(message)
				if err != nil {
					log.Println("Failed to marshal message:", err)
					continue
				}
				for client := range clients {
					select {
					case client.Send <- msgBytes:
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}
