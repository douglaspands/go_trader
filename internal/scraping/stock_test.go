package scraping_test

import (
	"slices"
	"testing"
	"trader/internal/scraping"
)

func TestGetStockByTickerOk(t *testing.T) {
	ticker := "PETR4"
	result, _ := scraping.GetStockByTicker(ticker)
	if result.Ticker != ticker {
		t.Errorf(`expected at "%s" and received at %s`, ticker, result.Ticker)
	}
}

func TestListStocksByTickersOk(t *testing.T) {
	tickers := []string{"ITSA3", "BBDC3", "VALE3", "ABEV3", "PETR4", "WEGE3", "IGTA3", "B3SA3"}
	result := scraping.ListStocksByTickers(tickers)
	for _, stock := range result {
		if !slices.Contains(tickers, stock.Ticker) {
			t.Errorf(`not found %s in %v`, stock.Ticker, tickers)
		}
	}
}
