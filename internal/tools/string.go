package tools

import (
	"strconv"
	"strings"
)

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
