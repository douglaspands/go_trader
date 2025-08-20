package tools

import (
	"strconv"
	"strings"
	"unicode"
)

func ToSnakeCase(text string) string {
	var result strings.Builder
	var char rune
	for i, r := range text {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				if i > 0 && char != 0 && char != '_' {
					result.WriteRune('_')
				}
				char = unicode.ToLower(r)
			} else {
				char = r
			}
			result.WriteRune(char)
		} else {
			if unicode.IsNumber(r) {
				if unicode.IsLetter(char) {
					result.WriteRune('_')
				}
				char = r
				result.WriteRune(char)
			} else {
				if i > 0 && char != 0 && char != '_' {
					char = '_'
					result.WriteRune(char)
				}
			}
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
