package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", app.health)
	mux.HandleFunc("POST /api/containers", app.insertContainer)
	mux.HandleFunc("GET /api/containers/{id}", app.getContainer)
	mux.HandleFunc("PATCH /api/containers/{id}", app.updateContainer)

	return mux
}
