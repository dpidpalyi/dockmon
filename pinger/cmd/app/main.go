package main

import (
	"log"
	"net/http"
	"os"
	"pinger/internal/config"
	"sync"
	"time"
)

type application struct {
	config      config.Config
	infoLogger  *log.Logger
	errorLogger *log.Logger
	client      http.Client
}

const (
	StatusUP   = "up"
	StatusDOWN = "down"
)

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

	app := &application{
		config:      cfg,
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		client:      client,
	}

	app.infoLogger.Printf("starting to work with api server: %v", app.config.APIurl)
	app.Run()
}

func (app *application) Run() {
	for {
		time.Sleep(app.config.IterationDelay)

		app.infoLogger.Printf("requesting data from server...")
		cs, err := app.Get()
		if err != nil {
			app.errorLogger.Print(err)
			continue
		}

		app.infoLogger.Printf("got %d containers to check\n", len(cs))

		var wg sync.WaitGroup
		for _, c := range cs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ping, err := app.Ping(c.IP)
				if err != nil {
					app.errorLogger.Printf("[DOWN] %10v:%-12v %v", c.Name, c.IP, err)
					if c.Status == StatusUP {
						c.Ping = 0.0
						c.Status = StatusDOWN
					} else {
						return
					}
				} else {
					c.Ping = ping
					c.Status = StatusUP
					app.infoLogger.Printf("  [UP] %10v:%-12v %.3fms", c.Name, c.IP, c.Ping)
				}

				err = app.Send(c)
				if err != nil {
					app.errorLogger.Print(err)
					return
				}
			}()
		}
		wg.Wait()
	}
}
