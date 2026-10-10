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

func TestCleanText(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect string
	}{
		{"plain text", "PETROBRAS S.A.", "PETROBRAS S.A."},
		{"accents", "Híbrido Logística", "Híbrido Logística"},
		{"empty", "", ""},
		{"clear screen sequence", "A\x1b[2JB", "A[2JB"},
		{"color sequence", "\x1b[31mRED\x1b[0m", "[31mRED[0m"},
		{"bell", "A\aB", "AB"},
		{"newline and tab", "A\nB\tC\r", "ABC"},
		{"C1 control", "A\u009b2JB", "A2JB"},
		{"delete", "A\x7fB", "AB"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// WHEN
			result := tools.CleanText(c.input)

			// THEN
			if result != c.expect {
				t.Errorf("expected %q and received %q", c.expect, result)
			}
		})
	}
}

func TestCSVSafe(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect string
	}{
		{"equals", `=HYPERLINK("http://x")`, `'=HYPERLINK("http://x")`},
		{"plus", "+1", "'+1"},
		{"minus", "-1", "'-1"},
		{"at", "@SUM(A1)", "'@SUM(A1)"},
		{"plain text", "PETROBRAS", "PETROBRAS"},
		{"prefix char in the middle", "A=B", "A=B"},
		{"document", "33.000.167/0001-01", "33.000.167/0001-01"},
		{"empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// WHEN
			result := tools.CSVSafe(c.input)

			// THEN
			if result != c.expect {
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
