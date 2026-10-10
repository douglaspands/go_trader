package cmd_test

import (
	"errors"
	"strings"
	"testing"
	"trader/internal/resource"
)

func failure(ticker string, securityType string, reason string) *resource.TickerFailure {
	return &resource.TickerFailure{Ticker: ticker, Type: securityType, Err: errors.New(reason)}
}

func TestStockListWithOneFailure(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.stock.list = func(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
		return []*resource.Security{newStock("ITSA3", 10.5)}, []*resource.TickerFailure{failure("XXXX0", resource.STOCK_TYPE, `status="404"`)}
	}
	tr := newTree(f)

	// WHEN
	err := tr.run("stock", "list", "ITSA3", "XXXX0", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	expected := `TICKER,NAME,DOCUMENT,PRICE,CURRENCY,CAPTURED AT
ITSA3,Name ITSA3,11.111.111/0001-11,10.50,R$ BRL,2025-06-08 22:59:37
`
	if tr.stdout.String() != expected {
		t.Errorf("expected stdout %q and received %q", expected, tr.stdout.String())
	}
	if tr.stderr.String() != "Warning: XXXX0: status=\"404\"\n" {
		t.Errorf("unexpected stderr %q", tr.stderr.String())
	}
}

func TestReitListWithFailures(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.reit.list = func(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
		return []*resource.Security{newReit("MXRF11", 10)}, []*resource.TickerFailure{
			failure("XXXX00", resource.REIT_TYPE, "ticker not found"),
			failure("YYYY11", resource.REIT_TYPE, "timeout\x1b[2J"),
		}
	}
	tr := newTree(f)

	// WHEN
	err := tr.run("reit", "list", "MXRF11", "XXXX00", "YYYY11", "--no-color")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if !strings.Contains(tr.stdout.String(), "MXRF11") {
		t.Errorf("expected the table in stdout and received %q", tr.stdout.String())
	}
	expected := "Warning: XXXX00: ticker not found\nWarning: YYYY11: timeout[2J\n"
	if tr.stderr.String() != expected {
		t.Errorf("expected stderr %q and received %q", expected, tr.stderr.String())
	}
}

func quoteFailures(failures ...*resource.TickerFailure) error {
	return &resource.QuoteFailuresError{Failures: failures}
}

func TestPurchaseBalanceQuoteFailures(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		err      error
		expected string
	}{
		{
			"stock with one failure",
			[]string{"stock", "purchase-balance", "ITSA3", "XXXX0", "--amount", "1000"},
			quoteFailures(failure("XXXX0", resource.STOCK_TYPE, `status="404"`)),
			"Error: could not get the quote of 1 ticker(s):\nXXXX0: status=\"404\"\n",
		},
		{
			"reit with one failure",
			[]string{"reit", "purchase-balance", "MXRF11", "XXXX00", "--amount", "1000"},
			quoteFailures(failure("XXXX00", resource.REIT_TYPE, "ticker not found")),
			"Error: could not get the quote of 1 ticker(s):\nXXXX00: ticker not found\n",
		},
		{
			"stocks and reits",
			[]string{"security", "purchase-balance", "-s", "XXXX0", "-r", "MXRF11,XXXX00", "-a", "1000"},
			quoteFailures(failure("XXXX0", resource.STOCK_TYPE, `status="404"`), failure("XXXX00", resource.REIT_TYPE, "ticker not found\a")),
			"Error: could not get the quote of 2 ticker(s):\nXXXX0: status=\"404\"\nXXXX00: ticker not found\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// GIVEN
			f := newFakes()
			f.purchase.err = c.err
			tr := newTree(f)

			// WHEN
			err := tr.run(c.args...)

			// THEN
			if err == nil {
				t.Fatal("expected an error, received nil")
			}
			if tr.stdout.Len() != 0 {
				t.Errorf("expected empty stdout and received %q", tr.stdout.String())
			}
			if tr.stderr.String() != c.expected {
				t.Errorf("expected stderr %q and received %q", c.expected, tr.stderr.String())
			}
		})
	}
}
