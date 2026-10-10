package resource_test

import (
	"errors"
	"testing"
	"trader/internal/resource"
)

func TestQuoteFailuresErrorMessage(t *testing.T) {
	// GIVEN
	err := &resource.QuoteFailuresError{Failures: []*resource.TickerFailure{
		{Ticker: "XXXX0", Type: resource.STOCK_TYPE, Err: errors.New(`status="404"`)},
		{Ticker: "XXXX00", Type: resource.REIT_TYPE, Err: errors.New("ticker not found")},
	}}

	// WHEN
	message := err.Error()

	// THEN
	expected := "could not get the quote of 2 ticker(s):\nXXXX0: status=\"404\"\nXXXX00: ticker not found"
	if message != expected {
		t.Errorf("expected %q and received %q", expected, message)
	}
}
