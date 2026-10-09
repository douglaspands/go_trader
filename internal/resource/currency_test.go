package resource_test

import (
	"testing"
	"trader/internal/resource"
)

func TestCurrencyString(t *testing.T) {
	// GIVEN
	currency := &resource.Currency{Code: "BRL", Sign: "R$", Description: "Brazilian Real"}

	// WHEN
	result := currency.String()

	// THEN
	if result != "R$ BRL" {
		t.Errorf(`expected "R$ BRL" and received "%s"`, result)
	}
}
