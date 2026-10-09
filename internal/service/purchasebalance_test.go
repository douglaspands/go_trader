package service_test

import (
	"fmt"
	"reflect"
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

func tickersOf(result *resource.PurchaseBalance) []string {
	tickers := make([]string, 0)
	for _, p := range result.SecuritiesBalance {
		tickers = append(tickers, p.Security.Ticker)
	}
	return tickers
}

func countOf(result *resource.PurchaseBalance, ticker string) int {
	for _, p := range result.SecuritiesBalance {
		if p.Security.Ticker == ticker {
			return p.Count
		}
	}
	return 0
}

func TestPurchaseBalance_UnevenDistribution_HighPrice(t *testing.T) {
	// GIVEN
	securities := []*resource.Security{
		{Ticker: "High", Price: 600},
		{Ticker: "Low", Price: 100},
	}
	svc := service.NewPurchaseBalanceService(&mockStockService{}, &mockReitService{})

	// WHEN
	result := svc.PurchaseBalance(securities, 1000)

	// THEN
	if !reflect.DeepEqual(tickersOf(result), []string{"Low", "High"}) {
		t.Errorf("Expected order [Low High], got %v", tickersOf(result))
	}
	if countOf(result, "Low") != 4 || countOf(result, "High") != 1 {
		t.Errorf("Expected Low=4 and High=1, got Low=%d High=%d", countOf(result, "Low"), countOf(result, "High"))
	}
	if result.RemainingBalance() != 0 {
		t.Errorf("Expected remaining balance 0, got %.2f", result.RemainingBalance())
	}
}

func TestPurchaseBalance_ExpensiveDoesNotFit(t *testing.T) {
	// GIVEN
	securities := []*resource.Security{
		{Ticker: "High", Price: 900},
		{Ticker: "Low", Price: 100},
	}
	svc := service.NewPurchaseBalanceService(&mockStockService{}, &mockReitService{})

	// WHEN
	result := svc.PurchaseBalance(securities, 1000)

	// THEN
	if !reflect.DeepEqual(tickersOf(result), []string{"Low"}) {
		t.Errorf("Expected only Low, got %v", tickersOf(result))
	}
	if countOf(result, "Low") != 10 {
		t.Errorf("Expected 10 Low, got %d", countOf(result, "Low"))
	}
	if result.RemainingBalance() != 0 {
		t.Errorf("Expected remaining balance 0, got %.2f", result.RemainingBalance())
	}
}

func TestPurchaseBalance_ZeroPriceIgnored(t *testing.T) {
	// GIVEN
	securities := []*resource.Security{
		{Ticker: "ZERO", Price: 0},
		nil,
		{Ticker: "A", Price: 10},
	}
	svc := service.NewPurchaseBalanceService(&mockStockService{}, &mockReitService{})

	// WHEN
	result := svc.PurchaseBalance(securities, 100)

	// THEN
	if !reflect.DeepEqual(tickersOf(result), []string{"A"}) {
		t.Errorf("Expected only A, got %v", tickersOf(result))
	}
	if countOf(result, "A") != 10 {
		t.Errorf("Expected 10 A, got %d", countOf(result, "A"))
	}
	if result.RemainingBalance() != 0 {
		t.Errorf("Expected remaining balance 0, got %.2f", result.RemainingBalance())
	}
}

func TestPurchaseBalance_EmptyList(t *testing.T) {
	// GIVEN
	svc := service.NewPurchaseBalanceService(&mockStockService{}, &mockReitService{})

	// WHEN
	result := svc.PurchaseBalance(nil, 100)

	// THEN
	if len(result.SecuritiesBalance) != 0 {
		t.Errorf("Expected no securities, got %v", tickersOf(result))
	}
	if result.RemainingBalance() != 100 {
		t.Errorf("Expected remaining balance 100, got %.2f", result.RemainingBalance())
	}
}

func TestPurchaseBalance_RemainingBalanceProperty(t *testing.T) {
	cases := []struct {
		prices []float64
		amount float64
	}{
		{[]float64{10, 10}, 100},
		{[]float64{600, 100}, 1000},
		{[]float64{900, 100}, 1000},
		{[]float64{37.5, 12.2, 99.9}, 1000},
		{[]float64{5.25, 7.75, 130, 48.1}, 2500.55},
		{[]float64{1000, 2000, 3000}, 2500},
		{[]float64{0.5, 333.33, 71.4}, 777.77},
		{[]float64{25}, 100},
		{[]float64{3, 0, 7, -1}, 50},
		{[]float64{150.9, 151.1}, 300},
	}
	svc := service.NewPurchaseBalanceService(&mockStockService{}, &mockReitService{})
	for _, c := range cases {
		// GIVEN
		securities := make([]*resource.Security, 0)
		for i, price := range c.prices {
			securities = append(securities, &resource.Security{Ticker: fmt.Sprintf("S%d", i), Price: price})
		}

		// WHEN
		result := svc.PurchaseBalance(securities, c.amount)

		// THEN
		remaining := result.RemainingBalance()
		if remaining < -1e-9 {
			t.Errorf("prices %v amount %.2f: negative remaining balance %.4f", c.prices, c.amount, remaining)
		}
		for _, security := range securities {
			if security.Price > 0 && remaining >= security.Price {
				t.Errorf("prices %v amount %.2f: remaining balance %.4f is not lower than price %.2f", c.prices, c.amount, remaining, security.Price)
			}
		}
	}
}

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

func TestPurchaseBalancesBySecurities_TickerWithoutQuote(t *testing.T) {
	// GIVEN
	mockStockSvc := &mockStockService{
		ListStocksByTickersFunc: func(tickers []string) []*resource.Security { return nil },
	}
	mockReitSvc := &mockReitService{
		ListReitsByTickersFunc: func(tickers []string) []*resource.Security {
			return []*resource.Security{{Ticker: "HGLG11", Price: 100}}
		},
	}
	svc := service.NewPurchaseBalanceService(mockStockSvc, mockReitSvc)

	// WHEN
	result := svc.PurchaseBalancesBySecurities([]string{"XXXX0"}, []string{"HGLG11"}, 1000)

	// THEN
	if !reflect.DeepEqual(tickersOf(result), []string{"HGLG11"}) {
		t.Errorf("Expected only HGLG11, got %v", tickersOf(result))
	}
	if countOf(result, "HGLG11") != 10 {
		t.Errorf("Expected 10 HGLG11, got %d", countOf(result, "HGLG11"))
	}
}
