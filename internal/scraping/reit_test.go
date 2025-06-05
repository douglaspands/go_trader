package scraping_test

import (
	"slices"
	"testing"
	"trader/internal/scraping"
)

func TestGetReitByTickerOk(t *testing.T) {
	ticker := "BPML11"
	result, _ := scraping.GetReitByTicker(ticker)
	if result.Ticker != ticker {
		t.Errorf(`expected at "%s" and received at %s`, ticker, result.Ticker)
	}
}

func TestListReitsByTickersOk(t *testing.T) {
	tickers := []string{"HTMX11", "PORD11"}
	result := scraping.ListReitsByTickers(tickers)
	for _, reit := range result {
		if !slices.Contains(tickers, reit.Ticker) {
			t.Errorf(`not found %s in %v`, reit.Ticker, tickers)
		}
	}
}
