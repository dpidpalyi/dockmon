package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pinger/internal/data"
)

func (app *application) Get() ([]*data.Container, error) {
	req, err := http.NewRequest("GET", app.config.APIurl, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := app.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("internal api server error")
	}

	var containers []*data.Container

	err = json.NewDecoder(resp.Body).Decode(&containers)
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (app *application) Send(container *data.Container) error {
	js, err := json.Marshal(container)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%v/%d", app.config.APIurl, container.ID)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(js))
	req.Header.Set("Content-type", "application/json")
	resp, err := app.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return fmt.Errorf("internal api server error")
	case http.StatusNotFound:
		return fmt.Errorf("the requested container ID:%v could not be found", container.ID)
	case http.StatusUnprocessableEntity:
		return fmt.Errorf("invalidated data provided to server")
	case http.StatusBadRequest:
		return fmt.Errorf("wrong JSON data provided to server")
	}

	return nil
}
