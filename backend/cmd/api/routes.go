package main

import (
	"net/http"

	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", app.health)
	mux.HandleFunc("POST /api/containers", app.insertContainer)
	mux.HandleFunc("GET /api/containers", app.listContainer)
	mux.HandleFunc("GET /api/containers/{id}", app.getContainer)
	mux.HandleFunc("PATCH /api/containers/{id}", app.updateContainer)
	mux.HandleFunc("DELETE /api/containers/{id}", app.deleteContainer)

	cMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	return cMiddleware.Handler(mux)
}
