package service_test

import (
	"fmt"
	"testing"
	"trader/internal/resource"
	"trader/internal/service"
)

type stockScrapingMock struct {
	ReturnValue map[string]interface{}
	ReturnError error
}

func (m *stockScrapingMock) GetStockByTicker(ticket string) (*resource.Security, error) {
	if m.ReturnError != nil {
		return nil, m.ReturnError
	}
	return m.ReturnValue["GetStockByTicker"].(*resource.Security), nil
}

func (m *stockScrapingMock) ListStocksByTickers(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
	failures, _ := m.ReturnValue["ListStocksByTickersFailures"].([]*resource.TickerFailure)
	return m.ReturnValue["ListStocksByTickers"].([]*resource.Security), failures
}

func newStockScrapingMock() *stockScrapingMock {
	securityUnit := &resource.Security{
		Ticker: "PETR4",
		Name:   "PETR4",
		Currency: &resource.Currency{
			Code:        "BRL",
			Description: "Real Brasileiro",
			Sign:        "R$",
		},
	}
	return &stockScrapingMock{
		ReturnValue: map[string]interface{}{
			"GetStockByTicker": securityUnit,
			"ListStocksByTickers": []*resource.Security{
				securityUnit,
			},
		},
		ReturnError: nil,
	}
}

func TestGetStockFound(t *testing.T) {
	// GIVEN
	ticker := "PETR4"

	// MOCK
	mock := newStockScrapingMock()

	// THEN
	stockService := service.NewStockService(mock)
	result := stockService.GetStockByTicker(ticker)

	// WHEN
	if result.Ticker != ticker {
		t.Errorf(`expected at "%s" and received at %s`, ticker, result.Ticker)
	}
}

func TestGetStockNotFound(t *testing.T) {
	// GIVEN
	ticker := "XPTO"

	// MOCK
	mock := newStockScrapingMock()
	mock.ReturnValue["GetStockByTicker"] = nil
	mock.ReturnError = fmt.Errorf("mock error")

	// THEN
	stockService := service.NewStockService(mock)
	result := stockService.GetStockByTicker(ticker)

	// WHEN
	if result != nil {
		t.Errorf(`expected at "%s" and received at %v`, "nil", result)
	}
}

func TestListStockFound(t *testing.T) {
	// GIVEN
	tickers := []string{"PETR4"}

	// MOCK
	mock := newStockScrapingMock()

	// THEN
	stockService := service.NewStockService(mock)
	result, failures := stockService.ListStocksByTickers(tickers)

	// WHEN
	if len(result) < 1 {
		t.Errorf(`expected at "%v" and received at %v`, tickers, result)
	}
	if len(failures) != 0 {
		t.Errorf("expected no failures and received %v", failures)
	}
}

func TestListStockWithFailures(t *testing.T) {
	// GIVEN
	tickers := []string{"PETR4", "XXXX0"}

	// MOCK
	mock := newStockScrapingMock()
	expected := []*resource.TickerFailure{{Ticker: "XXXX0", Type: resource.STOCK_TYPE, Err: fmt.Errorf("mock error")}}
	mock.ReturnValue["ListStocksByTickersFailures"] = expected

	// THEN
	stockService := service.NewStockService(mock)
	result, failures := stockService.ListStocksByTickers(tickers)

	// WHEN
	if len(result) != 1 || result[0].Ticker != "PETR4" {
		t.Errorf("expected [PETR4] and received %v", result)
	}
	if len(failures) != 1 || failures[0] != expected[0] {
		t.Errorf("expected the scraping failures and received %v", failures)
	}
}
