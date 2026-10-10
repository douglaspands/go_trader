package tools

import (
	"strconv"
	"strings"
	"unicode"
)

// CleanText removes every control character, which breaks any terminal escape sequence.
func CleanText(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, input)
}

// CSVSafe prefixes a single quote to text that a spreadsheet would evaluate as a formula.
func CSVSafe(input string) string {
	if input != "" && strings.ContainsRune("=+-@", rune(input[0])) {
		return "'" + input
	}
	return input
}

func SplitList(input string, sep string) []string {
	result := make([]string, 0)
	for _, part := range strings.Split(input, sep) {
		if part = strings.TrimSpace(part); part != "" {
			result = append(result, part)
		}
	}
	return result
}

func ToFloat(input string, decimalSeparator string) float64 {
	thousandSeparator := ","
	if strings.Contains(decimalSeparator, ",") {
		thousandSeparator = "."
	}
	value, err := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(input, thousandSeparator, ""), decimalSeparator, "."), 64)
	if err != nil {
		return 0.0
	}
	return value
}
