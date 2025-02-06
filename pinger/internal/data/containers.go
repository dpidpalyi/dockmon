package data

import "time"

type Container struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Status  string `json:"status"`
	Version int    `json:"version"`
	// in ns
	Ping      float64    `json:"ping"`
	UpdatedAt *time.Time `json:"updated_at"`
}
