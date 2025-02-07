package main

import (
	"errors"

	probing "github.com/prometheus-community/pro-bing"
)

var ErrUnreachable = errors.New("host is unreachable")

func (app *application) Ping(addr string) (float64, error) {
	pinger := probing.New(addr)
	pinger.Count = app.config.Count
	pinger.Interval = app.config.Interval
	pinger.Timeout = app.config.Timeout
	err := pinger.Run()
	if err != nil {
		return 0, err
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		return 0, ErrUnreachable
	}
	result := float64(stats.AvgRtt.Microseconds()) / 1000

	return result, nil
}
