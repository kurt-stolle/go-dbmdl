package dbmdl

import "reflect"

// The table type
type table struct {
	dialect *Dialect
	name    string
}

// Maps for storing data about the server environment
var tables map[reflect.Type]*table
var dialects map[string]*Dialect

// Initialize maps
func init() {
	tables = make(map[reflect.Type]*table)
	dialects = make(map[string]*Dialect)
}
