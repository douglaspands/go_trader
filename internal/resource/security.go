package resource

import "time"

const (
	STOCK_TYPE string = "STOCK"
	REIT_TYPE  string = "REIT"
)

type Security struct {
	Ticker      string    `json:"ticker"`
	Type        string    `json:"type"` // Stock, REIT
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Document    string    `json:"document"`
	Admin       string    `json:"admin"`
	Segment     string    `json:"segment"`
	Currency    *Currency `json:"currency"`
	Price       float64   `json:"price"`
	Origin      string    `json:"origin"`
	CapturedAt  time.Time `json:"capturedAt"`
}
