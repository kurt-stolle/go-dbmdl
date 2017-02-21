package dbmdl

import (
	"reflect"
	"strings"
)

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
