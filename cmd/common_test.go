package cmd_test

import (
	"strings"
	"testing"
	"time"
)

func TestPurchaseBalanceInvalidAmount(t *testing.T) {
	commands := [][]string{
		{"stock", "purchase-balance", "PETR4"},
		{"reit", "purchase-balance", "MXRF11"},
		{"security", "purchase-balance", "-s", "PETR4"},
	}
	for _, amount := range []string{"inf", "NaN", "0", "-100", "1e13"} {
		for _, command := range commands {
			t.Run(command[0]+" "+amount, func(t *testing.T) {
				// GIVEN
				f := newFakes()
				f.purchase.result = stockPurchaseResult()
				tr := newTree(f)
				args := append(append([]string{}, command...), "--amount", amount)

				// WHEN
				done := make(chan error, 1)
				go func() { done <- tr.run(args...) }()
				var err error
				select {
				case err = <-done:
				case <-time.After(5 * time.Second):
					t.Fatal("the command did not finish")
				}

				// THEN
				if err == nil {
					t.Fatal("expected an error, received nil")
				}
				if !strings.Contains(tr.stderr.String(), "Error: invalid amount") {
					t.Errorf("expected an invalid amount error in stderr and received %q", tr.stderr.String())
				}
				if strings.Contains(tr.stdout.String(), "TICKER") {
					t.Errorf("expected no table and received %q", tr.stdout.String())
				}
				if len(f.purchase.calls) != 0 {
					t.Errorf("expected no purchase balance call and received %v", f.purchase.calls)
				}
			})
		}
	}
}

func TestPurchaseBalanceAmountAtCeiling(t *testing.T) {
	// GIVEN
	f := newFakes()
	f.purchase.result = stockPurchaseResult()
	tr := newTree(f)

	// WHEN
	err := tr.run("stock", "purchase-balance", "ITSA3", "--amount", "1e12", "--csv")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if len(f.purchase.calls) != 1 || f.purchase.calls[0].amount != 1e12 {
		t.Errorf("expected one call with amount 1e12 and received %v", f.purchase.calls)
	}
}
