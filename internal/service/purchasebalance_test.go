package service_test

import (
	"testing"
	"trader/internal/resource"
	"trader/internal/service"
)

// Mocks

type mockStockService struct {
	GetStockByTickerFunc    func(string) *resource.Security
	ListStocksByTickersFunc func([]string) []*resource.Security
}

func (m *mockStockService) GetStockByTicker(ticker string) *resource.Security {
	return m.GetStockByTickerFunc(ticker)
}

func (m *mockStockService) ListStocksByTickers(tickers []string) []*resource.Security {
	return m.ListStocksByTickersFunc(tickers)
}

type mockReitService struct {
	GetReitByTickerFunc    func(string) *resource.Security
	ListReitsByTickersFunc func([]string) []*resource.Security
}

func (m *mockReitService) GetReitByTicker(ticker string) *resource.Security {
	return m.GetReitByTickerFunc(ticker)
}

func (m *mockReitService) ListReitsByTickers(tickers []string) []*resource.Security {
	return m.ListReitsByTickersFunc(tickers)
}

// Tests

func TestPurchaseBalance_LowInvestment(t *testing.T) {
	// GIVEN
	securities := []*resource.Security{
		{Ticker: "S1", Price: 100},
		{Ticker: "S2", Price: 200},
	}
	amountInvested := 50.0

	svc := service.NewPurchaseBalanceService(&mockStockService{}, &mockReitService{})

	// WHEN
	result := svc.PurchaseBalance(securities, amountInvested)

	// THEN
	if result.TotalCount() != 0 {
		t.Errorf("Expected total count 0, got %d", result.TotalCount())
	}
	if result.AmountInvested != amountInvested {
		t.Errorf("Expected amount invested %.2f, got %.2f", amountInvested, result.AmountInvested)
	}
}

func TestPurchaseBalance_StandardDistribution(t *testing.T) {
	// GIVEN
	securities := []*resource.Security{
		{Ticker: "S1", Price: 10},
		{Ticker: "S2", Price: 10},
	}
	amountInvested := 100.0 // 50 per security -> 5 each

	svc := service.NewPurchaseBalanceService(&mockStockService{}, &mockReitService{})

	// WHEN
	result := svc.PurchaseBalance(securities, amountInvested)

	// THEN
	if result.TotalCount() != 10 {
		t.Errorf("Expected total count 10, got %d", result.TotalCount())
	}
	for _, p := range result.SecuritiesBalance {
		if p.Count != 5 {
			t.Errorf("Expected 5 units for %s, got %d", p.Security.Ticker, p.Count)
		}
	}
}

// func TestPurchaseBalance_UnevenDistribution_HighPrice(t *testing.T) {
// 	// GIVEN
// 	// Split is 500 each.
// 	// S1 needs 500, price 600 -> count 0 by split, goes to expensive list?
// 	// Wait logic: securityValue = 1000/2 = 500.
// 	// S1 (600): 500/600 < 1 -> expensive list.
// 	// S2 (100): 500/100 = 5 -> normal buy. remaining from split? logic doesn't subtract from split, it subtracts from total remaining.
// 	// Logic trace:
// 	// Init: remaining = 1000.
// 	// Loop 1:
// 	// S1: < 1 -> expensive list.
// 	// S2: count 5 (500). remaining = 1000 - 500 = 500.
// 	// Loop 2 (expensive):
// 	// S1: remaining (500) < 600 -> can't buy.
// 	// Loop 3 (remaining balance re-loop):
// 	// remaining 500.
// 	// S2 (100) < 500 -> buy 1 more. rem=400.
// 	// ... buys S2 until rem=0?
// 	// S2 is in securitiesPurchaseSort.
// 	// Logic line 75: sort by price. S2 checked.
// 	// It will keep buying S2.

// 	securities := []*resource.Security{
// 		{Ticker: "High", Price: 600},
// 		{Ticker: "Low", Price: 100},
// 	}
// 	amountInvested := 1000.0

// 	svc := service.NewPurchaseBalanceService(&mockStockService{}, &mockReitService{})

// 	// WHEN
// 	result := svc.PurchaseBalance(securities, amountInvested)

// 	// THEN
// 	// S2 should soak up the rest?
// 	// S1: 0
// 	// S2: 5 initially + 5 from remaining loop? = 10 total?
// 	// Total spent: 10 * 100 = 1000.

// 	countHigh := 0
// 	countLow := 0
// 	for _, p := range result.SecuritiesBalance {
// 		if p.Security.Ticker == "High" {
// 			countHigh = p.Count
// 		}
// 		if p.Security.Ticker == "Low" {
// 			countLow = p.Count
// 		}
// 	}

// 	if countHigh != 0 {
// 		t.Errorf("Expected 0 High, got %d", countHigh)
// 	}
// 	if countLow != 10 {
// 		t.Errorf("Expected 10 Low, got %d", countLow)
// 	}
// }

func TestPurchaseBalancesBySecurities(t *testing.T) {
	// GIVEN
	stockTickers := []string{"PETR4"}
	reitTickers := []string{"HGLG11"}
	amountInvested := 1000.0

	mockStocks := []*resource.Security{{Ticker: "PETR4", Price: 50}}
	mockReits := []*resource.Security{{Ticker: "HGLG11", Price: 100}}

	mockStockSvc := &mockStockService{
		ListStocksByTickersFunc: func(tickers []string) []*resource.Security {
			if len(tickers) == 1 && tickers[0] == "PETR4" {
				return mockStocks
			}
			return nil
		},
	}
	mockReitSvc := &mockReitService{
		ListReitsByTickersFunc: func(tickers []string) []*resource.Security {
			if len(tickers) == 1 && tickers[0] == "HGLG11" {
				return mockReits
			}
			return nil
		},
	}

	svc := service.NewPurchaseBalanceService(mockStockSvc, mockReitSvc)

	// WHEN
	result := svc.PurchaseBalancesBySecurities(stockTickers, reitTickers, amountInvested)

	// THEN
	// 2 securities. 500 each.
	// PETR4 (50): 10 shares.
	// HGLG11 (100): 5 shares.
	if result.TotalCount() != 15 {
		t.Errorf("Expected 15 total, got %d", result.TotalCount())
	}
}
