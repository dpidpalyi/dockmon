package main

import (
	"backend/internal/data"
	"backend/internal/validator"
	"context"
	"fmt"
	"net/http"
	"time"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	if err := app.models.Containers.DB.Ping(); err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "database is down")
		return
	}

	fmt.Fprintln(w, "ok")
}

func (app *application) insertContainer(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
	}

	err := app.readJSON(r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	container := &data.Container{
		Name: input.Name,
		IP:   input.IP,
	}

	v := validator.New()
	if data.ValidateContainer(v, container); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = app.models.Containers.Insert(ctx, container)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("api/containers/%d", container.ID))

	err = app.writeJSON(w, http.StatusCreated, container, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
