package cmd_test

import (
	"reflect"
	"strings"
	"testing"
	"trader/internal/resource"
)

func securityPurchaseResult() *resource.PurchaseBalance {
	return &resource.PurchaseBalance{
		SecuritiesBalance: []*resource.SecurityPurchase{
			{Security: newStock("ITSA3", 10), Count: 5},
			{Security: newReit("MXRF11", 10), Count: 4},
		},
		AmountInvested: 100,
	}
}

func TestSecurityPurchaseBalanceOnlyStocks(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.purchase.result = securityPurchaseResult()
	tr := newTree(f)

	// WHEN
	err := tr.run("security", "purchase-balance", "--stocks", "ITSA3, BBDC3,", "--amount", "100", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	expectedCalls := []purchaseCall{{stocks: []string{"ITSA3", "BBDC3"}, reits: []string{}, amount: 100}}
	if !reflect.DeepEqual(f.purchase.calls, expectedCalls) {
		t.Errorf("expected calls %v and received %v", expectedCalls, f.purchase.calls)
	}
	if tr.stderr.Len() != 0 {
		t.Errorf("expected empty stderr and received %q", tr.stderr.String())
	}
}

func TestSecurityPurchaseBalanceStocksAndReits(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.purchase.result = securityPurchaseResult()
	tr := newTree(f)

	// WHEN
	err := tr.run("security", "purchase-balance", "-s", "ITSA3", "-r", "MXRF11", "-a", "100", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	expectedCalls := []purchaseCall{{stocks: []string{"ITSA3"}, reits: []string{"MXRF11"}, amount: 100}}
	if !reflect.DeepEqual(f.purchase.calls, expectedCalls) {
		t.Errorf("expected calls %v and received %v", expectedCalls, f.purchase.calls)
	}
	expected := `TICKER,TYPE,PRICE,COUNT,TOTAL,CURRENCY,CAPTURED AT
ITSA3,Stock,10.00,5,50.00,R$ BRL,2025-06-08 22:59:37
MXRF11,REIT,10.00,4,40.00,R$ BRL,2025-06-08 22:59:37
,,,9,90.00,R$ BRL,SPENT AMOUNT
,,,,10.00,R$ BRL,REMAINING AMOUNT
`
	if tr.stdout.String() != expected {
		t.Errorf("expected stdout %q and received %q", expected, tr.stdout.String())
	}
	if tr.stderr.Len() != 0 {
		t.Errorf("expected empty stderr and received %q", tr.stderr.String())
	}
}

func TestSecurityPurchaseBalanceWithoutAmount(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.purchase.result = securityPurchaseResult()
	tr := newTree(f)

	// WHEN
	err := tr.run("security", "purchase-balance", "-s", "ITSA3")

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

func TestSecurityPurchaseBalanceNothingPurchased(t *testing.T) {
	// GIVEN
	tr := newTree(newFakes())

	// WHEN
	err := tr.run("security", "purchase-balance", "-s", "ITSA3", "-a", "1")

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

func TestSecurityPurchaseBalanceNoColor(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.purchase.result = securityPurchaseResult()
	tr := newTree(f)

	// WHEN
	err := tr.run("security", "purchase-balance", "-s", "ITSA3", "-a", "100", "--no-color")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if strings.Contains(tr.stdout.String(), "\x1b[") {
		t.Errorf("expected no ANSI codes and received %q", tr.stdout.String())
	}
	if !strings.Contains(tr.stdout.String(), "REMAINING AMOUNT") {
		t.Errorf("expected the table in the output and received %q", tr.stdout.String())
	}
}
