package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	url := fmt.Sprintf("%v/%d", app.url, container.ID)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(js))
	req.Header.Set("Content-type", "application/json")
	resp, err := app.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	app.infoLogger.Printf("%s", data)

	return nil
}
