package service_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"trader/internal/config"
	"trader/internal/resource"
	"trader/internal/scraping"
	"trader/internal/service"
)

// Mocks

type mockStockService struct {
	GetStockByTickerFunc    func(string) *resource.Security
	ListStocksByTickersFunc func([]string) ([]*resource.Security, []*resource.TickerFailure)
}

func (m *mockStockService) GetStockByTicker(ticker string) *resource.Security {
	return m.GetStockByTickerFunc(ticker)
}

func (m *mockStockService) ListStocksByTickers(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
	return m.ListStocksByTickersFunc(tickers)
}

type mockReitService struct {
	GetReitByTickerFunc    func(string) *resource.Security
	ListReitsByTickersFunc func([]string) ([]*resource.Security, []*resource.TickerFailure)
}

func (m *mockReitService) GetReitByTicker(ticker string) *resource.Security {
	return m.GetReitByTickerFunc(ticker)
}

func (m *mockReitService) ListReitsByTickers(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
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
		ListStocksByTickersFunc: func(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
			if len(tickers) == 1 && tickers[0] == "PETR4" {
				return mockStocks, nil
			}
			return nil, nil
		},
	}
	mockReitSvc := &mockReitService{
		ListReitsByTickersFunc: func(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
			if len(tickers) == 1 && tickers[0] == "HGLG11" {
				return mockReits, nil
			}
			return nil, nil
		},
	}

	svc := service.NewPurchaseBalanceService(mockStockSvc, mockReitSvc)

	// WHEN
	result, err := svc.PurchaseBalancesBySecurities(stockTickers, reitTickers, amountInvested)

	// THEN
	// 2 securities. 500 each.
	// PETR4 (50): 10 shares.
	// HGLG11 (100): 5 shares.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.TotalCount() != 15 {
		t.Errorf("Expected 15 total, got %d", result.TotalCount())
	}
	if countOf(result, "PETR4") != 10 || countOf(result, "HGLG11") != 5 {
		t.Errorf("Expected PETR4=10 and HGLG11=5, got PETR4=%d HGLG11=%d", countOf(result, "PETR4"), countOf(result, "HGLG11"))
	}
}

func TestPurchaseBalancesBySecurities_StocksBeforeReits(t *testing.T) {
	// GIVEN
	mockStockSvc := &mockStockService{
		ListStocksByTickersFunc: func(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
			return []*resource.Security{{Ticker: "S1", Price: 10}, {Ticker: "S2", Price: 10}}, nil
		},
	}
	mockReitSvc := &mockReitService{
		ListReitsByTickersFunc: func(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
			return []*resource.Security{{Ticker: "R1", Price: 10}, {Ticker: "R2", Price: 10}}, nil
		},
	}
	svc := service.NewPurchaseBalanceService(mockStockSvc, mockReitSvc)

	// WHEN
	result, err := svc.PurchaseBalancesBySecurities([]string{"S1", "S2"}, []string{"R1", "R2"}, 400)

	// THEN
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(tickersOf(result), []string{"S1", "S2", "R1", "R2"}) {
		t.Errorf("Expected [S1 S2 R1 R2], got %v", tickersOf(result))
	}
}

func failingList(securities []*resource.Security, failures ...*resource.TickerFailure) func([]string) ([]*resource.Security, []*resource.TickerFailure) {
	return func([]string) ([]*resource.Security, []*resource.TickerFailure) { return securities, failures }
}

func TestPurchaseBalancesBySecurities_TickerWithoutQuote(t *testing.T) {
	// GIVEN
	failure := &resource.TickerFailure{Ticker: "XXXX0", Type: resource.STOCK_TYPE, Err: errors.New(`status="404"`)}
	svc := service.NewPurchaseBalanceService(
		&mockStockService{ListStocksByTickersFunc: failingList([]*resource.Security{{Ticker: "PETR4", Price: 50}}, failure)},
		&mockReitService{ListReitsByTickersFunc: failingList(nil)},
	)

	// WHEN
	result, err := svc.PurchaseBalancesBySecurities([]string{"PETR4", "XXXX0"}, []string{}, 1000)

	// THEN
	if result != nil {
		t.Errorf("Expected no balance, got %v", tickersOf(result))
	}
	var quoteErr *resource.QuoteFailuresError
	if !errors.As(err, &quoteErr) {
		t.Fatalf("Expected a QuoteFailuresError, got %v", err)
	}
	if !reflect.DeepEqual(quoteErr.Failures, []*resource.TickerFailure{failure}) {
		t.Errorf("Expected the XXXX0 failure, got %v", quoteErr.Failures)
	}
	if !strings.Contains(err.Error(), `XXXX0: status="404"`) {
		t.Errorf("Expected the error to name XXXX0 and its reason, got %q", err.Error())
	}
}

func TestPurchaseBalancesBySecurities_FailuresOnBothSides(t *testing.T) {
	// GIVEN
	stockFailure := &resource.TickerFailure{Ticker: "XXXX0", Type: resource.STOCK_TYPE, Err: errors.New(`status="404"`)}
	reitFailure := &resource.TickerFailure{Ticker: "XXXX00", Type: resource.REIT_TYPE, Err: errors.New("ticker not found")}
	svc := service.NewPurchaseBalanceService(
		&mockStockService{ListStocksByTickersFunc: failingList(nil, stockFailure)},
		&mockReitService{ListReitsByTickersFunc: failingList([]*resource.Security{{Ticker: "MXRF11", Price: 10}}, reitFailure)},
	)

	// WHEN
	result, err := svc.PurchaseBalancesBySecurities([]string{"XXXX0"}, []string{"MXRF11", "XXXX00"}, 1000)

	// THEN
	if result != nil {
		t.Errorf("Expected no balance, got %v", tickersOf(result))
	}
	var quoteErr *resource.QuoteFailuresError
	if !errors.As(err, &quoteErr) {
		t.Fatalf("Expected a QuoteFailuresError, got %v", err)
	}
	if !reflect.DeepEqual(quoteErr.Failures, []*resource.TickerFailure{stockFailure, reitFailure}) {
		t.Errorf("Expected the stock failure then the REIT failure, got %v", quoteErr.Failures)
	}
	for _, expected := range []string{`XXXX0: status="404"`, "XXXX00: ticker not found"} {
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected the error to contain %q, got %q", expected, err.Error())
		}
	}
}

// countingTransport answers every request with a page priced at 10 and records the peak number
// of requests in flight. The first requests wait until 4 are in flight, or for a short timeout.
type countingTransport struct {
	current atomic.Int32
	max     atomic.Int32
	barrier chan struct{}
	once    sync.Once
}

func (c *countingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	n := c.current.Add(1)
	defer c.current.Add(-1)
	for {
		peak := c.max.Load()
		if n <= peak || c.max.CompareAndSwap(peak, n) {
			break
		}
	}
	if n >= 4 {
		c.once.Do(func() { close(c.barrier) })
	}
	select {
	case <-c.barrier:
	case <-time.After(2 * time.Second):
	}
	page := `<html><body><div title="Valor atual do ativo"><strong>10,00</strong></div></body></html>`
	return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(page)), Request: r}, nil
}

func TestPurchaseBalancesBySecurities_SharedRequestLimit(t *testing.T) {
	// GIVEN real scrapers sharing one fetcher
	transport := &countingTransport{barrier: make(chan struct{})}
	cfg := config.NewConfig()
	fetcher := scraping.NewFetcher(cfg, scraping.WithTransport(transport))
	svc := service.NewPurchaseBalanceService(
		service.NewStockService(scraping.NewStockScraping(cfg, fetcher)),
		service.NewReitService(scraping.NewReitScraping(cfg, fetcher)),
	)

	// WHEN
	result, err := svc.PurchaseBalancesBySecurities([]string{"S1", "S2", "S3"}, []string{"R1", "R2", "R3"}, 600)

	// THEN
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(tickersOf(result), []string{"S1", "S2", "S3", "R1", "R2", "R3"}) {
		t.Errorf("Expected stocks before REITs in input order, got %v", tickersOf(result))
	}
	if transport.max.Load() != 4 {
		t.Errorf("Expected at most and at least 4 requests in flight, got %d", transport.max.Load())
	}
}
