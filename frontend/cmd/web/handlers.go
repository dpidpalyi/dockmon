package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frontend/internal/data"
	"net/http"
	"path/filepath"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", app.config.APIurl, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := app.client.Do(req)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		app.serverError(w, err)
		return
	}

	var containers []*data.Container

	err = json.NewDecoder(resp.Body).Decode(&containers)
	if err != nil {
		app.serverError(w, err)
		return
	}

	d := struct {
		Containers []*data.Container
	}{
		Containers: containers,
	}

	app.render(w, http.StatusOK, "home.html", d)
}

func (app *application) view(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/%d", app.config.APIurl, id)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := app.client.Do(req)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		app.serverError(w, err)
		return
	}

	var container data.Container

	err = json.NewDecoder(resp.Body).Decode(&container)
	if err != nil {
		app.serverError(w, err)
		return
	}

	d := struct {
		Container data.Container
	}{
		Container: container,
	}

	app.render(w, http.StatusOK, "view.html", d)
}

func (app *application) addGet(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "add.html", nil)
}

func (app *application) addPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	ip := r.PostForm.Get("ip")

	container := data.Container{
		Name: name,
		IP:   ip,
	}

	js, err := json.Marshal(container)
	if err != nil {
		app.serverError(w, err)
		return
	}

	req, err := http.NewRequest("POST", app.config.APIurl, bytes.NewBuffer(js))
	req.Header.Add("Content-type", "application/json")
	resp, err := app.client.Do(req)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		app.clientError(w, http.StatusBadRequest)
		return
	case http.StatusInternalServerError:
		app.serverError(w, err)
		return
	}

	id := filepath.Base(resp.Header.Get("Location"))
	if id == "" {
		app.serverError(w, fmt.Errorf("unknown error"))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/view?id=%s", id), http.StatusSeeOther)
}
