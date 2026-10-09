package tools

import (
	"fmt"
	"time"
)

func TableRowValue(value interface{}) interface{} {
	var result interface{}
	switch v := value.(type) {
	case time.Time:
		result = v.Format("2006-01-02 15:04:05")
	case float64:
		result = fmt.Sprintf("%.2f", v)
	default:
		result = v
	}
	return result
}
