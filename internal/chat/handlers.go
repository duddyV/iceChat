package chat

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

// upgrader configures the WebSocket upgrader with buffer sizes and an origin check function.
// The CheckOrigin function here allows all origins; this should be customized for production security.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Server holds the dependencies required for running the chat server,
// including a Hub for managing client connections and a Jet template set for rendering views.
type Server struct {
	Hub  *Hub
	View *jet.Set
}

// HandleChat renders the chat page when the root URL ("/") is accessed.
// It uses the Jet template engine to render the "chat.jet" template.
func (s *Server) HandleChat(w http.ResponseWriter, r *http.Request) {
	err := s.renderPage(w, "chat.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

// renderPage renders a Jet template with the given name and data.
// It retrieves the template from the template set and executes it, writing the output to the response writer.
func (s *Server) renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := s.View.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// HandleWs upgrades an HTTP connection to a WebSocket connection.
// It creates a new client instance and registers it with the Hub, starting the read and write pumps.
func (s *Server) HandleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	log.Println("Client connected to WebSocket")

	client := NewClient(s.Hub, conn, "Anonymous")
	s.Hub.register <- client

	go client.ReadPump()
	go client.WritePump()
}

// Routes sets up the HTTP routes for the chat server.
// It defines handlers for the chat page and the WebSocket endpoint.
func (s *Server) Routes(r chi.Router) {
	r.Get("/", s.HandleChat)
	r.Get("/ws", s.HandleWs)
}
