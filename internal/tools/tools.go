package tools

import (
	"reflect"
)

func StructToMap(s interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)
			if fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct {
				result[field.Name] = StructToMap(fieldValue.Interface())
			} else {
				result[field.Name] = fieldValue.Interface()
			}
		}
	}
	return result
}
