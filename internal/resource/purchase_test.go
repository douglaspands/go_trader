package resource_test

import (
	"testing"
	"trader/internal/resource"
)

func TestSecurityPurchaseTotalAmount(t *testing.T) {
	// GIVEN
	purchase := &resource.SecurityPurchase{
		Security: &resource.Security{Ticker: "S1", Price: 12.5},
		Count:    4,
	}

	// WHEN
	result := purchase.TotalAmount()

	// THEN
	if result != 50.0 {
		t.Errorf("expected total amount 50.00 and received %.2f", result)
	}
}

func TestPurchaseBalanceTotals(t *testing.T) {
	// GIVEN
	balance := &resource.PurchaseBalance{
		SecuritiesBalance: []*resource.SecurityPurchase{
			{Security: &resource.Security{Ticker: "S1", Price: 10}, Count: 3},
			{Security: &resource.Security{Ticker: "S2", Price: 25}, Count: 2},
		},
		AmountInvested: 100,
	}

	// THEN
	if result := balance.TotalCount(); result != 5 {
		t.Errorf("expected total count 5 and received %d", result)
	}
	if result := balance.AmountSpent(); result != 80.0 {
		t.Errorf("expected amount spent 80.00 and received %.2f", result)
	}
	if result := balance.RemainingBalance(); result != 20.0 {
		t.Errorf("expected remaining balance 20.00 and received %.2f", result)
	}
}

func TestPurchaseBalanceEmpty(t *testing.T) {
	// GIVEN
	balance := &resource.PurchaseBalance{AmountInvested: 100}

	// THEN
	if result := balance.TotalCount(); result != 0 {
		t.Errorf("expected total count 0 and received %d", result)
	}
	if result := balance.AmountSpent(); result != 0.0 {
		t.Errorf("expected amount spent 0.00 and received %.2f", result)
	}
	if result := balance.RemainingBalance(); result != 100.0 {
		t.Errorf("expected remaining balance 100.00 and received %.2f", result)
	}
}
