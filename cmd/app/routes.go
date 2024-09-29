package main

import (
	"net/http"

	"github.com/duddyV/iceChat/internal/server"
	"github.com/go-chi/chi/v5"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", server.HomePage)

	return mux
}
