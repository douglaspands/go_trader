package service

import (
	"trader/internal/resource"
	"trader/internal/scraping"
)

func GetStock(ticker string) *resource.Stock {
	stock, err := scraping.GetStockByTicker(ticker)
	if err != nil {
		return nil
	}
	return stock
}

func ListStocks(tickers []string) []*resource.Stock {
	return scraping.ListStocksByTickers(tickers)
}
