package dbmdl

import (
	"errors"
	"reflect"
)

// StructInterface for all dbmodels
type StructInterface interface {
	CreateTables(ref reflect.Type) error
}

// Struct is a struct that can be inherited for use with dbmdl
type Struct struct{}

// CreateTables will register the struct in the database
func (s *Struct) CreateTables(ref reflect.Type) error {
	var t, ok = tables[ref]
	if !ok {
		return errors.New("[dbmdl] Type not in tables map: " + ref.Name())
	}

	// Build fields list
	var fields []string

	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)              // Get the field at index i
		dataType := field.Tag.Get("dbmdl") // Find the datatype from the dbmdl tag

		if dataType == "" {
			continue
		}

		fields = append(fields, field.Name+" "+dataType)
	}

	// Query
	q := t.dialect.CreateTable(t.name, fields...)
	query(q)

	return nil
}
