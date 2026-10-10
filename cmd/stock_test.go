package cmd_test

import (
	"reflect"
	"strings"
	"testing"
	"trader/internal/resource"
)

func stockPurchaseResult() *resource.PurchaseBalance {
	return &resource.PurchaseBalance{
		SecuritiesBalance: []*resource.SecurityPurchase{
			{Security: newStock("ITSA3", 10.5), Count: 5},
			{Security: newStock("BBDC3", 20), Count: 2},
		},
		AmountInvested: 100,
	}
}

func TestStockGetFound(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.stock.get = func(ticker string) *resource.Security { return newStock(ticker, 10.5) }
	tr := newTree(f)

	// WHEN
	err := tr.run("stock", "get", "ITSA3", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	expected := `FIELD,VALUE
Ticker,ITSA3
Name,Name ITSA3
Document,11.111.111/0001-11
Currency,R$ BRL
Price,10.50
CapturedAt,2025-06-08 22:59:37
Origin,https://example.com/acoes/ITSA3
`
	if tr.stdout.String() != expected {
		t.Errorf("expected stdout %q and received %q", expected, tr.stdout.String())
	}
	if tr.stderr.Len() != 0 {
		t.Errorf("expected empty stderr and received %q", tr.stderr.String())
	}
}

func TestStockGetNotFound(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run("stock", "get", "XXXX0")

	// THEN
	if err == nil {
		t.Fatal("expected an error, received nil")
	}
	if tr.stderr.String() != "Error: ticker \"XXXX0\" not found!\n" {
		t.Errorf("unexpected stderr %q", tr.stderr.String())
	}
	if tr.stdout.Len() != 0 {
		t.Errorf("expected empty stdout and received %q", tr.stdout.String())
	}
}

func TestStockGetWithoutArgument(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run("stock", "get")

	// THEN
	if err == nil {
		t.Fatal("expected an error, received nil")
	}
	if !strings.Contains(tr.stderr.String(), "Error: accepts 1 arg(s), received 0") {
		t.Errorf("unexpected stderr %q", tr.stderr.String())
	}
}

func TestStockListFound(t *testing.T) {
	// GIVEN
	f := newFakes()
	var received []string
	f.stock.list = func(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
		received = tickers
		return []*resource.Security{newStock("ITSA3", 10.5), newStock("BBDC3", 1234.5)}, nil
	}
	tr := newTree(f)

	// WHEN
	err := tr.run("stock", "list", "ITSA3", "BBDC3", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if !reflect.DeepEqual(received, []string{"ITSA3", "BBDC3"}) {
		t.Errorf("unexpected tickers %v", received)
	}
	expected := `TICKER,NAME,DOCUMENT,PRICE,CURRENCY,CAPTURED AT
ITSA3,Name ITSA3,11.111.111/0001-11,10.50,R$ BRL,2025-06-08 22:59:37
BBDC3,Name BBDC3,11.111.111/0001-11,1234.50,R$ BRL,2025-06-08 22:59:37
`
	if tr.stdout.String() != expected {
		t.Errorf("expected stdout %q and received %q", expected, tr.stdout.String())
	}
	if tr.stderr.Len() != 0 {
		t.Errorf("expected empty stderr and received %q", tr.stderr.String())
	}
}

func TestStockListNotFound(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run("stock", "list", "XXXX0")

	// THEN
	if err == nil {
		t.Fatal("expected an error, received nil")
	}
	if tr.stderr.String() != "Error: tickers not found!\n" {
		t.Errorf("unexpected stderr %q", tr.stderr.String())
	}
	if tr.stdout.Len() != 0 {
		t.Errorf("expected empty stdout and received %q", tr.stdout.String())
	}
}

func TestStockPurchaseBalance(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.purchase.result = stockPurchaseResult()
	tr := newTree(f)

	// WHEN
	err := tr.run("stock", "purchase-balance", "ITSA3", "BBDC3", "--amount", "100", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	expectedCalls := []purchaseCall{{stocks: []string{"ITSA3", "BBDC3"}, reits: []string{}, amount: 100}}
	if !reflect.DeepEqual(f.purchase.calls, expectedCalls) {
		t.Errorf("expected calls %v and received %v", expectedCalls, f.purchase.calls)
	}
	expected := `TICKER,PRICE,COUNT,TOTAL,CURRENCY,CAPTURED AT
ITSA3,10.50,5,52.50,R$ BRL,2025-06-08 22:59:37
BBDC3,20.00,2,40.00,R$ BRL,2025-06-08 22:59:37
,,7,92.50,R$ BRL,SPENT AMOUNT
,,,7.50,R$ BRL,REMAINING AMOUNT
`
	if tr.stdout.String() != expected {
		t.Errorf("expected stdout %q and received %q", expected, tr.stdout.String())
	}
	if tr.stderr.Len() != 0 {
		t.Errorf("expected empty stderr and received %q", tr.stderr.String())
	}
}

func TestStockPurchaseBalanceNotFound(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run("stock", "purchase-balance", "ITSA3", "-a", "1")

	// THEN
	if err == nil {
		t.Fatal("expected an error, received nil")
	}
	if tr.stderr.String() != "Error: tickers not found!\n" {
		t.Errorf("unexpected stderr %q", tr.stderr.String())
	}
	if tr.stdout.Len() != 0 {
		t.Errorf("expected empty stdout and received %q", tr.stdout.String())
	}
}

func TestStockPurchaseBalanceWithoutAmount(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.purchase.result = stockPurchaseResult()
	tr := newTree(f)

	// WHEN
	err := tr.run("stock", "purchase-balance", "ITSA3")

	// THEN
	if err == nil {
		t.Fatal("expected an error, received nil")
	}
	if !strings.Contains(tr.stderr.String(), `Error: required flag(s) "amount" not set`) {
		t.Errorf("unexpected stderr %q", tr.stderr.String())
	}
	if len(f.purchase.calls) != 0 {
		t.Errorf("expected no purchase balance call and received %v", f.purchase.calls)
	}
}

func TestStockNoColor(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.stock.list = listing(newStock("ITSA3", 10.5))
	tr := newTree(f)

	// WHEN
	colored := tr.run("stock", "list", "ITSA3")
	coloredOut := tr.stdout.String()
	plain := tr.run("stock", "list", "ITSA3", "--no-color")
	plainOut := tr.stdout.String()

	// THEN
	if colored != nil || plain != nil {
		t.Fatalf("expected no errors, received %v and %v", colored, plain)
	}
	if !strings.Contains(coloredOut, "\x1b[") {
		t.Errorf("expected ANSI codes in the default output and received %q", coloredOut)
	}
	if strings.Contains(plainOut, "\x1b[") {
		t.Errorf("expected no ANSI codes with --no-color and received %q", plainOut)
	}
	if !strings.Contains(plainOut, "ITSA3") {
		t.Errorf("expected the table in the output and received %q", plainOut)
	}
}

func TestOptionDoesNotLeakToAnotherCommand(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.stock.list = listing(newStock("ITSA3", 10.5))
	f.reit.list = listing(newReit("MXRF11", 10))
	tr := newTree(f)

	// WHEN
	first := tr.run("stock", "list", "ITSA3", "--csv")
	firstOut := tr.stdout.String()
	second := tr.run("reit", "list", "MXRF11")
	secondOut := tr.stdout.String()

	// THEN
	if first != nil || second != nil {
		t.Fatalf("expected no errors, received %v and %v", first, second)
	}
	if !strings.HasPrefix(firstOut, "TICKER,NAME") {
		t.Errorf("expected csv in the first run and received %q", firstOut)
	}
	if strings.Contains(secondOut, "TICKER,NAME") || !strings.Contains(secondOut, "MXRF11") {
		t.Errorf("expected a table in the second run and received %q", secondOut)
	}
}

func TestStockListFormatting(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.stock.list = listing(newStock("ITSA3", 10.5))
	tr := newTree(f)

	// WHEN
	err := tr.run("stock", "list", "ITSA3", "--no-color")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	for _, value := range []string{"10.50", "2025-06-08 22:59:37", "R$ BRL"} {
		if !strings.Contains(tr.stdout.String(), value) {
			t.Errorf("expected %q in the table and received %q", value, tr.stdout.String())
		}
	}
}
