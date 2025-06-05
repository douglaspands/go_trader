package resource

import "time"

type Reit struct {
	Ticker     string    `json:"ticker"`
	Name       string    `json:"name"`
	Document   string    `json:"document"`
	Admin      string    `json:"admin"`
	Segment    string    `json:"segment"`
	Price      float64   `json:"price"`
	Origin     string    `json:"origin"`
	CapturedAt time.Time `json:"capturedAt"`
}
