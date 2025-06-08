package service

import (
	"trader/internal/resource"
	"trader/internal/scraping"
)

func GetStock(ticker string) *resource.Security {
	stock, err := scraping.GetStockByTicker(ticker)
	if err != nil {
		return nil
	}
	return stock
}

func ListStocks(tickers []string) []*resource.Security {
	return scraping.ListStocksByTickers(tickers)
}
