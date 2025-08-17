package service_test

import (
	"fmt"
	"testing"
	"trader/internal/resource"
	"trader/internal/service"
)

type reitScrapingMock struct {
	ReturnValue map[string]interface{}
	ReturnError error
}

func (m *reitScrapingMock) GetReitByTicker(ticket string) (*resource.Security, error) {
	if m.ReturnError != nil {
		return nil, m.ReturnError
	}
	return m.ReturnValue["GetReitByTicker"].(*resource.Security), nil
}

func (m *reitScrapingMock) ListReitsByTickers(tickers []string) []*resource.Security {
	return m.ReturnValue["ListReitsByTickers"].([]*resource.Security)
}

func newReitScrapingMock() *reitScrapingMock {
	securityUnit := &resource.Security{
		Ticker: "BPML11",
		Name:   "BPML11",
		Currency: &resource.Currency{
			Code:        "BRL",
			Description: "Real Brasileiro",
			Sign:        "R$",
		},
	}
	return &reitScrapingMock{
		ReturnValue: map[string]interface{}{
			"GetReitByTicker": securityUnit,
			"ListReitsByTickers": []*resource.Security{
				securityUnit,
			},
		},
		ReturnError: nil,
	}
}

func TestGetReitFound(t *testing.T) {
	// GIVEN
	ticker := "BPML11"

	// MOCK
	mock := newReitScrapingMock()

	// THEN
	reitService := service.NewReitService(mock)
	result := reitService.GetReitByTicker(ticker)

	// WHEN
	if result.Ticker != ticker {
		t.Errorf(`expected at "%s" and received at %s`, ticker, result.Ticker)
	}
}

func TestGetReitNotFound(t *testing.T) {
	// GIVEN
	ticker := "XPTO"

	// MOCK
	mock := newReitScrapingMock()
	mock.ReturnValue["GetReitByTicker"] = nil
	mock.ReturnError = fmt.Errorf("mock error")

	// THEN
	reitService := service.NewReitService(mock)
	result := reitService.GetReitByTicker(ticker)

	// WHEN
	if result != nil {
		t.Errorf(`expected at "%s" and received at %v`, "nil", result)
	}
}

func TestListReitFound(t *testing.T) {
	// GIVEN
	tickers := []string{"BPML11"}

	// MOCK
	mock := newReitScrapingMock()

	// THEN
	reitService := service.NewReitService(mock)
	result := reitService.ListReitsByTickers(tickers)

	// WHEN
	if len(result) < 1 {
		t.Errorf(`expected at "%v" and received at %v`, tickers, result)
	}
}
