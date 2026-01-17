package tools_test

import (
	"testing"
	"trader/internal/tools"
)

func TestTableRowValue(t *testing.T) {
	// GIVEN
	value := map[string]interface{}{
		"price": 10.50,
		"name":  "Test",
	}

	// THEN
	result := tools.TableRowValue(value)

	// WHEN
	if result.(map[string]interface{})["price"] != 10.50 {
		t.Errorf(`expected at %f and received at %f`, 10.50, result.(map[string]interface{})["price"])
	}

	if result.(map[string]interface{})["name"] != "Test" {
		t.Errorf(`expected at "%s" and received at %s`, "Test", result.(map[string]interface{})["name"])
	}
}
