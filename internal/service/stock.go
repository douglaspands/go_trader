package service

import (
	"trader/internal/resource"
	"trader/internal/scraping"
)

type StockService interface {
	GetStockByTicker(ticker string) *resource.Security
	ListStocksByTickers(tickers []string) ([]*resource.Security, []*resource.TickerFailure)
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

func (ss *stockService) ListStocksByTickers(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
	return ss.stockScraping.ListStocksByTickers(tickers)
}

func NewStockService(stockScraping scraping.StockScraping) StockService {
	return &stockService{
		stockScraping: stockScraping,
	}
}
