package main

import (
	"backend/internal/data"
	"backend/internal/validator"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	responseTimeout = 3
	statusUp        = "up"
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

	ctx, cancel := context.WithTimeout(context.Background(), responseTimeout*time.Second)
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

func (app *application) getContainer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), responseTimeout*time.Second)
	defer cancel()

	container, err := app.models.Containers.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, container, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateContainer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), responseTimeout*time.Second)
	defer cancel()

	container, err := app.models.Containers.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name   *string  `json:"name"`
		IP     *string  `json:"ip"`
		Status *string  `json:"status"`
		Ping   *float64 `json:"ping"`
	}

	err = app.readJSON(r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		container.Name = *input.Name
	}

	if input.IP != nil {
		container.IP = *input.IP
	}

	if input.Status != nil {
		container.Status = *input.Status
	}

	if input.Ping != nil {
		container.Ping = *input.Ping
	}

	if container.Status == statusUp && container.Ping > 0 {
		t := time.Now()
		container.UpdatedAt = &t
	}

	v := validator.New()
	if data.ValidateContainer(v, container); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Containers.Update(ctx, container)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, container, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteContainer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), responseTimeout*time.Second)
	defer cancel()

	err = app.models.Containers.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "container successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listContainer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), responseTimeout*time.Second)
	defer cancel()

	containers, err := app.models.Containers.List(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, containers, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
