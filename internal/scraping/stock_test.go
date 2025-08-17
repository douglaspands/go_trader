package scraping_test

import (
	"slices"
	"testing"
	"trader/internal/config"
	"trader/internal/scraping"
)

func TestGetStockByTickerOk(t *testing.T) {
	// GIVEN
	ticker := "PETR4"

	// THEN
	stockScraping := scraping.NewStockScraping(config.NewConfig())
	result, _ := stockScraping.GetStockByTicker(ticker)

	// WHEN
	if result.Ticker != ticker {
		t.Errorf(`expected at "%s" and received at %s`, ticker, result.Ticker)
	}
}

func TestListStocksByTickersOk(t *testing.T) {
	// GIVEN
	tickers := []string{"ITSA3", "BBDC3", "VALE3", "ABEV3", "PETR4", "WEGE3", "IGTA3", "B3SA3"}

	// THEN
	stockScraping := scraping.NewStockScraping(config.NewConfig())
	result := stockScraping.ListStocksByTickers(tickers)

	// WHEN
	for _, stock := range result {
		if !slices.Contains(tickers, stock.Ticker) {
			t.Errorf(`not found %s in %v`, stock.Ticker, tickers)
		}
	}
}
