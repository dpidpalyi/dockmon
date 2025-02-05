package main

import (
	"encoding/json"
	"net/http"
	"pinger/internal/data"
)

func (app *application) Get() ([]*data.Container, error) {
	req, err := http.NewRequest("GET", app.url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := app.client.Do(req)
	if err != nil {
		return nil, err
	}

	var containers []*data.Container

	err = json.NewDecoder(resp.Body).Decode(&containers)
	if err != nil {
		return nil, err
	}

	return containers, nil
}
