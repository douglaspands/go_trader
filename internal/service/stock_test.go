package service

import (
	"testing"
	"trader/internal/resource"
	"trader/internal/scraping"
)

func TestGetStockFound(t *testing.T) {
	// EXPECT
	var beforeScrapping = scraping.GetStockByTicker
	scraping.GetStockByTicker = func(ticket string) *resource.Security {
		return nil
	}
	// THEN

	// WHEN
	scraping.GetStockByTicker = beforeScrapping
	if outString != expected {
		t.Errorf("Output expected: \"%s\", but received: \"%s\"", expected, outString)
	}
}
