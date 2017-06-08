package dbmdl

import "regexp"

var (
	regDefault = regexp.MustCompile("default .+")
	regExtern  = regexp.MustCompile("extern (.*) from (.*) on (.*)")
)