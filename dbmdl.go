package dbmdl

import (
	"errors"
	"log"
	"reflect"
)

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

// RegisterDialect will add a dialect so that it can be used later
func RegisterDialect(d string, strct *Dialect) error {
	dialects[d] = strct

	return nil
}

// RegisterStruct registers a struct for use with dbmdl
func RegisterStruct(dlct string, t string, s StructInterface) error {
	d, ok := dialects[dlct]
	if !ok {
		return errors.New("Failed to register struct; dialect unknown!")
	}

	refType := reflect.TypeOf(s)
	tables[refType] = &table{dialect: d, name: t}

	log.Println("[dbmdl] Registered struct: " + refType.Name())

	s.CreateTables()

	return nil
}
