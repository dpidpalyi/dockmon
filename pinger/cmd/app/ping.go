package main

import (
	"context"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func (app *application) Ping(addr string) (float64, error) {
	pinger := probing.New(addr)
	app.infoLogger.Print(pinger.Addr())
	pinger.PacketsSent = 1
	pinger.Count = 1
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := pinger.RunWithContext(ctx)
	if err != nil {
		return 0, err
	}

	stats := pinger.Statistics()
	result := float64(stats.AvgRtt.Microseconds()) / 1000
	app.infoLogger.Print(result)

	return result, nil
}
