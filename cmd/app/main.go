package main

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/duddyV/iceChat/internal/chat"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialize and start a hub
	hub := chat.NewHub()
	go hub.Run()

	// Initialize the Jet template engine
	viewSet := jet.NewSet(
		jet.NewOSFileSystemLoader("./web/templates"),
		jet.InDevelopmentMode(),
	)

	// Set up the server
	r := chi.NewRouter()
	server := &chat.Server{Hub: hub, View: viewSet}
	server.Routes(r)

	// Serve static files for CSS and JS
	r.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))).ServeHTTP)
	r.Get("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js/"))).ServeHTTP)

	// Start the HTTP server
	log.Println("Starting web server on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic("Error starting web server: " + err.Error())
	}
}
