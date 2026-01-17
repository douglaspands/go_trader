package scraping_test

import (
	"slices"
	"testing"
	"trader/internal/config"
	"trader/internal/scraping"
)

func TestGetReitByTickerOk(t *testing.T) {
	// GIVEN
	ticker := "BPML11"

	// THEN
	reitScraping := scraping.NewReitScraping(config.NewConfig())
	result, _ := reitScraping.GetReitByTicker(ticker)

	// WHEN
	if result.Ticker != ticker {
		t.Errorf(`expected at "%s" and received at %s`, ticker, result.Ticker)
	}
}

func TestGetReitByTickerNotFound(t *testing.T) {
	// GIVEN
	ticker := "XXXX00"

	// THEN
	reitScraping := scraping.NewReitScraping(config.NewConfig())
	_, err := reitScraping.GetReitByTicker(ticker)

	// WHEN
	if err == nil {
		t.Errorf(`expected error but received nil`)
	}
}

func TestListReitsByTickersOk(t *testing.T) {
	// GIVEN
	tickers := []string{"HTMX11", "PORD11"}

	// THEN
	reitScraping := scraping.NewReitScraping(config.NewConfig())
	result := reitScraping.ListReitsByTickers(tickers)

	// WHEN
	for _, reit := range result {
		if !slices.Contains(tickers, reit.Ticker) {
			t.Errorf(`not found %s in %v`, reit.Ticker, tickers)
		}
	}
}
