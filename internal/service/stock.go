package service

import (
	"trader/internal/resource"
	"trader/internal/scraping"
)

type StockService interface {
	GetStockByTicker(ticker string) *resource.Security
	ListStocksByTickers(tickers []string) []*resource.Security
}

type stockService struct {
	stockScraping scraping.StockScraping
}

func (ss *stockService) GetStockByTicker(ticker string) *resource.Security {
	stock, err := ss.stockScraping.GetStockByTicker(ticker)
	if err != nil {
		return nil
	}
	return stock
}

func (ss *stockService) ListStocksByTickers(tickers []string) []*resource.Security {
	result := ss.stockScraping.ListStocksByTickers(tickers)
	return result
}

func NewStockService(stockScraping scraping.StockScraping) StockService {
	return &stockService{
		stockScraping: stockScraping,
	}
}
