package tools_test

import (
	"reflect"
	"testing"
	"trader/internal/tools"
)

func TestSplitList(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect []string
	}{
		{"empty", "", []string{}},
		{"single", "A", []string{"A"}},
		{"two", "A,B", []string{"A", "B"}},
		{"spaces and empty entries", " A , ,B ", []string{"A", "B"}},
		{"only separator", ",", []string{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// WHEN
			result := tools.SplitList(c.input, ",")

			// THEN
			if !reflect.DeepEqual(result, c.expect) {
				t.Errorf("expected %q and received %q", c.expect, result)
			}
		})
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
