package resource

import "time"

type Stock struct {
	Ticker      string    `json:"ticker"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Document    string    `json:"document"`
	Price       float64   `json:"price"`
	Origin      string    `json:"origin"`
	CapturedAt  time.Time `json:"capturedAt"`
}
