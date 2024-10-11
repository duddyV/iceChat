package chat

import "log"

// Hub maintains the set of active clients and broadcasts messages to the clients.
// It is the central hub for managing client registration, unregistration, and message broadcasting.
type Hub struct {
	clients    map[*Client]bool // All connected clients
	broadcast  chan *Message    // Channel for broadcasting messages
	register   chan *Client     // Register new clients
	unregister chan *Client     // Unregister clients
}

// NewHub initializes and returns a pointer to a new Hub instance.
// It sets up all necessary channels and the map to track active clients.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the Hub's main event loop, handling client registration, unregistration, and message broadcasting.
// This function should be run in a goroutine as it blocks indefinitely, processing events.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("client registered")
			h.broadcastOnlineClients() // Broadcast updated online clients list
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("client removed")
				h.broadcastOnlineClients() // Broadcast updated online clients list
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					log.Println("msg sent")
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// broadcastOnlineClients sends an updated list of online clients to all connected clients.
// It creates a system message containing the list of online users and broadcasts it.
func (h *Hub) broadcastOnlineClients() {
	var onlineUsers []string
	for client := range h.clients {
		onlineUsers = append(onlineUsers, client.username) // Collect usernames directly
	}

	// Create the message containing the online users list
	message := Message{
		Type:        "update-username", // Message type
		Username:    "system",          // System-generated message
		Message:     "online users",    // Default message text
		OnlineUsers: onlineUsers,       // Add online users list to message
	}

	// Broadcast the message to all connected clients
	for client := range h.clients {
		client.send <- &message
	}
}
