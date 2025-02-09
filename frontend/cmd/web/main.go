package main

import (
	"frontend/internal/config"
	"log"
	"net/http"
	"os"
	"text/template"
)

type application struct {
	config        config.Config
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	client        http.Client
	templateCache map[string]*template.Template
}

func main() {
	infoLogger := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime)

	cfg, err := config.New(".")
	if err != nil {
		errorLogger.Fatal(err)
	}

	client := http.Client{
		Timeout: cfg.APIRequestTimeout,
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLogger.Fatal(err)
	}

	app := &application{
		config:        cfg,
		infoLogger:    infoLogger,
		errorLogger:   errorLogger,
		client:        client,
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: app.routes(),
	}

	infoLogger.Printf("starting frontend on %s", srv.Addr)
	errorLogger.Fatal(srv.ListenAndServe())
}
