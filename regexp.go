package dbmdl

import "regexp"

var (
	regDefault = regexp.MustCompile("^default .+$")
)
