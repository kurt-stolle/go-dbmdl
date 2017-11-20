package dbmdl

import (
	"reflect"
	"strings"
)

func getTagParameters(rawTag string) []string {
	tag := strings.Split(rawTag, ",")
	for i, v := range tag {
		tag[i] = strings.Trim(v, " \t\n")
	}

	return tag
}