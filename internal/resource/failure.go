package resource

import (
	"fmt"
	"strings"
)

// TickerFailure tells which ticker could not be quoted and why.
type TickerFailure struct {
	Ticker string
	Type   string // Stock, REIT
	Err    error
}

// QuoteFailuresError is returned when the quote of one or more tickers cannot be obtained.
type QuoteFailuresError struct {
	Failures []*TickerFailure
}

func (e *QuoteFailuresError) Error() string {
	lines := make([]string, 0, len(e.Failures)+1)
	lines = append(lines, fmt.Sprintf("could not get the quote of %d ticker(s):", len(e.Failures)))
	for _, failure := range e.Failures {
		lines = append(lines, fmt.Sprintf("%s: %v", failure.Ticker, failure.Err))
	}
	return strings.Join(lines, "\n")
}
