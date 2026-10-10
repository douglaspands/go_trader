package cmd_test

import (
	"reflect"
	"strings"
	"testing"
	"trader/internal/resource"
)

func reitPurchaseResult() *resource.PurchaseBalance {
	return &resource.PurchaseBalance{
		SecuritiesBalance: []*resource.SecurityPurchase{
			{Security: newReit("MXRF11", 10), Count: 5},
			{Security: newReit("HGLG11", 100), Count: 1},
		},
		AmountInvested: 200,
	}
}

func TestReitGetFound(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.reit.get = func(ticker string) *resource.Security { return newReit(ticker, 10.5) }
	tr := newTree(f)

	// WHEN
	err := tr.run("reit", "get", "MXRF11", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	expected := `FIELD,VALUE
Ticker,MXRF11
Name,Name MXRF11
Admin,Admin MXRF11
Document,22.222.222/0001-22
Segment,Segment MXRF11
Currency,R$ BRL
Price,10.50
CapturedAt,2025-06-08 22:59:37
Origin,https://example.com/fundos-imobiliarios/MXRF11
`
	if tr.stdout.String() != expected {
		t.Errorf("expected stdout %q and received %q", expected, tr.stdout.String())
	}
	if tr.stderr.Len() != 0 {
		t.Errorf("expected empty stderr and received %q", tr.stderr.String())
	}
}

func TestReitGetNotFound(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run("reit", "get", "XXXX00")

	// THEN
	if err == nil {
		t.Fatal("expected an error, received nil")
	}
	if tr.stderr.String() != "Error: ticker \"XXXX00\" not found!\n" {
		t.Errorf("unexpected stderr %q", tr.stderr.String())
	}
	if tr.stdout.Len() != 0 {
		t.Errorf("expected empty stdout and received %q", tr.stdout.String())
	}
}

func TestReitGetWithoutArgument(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run("reit", "get")

	// THEN
	if err == nil {
		t.Fatal("expected an error, received nil")
	}
	if !strings.Contains(tr.stderr.String(), "Error: accepts 1 arg(s), received 0") {
		t.Errorf("unexpected stderr %q", tr.stderr.String())
	}
}

func TestReitListFound(t *testing.T) {
	// GIVEN
	f := newFakes()
	var received []string
	f.reit.list = func(tickers []string) ([]*resource.Security, []*resource.TickerFailure) {
		received = tickers
		return []*resource.Security{newReit("MXRF11", 10.5), newReit("HGLG11", 1050.75)}, nil
	}
	tr := newTree(f)

	// WHEN
	err := tr.run("reit", "list", "MXRF11", "HGLG11", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if !reflect.DeepEqual(received, []string{"MXRF11", "HGLG11"}) {
		t.Errorf("unexpected tickers %v", received)
	}
	expected := `TICKER,NAME,DOCUMENT,PRICE,CURRENCY,CAPTURED AT
MXRF11,Name MXRF11,22.222.222/0001-22,10.50,R$ BRL,2025-06-08 22:59:37
HGLG11,Name HGLG11,22.222.222/0001-22,1050.75,R$ BRL,2025-06-08 22:59:37
`
	if tr.stdout.String() != expected {
		t.Errorf("expected stdout %q and received %q", expected, tr.stdout.String())
	}
	if tr.stderr.Len() != 0 {
		t.Errorf("expected empty stderr and received %q", tr.stderr.String())
	}
}

func TestReitListNotFound(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run("reit", "list", "XXXX00")

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

func TestReitPurchaseBalance(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.purchase.result = reitPurchaseResult()
	tr := newTree(f)

	// WHEN
	err := tr.run("reit", "purchase-balance", "MXRF11", "HGLG11", "--amount", "200", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	expectedCalls := []purchaseCall{{stocks: []string{}, reits: []string{"MXRF11", "HGLG11"}, amount: 200}}
	if !reflect.DeepEqual(f.purchase.calls, expectedCalls) {
		t.Errorf("expected calls %v and received %v", expectedCalls, f.purchase.calls)
	}
	expected := `TICKER,PRICE,COUNT,TOTAL,CURRENCY,CAPTURED AT
MXRF11,10.00,5,50.00,R$ BRL,2025-06-08 22:59:37
HGLG11,100.00,1,100.00,R$ BRL,2025-06-08 22:59:37
,,6,150.00,R$ BRL,SPENT AMOUNT
,,,50.00,R$ BRL,REMAINING AMOUNT
`
	if tr.stdout.String() != expected {
		t.Errorf("expected stdout %q and received %q", expected, tr.stdout.String())
	}
	if tr.stderr.Len() != 0 {
		t.Errorf("expected empty stderr and received %q", tr.stderr.String())
	}
}

func TestReitPurchaseBalanceNotFound(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run("reit", "purchase-balance", "HGLG11", "--amount", "10")

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

func TestReitPurchaseBalanceWithoutAmount(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.purchase.result = reitPurchaseResult()
	tr := newTree(f)

	// WHEN
	err := tr.run("reit", "purchase-balance", "MXRF11")

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

func TestReitNoColor(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.reit.get = func(ticker string) *resource.Security { return newReit(ticker, 10.5) }
	tr := newTree(f)

	// WHEN
	colored := tr.run("reit", "get", "MXRF11")
	coloredOut := tr.stdout.String()
	plain := tr.run("reit", "get", "MXRF11", "--no-color")
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
	if !strings.Contains(plainOut, "MXRF11") {
		t.Errorf("expected the table in the output and received %q", plainOut)
	}
}

func TestReitOptionDoesNotLeakToAnotherCommand(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.stock.list = listing(newStock("ITSA3", 10.5))
	f.reit.list = listing(newReit("MXRF11", 10))
	tr := newTree(f)

	// WHEN
	first := tr.run("reit", "list", "MXRF11", "--csv")
	firstOut := tr.stdout.String()
	second := tr.run("stock", "list", "ITSA3")
	secondOut := tr.stdout.String()

	// THEN
	if first != nil || second != nil {
		t.Fatalf("expected no errors, received %v and %v", first, second)
	}
	if !strings.HasPrefix(firstOut, "TICKER,NAME") {
		t.Errorf("expected csv in the first run and received %q", firstOut)
	}
	if strings.Contains(secondOut, "TICKER,NAME") || !strings.Contains(secondOut, "ITSA3") {
		t.Errorf("expected a table in the second run and received %q", secondOut)
	}
}
