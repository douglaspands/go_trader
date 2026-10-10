//go:build integration

package scraping_test

import (
	"slices"
	"testing"
	"trader/internal/config"
	"trader/internal/scraping"
)

func newStockScraping() scraping.StockScraping {
	cfg := config.NewConfig()
	return scraping.NewStockScraping(cfg, scraping.NewFetcher(cfg))
}

func TestGetStockByTickerOk(t *testing.T) {
	// GIVEN
	ticker := "PETR4"

	// THEN
	stockScraping := newStockScraping()
	result, _ := stockScraping.GetStockByTicker(ticker)

	// WHEN
	if result.Ticker != ticker {
		t.Errorf(`expected at "%s" and received at %s`, ticker, result.Ticker)
	}
}

func TestGetStockByTickerNotFound(t *testing.T) {
	// GIVEN
	ticker := "XXXX0"

	// THEN
	stockScraping := newStockScraping()
	_, err := stockScraping.GetStockByTicker(ticker)

	// WHEN
	if err == nil {
		t.Errorf(`expected error but received nil`)
	}
}

func TestListStocksByTickersOk(t *testing.T) {
	// GIVEN
	tickers := []string{"ITSA3", "BBDC3", "VALE3", "ABEV3", "PETR4", "WEGE3", "IGTA3", "B3SA3"}

	// THEN
	stockScraping := newStockScraping()
	result, failures := stockScraping.ListStocksByTickers(tickers)
	for _, failure := range failures {
		t.Logf("failed %s: %v", failure.Ticker, failure.Err)
	}

	// WHEN
	for _, stock := range result {
		if !slices.Contains(tickers, stock.Ticker) {
			t.Errorf(`not found %s in %v`, stock.Ticker, tickers)
		}
	}
}
