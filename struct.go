package dbmdl

import (
	"errors"
	"reflect"
)

// createTables will register the struct in the database
func createTables(ref reflect.Type) error {
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
	q := t.dialect.CreateTable(t.name, fields)
	query(nil, q)

	return nil
}

//
