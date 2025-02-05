package main

import (
	"log"
	"net/http"
	"os"
	"pinger/internal/config"

	probing "github.com/prometheus-community/pro-bing"
)

type application struct {
	url    string
	logger *log.Logger
	client http.Client
	pinger *probing.Pinger
}

func main() {
	logger := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)

	cfg, err := config.New(".")
	if err != nil {
		logger.Fatal(err)
	}

	client := http.Client{
		Timeout: cfg.RequestTimeout,
	}

	pinger := probing.New("127.0.0.1")

	app := &application{
		url:    cfg.APIurl,
		logger: logger,
		client: client,
		pinger: pinger,
	}

	app.Run()
}

func (app *application) Run() {
	cs, err := app.Get()
	if err != nil {
		app.logger.Print(err)
	}

	for _, c := range cs {
		app.logger.Print(*c)
	}
}
