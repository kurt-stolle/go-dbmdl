package dbmdl

import "reflect"

// Constants
const (
	Version = "1.0.0"
)

// Privates
var dialects map[string]*Dialect
var tables map[reflect.Type]string

// RegisterDialect will add a dialect so that it can be used later
func RegisterDialect(d string, strct *Dialect) {
	dialects[d] = strct
}

// RegisterStruct registers a struct for use with dbmdl
func RegisterStruct(d string, t string, s ifc) {
	refType := reflect.TypeOf(s)
	tables[refType] = t
}
