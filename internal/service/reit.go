package service

import (
	"trader/internal/resource"
	"trader/internal/scraping"
)

func GetReit(ticker string) *resource.Reit {
	stock, err := scraping.GetReitByTicker(ticker)
	if err != nil {
		return nil
	}
	return stock
}

func ListReits(tickers []string) []*resource.Reit {
	return scraping.ListReitsByTickers(tickers)
}
