package dbmdl

import "reflect"

// Constants
const (
	Version = "1.0.0"
)

// Privates
type table struct {
	dialect *Dialect
	name    string
}

var tables map[reflect.Type]*table
var dialects map[string]*Dialect

func init() {
	tables = make(map[reflect.Type]*table)
	dialects = make(map[string]*Dialect)
}
