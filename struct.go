package dbmdl

import (
	"errors"
	"reflect"
)

// StructInterface for all dbmodels
type StructInterface interface {
	CreateTables() error
}

// Struct is a struct that can be inherited for use with dbmdl
type Struct struct{}

// CreateTables will register the struct in the database
func (s *Struct) CreateTables() error {
	var ref = reflect.TypeOf(s)
	var fields []string

	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)              // Get the field at index i
		dataType := field.Tag.Get("dbmdl") // Find the datatype from the dbmdl tag

		if dataType == "" {
			return errors.New("Failed to create tables for struct " + ref.Name() + ", field " + field.Name + " does not have a `dbmdl` tag")
		}

		fields = append(fields, field.Name+" "+dataType)
	}

	t := tables[ref]
	query(t.dialect.CreateTable(t.name, fields...))

	return nil
}
