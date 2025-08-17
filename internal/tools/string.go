package tools

import (
	"strconv"
	"strings"
	"unicode"
)

func ToSnakeCase(input string) string {
	var result strings.Builder
	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
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
