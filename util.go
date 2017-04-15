package dbmdl

import (
	"reflect"
	"strings"
)

func getFields(targetType reflect.Type) []string {
	var fields []string
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i) // Get the field at index i
		if field.Tag.Get("dbmdl") == "" {
			continue
		}

		for _, tag := range getTagParameters(field) {
			if tag == omit {
				continue
			}
		}

		fields = append(fields, field.Name)
	}

	return fields
}

func getTagParameters(field reflect.StructField) []string {
	tag := strings.Split(field.Tag.Get("dbmdl"), ",")
	for i, v := range tag {
		tag[i] = strings.Trim(v, " \t\n")
	}

	return tag
}

func getReflectType(ifc interface{}) reflect.Type {
	switch v := ifc.(type) {
	case reflect.Type:
		if v.Kind() == reflect.Ptr {
			return v.Elem()
		}

		return v
	default:
		return getReflectType(reflect.TypeOf(v))
	}
}
