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

func (m *stockScrapingMock) ListStocksByTickers(tickers []string) []*resource.Security {
	return m.ReturnValue["ListStocksByTickers"].([]*resource.Security)
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
	result := stockService.ListStocksByTickers(tickers)

	// WHEN
	if len(result) < 1 {
		t.Errorf(`expected at "%v" and received at %v`, tickers, result)
	}
}
