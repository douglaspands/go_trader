package tools_test

import (
	"testing"
	"time"
	"trader/internal/tools"
)

func TestTableRowValueFloatFormat(t *testing.T) {

	// EXPECTED
	expected := "10.50"

	// GIVEN
	price := 10.50

	// THEN
	result := tools.TableRowValue(price)

	// WHEN
	if result.(string) != expected {
		t.Errorf(`expected at %s and received at %s`, expected, result.(string))
	}
}

func TestTableRowValueTimeFormat(t *testing.T) {

	// EXPECTED
	expected := "2022-01-01 00:00:00"

	// GIVEN
	date := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	// THEN
	result := tools.TableRowValue(date)

	// WHEN
	if result.(string) != expected {
		t.Errorf(`expected at %s and received at %s`, expected, result.(string))
	}
}

func TestTableRowValueStringFormat(t *testing.T) {

	// EXPECTED
	expected := "Test"

	// GIVEN
	value := "Test"

	// THEN
	result := tools.TableRowValue(value)

	// WHEN
	if result.(string) != expected {
		t.Errorf(`expected at %s and received at %s`, expected, result.(string))
	}
}

func TestTableRowValueIntFormat(t *testing.T) {

	// EXPECTED
	expected := 10

	// GIVEN
	value := 10

	// THEN
	result := tools.TableRowValue(value)

	// WHEN
	if result.(int) != expected {
		t.Errorf(`expected at %d and received at %d`, expected, result.(int))
	}
}
