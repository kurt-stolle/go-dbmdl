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
