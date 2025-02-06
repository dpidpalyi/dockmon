package main

import (
	"log"
	"net/http"
	"os"
	"pinger/internal/config"

	probing "github.com/prometheus-community/pro-bing"
)

type application struct {
	url         string
	infoLogger  *log.Logger
	errorLogger *log.Logger
	client      http.Client
	pinger      *probing.Pinger
}

func main() {
	infoLogger := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime)

	cfg, err := config.New(".")
	if err != nil {
		errorLogger.Fatal(err)
	}

	client := http.Client{
		Timeout: cfg.RequestTimeout,
	}

	pinger := probing.New("127.0.0.1")

	app := &application{
		url:         cfg.APIurl,
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		client:      client,
		pinger:      pinger,
	}

	app.Run()
}

func (app *application) Run() {
	cs, err := app.Get()
	if err != nil {
		app.errorLogger.Print(err)
	}

	for _, c := range cs {
		c.Ping = 100
		c.Status = "up"
		app.Send(c)
	}
}
