package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", app.health)
	mux.HandleFunc("POST /api/containers", app.insertContainer)

	return mux
}
