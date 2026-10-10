package cmd_test

import (
	"strings"
	"testing"
	"trader/internal/resource"
)

func TestStockListRemovesEscapeSequences(t *testing.T) {
	// GIVEN
	f := newFakes()
	stock := newStock("ITSA3", 10.5)
	stock.Name = "ITAU\x1b[2JSA"
	f.stock.list = listing(stock)
	tr := newTree(f)

	// WHEN
	err := tr.run("stock", "list", "ITSA3", "--no-color")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if strings.Contains(tr.stdout.String(), "\x1b") {
		t.Errorf("expected no ESC character and received %q", tr.stdout.String())
	}
	if !strings.Contains(tr.stdout.String(), "ITAU[2JSA") {
		t.Errorf("expected the cleaned name and received %q", tr.stdout.String())
	}
}

func TestReitGetRemovesControlCharacters(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.reit.get = func(ticker string) *resource.Security {
		reit := newReit(ticker, 10)
		reit.Admin = "ADMIN\a\x1b]0;title\a"
		reit.Segment = "SEG\u009bMENT"
		return reit
	}
	tr := newTree(f)

	// WHEN
	err := tr.run("reit", "get", "MXRF11", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if !strings.Contains(tr.stdout.String(), "Admin,ADMIN]0;title\n") || !strings.Contains(tr.stdout.String(), "Segment,SEGMENT\n") {
		t.Errorf("expected cleaned cells and received %q", tr.stdout.String())
	}
}

func formulaStock(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
	stock := newStock("ITSA3", 10.5)
	stock.Name = `=HYPERLINK("http://x")`
	return []*resource.Security{stock}, nil
}

func TestStockListCSVNeutralizesFormulas(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.stock.list = formulaStock
	tr := newTree(f)

	// WHEN
	err := tr.run("stock", "list", "ITSA3", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	expected := `TICKER,NAME,DOCUMENT,PRICE,CURRENCY,CAPTURED AT
ITSA3,"'=HYPERLINK(""http://x"")",11.111.111/0001-11,10.50,R$ BRL,2025-06-08 22:59:37
`
	if tr.stdout.String() != expected {
		t.Errorf("expected stdout %q and received %q", expected, tr.stdout.String())
	}
}

func TestStockListTableKeepsFormulaText(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.stock.list = formulaStock
	tr := newTree(f)

	// WHEN
	err := tr.run("stock", "list", "ITSA3", "--no-color")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if !strings.Contains(tr.stdout.String(), ` =HYPERLINK("http://x")`) || strings.Contains(tr.stdout.String(), "'=") {
		t.Errorf("expected the name without a prefix and received %q", tr.stdout.String())
	}
}
