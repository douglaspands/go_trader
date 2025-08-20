package tools_test

import (
	"testing"
	"trader/internal/tools"
)

func TestToSnakeCaseOk01(t *testing.T) {
	// EXPECT
	expect_text := "get_cheese"

	// GIVEN
	text := "getCheese"

	// THEN
	result := tools.ToSnakeCase(text)

	// WHEN
	if result != expect_text {
		t.Errorf(`expected at "%s" and received at %s`, expect_text, result)
	}
}

func TestToSnakeCaseOk02(t *testing.T) {
	// EXPECT
	expect_text := "get_cheese"

	// GIVEN
	text := "Get Cheese"

	// THEN
	result := tools.ToSnakeCase(text)

	// WHEN
	if result != expect_text {
		t.Errorf(`expected at "%s" and received at %s`, expect_text, result)
	}
}

func TestToSnakeCaseOk03(t *testing.T) {
	// EXPECT
	expect_text := "get_cheese_01_02"

	// GIVEN
	text := " Get Cheese01 02"

	// THEN
	result := tools.ToSnakeCase(text)

	// WHEN
	if result != expect_text {
		t.Errorf(`expected at "%s" and received at %s`, expect_text, result)
	}
}

func TestToFloatOk01(t *testing.T) {
	// EXPECT
	expect_value := 1.2

	// GIVEN
	value := "1.20"

	// THEN
	result := tools.ToFloat(value, ".")

	// WHEN
	if result != expect_value {
		t.Errorf(`expected at "%f" and received at %f`, expect_value, result)
	}
}

func TestToFloatOk02(t *testing.T) {
	// EXPECT
	expect_value := 1.2

	// GIVEN
	value := "1,20"

	// THEN
	result := tools.ToFloat(value, ",")

	// WHEN
	if result != expect_value {
		t.Errorf(`expected at "%f" and received at %f`, expect_value, result)
	}
}

func TestToFloatOk03(t *testing.T) {
	// EXPECT
	expect_value := 0.0

	// GIVEN
	value := "1@20"

	// THEN
	result := tools.ToFloat(value, ",")

	// WHEN
	if result != expect_value {
		t.Errorf(`expected at "%f" and received at %f`, expect_value, result)
	}
}
