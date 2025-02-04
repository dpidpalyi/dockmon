package main

import (
	"backend/internal/config"
	"backend/internal/data"
	"backend/internal/dbinit"
	"log"
	"net/http"
	"os"
)

type application struct {
	config config.Config
	logger *log.Logger
	models *data.Models
}

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	cfg, err := config.New(".")
	if err != nil {
		logger.Fatal(err)
	}

	db, err := dbinit.OpenDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	err = dbinit.RunMigrate(db, cfg.MigratePath)
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: app.routes(),
	}

	logger.Printf("starting backend on %s", srv.Addr)
	logger.Fatal(srv.ListenAndServe())
}
