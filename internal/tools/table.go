package tools

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

func StructToTableHeader(s interface{}, ignore []string) table.Row {
	header := make([]interface{}, 0)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if !slices.Contains(ignore, field.Name) {
				header = append(header, strings.ToUpper(ToSnakeCase(field.Name)))
			}
		}
	}
	return header
}

func StructToTableRow(s interface{}, ignore []string) table.Row {
	row := make([]interface{}, 0)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if !slices.Contains(ignore, field.Name) {
				fieldValue := v.Field(i)
				row = append(row, tableRowValue(fieldValue.Interface()))
			}
		}
	}
	return row
}

func StructToTableRowsFieldValue(s interface{}, ignore []string) []table.Row {
	rows := make([]table.Row, 0)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			row := make([]interface{}, 0)
			field := t.Field(i)
			if !slices.Contains(ignore, field.Name) {
				fieldValue := v.Field(i)
				row = append(row, strings.ToUpper(ToSnakeCase(field.Name)), tableRowValue(fieldValue.Interface()))
			}
			if len(row) > 0 {
				rows = append(rows, row)
			}
		}
	}
	return rows
}

func tableRowValue(value interface{}) interface{} {
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
