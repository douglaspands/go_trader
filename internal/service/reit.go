package service

import (
	"trader/internal/resource"
	"trader/internal/scraping"
)

func GetReit(ticker string) *resource.Security {
	stock, err := scraping.GetReitByTicker(ticker)
	if err != nil {
		return nil
	}
	return stock
}

func ListReits(tickers []string) []*resource.Security {
	return scraping.ListReitsByTickers(tickers)
}
