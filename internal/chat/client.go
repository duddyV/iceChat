package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a single chat client connected to the Hub.
// It holds a reference to the Hub, a WebSocket connection, a send channel for outgoing messages, and the client's username.
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan *Message
	username string
}

// Message represents the structure of messages exchanged between clients.
// It includes the message type, the username of the sender, the message content, and an optional list of online users.
type Message struct {
	Type        string   `json:"type"`
	Username    string   `json:"username"`
	Message     string   `json:"message"`
	OnlineUsers []string `json:"online_users,omitempty"`
}

// NewClient initializes and returns a pointer to a new Client instance.
// It takes the Hub, a WebSocket connection, and a username as arguments.
func NewClient(hub *Hub, conn *websocket.Conn, username string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan *Message),
		username: username,
	}
}

// ReadPump listens for incoming messages from the WebSocket connection and processes them.
// This function should be run as a goroutine since it blocks while listening for messages.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		var msg Message
		if err := c.conn.ReadJSON(&msg); err != nil {
			break
		}

		switch msg.Type {
		case "chat":
			// Handle regular chat messages
			c.hub.broadcast <- &msg
		case "update-username":
			// Update the client's username and notify others
			c.username = msg.Username
			c.hub.broadcastOnlineClients() // Notify others of the updated online list
		}
	}
}

// WritePump listens for messages to send to the client via the WebSocket connection.
// It should be run as a goroutine to allow asynchronous sending of messages.
func (c *Client) WritePump() {
	defer c.conn.Close()

	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			log.Println(err)
			return
		}
	}

	// Send the close message to the client
	if err := c.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(
			websocket.CloseNormalClosure,
			"Normal closure",
		)); err != nil {
		log.Println(err)
		return
	}
}
